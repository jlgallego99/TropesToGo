package tropestogo

import (
	"errors"
	"fmt"
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

	// A Page in TvTropes can be, mainly, a main page, a work page or an index page
	pageType PageType
}

// NewPage creates a valid Page value-object that represents a generic and immutable TvTropes web page
// It accepts a URL object and checks if it belongs to TvTropes and extracts the type of the page from it
// (main page, work page, index page, etc.)
func NewPage(URL *url.URL) (Page, error) {
	if URL == nil {
		return Page{}, fmt.Errorf("%w: URL object is null", ErrBadUrl)
	}

	if URL.Hostname() != TvTropesHostname {
		return Page{}, ErrNotTvTropes
	}

	var pageType PageType
	if strings.HasPrefix(URL.Path, TvTropesMainPath) {
		pageType = MainPage
	} else if strings.HasPrefix(URL.Path, TvTropesPmwiki) {
		pageType = WorkPage
	} else if strings.HasPrefix(URL.Path, TvTropesIndexPath) {
		pageType = IndexPage
	} else {
		pageType = UnknownPageType
	}

	return Page{
		url:      URL,
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
