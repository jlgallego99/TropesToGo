package tropestogo

import (
	"github.com/google/uuid"
	"net/url"
	"time"
)

// Page is a TvTropes raw page for later scraping
type Page struct {
	ID          uuid.UUID
	URL         url.URL
	RawHTML     string
	LastUpdated time.Time
}
