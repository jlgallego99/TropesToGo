package tropestogo

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"time"
)

// Page is a TvTropes Work page for later scraping
type Page struct {
	// URL defines the identity of the Page entity
	URL *url.URL

	// DOMTree holds a traversable tree of HTML elements comprised of only the main body of a Work page
	// This has all the information the Scraper needs
	DOMTree *goquery.Document

	// LastUpdated is the last time the page was updated, for helping with maintaining information updated
	LastUpdated time.Time
}
