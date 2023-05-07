package scraper

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/index"
	"github.com/jlgallego99/TropesToGo/media"
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

	// Check the main article structure
	title := page.Document.Find("h1.entry-title")
	index := title.Find("strong")
	if strings.Trim(index.Text(), " /") != "Film" {
		return false, ErrNotWorkPage
	}

	// Check if the tropes part of the page has a known structure

	return true, nil
}
