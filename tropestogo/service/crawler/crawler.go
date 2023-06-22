package crawler

import (
	"errors"
	"fmt"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"net/url"
)

const (
	Pagelist = "https://tvtropes.org/pmwiki/pagelist_having_pagetype_in_namespace.php?t=work"
)

var (
	ErrNotFound    = errors.New("couldn't request the URL")
	ErrBadUrl      = errors.New("invalid URL")
	ErrInvalidPage = errors.New("couldn't crawl in page")
)

type ServiceCrawler struct {
	// The seed is the starting URL of the crawler
	seed tropestogo.Page
}

func NewCrawler(mediaTypeString string) (*ServiceCrawler, error) {
	crawler := &ServiceCrawler{}
	mediaType, errMediaType := media.ToMediaType(mediaTypeString)
	if errMediaType != nil {
		return nil, errMediaType
	}

	crawler.SetMediaSeed(mediaType)

	return crawler, nil
}

// SetMediaSeed sets a mediaType for the crawler seed (starting page) for crawling all pages of that medium
// It returns an error if the page URL or the mediaType isn't valid/doesn't exist on TvTropes
func (crawler *ServiceCrawler) SetMediaSeed(mediaType media.MediaType) error {
	seedUrl, errParse := url.Parse(Pagelist)
	if errParse != nil {
		return fmt.Errorf("%w: "+Pagelist+"\n%w", ErrBadUrl, errParse)
	}

	values := seedUrl.Query()
	values.Add("n", mediaType.String())
	seedUrl.RawQuery = values.Encode()

	seedPage, errNewPage := tropestogo.NewPage(seedUrl.String())
	if errNewPage != nil {
		return errNewPage
	}

	crawler.seed = seedPage

	return nil
}
