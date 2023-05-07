package scraper

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/index"
	"github.com/jlgallego99/TropesToGo/media"
	"regexp"
	"strings"
)

var (
	ErrInvalidField         = errors.New("one or more fields for the Scraper are invalid")
	ErrNotTvTropes          = errors.New("the URL does not belong to a TvTropes page")
	ErrNotWorkPage          = errors.New("the page isn't a TvTropes Work page")
	ErrUnknownPageStructure = errors.New("the scraper doesn't recognize the page structure")
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
		ss.index = ir
		return nil
	}
}

func ConfigRepository(mr media.RepositoryMedia) ScraperConfig {
	return func(ss *ServiceScraper) error {
		ss.data = mr
		return nil
	}
}

// CheckValidWorkPage checks if a TvTropes Work page has a valid structure in which the scraper can extract data
// This allows the scraper to check if TvTropes template has somewhat changed
func (*ServiceScraper) CheckValidWorkPage(page *tropestogo.Page) (bool, error) {
	// First check if the domain is TvTropes
	if page.URL.Hostname() != "tvtropes.org" {
		return false, ErrNotTvTropes
	}

	// Check if it's a Film Work page
	splitPath := strings.Split(page.URL.Path, "/")
	if !strings.HasPrefix(page.URL.Path, "/pmwiki/pmwiki.php") || splitPath[3] != "Film" {
		return false, ErrNotWorkPage
	}

	// Check if the main article structure has all known ids and elements that comprise a TvTropes work page
	if page.Document.Find("#main-article").Length() == 0 ||
		page.Document.Find("nav.body-options").Find("ul.subpage-links").Find("a.subpage-link").Length() == 0 {
		return false, ErrUnknownPageStructure
	}

	// Check the title
	title := page.Document.Find("h1.entry-title")
	index := title.Find("strong")
	if strings.Trim(index.Text(), " /") != "Film" {
		return false, ErrNotWorkPage
	}

	// Look for the tropes section
	if page.Document.Find("#main-article ul").Length() == 0 &&
		strings.Contains(strings.ToLower(page.Document.Find("#main-article h2").First().Text()), "tropes") {

		return false, ErrUnknownPageStructure
	} else {
		// Check if the list is a) a simple trope list or c) a list of subpages with tropes
		if page.Document.Find("#main-article ul li a.twikilink").Length() != 0 {
			// Get the first word of the first element of the list, check if is an anchor to a trope page or a sub page
			tropeHref, exists := page.Document.Find("#main-article ul li a.twikilink").First().Attr("href")

			// a) Tropes are presented on a single list
			// Checks if it's a Main page
			if exists && strings.HasPrefix(tropeHref, "/pmwiki/pmwiki.php/Main/") {
				return true, nil
			}

			// c) Tropes are on subpages
			// Checks if the elements of the list are anchors to a subpage inside the work
			// A regex matches if the last part of the URL is of the type TropesXtoY
			hrefSplit := strings.Split(tropeHref, "/")
			r, _ := regexp.Compile("Tropes[A-Z]To[A-Z]")
			match := r.MatchString(hrefSplit[len(hrefSplit)-1])
			if exists && strings.HasPrefix(tropeHref, "/pmwiki/pmwiki.php/") && match {
				return true, nil
			}
		}

		// b) Check if tropes are on folders
		// If there's a close all folders button, then the tropes are on folders
		folderFunctionName, existsFolderButton := page.Document.Find("#main-article div.folderlabel").Attr("onclick")
		if existsFolderButton && folderFunctionName == "toggleAllFolders()" {
			return true, nil
		}

		// Tropes are presented in an unknown form, so data can't be extracted
		return false, ErrUnknownPageStructure
	}

	// If it isn't any of the know formats for the trope list, check if there are tropes references
}
