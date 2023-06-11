package scraper_test

import (
	"github.com/jlgallego99/TropesToGo/media/csv_dataset"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	"net/url"
	"time"

	"github.com/jlgallego99/TropesToGo/media"

	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/service/scraper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// A scraper service for test purposes
var serviceScraperJson, serviceScraperCsv, invalidScraper *scraper.ServiceScraper
var newScraperJsonErr, newScraperCsvErr, invalidScraperErr error
var csvRepositoryErr, jsonRepositoryErr error

var tvTropesPage, tvTropesPage2, tvTropesPage3, notTvTropesPage, notWorkPage, unknownMediaPage *tropestogo.Page

var _ = BeforeSuite(func() {
	// Create two scrapers, one for the JSON dataset and the other for the CSV dataset
	var csvRepository, jsonRepository media.RepositoryMedia
	csvRepository, csvRepositoryErr = csv_dataset.NewCSVRepository("dataset", ',')
	jsonRepository, jsonRepositoryErr = json_dataset.NewJSONRepository("dataset")
	serviceScraperJson, newScraperJsonErr = scraper.NewServiceScraper(scraper.ConfigMediaRepository(csvRepository))
	serviceScraperCsv, newScraperCsvErr = scraper.NewServiceScraper(scraper.ConfigMediaRepository(jsonRepository))

	// Create invalid scraper
	invalidScraper, invalidScraperErr = scraper.NewServiceScraper(scraper.ConfigIndexRepository(nil))

	tvTropesUrl, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003")
	tvTropesUrl2, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/TheAvengers2012")
	tvTropesUrl3, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/ANewHope")

	tvTropesUrlUnknown, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Manga/AttackOnTitan")

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

	unknownMediaPage = &tropestogo.Page{
		URL:         tvTropesUrlUnknown,
		LastUpdated: time.Now(),
	}
})

var _ = Describe("Scraper", func() {
	Describe("Create the scraper services", func() {
		Context("The services are created correctly", func() {
			It("Shouldn't return an error", func() {
				Expect(newScraperJsonErr).To(BeNil())
				Expect(newScraperCsvErr).To(BeNil())
			})

			It("Should have a correct media repository", func() {

			})
		})

		Context("The service is created incorrectly", func() {
			It("Should return an empty ServiceScraper", func() {
				Expect(invalidScraper).To(BeNil())
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
				validTvTropesPage, errTvTropes = serviceScraperJson.CheckValidWorkPage(tvTropesPage)
				validTvTropesPage, errTvTropes = serviceScraperCsv.CheckValidWorkPage(tvTropesPage)
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
				validTvTropesPage2, errTvTropes2 = serviceScraperJson.CheckValidWorkPage(tvTropesPage2)
				validTvTropesPage2, errTvTropes2 = serviceScraperCsv.CheckValidWorkPage(tvTropesPage2)
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
				validTvTropesPage3, errTvTropes3 = serviceScraperJson.CheckValidWorkPage(tvTropesPage3)
				validTvTropesPage3, errTvTropes3 = serviceScraperCsv.CheckValidWorkPage(tvTropesPage3)
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
				validNotWorkPage, errNotWorkPage = serviceScraperJson.CheckValidWorkPage(notWorkPage)
				validNotWorkPage, errNotWorkPage = serviceScraperCsv.CheckValidWorkPage(notWorkPage)
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
				validDifferentPage, errDifferent = serviceScraperJson.CheckValidWorkPage(notTvTropesPage)
				validDifferentPage, errDifferent = serviceScraperCsv.CheckValidWorkPage(notTvTropesPage)
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
				validFilm1, errorFilm1 = serviceScraperJson.ScrapeWorkPage(tvTropesPage)
				validFilm1, errorFilm1 = serviceScraperCsv.ScrapeWorkPage(tvTropesPage)
			})

			It("Should have correct fields", func() {
				testValidScrapedMedia(validFilm1, "Oldboy (2003)", "2003", media.Film)
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
				validFilm3, errorFilm3 = serviceScraperJson.ScrapeWorkPage(tvTropesPage3)
				validFilm3, errorFilm3 = serviceScraperCsv.ScrapeWorkPage(tvTropesPage3)
			})

			It("Should have correct fields", func() {
				testValidScrapedMedia(validFilm3, "A New Hope", "", media.Film)
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
				filmInvalidType, errorFilmInvalidType = serviceScraperJson.ScrapeWorkPage(unknownMediaPage)
				filmInvalidType, errorFilmInvalidType = serviceScraperCsv.ScrapeWorkPage(unknownMediaPage)
			})

			It("Should return an empty media object", func() {
				Expect(filmInvalidType.GetWork()).To(BeNil())
				Expect(filmInvalidType.GetPage()).To(BeNil())
				Expect(filmInvalidType.GetMediaType()).To(Equal(media.UnknownMediaType))
			})

			It("Should return an appropriate error", func() {
				Expect(errorFilmInvalidType).To(Equal(media.ErrUnknownMediaType))
			})
		})
	})
})

func testValidScrapedMedia(validMedia media.Media, title, year string, mediaType media.MediaType) {
	Expect(validMedia.GetWork().Title).To(Equal(title))
	Expect(validMedia.GetWork().Year).To(Equal(year))
	Expect(validMedia.GetMediaType()).To(Equal(mediaType))
	Expect(validMedia.GetWork().Tropes).To(Not(BeEmpty()))
}

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
