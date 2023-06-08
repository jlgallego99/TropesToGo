package scraper

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

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
	TvTropesHostname         = "tvtropes.org"
	TvTropesPmwiki           = "/pmwiki/pmwiki.php/"
	TvTropesMainPath         = TvTropesPmwiki + "Main/"
	WorkTitleSelector        = "h1.entry-title"
	WorkIndexSelector        = WorkTitleSelector + " strong"
	MainArticleSelector      = "#main-article"
	TropeListSelector        = MainArticleSelector + " ul"
	TropeListHeaderSelector  = MainArticleSelector + " h2"
	SubPagesNavSelector      = "nav.body-options"
	SubPageListSelector      = "ul.subpage-links"
	SubPageLinkSelector      = "a.subpage-link"
	TropeTag                 = "a.twikilink"
	TropeLinkSelector        = "#main-article ul li " + TropeTag
	TropeFolderSelector      = MainArticleSelector + " div.folderlabel"
	FolderToggleFunction     = "toggleAllFolders();"
	MainTropesSelector       = TropeListHeaderSelector + " ~ ul > li > " + TropeTag + ":first-child"
	MainTropesFolderSelector = MainArticleSelector + " .folder > ul > li > " + TropeTag + ":first-child"
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

func ConfigMediaRepository(mr media.RepositoryMedia) ScraperConfig {
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

	validWorkPage, errWorkPage := checkTvTropesWorkPage(page)
	if !validWorkPage {
		return false, errWorkPage
	}

	validMainArticle, errMainArticle := checkMainArticle(doc)
	if !validMainArticle {
		return false, errMainArticle
	}

	validTropeSection, errTropeSection := checkTropeSection(doc)
	if !validTropeSection {
		return false, errTropeSection
	}

	return true, nil
}

// CheckTvTropesWorkPage checks if the received page belongs to a tvtropes.org Work page
func checkTvTropesWorkPage(page *tropestogo.Page) (bool, error) {
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
func checkMainArticle(doc *goquery.Document) (bool, error) {
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

func checkTropeSection(doc *goquery.Document) (bool, error) {
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
		if checkTropesOnFolders(doc) {
			return true, nil
		}
	}

	// Tropes are presented in an unknown form, so data can't be extracted
	return false, ErrUnknownPageStructure
}

// checkTropesOnFolders validates whether tropes are presented on folders
func checkTropesOnFolders(doc *goquery.Document) bool {
	// If there's a close all folders button, then the tropes are on folders
	folderFunctionName, existsFolderButton := doc.Find(TropeFolderSelector).Attr("onclick")
	if existsFolderButton && folderFunctionName == FolderToggleFunction {
		return true
	}

	return false
}

// ScrapeWorkPage extracts all the relevant information from a TvTropes Work Page
func (scraper *ServiceScraper) ScrapeWorkPage(page *tropestogo.Page) (media.Media, error) {
	res, _ := http.Get(page.URL.String())
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Url = page.URL

	title, year, mediaIndex, errMediaIndex := scraper.ScrapeWorkTitleAndYear(doc)
	if errMediaIndex != nil {
		return media.Media{}, errMediaIndex
	}

	tropes, errTropes := scraper.ScrapeWorkTropes(doc)
	if errTropes != nil {
		return media.Media{}, errTropes
	}

	media, error := media.NewMedia(title, year, time.Now(), tropes, page, mediaIndex)
	return media, error
}

// ScrapeWorkTitleAndYear extracts the title, the year on the title/URL if it's there and the media index from the HTML document of a Work Page
func (scraper *ServiceScraper) ScrapeWorkTitleAndYear(doc *goquery.Document) (string, string, media.MediaType, error) {
	var title, year string

	// Get Title of the Work page by extracting the title and discarding the index name
	title = strings.TrimSpace(strings.Split(doc.Find(WorkTitleSelector).Text(), "/")[1])
	mediaIndex, errMediaIndex := media.ToMediaType(strings.Trim(doc.Find(WorkIndexSelector).First().Text(), " /"))

	// Search for the year of the Work in the URI
	r, _ := regexp.Compile(`(19|20)\d{2}`)
	matchedString := r.FindStringSubmatch(doc.Url.String())
	if len(matchedString) > 0 {
		year = matchedString[0]
	}

	return title, year, mediaIndex, errMediaIndex
}

// ScrapeWorkTropes extracts all the tropes from the HTML document of a Work Page
func (scraper *ServiceScraper) ScrapeWorkTropes(doc *goquery.Document) (map[tropestogo.Trope]struct{}, error) {
	tropes := make(map[tropestogo.Trope]struct{}, 0)
	var newTrope tropestogo.Trope
	var newTropeError error

	// Use a different selector if it's a simple list or folders
	var selector string
	if checkTropesOnFolders(doc) {
		selector = MainTropesFolderSelector
	} else {
		selector = MainTropesSelector
	}

	// Searches for all trope tags on the main list
	// Extracts only the main tropes on the list, that is, the ones that start an element on the list
	// Ignores referenced tropes on the description as those may not belong to the Work
	doc.Find(selector).EachWithBreak(func(_ int, selection *goquery.Selection) bool {
		// Get trope name from the last part of the URI
		tropeUri, tropeUriExists := selection.Attr("href")
		if tropeUriExists {
			// Insert tropes without repeating
			newTrope, newTropeError = tropestogo.NewTrope(strings.Split(tropeUri, "/")[4], tropestogo.TropeIndex(0))
			if newTropeError != nil {
				return false
			}

			// Check if it's really a trope by checking its URI
			if tropeUri == TvTropesMainPath+newTrope.GetTitle() {
				// Add trope to the set. If the trope is already there then it's ignored
				tropes[newTrope] = struct{}{}
			}
		}

		return true
	})

	if newTropeError != nil {
		return make(map[tropestogo.Trope]struct{}), newTropeError
	}

	return tropes, nil
}
