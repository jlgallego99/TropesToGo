package tvtropespages

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
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

	httpClient = &http.Client{}
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
func NewPage(pageUrl string, requestPage bool, req *http.Request) (Page, error) {
	if pageUrl == "" {
		return Page{}, ErrEmptyUrl
	}

	parsedUrl, errParse := parseTvTropesUrl(pageUrl)
	if errParse != nil {
		return Page{}, errParse
	}

	var doc *goquery.Document = nil
	if requestPage {
		httpResponse, errRequest := doRequest(req)
		if errRequest != nil {
			return Page{}, errRequest
		}

		if httpResponse.StatusCode == 403 || httpResponse.StatusCode == 429 {
			return Page{}, fmt.Errorf("%w: "+pageUrl, ErrForbidden)
		}

		var errParseDocument error
		doc, errParseDocument = parsePageDocument(httpResponse.Body)
		if errParseDocument != nil {
			return Page{}, errParseDocument
		}
	}

	pageType := inferPageType(parsedUrl)

	return Page{
		url:      parsedUrl,
		document: doc,
		pageType: pageType,
	}, nil
}

// NewPageWithDocument creates a valid Page value-object with a custom document that represents a generic and immutable TvTropes web page
func NewPageWithDocument(pageUrl string, doc *goquery.Document) (Page, error) {
	page, errPage := NewPage(pageUrl, false, nil)
	if errPage != nil {
		return Page{}, errPage
	}

	page.document = doc

	return page, nil
}

// parseTvTropesUrl accepts a pageUrl string and parses it to a valid URL object, only if it belongs to TvTropes
// If it can't be parsed it returns an ErrBadUrl error and if it's not a TvTropes page it returns an ErrNotTvTropes error
func parseTvTropesUrl(pageUrl string) (*url.URL, error) {
	newUrl, errParse := url.Parse(pageUrl)
	if errParse != nil {
		return nil, fmt.Errorf("%w: "+pageUrl+"\n%w", ErrBadUrl, errParse)
	}

	if newUrl.Hostname() != TvTropesHostname {
		return nil, ErrNotTvTropes
	}

	return newUrl, nil
}

// doRequest tries to make an HTTP request and returns its contents
// If the URL isn't available for retrieving its content will return an ErrNotFound or an ErrForbidden error
func doRequest(request *http.Request) (*http.Response, error) {
	httpResponse, errDoRequest := httpClient.Do(request)
	if errDoRequest != nil {
		return nil, fmt.Errorf("%w: "+request.URL.String(), ErrNotFound)
	}

	return httpResponse, nil
}

// parsePageDocument accepts a reader object containing the contents of a web page and returns a goquery Document with them
// if it can't be parsed, it will return an ErrParsing error
func parsePageDocument(reader io.Reader) (*goquery.Document, error) {
	doc, errDoc := goquery.NewDocumentFromReader(reader)
	if errDoc != nil {
		return nil, ErrParsing
	}

	return doc, nil
}

// inferPageType accepts a URL object and analyzes it to infer what TvTropes type page it is
// If it isn't a known type, it will return an UnknownPageType
func inferPageType(url *url.URL) PageType {
	var pageType PageType
	if strings.HasPrefix(url.Path, TvTropesMainPath) {
		pageType = MainPage
	} else if strings.HasPrefix(url.Path, TvTropesPmwiki) {
		pageType = WorkPage
	} else if strings.HasPrefix(url.Path, TvTropesIndexPath) {
		pageType = IndexPage
	} else {
		pageType = UnknownPageType
	}

	return pageType
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
