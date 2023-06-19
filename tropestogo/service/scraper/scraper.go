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
	ErrInvalidSubpage       = errors.New("couldn't scrape tropes in subpage")
)

const (
	TvTropesHostname         = "tvtropes.org"
	TvTropesWeb              = "https://" + TvTropesHostname
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
func (scraper *ServiceScraper) CheckValidWorkPage(reader io.Reader, checkUrl *url.URL) (bool, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return false, ErrNotTvTropes
	}

	validWorkPage, errWorkPage := scraper.CheckIsWorkPage(doc, checkUrl)
	if !validWorkPage {
		return false, errWorkPage
	}

	validTropeSection, errTropeSection := scraper.CheckTropeSection(doc)
	if !validTropeSection {
		return false, errTropeSection
	}

	return true, nil
}

// CheckIsWorkPage checks if the received url belongs to a correct tvtropes.org Work Page
// It returns an ErrNotTvTropes error if it doesn't belong to TvTropes, an ErrNotWorkPage if it's from TvTropes but of any other type
// or an ErrUnknownPageStructure if there are strange elements on it
// Returns true if all checks passes
func (scraper *ServiceScraper) CheckIsWorkPage(doc *goquery.Document, url *url.URL) (bool, error) {
	if doc == nil || url == nil {
		return false, ErrInvalidField
	}

	if url.Hostname() != TvTropesHostname {
		return false, fmt.Errorf("%w: "+url.String(), ErrNotTvTropes)
	}

	splitPath := strings.Split(url.Path, "/")
	if !strings.HasPrefix(url.Path, TvTropesPmwiki) || splitPath[3] != media.Film.String() {
		return false, fmt.Errorf("%w: "+url.String(), ErrNotWorkPage)
	}

	if doc.Find(MainArticleSelector).Length() == 0 ||
		doc.Find(SubPagesNavSelector).Find(SubPageListSelector).Find(SubPageLinkSelector).Length() == 0 ||
		doc.Find(TropeListSelector).Length() == 0 ||
		doc.Find(TropeLinkSelector).Length() == 0 {
		return false, ErrUnknownPageStructure
	}

	tropeIndex := strings.Trim(doc.Find(WorkIndexSelector).Text(), " /")
	if tropeIndex != media.Film.String() {
		return false, fmt.Errorf("%w: the index is"+tropeIndex, ErrNotWorkPage)
	}

	return true, nil
}

