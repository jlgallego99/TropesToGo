package crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	Pagelist = "https://tvtropes.org/pmwiki/pagelist_having_pagetype_in_namespace.php?t=work"

	WorkPageSelector = "table a"
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

// CrawlWorkPages searches all Work pages from the defined seed starting page
// It returns an array of Page objects that are all crawled TvTropes Work pages
func (crawler *ServiceCrawler) CrawlWorkPages() ([]tropestogo.Page, error) {
	var crawledPages []tropestogo.Page
	pageNumber := 1

	for {
		workPageList := crawler.seed.GetUrl()
		values := workPageList.Query()
		values.Add("page", strconv.Itoa(pageNumber))
		workPageList.RawQuery = values.Encode()

		resp, errGet := http.Get(workPageList.String())
		if errGet != nil {
			return []tropestogo.Page{}, fmt.Errorf("%w: "+crawler.seed.GetUrl().String(), ErrNotFound)
		}

		doc, errDocument := goquery.NewDocumentFromReader(resp.Body)
		if errDocument != nil {
			return []tropestogo.Page{}, fmt.Errorf("%w: "+crawler.seed.GetUrl().String(), ErrInvalidPage)
		}

		pageSelector := doc.Find(WorkPageSelector)
		if pageSelector.Length() == 0 {
			break
		}

		pageSelector.Each(func(_ int, selection *goquery.Selection) {
			workUrl, urlExists := selection.Attr("href")

			if urlExists {
				newPage, errNewPage := tropestogo.NewPage(workUrl)

				if errNewPage == nil {
					crawledPages = append(crawledPages, newPage)
				}
			}
		})

		pageNumber += 1
		time.Sleep(time.Second / 2)
	}

	return crawledPages, nil
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
