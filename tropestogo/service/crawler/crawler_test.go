package crawler_test

import (
	tropestogo "github.com/jlgallego99/TropesToGo"
	crawler "github.com/jlgallego99/TropesToGo/service/crawler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// A crawler service for test purposes
var serviceCrawler *crawler.ServiceCrawler
var errNewCrawler, errCrawling error
var crawledPages *tropestogo.TvTropesPages

var _ = BeforeSuite(func() {
	serviceCrawler, errNewCrawler = crawler.NewCrawler("Film")
	Expect(errNewCrawler).To(BeNil())
	crawledPages, errCrawling = serviceCrawler.CrawlWorkPages(10)
})

var _ = Describe("Crawler", func() {
	Context("The Crawler gets all Films", func() {
		It("Shouldn't return an error", func() {
			Expect(errCrawling).To(BeNil())
		})

		It("Should have crawled web pages", func() {
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
