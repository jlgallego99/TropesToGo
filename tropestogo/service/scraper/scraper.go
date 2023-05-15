package scraper

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/index"
	"github.com/jlgallego99/TropesToGo/media"
)

var (
	ErrInvalidField         = errors.New("one or more fields for the Scraper are invalid")
	ErrNotTvTropes          = errors.New("the URL does not belong to a TvTropes page")
	ErrNotWorkPage          = errors.New("the page isn't a TvTropes Work page")
	ErrUnknownPageStructure = errors.New("the scraper doesn't recognize the page structure")
)

const (
	TvTropesHostname        = "tvtropes.org"
	TvTropesPmwiki          = "/pmwiki/pmwiki.php/"
	TvTropesMainPath        = TvTropesPmwiki + "Main/"
	WorkTitleSelector       = "h1.entry-title"
	WorkIndexSelector       = WorkTitleSelector + " strong"
	MainArticleSelector     = "#main-article"
	TropeListSelector       = "#main-article ul"
	TropeListHeaderSelector = "#main-article h2"
	SubPagesNavSelector     = "nav.body-options"
	SubPageListSelector     = "ul.subpage-links"
	SubPageLinkSelector     = "a.subpage-link"
	TropeTag                = "a.twikilink"
	TropeLinkSelector       = "#main-article ul li " + TropeTag
	TropeFolderSelector     = "#main-article div.folderlabel"
	FolderToggleFunction    = "toggleAllFolders()"
)

// ScraperConfig is an alias for a function that will accept a pointer to a ServiceScraper and modify its fields
// Each function acts as one configuration for the scraper
type ScraperConfig func(ss *ServiceScraper) error

type ServiceScraper struct {
	// TvTropes index
	index index.RepositoryIndex
	// TvTropes dataset
	data media.RepositoryMedia
}

// NewServiceScraper takes a variable amount of configuration functions and returns a ServiceScraper with all configs passed
func NewServiceScraper(cfgs ...ScraperConfig) (*ServiceScraper, error) {
	ss := &ServiceScraper{}
	// Apply all config functions
	for _, cfg := range cfgs {
		// Configure the service we are creating
		err := cfg(ss)
		if err != nil {
			return nil, err
		}
	}

	return ss, nil
}

func ConfigIndexRepository(ir index.RepositoryIndex) ScraperConfig {
	return func(ss *ServiceScraper) error {
		if ir == nil {
			return ErrInvalidField
		}

		ss.index = ir
		return nil
	}
}

func ConfigRepository(mr media.RepositoryMedia) ScraperConfig {
	return func(ss *ServiceScraper) error {
		if mr == nil {
			return ErrInvalidField
		}

		ss.data = mr
		return nil
	}
}

// CheckValidWorkPage checks if a TvTropes Work page has a valid structure in which the scraper can extract data
// This allows the scraper to check if TvTropes template has somewhat changed
func (scraper *ServiceScraper) CheckValidWorkPage(page *tropestogo.Page) (bool, error) {
	res, _ := http.Get(page.URL.String())
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	validWorkPage, errWorkPage := CheckTvTropesWorkPage(page)
	if !validWorkPage {
		return false, errWorkPage
	}

	validMainArticle, errMainArticle := CheckMainArticle(doc)
	if !validMainArticle {
		return false, errMainArticle
	}

	validTropeSection, errTropeSection := CheckTropeSection(doc)
	if !validTropeSection {
		return false, errTropeSection
	}

	return true, nil
}

// CheckTvTropesWorkPage checks if the received page belongs to a tvtropes.org Work page
func CheckTvTropesWorkPage(page *tropestogo.Page) (bool, error) {
	// First check if the domain is TvTropes
	if page.URL.Hostname() != TvTropesHostname {
		return false, ErrNotTvTropes
	}

	// Check if it's a Film Work page
	splitPath := strings.Split(page.URL.Path, "/")
	if !strings.HasPrefix(page.URL.Path, TvTropesPmwiki) || splitPath[3] != media.Film.String() {
		return false, ErrNotWorkPage
	}

	return true, nil
}

// CheckMainArticle checks if the tvtropes work main article page has a known structure which can be extracted later
func CheckMainArticle(doc *goquery.Document) (bool, error) {
	// Check if the main article structure has all known ids and elements that comprise a TvTropes work page
	if doc.Find(MainArticleSelector).Length() == 0 ||
		doc.Find(SubPagesNavSelector).Find(SubPageListSelector).Find(SubPageLinkSelector).Length() == 0 {
		return false, ErrUnknownPageStructure
	}

	// Check the index of the work
	index := doc.Find(WorkIndexSelector)
	if strings.Trim(index.Text(), " /") != media.Film.String() {
		return false, ErrNotWorkPage
	}

	return true, nil
}

func CheckTropeSection(doc *goquery.Document) (bool, error) {
	// Look for the tropes section
	if doc.Find(TropeListSelector).Length() != 0 {
		// Check if the list is a simple trope list or a list of subpages with tropes
		if doc.Find(TropeLinkSelector).Length() != 0 {
			// Get the first word of the first element of the list, check if is an anchor to a trope page or a sub page
			tropeHref, exists := doc.Find(TropeLinkSelector).First().Attr("href")

			// Checks if the elements of the list are anchors to a subpage inside the work
			// A regex matches if the last part of the URL is of the type TropesXtoY
			hrefSplit := strings.Split(tropeHref, "/")
			r, _ := regexp.Compile("Tropes[A-Z]To[A-Z]")
			match := r.MatchString(hrefSplit[len(hrefSplit)-1])

			// Check if the trope link directs to a Main trope page or a subpage with tropes
			if exists && strings.HasPrefix(tropeHref, TvTropesPmwiki) {
				// If the page isn't a main trope page and doesn't match the subpage regex, we don't know what that is
				if !strings.HasPrefix(tropeHref, TvTropesMainPath) && !match {
					return false, ErrUnknownPageStructure
				}

				return true, nil
			}
		}

		// Check if tropes are on folders
		// If there's a close all folders button, then the tropes are on folders
		folderFunctionName, existsFolderButton := doc.Find(TropeFolderSelector).Attr("onclick")
		if existsFolderButton && folderFunctionName == FolderToggleFunction {
			return true, nil
		}
	}

	// Tropes are presented in an unknown form, so data can't be extracted
	return false, ErrUnknownPageStructure
}
