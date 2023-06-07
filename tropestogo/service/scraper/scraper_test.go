package scraper_test

import (
	"github.com/jlgallego99/TropesToGo/media"
	"net/url"
	"time"

	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/service/scraper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// A scraper service for test purposes
var serviceScraper *scraper.ServiceScraper
var newScraperErr, invalidScraperErr error

var tvTropesPage, tvTropesPage2, tvTropesPage3, notTvTropesPage, notWorkPage, unsupportedMediaPage *tropestogo.Page

var _ = BeforeSuite(func() {
	serviceScraper, newScraperErr = scraper.NewServiceScraper()
	// Create invalid scraper
	_, invalidScraperErr = scraper.NewServiceScraper(scraper.ConfigIndexRepository(nil))

	tvTropesUrl, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003")
	tvTropesUrl2, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/TheAvengers2012")
	tvTropesUrl3, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/ANewHope")

	tvTropesUrlUnsupported, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Manga/AttackOnTitan")

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

	unsupportedMediaPage = &tropestogo.Page{
		URL:         tvTropesUrlUnsupported,
		LastUpdated: time.Now(),
	}
})

var _ = Describe("Scraper", func() {
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
				Expect(invalidScraperErr).To(Equal(scraper.ErrInvalidField))
			})
		})
	})

	Describe("Check if page can be scraped", func() {
		var validTvTropesPage, validTvTropesPage2, validTvTropesPage3, validDifferentPage, validNotWorkPage bool
		var errTvTropes, errTvTropes2, errTvTropes3, errDifferent, errNotWorkPage error

		Context("URL belongs to a TvTropes Work page with tropes on a list", func() {
			BeforeEach(func() {
				validTvTropesPage, errTvTropes = serviceScraper.CheckValidWorkPage(tvTropesPage)
			})

			It("Should mark the page as valid", func() {
				Expect(validTvTropesPage).To(BeTrue())
			})

			It("Shouldn't return an error", func() {
				Expect(errTvTropes).To(BeNil())
			})
		})

		Context("URL belongs to a TvTropes Work page with tropes on subpages", func() {
			BeforeEach(func() {
				validTvTropesPage2, errTvTropes2 = serviceScraper.CheckValidWorkPage(tvTropesPage2)
			})

			It("Should mark the page as valid", func() {
				Expect(validTvTropesPage2).To(BeTrue())
			})

			It("Shouldn't return an error", func() {
				Expect(errTvTropes2).To(BeNil())
			})
		})

		Context("URL belongs to a TvTropes Work page with tropes on folders", func() {
			BeforeEach(func() {
				validTvTropesPage3, errTvTropes3 = serviceScraper.CheckValidWorkPage(tvTropesPage3)
			})

			It("Should mark the page as valid", func() {
				Expect(validTvTropesPage3).To(BeTrue())
			})

			It("Shouldn't return an error", func() {
				Expect(errTvTropes3).To(BeNil())
			})
		})

		Context("URL belongs to TvTropes but isn't from a Work page", func() {
			BeforeEach(func() {
				validNotWorkPage, errNotWorkPage = serviceScraper.CheckValidWorkPage(notWorkPage)
			})

			It("Should mark the page as invalid", func() {
				Expect(validNotWorkPage).To(BeFalse())
			})

			It("Should return an appropriate error", func() {
				Expect(errNotWorkPage).To(Equal(scraper.ErrNotWorkPage))
			})
		})

		Context("URL isn't from TvTropes", func() {
			BeforeEach(func() {
				validDifferentPage, errDifferent = serviceScraper.CheckValidWorkPage(notTvTropesPage)
			})

			It("Should mark the page as invalid", func() {
				Expect(validDifferentPage).To(BeFalse())
			})

			It("Should return an appropriate error", func() {
				Expect(errDifferent).To(Equal(scraper.ErrNotTvTropes))
			})
		})
	})

	Describe("Scrape Film Page", func() {
		var validFilm1, validFilm3, filmInvalidType media.Media
		var errorFilm1, errorFilm3, errorFilmInvalidType error

		Context("Valid Film Page with tropes on a simple list", func() {
			BeforeEach(func() {
				validFilm1, errorFilm1 = serviceScraper.ScrapeWorkPage(tvTropesPage)
			})

			It("Should have correct Work fields", func() {
				Expect(validFilm1.GetWork().Title).To(Equal("Oldboy (2003)"))
				Expect(validFilm1.GetWork().Year).To(Equal("2003"))
				Expect(validFilm1.GetMediaType()).To(Equal(media.Film))
				Expect(validFilm1.GetWork().Tropes).To(Not(BeEmpty()))
			})

			It("Shouldn't return an error", func() {
				Expect(errorFilm1).To(BeNil())
			})

			It("Shouldn't have repeated tropes", func() {
				unique := areTropesUnique(validFilm1.GetWork().Tropes)

				Expect(unique).To(BeTrue())
			})
		})

		Context("Valid Film Page with tropes on folders", func() {
			BeforeEach(func() {
				validFilm3, errorFilm3 = serviceScraper.ScrapeWorkPage(tvTropesPage3)
			})

			It("Should have correct Work fields", func() {
				Expect(validFilm3.GetWork().Title).To(Equal("A New Hope"))
				Expect(validFilm3.GetWork().Year).To(BeEmpty())
				Expect(validFilm3.GetMediaType()).To(Equal(media.Film))
				Expect(validFilm3.GetWork().Tropes).To(Not(BeEmpty()))
			})

			It("Shouldn't return an error", func() {
				Expect(errorFilm3).To(BeNil())
			})

			It("Shouldn't have repeated tropes", func() {
				unique := areTropesUnique(validFilm3.GetWork().Tropes)

				Expect(unique).To(BeTrue())
			})
		})

		Context("Invalid Film because the media type isn't supported", func() {
			BeforeEach(func() {
				filmInvalidType, errorFilmInvalidType = serviceScraper.ScrapeWorkPage(unsupportedMediaPage)
			})

			It("Should return an empty media object", func() {
				Expect(filmInvalidType.GetWork()).To(BeNil())
				Expect(filmInvalidType.GetPage()).To(BeNil())
				Expect(filmInvalidType.GetMediaType()).To(Equal(media.UnknownMediaType))
			})

			It("Should return an appropriate error", func() {
				Expect(errorFilmInvalidType).To(Equal(media.ErrUnsupportedMediaType))
			})
		})
	})
})

func areTropesUnique(tropes map[tropestogo.Trope]struct{}) bool {
	visited := make(map[string]bool, 0)
	for trope := range tropes {
		if visited[trope.GetTitle()] == true {
			return false
		} else {
			visited[trope.GetTitle()] = true
		}
	}

	return true
}
