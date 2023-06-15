package scraper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	ErrNotFound             = errors.New("couldn't request the URL")
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

// ServiceScraper manages the TropesToGo scraper for checking TvTropes pages, extracting/cleaning its information
// and persisting the data on a RepositoryMedia
type ServiceScraper struct {
	// TvTropes index
	index index.RepositoryIndex
	// TvTropes dataset
	data media.RepositoryMedia
}

// NewServiceScraper takes a variable amount of configuration functions, applies them and returns a ServiceScraper with all configs passed
func NewServiceScraper(cfgs ...ScraperConfig) (*ServiceScraper, error) {
	ss := &ServiceScraper{}
	for _, cfg := range cfgs {
		err := cfg(ss)
		if err != nil {
			return nil, err
		}
	}

	return ss, nil
}

// ConfigIndexRepository defines a function that applies a RepositoryIndex so it can be used as a config when creating a ServiceScraper
func ConfigIndexRepository(ir index.RepositoryIndex) ScraperConfig {
	return func(ss *ServiceScraper) error {
		if ir == nil {
			return ErrInvalidField
		}

		ss.index = ir
		return nil
	}
}

// ConfigMediaRepository defines a function that applies a RepositoryMedia so it can be used as a config when creating a ServiceScraper
// It accepts any implementation of a RepositoryMedia, so it can accept either a CSVRepository or a JSONRepository for defining the persistence model
func ConfigMediaRepository(mr media.RepositoryMedia) ScraperConfig {
	return func(ss *ServiceScraper) error {
		if mr == nil {
			return ErrInvalidField
		}

		ss.data = mr
		return nil
	}
}

// CheckTvTropesPage makes an HTTP request to a TvTropes web page and checks if it's valid for scraping
// If the page doesn't exist, it returns an ErrNotFound error
// It returns true if all checks passes
func (scraper *ServiceScraper) CheckTvTropesPage(page tropestogo.Page) (bool, error) {
	res, err := http.Get(page.GetUrl().String())
	if err != nil {
		return false, fmt.Errorf("%w: "+page.GetUrl().String(), ErrNotFound)
	}

	return scraper.CheckValidWorkPage(res.Body, page.GetUrl())
}

// CheckValidWorkPage accepts any reader (whether it's from an URL or a local HTML file) with a webpage contents
// and checks if it's a valid TvTropes Work page
// This allows the scraper to check if TvTropes template has somewhat changed and if the scraper can extract its data
// It full-checks a TvTropes page, validating if its url is one of a TvTropes Work Page,
// if it's main article has a known structure that can be scraped and if trope section can also be scraped and contains valid tropes
// It returns true if all checks passes
func (scraper *ServiceScraper) CheckValidWorkPage(reader io.Reader, url *url.URL) (bool, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return false, ErrNotTvTropes
	}

	validWorkPage, errWorkPage := scraper.CheckIsWorkPage(url)
	if !validWorkPage {
		return false, errWorkPage
	}

	validMainArticle, errMainArticle := scraper.CheckMainArticle(doc)
	if !validMainArticle {
		return false, errMainArticle
	}

	validTropeSection, errTropeSection := scraper.CheckTropeSection(doc)
	if !validTropeSection {
		return false, errTropeSection
	}

	return true, nil
}

// CheckIsWorkPage checks if the received url belongs to a tvtropes.org Work Page
// It returns an ErrNotTvTropes error if it doesn't belong to TvTropes or a ErrNotWorkPage if it's from TvTropes but of any other type
// It returns true if all checks passes
func (scraper *ServiceScraper) CheckIsWorkPage(url *url.URL) (bool, error) {
	if url.Hostname() != TvTropesHostname {
		return false, fmt.Errorf("%w: "+url.String(), ErrNotTvTropes)
	}

	splitPath := strings.Split(url.Path, "/")
	if !strings.HasPrefix(url.Path, TvTropesPmwiki) || splitPath[3] != media.Film.String() {
		return false, fmt.Errorf("%w: "+url.String(), ErrNotWorkPage)
	}

	return true, nil
}

// CheckMainArticle checks the received goquery Document DOM Tree
// and validates if the work main article page has a known structure which can be extracted later
// If it doesn't recognize something in it, it returns an ErrUnknownPageStructure or an ErrNotWorkPage, so it can't be scraped
// It returns true if all checks passes
func (scraper *ServiceScraper) CheckMainArticle(doc *goquery.Document) (bool, error) {
	if doc.Find(MainArticleSelector).Length() == 0 ||
		doc.Find(SubPagesNavSelector).Find(SubPageListSelector).Find(SubPageLinkSelector).Length() == 0 {
		return false, ErrUnknownPageStructure
	}

	tropeIndex := doc.Find(WorkIndexSelector)
	if strings.Trim(tropeIndex.Text(), " /") != media.Film.String() {
		return false, fmt.Errorf("%w: "+doc.Url.String(), ErrNotWorkPage)
	}

	return true, nil
}

// CheckTropeSection checks the received goquery Document DOM Tree
// and validates if the tropes on the TvTropes Work Page are arranged in a known way that can be scraped
// If the trope section isn't of the recognized types, it returns an ErrUnknownPageStructure, so it can't be scraped
// It returns true if all checks passes
func (scraper *ServiceScraper) CheckTropeSection(doc *goquery.Document) (bool, error) {
	if doc.Find(TropeListSelector).Length() != 0 {
		if doc.Find(TropeLinkSelector).Length() != 0 {
			tropeHref, exists := doc.Find(TropeLinkSelector).First().Attr("href")

			hrefSplit := strings.Split(tropeHref, "/")
			r, _ := regexp.Compile("Tropes[A-Z]To[A-Z]")
			match := r.MatchString(hrefSplit[len(hrefSplit)-1])

			if exists && strings.HasPrefix(tropeHref, TvTropesPmwiki) {
				if !strings.HasPrefix(tropeHref, TvTropesMainPath) && !match {
					return false, ErrUnknownPageStructure
				}

				return true, nil
			}
		}

		// Check if tropes are on folders
		if scraper.CheckTropesOnFolders(doc) {
			return true, nil
		}
	}

	return false, ErrUnknownPageStructure
}

