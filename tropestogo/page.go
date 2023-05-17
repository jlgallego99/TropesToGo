package tropestogo

import (
	"net/url"
	"time"
)

// Page is a TvTropes Work page for later scraping
type Page struct {
	// URL defines the identity of the Page entity
	URL *url.URL

	// LastUpdated is the last time the page was updated, for helping with maintaining information updated
	LastUpdated time.Time
}
