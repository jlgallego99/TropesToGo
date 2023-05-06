package scraper

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/index"
	"github.com/jlgallego99/TropesToGo/media"
)

var (
	ErrInvalidField = errors.New("one or more fields for the Scraper are invalid")
	ErrNotTvTropes  = errors.New("the URL does not belong to a TvTropes page")
	ErrNotWorkPage  = errors.New("the page isn't a TvTropes Work page")
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
func (*ServiceScraper) CheckValidWorkPage(*tropestogo.Page) (bool, error) {
	return false, nil
}
