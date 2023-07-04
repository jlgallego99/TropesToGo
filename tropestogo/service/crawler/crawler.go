package crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	// The seed is the starting URL of the crawler
	seed = "https://tvtropes.org/pmwiki/pagelist_having_pagetype_in_namespace.php?t=work&n="

	// Common headers for a Firefox browser
	userAgentHeader               = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0"
	acceptHeader                  = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
	acceptLanguageHeader          = "es"
	upgradeInsecureRequestsHeader = "1"

	// TvTropes date formats
	tvTropesHistoryDateFormat = "Jan 2 2006 at 3:04:05 PM"

	// Referer header for a well-known and trusty webpage
	refererHeader = "https://www.google.com/"

	TvTropesHostname        = "tvtropes.org"
	TvTropesWeb             = "https://" + TvTropesHostname
	TvTropesPmwiki          = TvTropesWeb + "/pmwiki/"
	WorkPageSelector        = "table a"
	CurrentSubpageSelector  = ".curr-subpage"
	SubWikiSelector         = "a.subpage-link:not(" + CurrentSubpageSelector + ")"
	SubPageSelector         = "ul a.twikilink"
	PaginationNavSelector   = "nav.pagination-box > a"
	WorkHistoryPageSelector = "li.link-history a"
	LastUpdatedSelector     = "#main-article > div:first-of-type .pull-right a"
)

var (
	ErrNotFound    = errors.New("couldn't request the URL")
	ErrCrawling    = errors.New("there was an error crawling TvTropes")
	ErrEndIndex    = errors.New("there's no next page on the index")
	ErrParse       = errors.New("couldn't parse the HTML contents of the page")
	ErrLastUpdated = errors.New("couldn't retrieve o the last updated time")
	ErrParseTime   = errors.New("couldn't parse the TvTropes last updated time")

	httpClient = &http.Client{}

	// date ordinals for removing them on a date string
	dateOrdinals = []string{"st", "nd", "rd", "th"}

	// Seed for works that belongs to a certain MediaType
	mediaSeed = seed
)

type ServiceCrawler struct{}

func NewCrawler() *ServiceCrawler {
	crawler := &ServiceCrawler{}

	return crawler
}

// CrawlWorkPages searches crawlLimit number of Work pages belonging to a mediaType from the defined seed starting page
// if the crawlLimit is 0 or less, then it crawls all Work pages on the selected MediaType
// It returns a TvTropesPages object with all crawled pages and subpages from TvTropes
func (crawler *ServiceCrawler) CrawlWorkPages(crawlLimit int, mediaType media.MediaType) (*tropestogo.TvTropesPages, error) {
	crawledPages := tropestogo.NewTvTropesPages()

	mediaSeed = seed + mediaType.String()
	indexPage := mediaSeed

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
			log.Error().Err(errDoRequest).Msg("CRAWLING FAILED " + indexPage)
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

			log.Info().Msg("CRAWLING: " + workUrl)

			// Create the Work Page with its subpages
			workPage, errAddPage = crawler.createWorkPage(workUrl, crawledPages)
			if errAddPage != nil {
				log.Error().Err(errAddPage).Msg("CRAWLING WORK PAGE FAILED " + workUrl)
				return false
			}

			// Set LastUpdated time
			lastUpdated, errLastUpdated := crawler.getLastUpdated(workPage.GetDocument())
			if errLastUpdated != nil {
				errAddPage = errLastUpdated
				log.Error().Err(errAddPage).Msg("CRAWLING LAST UPDATE DATE FAILED " + workUrl)
				return false
			}
			crawledPages.Pages[workPage].LastUpdated = lastUpdated

			// Crawl Work subpages and add them
			errAddPage = crawler.addWorkSubpages(workPage, crawledPages)
			if errAddPage != nil {
				log.Error().Err(errAddPage).Msg("CRAWLING WORK SUBPAGES FAILED " + workUrl)
				return false
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
		indexPage, errAddPage = crawler.getNextPageUriFromDocument(doc)
		indexPage = TvTropesPmwiki + indexPage
		if errAddPage != nil {
			log.Error().Err(errAddPage).Msg("CRAWLING NEXT INDEX PAGE FAILED " + indexPage)
			break
		}
	}

	return crawledPages, nil
}