// CheckTropeSection checks the received goquery Document DOM Tree
// and validates if the tropes on the TvTropes Work Page are arranged in a known way that can be scraped
// First it checks if there are folders, then if the list redirects to trope subpages and last if there's a list with only tropes
// If the trope section isn't of the recognized types, it returns an ErrUnknownPageStructure, so it can't be scraped
// It returns true if all checks passes
func (scraper *ServiceScraper) CheckTropeSection(doc *goquery.Document) (bool, error) {
	if scraper.CheckTropesOnFolders(doc) {
		return true, nil
	}

	title, _, _, _ := scraper.ScrapeWorkTitleAndYear(doc)
	if scraper.CheckTropesOnSubpages(doc, title) {
		return true, nil
	}

	tropeHref, exists := doc.Find(TropeLinkSelector).First().Attr("href")
	if exists && strings.HasPrefix(tropeHref, TvTropesMainPath) {
		return true, nil
	}

	// Tropes are presented in an unknown form, so data can't be extracted
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

// CheckTropesOnSubpages validates whether main tropes are divided on subpages that holds other tropes of the same Work
// The method accepts the DOM tree goquery document and the title of the Work for checking if the subpages has a correct URI
// It returns a boolean meaning if the check passes or not
func (scraper *ServiceScraper) CheckTropesOnSubpages(doc *goquery.Document, workTitle string) bool {
	// Get the first word of the first element of the list, check if is an anchor to a subpage
	tropeHref, exists := doc.Find(TropeLinkSelector).First().Attr("href")

	// Check if the link directs to a subpage with tropes
	if exists && strings.HasPrefix(tropeHref, TvTropesPmwiki) && scraper.CheckSubpageUri(tropeHref, workTitle) {
		return true
	}

	return false
}

// CheckSubpageUri checks if a TvTropes URI belongs to a subpage with tropes on it
// Checks if the elements of the list are anchors to a subpage inside the work
// A regex matches if the last part of the URL is of the type TropesXtoY and is preceded by the work title name, returning true or false
func (scraper *ServiceScraper) CheckSubpageUri(URI, title string) bool {
	title = strings.ToLower(strings.ReplaceAll(title, " ", ""))
	r, _ := regexp.Compile(`\/` + title + `\/tropes[a-z]to[a-z]`)
	match := r.MatchString(strings.ToLower(URI))

	return match
}

// CheckSubwikiUri checks both the URI and the title for validating if the given subpage is a sub wiki of the Work that holds secondary tropes
// Returns a true boolean if it's a subwiki, a false if it's not
func (scraper *ServiceScraper) CheckSubwikiUri(URI, title string) bool {
	title = strings.ToLower(strings.ReplaceAll(title, " ", ""))
	r, _ := regexp.Compile(`\/\w*\/` + title)
	match := r.MatchString(strings.ToLower(URI))

	return match
}

// ScrapeTvTropes tries to scrape all pages and its subpages that are TvTropesPages by making HTTP requests to TvTropes
// It only returns an error if it can't write or read the dataset, if the page can't be scraped it skips to the next
func (scraper *ServiceScraper) ScrapeTvTropes(tvtropespages *tropestogo.TvTropesPages) error {
	for page, tvtropessubpages := range tvtropespages.Pages {
		var subPages []tropestogo.Page
		for subPage := range tvtropessubpages.Subpages {
			subPages = append(subPages, subPage)
		}

		scraper.ScrapeTvTropesPage(page, subPages)
	}

	errPersist := scraper.Persist()
	if errPersist != nil {
		return errPersist
	}

	return nil
}

// ScrapeTvTropesPage makes an HTTP request to a TvTropes page and its subpages and fully scrapes its contents
// calling all sub functions and returning a valid Media object
// If the url of the page isn't found, it returns an ErrNotFound error
func (scraper *ServiceScraper) ScrapeTvTropesPage(page tropestogo.Page, subPages []tropestogo.Page) (media.Media, error) {
	res, err := http.Get(page.GetUrl().String())
	if err != nil {
		return media.Media{}, fmt.Errorf("%w: "+page.GetUrl().String(), ErrNotFound)
	}

	var subReaders []io.Reader
	for _, subPage := range subPages {
		subPageRes, subPageErr := http.Get(subPage.GetUrl().String())
		if subPageErr != nil {
			return media.Media{}, fmt.Errorf("%w: "+page.GetUrl().String(), ErrNotFound)
		}

		subReaders = append(subReaders, subPageRes.Body)
	}

	return scraper.ScrapeWorkPage(res.Body, subReaders, page.GetUrl())
}

// ScrapeWorkPage accepts a reader with a TvTropes Work Page contents and an array of readers with all its subpages contents
// Extracts all the relevant information from it and the url from the last parameter
// It scrapes the title, year, media type and all tropes, finally returning a correctly formed media object with all the data
// It calls sub functions for scraping the multiple parts and returns an error if some scraping has failed
func (scraper *ServiceScraper) ScrapeWorkPage(reader io.Reader, subReaders []io.Reader, url *url.URL) (media.Media, error) {
	var tropes map[tropestogo.Trope]struct{}
	var errTropes error
	var subDocs []*goquery.Document
	doc, _ := goquery.NewDocumentFromReader(reader)

	page, errNewPage := tropestogo.NewPage(url.String())
	if errNewPage != nil {
		return media.Media{}, fmt.Errorf("Error creating Page object \n%w", errNewPage)
	}

	title, year, mediaIndex, errMediaIndex := scraper.ScrapeWorkTitleAndYear(doc)
	if errMediaIndex != nil {
		return media.Media{}, errMediaIndex
	}

	if scraper.CheckTropesOnSubpages(doc, title) {
		for _, subReader := range subReaders {
			subDoc, _ := goquery.NewDocumentFromReader(subReader)
			subDocs = append(subDocs, subDoc)
		}

		tropes, errTropes = scraper.ScrapeMainSubpageTropes(doc, subDocs, title)
	} else {
		tropes, errTropes = scraper.ScrapeTropes(doc)
	}

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

// ScrapeTropes traverses the received goquery Document DOM Tree and extracts all the tropes that are in a list or in folders
// It returns a set (map of trope keys and empty values) of all the unique tropes found on the web page
func (scraper *ServiceScraper) ScrapeTropes(doc *goquery.Document) (map[tropestogo.Trope]struct{}, error) {
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
			newTrope, newTropeError = tropestogo.NewTrope(strings.Split(tropeUri, "/")[4], tropestogo.TropeIndex(0), "")
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

// ScrapeSubpageFullTitle scrapes the full title of any Work subpage from a Goquery document
// (<Title>/<TropesXtoY> for main tropes subpages and <Namespace>/<Title>) for subwikis)
// Returns a correctly formated string without blanks for comparing with URIs
func (scraper *ServiceScraper) ScrapeSubpageFullTitle(subDoc *goquery.Document) string {
	subPageTitle := "/" + strings.ReplaceAll(strings.ReplaceAll(subDoc.Find(WorkTitleSelector).Text(), "\n", ""), " ", "")

	return subPageTitle
}

// ScrapeMainSubpageTropes extracts the main tropes (the ones that are on the main article) that are divided into subpages because they are many
// It traverses the DOM tree document and searches for subpages whose URI has a known structure and have the Work title string
// It performs various ScrapeTropes calls for each of the subpages, adding its tropes to the main trope list
// Returns a trope list of all tropes found on the different main trope subpages
// If the subpage can't be scraped, it returns an ErrInvalidSubpage error
func (scraper *ServiceScraper) ScrapeMainSubpageTropes(doc *goquery.Document, subDocs []*goquery.Document, title string) (map[tropestogo.Trope]struct{}, error) {
	var errSubpage error
	tropes := make(map[tropestogo.Trope]struct{})

	doc.Find(MainTropesSelector).EachWithBreak(func(i int, selection *goquery.Selection) bool {
		subpageUri, subpageUriExists := selection.Attr("href")
		if i >= len(subDocs) {
			return false
		}

		if subpageUriExists && scraper.CheckSubpageUri(subpageUri, title) &&
			scraper.CheckSubpageUri(scraper.ScrapeSubpageFullTitle(subDocs[i]), title) {

			subpageTropes, err := scraper.ScrapeTropes(subDocs[i])
			if err == nil {
				for subpageTrope := range subpageTropes {
					tropes[subpageTrope] = struct{}{}
				}

				return true
			}
		}

		errSubpage = fmt.Errorf("%w: "+subpageUri, ErrInvalidSubpage)
		return false
	})

	if errSubpage != nil {
		return make(map[tropestogo.Trope]struct{}), errSubpage
	}

	return tropes, nil
}

// Persist calls the same method on the RepositoryMedia that is defined for the scraper and writes all data in the repository file
// If the internal data structure is empty, it will do nothing and return an ErrPersist error
// or return the proper Reading/Writing errors depending on the implementation
func (scraper *ServiceScraper) Persist() error {
	return scraper.data.Persist()
}
