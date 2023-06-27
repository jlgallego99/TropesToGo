package tropestogo

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var (
	ErrDuplicatedPage = errors.New("the page already exists")
	seededRand        = rand.New(rand.NewSource(time.Now().UnixNano()))
)

const minWaitingSeconds = 0
const maxWaitingSeconds = 4

// TvTropesPages is an entity that manages all relevant pages in TvTropes for its extraction
// Each Page has a last updated date, for checking future TvTropes updates on its Pages
// and a list of its subpages, each with their own last updated time
type TvTropesPages struct {
	Pages map[Page]*TvTropesSubpages
}

// TvTropesSubpages is an entity that manages all subpages of a TvTropes page for its extraction
// It has the main page last updated time and also each Page has its own last updated date
type TvTropesSubpages struct {
	// LastUpdated is the last time the main page was updated
	LastUpdated time.Time

	// Subpages are pages that are inside the main page, and each has a last time when they were updated
	Subpages map[Page]time.Time
}

// NewTvTropesPages creates an empty object to which we can add valid Pages
func NewTvTropesPages() *TvTropesPages {
	return &TvTropesPages{
		Pages: make(map[Page]*TvTropesSubpages, 0),
	}
}

// AddTvTropesPage creates a valid TvTropes Page from a string pageUrl and adds it to the internal structure of all pages
// except if the page has already been added before, then it will return an ErrDuplicatedPage error
// If the requestPages argument is true, it makes http requests to all pages with a random waiting time between requests
// If the url is empty or has an invalid format, it will return either an ErrEmptyUrl or ErrBadUrl error
// If the url does not belong to a TvTropes page, it will return an ErrNotTvTropes error
// If TvTropes denies access because of too many requests, it will not create the Page and return an ErrForbidden error for the crawler to manage
func (tvtropespages *TvTropesPages) AddTvTropesPage(pageUrl string, subpageUrls []string, requestPages bool) error {
	newPage, errNewPage := NewPage(pageUrl, requestPages)
	if errNewPage != nil {
		return errNewPage
	}

	for tvtropesPage := range tvtropespages.Pages {
		if tvtropesPage.GetUrl().String() == pageUrl {
			return fmt.Errorf("%w: "+pageUrl, ErrDuplicatedPage)
		}
	}

	subPages := &TvTropesSubpages{
		LastUpdated: time.Now(),
		Subpages:    make(map[Page]time.Time, 0),
	}
	for _, subpageUrl := range subpageUrls {
		newSubpage, errSubpage := NewPage(subpageUrl, requestPages)
		if errSubpage != nil {
			return errSubpage
		}

		subPages.Subpages[newSubpage] = time.Now()

		// Wait random time between HTTP requests
		if requestPages {
			waitingFactor := seededRand.Intn(maxWaitingSeconds-minWaitingSeconds) + minWaitingSeconds
			time.Sleep(time.Second * time.Duration(waitingFactor))
		}
	}

	tvtropespages.Pages[newPage] = subPages
	return nil
}
