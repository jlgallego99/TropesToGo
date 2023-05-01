package crawler

import "net/url"

type ServiceCrawler struct {
	// Seeds are the base URLs of the crawler
	seeds []url.URL
}

func NewCrawler() (*ServiceCrawler, error) {
	return nil, nil
}
