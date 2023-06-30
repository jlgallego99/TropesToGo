package crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	// The seed is the starting URL of the crawler
	seed = "https://tvtropes.org/pmwiki/pagelist_having_pagetype_in_namespace.php?t=work&n=Film"

	// Common headers for a Firefox browser
	userAgentHeader               = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0"
	acceptHeader                  = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
	acceptLanguageHeader          = "es"
	upgradeInsecureRequestsHeader = "1"

	// Referer header for a well-known and trusty webpage
	refererHeader = "https://www.google.com/"

	TvTropesHostname       = "tvtropes.org"
	TvTropesWeb            = "https://" + TvTropesHostname
	TvTropesPmwiki         = TvTropesWeb + "/pmwiki/"
	WorkPageSelector       = "table a"
	CurrentSubpageSelector = ".curr-subpage"
	SubWikiSelector        = "a.subpage-link:not(" + CurrentSubpageSelector + ")"
	SubPageSelector        = "ul a.twikilink"
	PaginationNavSelector  = "nav.pagination-box > a"
)

var (
	ErrNotFound = errors.New("couldn't request the URL")
	ErrCrawling = errors.New("there was an error crawling TvTropes")
	ErrEndIndex = errors.New("there's no next page on the index")
	ErrParse    = errors.New("couldn't parse the HTML contents of the page")

	httpClient = &http.Client{}
)

type ServiceCrawler struct{}

func NewCrawler() *ServiceCrawler {
	crawler := &ServiceCrawler{}

	return crawler
}

// CrawlWorkPages searches crawlLimit number of Work pages from the defined seed starting page; if it's 0 or less, then it crawls all Work pages
// It returns a TvTropesPages object with all crawled pages and subpages from TvTropes
func (crawler *ServiceCrawler) CrawlWorkPages(crawlLimit int) (*tropestogo.TvTropesPages, error) {
	crawledPages := tropestogo.NewTvTropesPages()
	indexPage := seed

	limitedCrawling := true
	if crawlLimit <= 0 {
		limitedCrawling = false
	}

	for {
		request, errValidRequest := crawler.makeValidRequest(indexPage)
		if errValidRequest != nil {
			return nil, errValidRequest
		}

		resp, errDoRequest := httpClient.Do(request)
		if errDoRequest != nil {
			return nil, fmt.Errorf("%w: "+indexPage, ErrNotFound)
		}

		doc, errDocument := goquery.NewDocumentFromReader(resp.Body)
		if errDocument != nil {
			return nil, fmt.Errorf("%w: "+indexPage, ErrParse)
		}

		pageSelector := doc.Find(WorkPageSelector)
		if pageSelector.Length() == 0 {
			return nil, fmt.Errorf("%w: "+indexPage, ErrCrawling)
		}

		var errAddPage error
		var workPage tropestogo.Page
		pageSelector.EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if limitedCrawling && len(crawledPages.Pages) == crawlLimit {
				return false
			}

			workUrl, urlExists := selection.Attr("href")
			if !urlExists {
				return false
			}

			// Create the Work Page
			validRequest, errRequest := crawler.makeValidRequest(workUrl)
			if errRequest != nil {
				errAddPage = errRequest
				return false
			}
			workPage, errAddPage = crawledPages.AddTvTropesPage(workUrl, true, validRequest)
			if errors.Is(errAddPage, tropestogo.ErrForbidden) {
				time.Sleep(time.Minute)
			}

			// Search for subpages on the new Work Page
			subPagesUrls := crawler.CrawlWorkSubpages(workPage.GetDocument())

			// Add its subpages to the Work Page
			var requests []*http.Request
			for _, subPagesUrl := range subPagesUrls {
				validRequest, errRequest = crawler.makeValidRequest(subPagesUrl)
				if errRequest != nil {
					errAddPage = errRequest
					return false
				}

				requests = append(requests, validRequest)
			}
			errAddPage = crawledPages.AddSubpages(workUrl, subPagesUrls, true, requests)

			// If there's been too many requests to TvTropes, wait longer
			if errors.Is(errAddPage, tropestogo.ErrForbidden) {
				time.Sleep(time.Minute)
			}

			return true
		})

		if limitedCrawling && len(crawledPages.Pages) == crawlLimit {
			break
		}

		if errAddPage != nil {
			return nil, fmt.Errorf("error crawling: %w", errAddPage)
		}

		// Get next index page for crawling
		indexPage, errAddPage = crawler.getNextPageUrl(doc)
		if errAddPage != nil {
			break
		}
	}

	return crawledPages, nil
}

