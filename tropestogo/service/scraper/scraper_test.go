package scraper_test

import (
	"github.com/PuerkitoBio/goquery"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/service/scraper"
	"net/http"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scraper", func() {
	// A scraper service for test purposes
	var serviceScraper *scraper.ServiceScraper
	var newScraperErr error
	// A valid TvTropes page and one page that is from other website
	var tvTropesPage, notTvTropesPage, notWorkPage *tropestogo.Page
	// HTTP request to a TvTropes Work page
	var res *http.Response
	// DOM Tree of a TvTropes Work page
	var doc *goquery.Document

	BeforeEach(func() {
		serviceScraper, newScraperErr = scraper.NewServiceScraper()

		res, _ = http.Get("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003")
		doc, _ = goquery.NewDocumentFromReader(res.Body)

		tvTropesUrl, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003")
		differentUrl, _ := url.Parse("https://www.google.com/")
		notWorkUrl, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Main/Media")

		tvTropesPage = &tropestogo.Page{
			URL:         tvTropesUrl,
			Document:    doc,
			LastUpdated: time.Now(),
		}

		notTvTropesPage = &tropestogo.Page{
			URL:         differentUrl,
			Document:    &goquery.Document{},
			LastUpdated: time.Now(),
		}

		notWorkPage = &tropestogo.Page{
			URL:         notWorkUrl,
			Document:    &goquery.Document{},
			LastUpdated: time.Now(),
		}
	})

	Describe("Create the scraper service", func() {
		Context("The service is created correctly", func() {
			It("Shouldn't return an error", func() {
				Expect(newScraperErr).To(BeNil())
			})
		})

		Context("The service is created incorrectly", func() {
			It("Should return an empty ServiceScraper", func() {
				Expect(*serviceScraper).To(Equal(scraper.ServiceScraper{}))
			})

			It("Should return an appropriate error", func() {
				Expect(newScraperErr).To(Equal(scraper.ErrInvalidField))
			})
		})
	})

	Describe("Check page URL", func() {
		var validTvTropesPage, validDifferentPage, validNotWorkPage bool
		var errTvTropes, errDifferent, errNotWorkPage error

		BeforeEach(func() {
			validTvTropesPage, errTvTropes = serviceScraper.CheckValidWorkPage(tvTropesPage)
			validDifferentPage, errDifferent = serviceScraper.CheckValidWorkPage(notTvTropesPage)
			validNotWorkPage, errNotWorkPage = serviceScraper.CheckValidWorkPage(notWorkPage)
		})

		Context("URL belongs to a TvTropes Work page", func() {
			It("Should mark the page as valid", func() {
				Expect(validTvTropesPage).To(BeTrue())
			})

			It("Shouldn't return an error", func() {
				Expect(errTvTropes).To(BeNil())
			})
		})

		Context("URL belongs to TvTropes but isn't from a Work page", func() {
			It("Should mark the page as invalid", func() {
				Expect(validNotWorkPage).To(BeFalse())
			})

			It("Should return an appropriate error", func() {
				Expect(errNotWorkPage).To(Equal(scraper.ErrNotWorkPage))
			})
		})

		Context("URL isn't from TvTropes", func() {
			It("Should mark the page as invalid", func() {
				Expect(validDifferentPage).To(BeFalse())
			})

			It("Should return an appropriate error", func() {
				Expect(errDifferent).To(Equal(scraper.ErrNotTvTropes))
			})
		})
	})
})
