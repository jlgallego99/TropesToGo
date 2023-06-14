package tropestogo

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

var (
	ErrDuplicatedPage = errors.New("the page already exists")
	ErrEmptyUrl       = errors.New("the provided URL string is empty")
)

// TvTropesPages is an entity that manages all relevant pages in TvTropes for its extraction
// Each Page has a last updated date, for checking future TvTropes updates on its Pages
type TvTropesPages struct {
	Pages map[Page]time.Time
}

// NewTvTropesPages creates an empty object to which we can add valid Pages
func NewTvTropesPages() *TvTropesPages {
	return &TvTropesPages{
		Pages: make(map[Page]time.Time, 0),
	}
}

// AddTvTropesPage creates a valid TvTropes Page from a string pageUrl and adds it to the internal structure of all pages
// except if the page has already been added before, then it will return an ErrDuplicatedPage error
// If the url is empty or has an invalid format, it will return an ErrBadUrl error
// If the url does not belong to a TvTropes page, it will return an ErrNotTvTropes error
func (tvtropespages *TvTropesPages) AddTvTropesPage(pageUrl string) error {
	if pageUrl == "" {
		return ErrEmptyUrl
	}

	newUrl, errParse := url.Parse(pageUrl)
	if errParse != nil {
		return fmt.Errorf("%w: "+pageUrl+"\n%w", ErrBadUrl, errParse)
	}

	newPage, errNewPage := NewPage(newUrl)
	if errNewPage != nil {
		return errNewPage
	}

	for tvtropesPage := range tvtropespages.Pages {
		if tvtropesPage.GetUrl().String() == pageUrl {
			return fmt.Errorf("%w: "+pageUrl, ErrDuplicatedPage)
		}
	}

	tvtropespages.Pages[newPage] = time.Now()
	return nil
}