// getNextPageUriFromDocument, internal function that looks for the next pagination URI on the current index
// It looks for a "Next" button on the pagination navigator, and returns an error if there's no next page
// It works for any page with a pagination navigator, and the path can be different, so it returns only the URI
func (crawler *ServiceCrawler) getNextPageUriFromDocument(doc *goquery.Document) (string, error) {
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
		return "", fmt.Errorf("%w: "+mediaSeed, ErrEndIndex)
	}

	return nextPageUri, nil
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
			return nil, fmt.Errorf("%w: "+mediaSeed, ErrParse)
		}

		listSelector := doc.Find(WorkPageSelector)
		if listSelector.Length() == 0 {
			return nil, fmt.Errorf("%w: "+mediaSeed, ErrCrawling)
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

// CrawlChanges crawls the latest changes on TvTropes Films and returns a TvTropesPages with all recently-updated Work Pages
// Receives a map of already crawled works, relating a name with its last updated time
// and only crawls them if there's record of them on the history page, and it's newer
// Returns a TvTropesPages containing the crawled Pages of the Media that needs to be updated
func (crawler *ServiceCrawler) CrawlChanges(crawledWorks map[string]time.Time) (*tropestogo.TvTropesPages, error) {
	crawledPages := tropestogo.NewTvTropesPages()

	for crawledUrl, lastUpdated := range crawledWorks {
		// Create the Work Page
		newPage, errNewPage := crawler.createWorkPage(crawledUrl, crawledPages)
		if errNewPage != nil {
			return nil, errNewPage
		}

		newLastUpdated, errGetLastUpdated := crawler.getLastUpdated(newPage.GetDocument())
		if errGetLastUpdated != nil {
			return nil, errGetLastUpdated
		}

		// Only crawl the page if it has been crawled before, and it's been updated
		if newLastUpdated.After(lastUpdated) {
			// Set LastUpdated time
			crawledPages.Pages[newPage].LastUpdated = newLastUpdated

			// Crawl Work subpages and add them
			errCrawlSubpages := crawler.addWorkSubpages(newPage, crawledPages)
			if errCrawlSubpages != nil {
				return nil, errCrawlSubpages
			}
		}
	}

	return crawledPages, nil
}

// createWorkPage forms a valid Work Page object and adds it to the crawledPages object
func (crawler *ServiceCrawler) createWorkPage(workUrl string, crawledPages *tropestogo.TvTropesPages) (tropestogo.Page, error) {
	validRequest, errRequest := crawler.makeValidRequest(workUrl)
	if errRequest != nil {
		return tropestogo.Page{}, errRequest
	}
	workPage, errAddPage := crawledPages.AddTvTropesPage(workUrl, true, validRequest)
	if errAddPage != nil {
		return tropestogo.Page{}, errAddPage
	}

	return workPage, nil
}

// addWorkSubpages crawls all Work subpages, creates them and adds them to the referenced crawledPages argument
func (crawler *ServiceCrawler) addWorkSubpages(workPage tropestogo.Page, crawledPages *tropestogo.TvTropesPages) error {
	// Search for subpages on the new Work Page
	subPagesUrls := crawler.CrawlWorkSubpages(workPage.GetDocument())

	// Add its subpages to the Work Page
	var requests []*http.Request
	for _, subPagesUrl := range subPagesUrls {
		validRequest, errRequest := crawler.makeValidRequest(subPagesUrl)
		if errRequest != nil {
			return errRequest
		}

		requests = append(requests, validRequest)
	}

	errSubpages := crawledPages.AddSubpages(workPage.GetUrl().String(), subPagesUrls, true, requests)

	// If there's been too many requests to TvTropes, wait longer
	if errors.Is(errSubpages, tropestogo.ErrForbidden) {
		time.Sleep(time.Minute)
		errSubpages = nil
	}

	return errSubpages
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

// GetLastUpdated retrieves the last updated date from the history page of a Work page and parses it to a valid time object
// If it couldn't be parsed or obtained, it will return an ErrLastUpdated error
func (crawler *ServiceCrawler) getLastUpdated(doc *goquery.Document) (time.Time, error) {
	historyPageUri, historyPageExists := doc.Find(WorkHistoryPageSelector).First().Attr("href")
	if !historyPageExists {
		return time.Time{}, nil
	}

	request, errRequest := crawler.makeValidRequest(TvTropesWeb + historyPageUri)
	if errRequest != nil {
		return time.Time{}, errRequest
	}

	resp, errDoRequest := httpClient.Do(request)
	if errDoRequest != nil {
		return time.Time{}, fmt.Errorf("%w because there was an error on the HTTP request to the history ", ErrLastUpdated)
	}

	historyDoc, _ := goquery.NewDocumentFromReader(resp.Body)
	lastUpdated, errLastUpdated := crawler.ParseTvTropesTime(historyDoc)
	if errLastUpdated != nil {
		return time.Time{}, errLastUpdated
	}

	return lastUpdated, nil
}

// ParseTvTropesTime searches for the last updated time in a work history page and parses it to a valid time object
// If it can't be parsed it will return an ErrParseTime error
func (crawler *ServiceCrawler) ParseTvTropesTime(historyDoc *goquery.Document) (time.Time, error) {
	lastUpdatedString := historyDoc.Find(LastUpdatedSelector).Text()
	for _, ordinal := range dateOrdinals {
		lastUpdatedString = strings.ReplaceAll(lastUpdatedString, ordinal, "")
	}

	lastUpdated, errParseTime := time.Parse(tvTropesHistoryDateFormat, lastUpdatedString)
	if errParseTime != nil {
		return time.Time{}, fmt.Errorf("%w: "+lastUpdatedString, ErrParseTime)
	}

	return lastUpdated, nil
}
