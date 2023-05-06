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

type ServiceScraper struct {
	// TvTropes index
	index index.RepositoryIndex
	// TvTropes dataset
	data media.RepositoryMedia
}

func NewServiceScraper() (*ServiceScraper, error) {
	return &ServiceScraper{}, nil
}

// CheckValidWorkPage checks if a TvTropes Work page has a valid structure in which the scraper can extract data
func (*ServiceScraper) CheckValidWorkPage(*tropestogo.Page) (bool, error) {
	return false, nil
}
