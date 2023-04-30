package tropestogo

import (
	"github.com/google/uuid"
	"time"
)

// Page is a TvTropes raw page for later scraping
type Page struct {
	ID          uuid.UUID
	URL         string
	RawHTML     string
	LastUpdated time.Time
}