// CheckTropesOnFolders checks the received goquery Document DOM Tree and validates whether tropes are presented on folders
// It returns true if all checks passes
func (scraper *ServiceScraper) CheckTropesOnFolders(doc *goquery.Document) bool {
	folderFunctionName, existsFolderButton := doc.Find(TropeFolderSelector).Attr("onclick")
	if existsFolderButton && folderFunctionName == FolderToggleFunction {
		return true
	}

	return false
}

// ScrapeTvTropesPage makes an HTTP request to a TvTropes page and fully scrapes its contents, calling all sub functions and returning a valid Media object
// If the url of the page isn't found, it returns an ErrNotFound error
func (scraper *ServiceScraper) ScrapeTvTropesPage(page tropestogo.Page) (media.Media, error) {
	res, err := http.Get(page.GetUrl().String())
	if err != nil {
		return media.Media{}, fmt.Errorf("%w: "+page.GetUrl().String(), ErrNotFound)
	}

	return scraper.ScrapeWorkPage(res.Body, page.GetUrl())
}

// ScrapeWorkPage accepts a reader with a TvTropes Work Page contents and extracts all the relevant information from it
// and the url from the second parameter
// It scrapes the title, year, media type and all tropes, finally returning a correctly formed media object with all the data
// It calls sub functions for scraping the multiple parts and returns an error if some scraping has failed
func (scraper *ServiceScraper) ScrapeWorkPage(reader io.Reader, url *url.URL) (media.Media, error) {
	doc, _ := goquery.NewDocumentFromReader(reader)
	page, errNewPage := tropestogo.NewPage(url.String())
	if errNewPage != nil {
		return media.Media{}, fmt.Errorf("Error creating Page object \n%w", errNewPage)
	}

	title, year, mediaIndex, errMediaIndex := scraper.ScrapeWorkTitleAndYear(doc)
	if errMediaIndex != nil {
		return media.Media{}, errMediaIndex
	}

	tropes, errTropes := scraper.ScrapeWorkTropes(doc)
	if errTropes != nil {
		return media.Media{}, errTropes
	}

	newMedia, errNewMedia := media.NewMedia(title, year, time.Now(), tropes, page, mediaIndex)
	if errNewMedia != nil {
		return media.Media{}, errNewMedia
	}

	errAddMedia := scraper.data.AddMedia(newMedia)
	return newMedia, errAddMedia
}

// ScrapeWorkTitleAndYear traverses the received goquery Document DOM Tree and extracts
// the title, the year and the media index from the title section of the article of a Work Page
// It returns the title, the year, the media type and an ErrUnknownMediaType error if the media type isn't known
func (scraper *ServiceScraper) ScrapeWorkTitleAndYear(doc *goquery.Document) (string, string, media.MediaType, error) {
	var title, year string
	var mediaIndex media.MediaType
	var errMediaIndex error

	r, _ := regexp.Compile(`\s\((19|20)\d{2}\)`)
	fullTitle := strings.TrimSpace(strings.Split(doc.Find(WorkTitleSelector).Text(), "/")[1])
	regexSubstringMatch := r.FindStringSubmatch(fullTitle)
	if len(regexSubstringMatch) > 0 {
		year = regexSubstringMatch[0]
	}

	if year != "" {
		title = strings.ReplaceAll(fullTitle, year, "")
	} else {
		title = strings.TrimRight(fullTitle, " ")
	}

	year = strings.ReplaceAll(year, "(", "")
	year = strings.ReplaceAll(year, ")", "")
	year = strings.TrimLeft(year, " ")

	mediaIndex, errMediaIndex = media.ToMediaType(strings.Trim(doc.Find(WorkIndexSelector).First().Text(), " /"))

	return title, year, mediaIndex, errMediaIndex
}

// ScrapeWorkTropes traverses the received goquery Document DOM Tree and extracts all the tropes from the main trope section of a Work Page
// It returns a set (map of trope keys and empty values) of all the unique tropes found on the web page
func (scraper *ServiceScraper) ScrapeWorkTropes(doc *goquery.Document) (map[tropestogo.Trope]struct{}, error) {
	tropes := make(map[tropestogo.Trope]struct{}, 0)
	var newTrope tropestogo.Trope
	var newTropeError error

	var selector string
	if scraper.CheckTropesOnFolders(doc) {
		selector = MainTropesFolderSelector
	} else {
		selector = MainTropesSelector
	}

	doc.Find(selector).EachWithBreak(func(_ int, selection *goquery.Selection) bool {
		tropeUri, tropeUriExists := selection.Attr("href")
		if tropeUriExists {
			newTrope, newTropeError = tropestogo.NewTrope(strings.Split(tropeUri, "/")[4], tropestogo.TropeIndex(0))
			if newTropeError != nil {
				return false
			}

			if tropeUri == TvTropesMainPath+newTrope.GetTitle() {
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

// Persist calls the same method on the RepositoryMedia that is defined for the scraper and writes all data in the repository file
// If the internal data structure is empty, it will do nothing and return an ErrPersist error
// or return the proper Reading/Writing errors depending on the implementation
func (scraper *ServiceScraper) Persist() error {
	return scraper.data.Persist()
}
