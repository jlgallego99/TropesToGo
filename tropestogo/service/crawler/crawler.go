package crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Pagelist = "https://tvtropes.org/pmwiki/pagelist_having_pagetype_in_namespace.php?t=work"

	TvTropesHostname       = "tvtropes.org"
	TvTropesWeb            = "https://" + TvTropesHostname
	WorkPageSelector       = "table a"
	CurrentSubpageSelector = ".curr-subpage"
	SubWikiSelector        = "a.subpage-link:not(" + CurrentSubpageSelector + ")"
	SubPageSelector        = "ul a.twikilink"
)

var (
	ErrNotFound    = errors.New("couldn't request the URL")
	ErrBadUrl      = errors.New("invalid URL")
	ErrInvalidPage = errors.New("couldn't crawl in page")
	ErrCrawling    = errors.New("there was an error crawling TvTropes")
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

// CrawlWorkPages searches crawlLimit number of Work pages from the defined seed starting page; if it's 0 or less, then it crawls all Work pages
// It returns a TvTropesPages object with all crawled pages and subpages from TvTropes
func (crawler *ServiceCrawler) CrawlWorkPages(crawlLimit int) (*tropestogo.TvTropesPages, error) {
	crawledPages := tropestogo.NewTvTropesPages()
	pageNumber := 1

	limitedCrawling := true
	if crawlLimit <= 0 {
		limitedCrawling = false
	}

	for {
		listSelector, errListSelector := crawler.GetWorkListSelector(pageNumber)
		if errListSelector != nil {
			return nil, errListSelector
		}

		var errAddPage error
		var pageReader *http.Response
		listSelector.EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if limitedCrawling && len(crawledPages.Pages) == crawlLimit {
				return false
			}

			workUrl, urlExists := selection.Attr("href")
			pageReader, errAddPage = http.Get(workUrl)
			if pageReader.StatusCode == 403 {
				time.Sleep(time.Minute)
			}

			if errAddPage != nil || !urlExists {
				return false
			}

			pageDoc, _ := goquery.NewDocumentFromReader(pageReader.Body)
			subPagesUrls := crawler.CrawlWorkSubpages(pageDoc)

			errAddPage = crawledPages.AddTvTropesPage(workUrl, subPagesUrls)
			time.Sleep(time.Second / 2)

			return true
		})

		if limitedCrawling && len(crawledPages.Pages) == crawlLimit {
			break
		}

		if errAddPage != nil {
			return nil, fmt.Errorf("error crawling: %w", errAddPage)
		}

		pageNumber += 1
	}

	return crawledPages, nil
}

// GetWorkListSelector returns the Nth index page selector for crawling Work Pages
// It returns an error if there's no page
func (crawler *ServiceCrawler) GetWorkListSelector(indexPage int) (*goquery.Selection, error) {
	workPageList := crawler.seed.GetUrl()
	values := workPageList.Query()
	values.Add("page", strconv.Itoa(indexPage))
	workPageList.RawQuery = values.Encode()

	resp, errGetIndex := http.Get(workPageList.String())
	if errGetIndex != nil {
		return nil, fmt.Errorf("%w: "+crawler.seed.GetUrl().String(), ErrNotFound)
	}

	doc, errDocument := goquery.NewDocumentFromReader(resp.Body)
	if errDocument != nil {
		return nil, fmt.Errorf("%w: "+crawler.seed.GetUrl().String(), ErrInvalidPage)
	}

	pageSelector := doc.Find(WorkPageSelector)
	if pageSelector.Length() == 0 {
		return nil, fmt.Errorf("%w: "+workPageList.String(), ErrInvalidPage)
	}

	return pageSelector, nil
}

// CrawlWorkSubpages searches all subpages (both with main tropes and SubWikis) on the goquery Document of a Work page
// It returns an array of string URLs that belong to all crawled TvTropes Work subpages
func (crawler *ServiceCrawler) CrawlWorkSubpages(doc *goquery.Document) []string {
	var subPagesUrls []string

	// Get all SubWikis
	doc.Find(SubWikiSelector).Each(func(_ int, selection *goquery.Selection) {
		subWikiUri, subWikiExists := selection.Attr("href")
		if subWikiExists {
			subPagesUrls = append(subPagesUrls, TvTropesWeb+subWikiUri)
		}
	})

	// Get all main trope subpages (if there are any)
	doc.Find(SubPageSelector).Each(func(_ int, selection *goquery.Selection) {
		subPageUri, subPageExists := selection.Attr("href")
		r, _ := regexp.Compile(`\/tropes[a-z]to[a-z]`)
		matchUri := r.MatchString(strings.ToLower(subPageUri))

		if subPageExists && matchUri {
			subPagesUrls = append(subPagesUrls, TvTropesWeb+subPageUri)
		}
	})

	return subPagesUrls
}

// SetMediaSeed sets a mediaType for the crawler seed (starting page) for crawling all pages of that medium
// It returns an error if the page URL or the mediaType isn't valid/doesn't exist on TvTropes
func (crawler *ServiceCrawler) SetMediaSeed(mediaType media.MediaType) error {
	seedUrl, errParse := url.Parse(Pagelist)
	if errParse != nil {
		return fmt.Errorf("%w: "+Pagelist, ErrBadUrl)
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

// CrawlWorkPagesFromReaders crawls all Work Pages and its subpages from an index reader and its pages readers. Only for test purposes
// It searches crawlLimit number of Work pages within the index
// It returns a TvTropesPages object with all crawled pages and subpages from TvTropes
func (crawler *ServiceCrawler) CrawlWorkPagesFromReaders(indexReader io.Reader, workReaders []io.Reader, crawlLimit int) (*tropestogo.TvTropesPages, error) {
	crawledPages := tropestogo.NewTvTropesPages()
	pageNumber := 1

	limitedCrawling := true
	if crawlLimit <= 0 {
		limitedCrawling = false
	}

	for {
		doc, errDocument := goquery.NewDocumentFromReader(indexReader)
		if errDocument != nil {
			return nil, fmt.Errorf("%w: "+crawler.seed.GetUrl().String(), ErrInvalidPage)
		}

		listSelector := doc.Find(WorkPageSelector)
		if listSelector.Length() == 0 {
			return nil, fmt.Errorf("%w: "+crawler.seed.GetUrl().String(), ErrInvalidPage)
		}

		var errAddPage error
		listSelector.EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if limitedCrawling && len(crawledPages.Pages) == crawlLimit {
				return false
			}

			workUrl, urlExists := selection.Attr("href")

			if !urlExists {
				return false
			}

			pageDoc, _ := goquery.NewDocumentFromReader(workReaders[i])
			subPagesUrls := crawler.CrawlWorkSubpages(pageDoc)

			errAddPage = crawledPages.AddTvTropesPage(workUrl, subPagesUrls)
			time.Sleep(time.Second / 2)

			return true
		})

		if limitedCrawling && len(crawledPages.Pages) == crawlLimit {
			break
		}

		if errAddPage != nil {
			return nil, fmt.Errorf("error crawling: %w", errAddPage)
		}

		pageNumber += 1
	}

	return crawledPages, nil
}