// getNextPageUrl, internal function that looks for the next Work index pagination URL on the current index page
// It looks for a "Next" button on the pagination navigator, and returns an error if there's no next page
func (crawler *ServiceCrawler) getNextPageUrl(doc *goquery.Document) (string, error) {
	// Search the "Next" button on the nav pagination
	nextPageUri := ""
	var nextPageExists bool
	doc.Find(PaginationNavSelector).EachWithBreak(func(_ int, selection *goquery.Selection) bool {
		nextPageUri, nextPageExists = selection.Attr("href")

		if nextPageExists && strings.EqualFold(selection.Find("a span.mobile-off").Text(), "Next") {
			return false
		} else {
			nextPageUri = ""
			return true
		}
	})

	if nextPageUri == "" {
		return "", fmt.Errorf("%w: "+seed, ErrEndIndex)
	}

	return TvTropesPmwiki + nextPageUri, nil
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
	doc.Find(SubPageSelector).EachWithBreak(func(_ int, selection *goquery.Selection) bool {
		subPageUri, subPageExists := selection.Attr("href")
		r, _ := regexp.Compile(`\/tropes[a-z]to[a-z]`)
		matchUri := r.MatchString(strings.ToLower(subPageUri))

		if subPageExists && matchUri {
			subPagesUrls = append(subPagesUrls, TvTropesWeb+subPageUri)

			return true
		} else {
			return false
		}
	})

	return subPagesUrls
}

// CrawlWorkPagesFromReaders crawls all Work Pages and its subpages from an index reader and its pages readers. Only for test purposes
// It searches crawlLimit number of Work pages within the index
// It returns a TvTropesPages object with all crawled pages and subpages from TvTropes
func (crawler *ServiceCrawler) CrawlWorkPagesFromReaders(indexReader io.Reader, workReaders []io.Reader, crawlLimit int) (*tropestogo.TvTropesPages, error) {
	crawledPages := tropestogo.NewTvTropesPages()

	limitedCrawling := true
	if crawlLimit <= 0 {
		limitedCrawling = false
	}

	for {
		doc, errDocument := goquery.NewDocumentFromReader(indexReader)
		if errDocument != nil {
			return nil, fmt.Errorf("%w: "+seed, ErrParse)
		}

		listSelector := doc.Find(WorkPageSelector)
		if listSelector.Length() == 0 {
			return nil, fmt.Errorf("%w: "+seed, ErrCrawling)
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

			_, errAddPage = crawledPages.AddTvTropesPage(workUrl, false, nil)
			errAddPage = crawledPages.AddSubpages(workUrl, subPagesUrls, false, nil)

			return true
		})

		if limitedCrawling && len(crawledPages.Pages) == crawlLimit {
			break
		}

		if errAddPage != nil {
			return nil, fmt.Errorf("error crawling: %w", errAddPage)
		}
	}

	return crawledPages, nil
}

// makeValidRequests builds an HTTP request to the url page and returns its contents
// The request sets very specific Headers to pass as a real browser, avoiding banning for being a bot
// It returns an ErrNotFound error if the request couldn't be made
func (crawler *ServiceCrawler) makeValidRequest(pageUrl string) (*http.Request, error) {
	request, errRequest := http.NewRequest("GET", pageUrl, nil)
	if errRequest != nil {
		return nil, fmt.Errorf("%w: "+pageUrl, ErrNotFound)
	}

	request.Header.Set("User-Agent", userAgentHeader)
	request.Header.Set("Referer", refererHeader)
	request.Header.Set("Accept", acceptHeader)
	request.Header.Set("Accept-Language", acceptLanguageHeader)
	request.Header.Set("Upgrade-Insecure-Requests", upgradeInsecureRequestsHeader)

	return request, nil
}
