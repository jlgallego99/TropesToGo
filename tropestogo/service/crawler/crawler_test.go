package crawler_test

import (
	tropestogo "github.com/jlgallego99/TropesToGo"
	crawler "github.com/jlgallego99/TropesToGo/service/crawler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"os"
)

const (
	indexResource = "resources/film_index_page1.html"
)

var filmResources = []string{"resources/film1.html", "resources/film2.html", "resources/film3.html",
	"resources/film4.html", "resources/film5.html"}

// A crawler service for test purposes
var serviceCrawler *crawler.ServiceCrawler
var errNewCrawler, errCrawling error
var crawledPages *tropestogo.TvTropesPages

var _ = BeforeSuite(func() {
	serviceCrawler, errNewCrawler = crawler.NewCrawler("Film")
	Expect(errNewCrawler).To(BeNil())

	indexReader, _ := os.Open(indexResource)
	workReaders := make([]io.Reader, 0)
	for _, filmResource := range filmResources {
		workResource, _ := os.Open(filmResource)
		workReaders = append(workReaders, workResource)
	}

	crawledPages, errCrawling = serviceCrawler.CrawlWorkPagesFromReaders(indexReader, workReaders, 5)
})

var _ = Describe("Crawler", func() {
	Context("Crawling a limited number of Work Pages from the Index", func() {
		It("Shouldn't return an error", func() {
			Expect(errCrawling).To(BeNil())
		})

		It("Should have crawled Work Pages and its subpages", func() {
			Expect(len(crawledPages.Pages) > 0).To(BeTrue())

			for crawledPage, crawledSubpages := range crawledPages.Pages {
				Expect(crawledPage.GetUrl()).To(Not(BeNil()))
				Expect(crawledPage.GetPageType()).To(Equal(tropestogo.WorkPage))
				Expect(len(crawledSubpages.Subpages) >= 0).To(BeTrue())

				for crawledSubpage := range crawledSubpages.Subpages {
					Expect(crawledSubpage.GetUrl()).To(Not(BeNil()))
				}
			}
		})
	})
})
