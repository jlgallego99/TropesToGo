package tropestogo

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

const (
	TvTropesHostname  = "tvtropes.org"
	TvTropesPmwiki    = "/pmwiki/pmwiki.php/"
	TvTropesMainPath  = TvTropesPmwiki + "Main/"
	TvTropesIndexPath = "/pmwiki/pagelist_having_pagetype_in_namespace.php"
)

var (
	ErrNotTvTropes = errors.New("the URL does not belong to a TvTropes web page")
	ErrBadUrl      = errors.New("invalid URL")
	ErrEmptyUrl    = errors.New("the provided URL string is empty")
	ErrNotFound    = errors.New("couldn't request the URL")
	ErrForbidden   = errors.New("http request denied, maybe there has been too many requests")
	ErrParsing     = errors.New("error parsing the web contents")
)

// PageType represents all the relevant types a TvTropes Page can be, so the scraper can know what it is traversing
type PageType int64

const (
	UnknownPageType PageType = iota
	WorkPage
	MainPage
	IndexPage
)

// Page is a value-object that represents a generic TvTropes web page
type Page struct {
	// A Page can be accessed only by its URL, which doesn't change
	url *url.URL

	// document is a reference to the parsed HTML contents of the page
	// to avoid duplicating HTTP requests to the url if it's been done before
	document *goquery.Document

	// A Page in TvTropes can be, mainly, a main page, a work page or an index page
	pageType PageType
}

// NewPage creates a valid Page value-object that represents a generic and immutable TvTropes web page
// It accepts a pageUrl string and checks if it belongs to TvTropes and extracts the type of the page from it
// If requestPage argument is true, it makes an HTTP request to the Page URL and parses its content to a Goquery document
// (main page, work page, index page, etc.)
// It returns an ErrEmptyUrl error if it's empty or an ErrBadUrl error if it's not properly represented
// It returns an ErrNotFound if the web page couldn't be retrieved or an ErrForbidden if it's access has been temporarily denied by a 403 error
func NewPage(pageUrl string, requestPage bool) (Page, error) {
	if pageUrl == "" {
		return Page{}, ErrEmptyUrl
	}

	newUrl, errParse := url.Parse(pageUrl)
	if errParse != nil {
		return Page{}, fmt.Errorf("%w: "+pageUrl+"\n%w", ErrBadUrl, errParse)
	}

	if newUrl.Hostname() != TvTropesHostname {
		return Page{}, ErrNotTvTropes
	}

	var doc *goquery.Document
	if requestPage {
		httpResponse, errResponse := http.Get(newUrl.String())
		if errResponse != nil {
			return Page{}, fmt.Errorf("%w: "+newUrl.String(), ErrNotFound)
		}

		if httpResponse.StatusCode == 403 {
			return Page{}, fmt.Errorf("%w: "+newUrl.String(), ErrForbidden)
		}

		var errDoc error
		doc, errDoc = goquery.NewDocumentFromReader(httpResponse.Body)
		if errDoc != nil {
			return Page{}, fmt.Errorf("%w: "+newUrl.String(), ErrParsing)
		}
	} else {
		doc = nil
	}

	var pageType PageType
	if strings.HasPrefix(newUrl.Path, TvTropesMainPath) {
		pageType = MainPage
	} else if strings.HasPrefix(newUrl.Path, TvTropesPmwiki) {
		pageType = WorkPage
	} else if strings.HasPrefix(newUrl.Path, TvTropesIndexPath) {
		pageType = IndexPage
	} else {
		pageType = UnknownPageType
	}

	return Page{
		url:      newUrl,
		document: doc,
		pageType: pageType,
	}, nil
}

// GetUrl returns the inmutable URL object that defines the Page
func (page Page) GetUrl() *url.URL {
	return page.url
}

// GetPageType returns the PageType enum that represents the type of the Page
func (page Page) GetPageType() PageType {
	return page.pageType
}

func (page Page) GetDocument() *goquery.Document {
	return page.document
}
