package scraper_test

import (
	"net/url"
	"time"

	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/service/scraper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scraper", func() {
	// A scraper service for test purposes
	var serviceScraper *scraper.ServiceScraper
	var newScraperErr error
	// A valid TvTropes page and one page that is from other website
	var tvTropesPage, tvTropesPage2, tvTropesPage3, notTvTropesPage, notWorkPage *tropestogo.Page

	BeforeEach(func() {
		serviceScraper, newScraperErr = scraper.NewServiceScraper()

		tvTropesUrl, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003")
		tvTropesUrl2, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/TheAvengers2012")
		tvTropesUrl3, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/ANewHope")

		differentUrl, _ := url.Parse("https://www.google.com/")
		notWorkUrl, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Main/Media")

		tvTropesPage = &tropestogo.Page{
			URL:         tvTropesUrl,
			LastUpdated: time.Now(),
		}

		tvTropesPage2 = &tropestogo.Page{
			URL:         tvTropesUrl2,
			LastUpdated: time.Now(),
		}

		tvTropesPage3 = &tropestogo.Page{
			URL:         tvTropesUrl3,
			LastUpdated: time.Now(),
		}

		notTvTropesPage = &tropestogo.Page{
			URL:         differentUrl,
			LastUpdated: time.Now(),
		}

		notWorkPage = &tropestogo.Page{
			URL:         notWorkUrl,
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
		var validTvTropesPage, validTvTropesPage2, validTvTropesPage3, validDifferentPage, validNotWorkPage bool
		var errTvTropes, errTvTropes2, errTvTropes3, errDifferent, errNotWorkPage error

		BeforeEach(func() {
			validTvTropesPage, errTvTropes = serviceScraper.CheckValidWorkPage(tvTropesPage)
			validTvTropesPage2, errTvTropes2 = serviceScraper.CheckValidWorkPage(tvTropesPage2)
			validTvTropesPage3, errTvTropes3 = serviceScraper.CheckValidWorkPage(tvTropesPage3)
			validDifferentPage, errDifferent = serviceScraper.CheckValidWorkPage(notTvTropesPage)
			validNotWorkPage, errNotWorkPage = serviceScraper.CheckValidWorkPage(notWorkPage)
		})

		Context("URL belongs to a TvTropes Work page with tropes on a list", func() {
			It("Should mark the page as valid", func() {
				Expect(validTvTropesPage).To(BeTrue())
			})

			It("Shouldn't return an error", func() {
				Expect(errTvTropes).To(BeNil())
			})
		})

		Context("URL belongs to a TvTropes Work page with tropes on subpages", func() {
			It("Should mark the page as valid", func() {
				Expect(validTvTropesPage2).To(BeTrue())
			})

			It("Shouldn't return an error", func() {
				Expect(errTvTropes2).To(BeNil())
			})
		})

		Context("URL belongs to a TvTropes Work page with tropes on folders", func() {
			It("Should mark the page as valid", func() {
				Expect(validTvTropesPage3).To(BeTrue())
			})

			It("Shouldn't return an error", func() {
				Expect(errTvTropes3).To(BeNil())
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
