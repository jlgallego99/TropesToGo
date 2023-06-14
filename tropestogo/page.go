package tropestogo

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

const TvTropesHostname = "tvtropes.org"

var (
	ErrNotTvTropes = errors.New("the URL does not belong to a TvTropes web page")
	ErrBadUrl      = errors.New("invalid URL")
)

// Page is a TvTropes Work page for later scraping
type Page struct {
	// URL defines the identity of the Page entity
	URL *url.URL

	// LastUpdated is the last time the page was updated, for helping with maintaining information updated
	LastUpdated time.Time
}

func NewPage(URL *url.URL) (*Page, error) {
	if URL == nil {
		return nil, fmt.Errorf("%w: URL object is null", ErrBadUrl)
	}

	if URL.Hostname() != TvTropesHostname {
		return nil, ErrNotTvTropes
	}

	return &Page{
		URL:         URL,
		LastUpdated: time.Now(),
	}, nil
}
