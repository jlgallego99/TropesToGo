package crawler

import (
	"errors"
	"fmt"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"net/url"
)

const Pagelist = "https://tvtropes.org/pmwiki/pagelist_having_pagetype_in_namespace.php?t=work"

var (
	ErrNotFound = errors.New("couldn't request the URL")
	ErrBadUrl   = errors.New("invalid URL")
)

type ServiceCrawler struct {
	// The seed is the starting URL of the crawler
	seed tropestogo.Page
}

func NewCrawler() (*ServiceCrawler, error) {
	return nil, nil
}

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
