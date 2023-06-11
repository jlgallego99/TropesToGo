package scraper_test

import (
	"encoding/csv"
	"encoding/json"
	"github.com/jlgallego99/TropesToGo/media/csv_dataset"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	"net/url"
	"os"
	"strings"
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
var csvRepository, jsonRepository media.RepositoryMedia

var tvTropesPage, tvTropesPage2, tvTropesPage3, notTvTropesPage, notWorkPage, unknownMediaPage *tropestogo.Page

var _ = BeforeSuite(func() {
	// Create two scrapers, one for the JSON dataset and the other for the CSV dataset
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
				Expect(csvRepositoryErr).To(BeNil())
				Expect(jsonRepositoryErr).To(BeNil())
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
		Context("Valid Film Page with tropes on a simple list", func() {
			var validfilm1Csv, validfilm1Json media.Media
			var errorfilm1Csv, errorfilm1Json error

			BeforeEach(func() {
				validfilm1Json, errorfilm1Json = serviceScraperJson.ScrapeWorkPage(tvTropesPage)
				validfilm1Csv, errorfilm1Csv = serviceScraperCsv.ScrapeWorkPage(tvTropesPage)
			})

			It("Shouldn't return an error", func() {
				Expect(errorfilm1Json).To(BeNil())
				Expect(errorfilm1Csv).To(BeNil())
			})

			It("Should have correct fields", func() {
				testValidScrapedMedia(validfilm1Csv, "Oldboy (2003)", "2003", media.Film)
				testValidScrapedMedia(validfilm1Json, "Oldboy (2003)", "2003", media.Film)

				Expect(errorfilm1Csv).To(Equal(csv_dataset.ErrDuplicatedMedia))
				Expect(errorfilm1Json).To(Equal(json_dataset.ErrDuplicatedMedia))
			})

			It("Shouldn't have repeated tropes", func() {
				uniqueCsv := areTropesUnique(validfilm1Csv.GetWork().Tropes)
				uniqueJson := areTropesUnique(validfilm1Json.GetWork().Tropes)

				Expect(uniqueCsv).To(BeTrue())
				Expect(uniqueJson).To(BeTrue())

				Expect(errorfilm1Csv).To(Equal(csv_dataset.ErrDuplicatedMedia))
				Expect(errorfilm1Json).To(Equal(json_dataset.ErrDuplicatedMedia))
			})

			It("Should have added a correct record on the JSON repository", func() {
				var dataset json_dataset.JSONDataset
				fileContents, _ := os.ReadFile("dataset.json")
				err := json.Unmarshal(fileContents, &dataset)

				Expect(err).To(BeNil())
				Expect(dataset.Tropestogo[0].Title).To(Equal("Oldboy (2003)"))
				Expect(dataset.Tropestogo[0].Year).To(Equal("2003"))
				Expect(dataset.Tropestogo[0].URL).To(Equal("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003"))
				Expect(dataset.Tropestogo[0].MediaType).To(Equal("Film"))
				Expect(len(dataset.Tropestogo[0].Tropes) > 0).To(BeTrue())
			})

			It("Should have added a correct record on the CSV repository", func() {
				f, errOpen := os.Open("dataset.csv")
				reader := csv.NewReader(f)
				records, errReadCSV := reader.ReadAll()

				Expect(errOpen).To(BeNil())
				Expect(errReadCSV).To(BeNil())

				Expect(len(records[0])).To(Equal(7))
				Expect(len(records[1])).To(Equal(7))
				Expect(records[1][0]).To(Equal("Oldboy (2003)"))
				Expect(records[1][1]).To(Equal("2003"))
				Expect(records[1][3]).To(Equal("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003"))
				Expect(records[1][4]).To(Equal("Film"))
				Expect(len(strings.Split(records[1][5], ";")) > 0).To(BeTrue())
			})
		})

		Context("Valid Film Page with tropes on folders", func() {
			var validfilm3Csv, validfilm3Json media.Media
			var errorfilm3Csv, errorfilm3Json error

			BeforeEach(func() {
				validfilm3Json, errorfilm3Json = serviceScraperJson.ScrapeWorkPage(tvTropesPage3)
				validfilm3Csv, errorfilm3Csv = serviceScraperCsv.ScrapeWorkPage(tvTropesPage3)
			})

			It("Shouldn't return an error", func() {
				Expect(errorfilm3Csv).To(BeNil())
				Expect(errorfilm3Json).To(BeNil())
			})

			It("Should have correct fields", func() {
				testValidScrapedMedia(validfilm3Csv, "A New Hope", "", media.Film)
				testValidScrapedMedia(validfilm3Json, "A New Hope", "", media.Film)

				Expect(errorfilm3Csv).To(Equal(csv_dataset.ErrDuplicatedMedia))
				Expect(errorfilm3Json).To(Equal(json_dataset.ErrDuplicatedMedia))
			})

			It("Shouldn't have repeated tropes", func() {
				uniqueCsv := areTropesUnique(validfilm3Csv.GetWork().Tropes)
				uniqueJson := areTropesUnique(validfilm3Json.GetWork().Tropes)

				Expect(uniqueCsv).To(BeTrue())
				Expect(uniqueJson).To(BeTrue())

				Expect(errorfilm3Csv).To(Equal(csv_dataset.ErrDuplicatedMedia))
				Expect(errorfilm3Json).To(Equal(json_dataset.ErrDuplicatedMedia))
			})
		})

		Context("Invalid Film because the media type isn't supported", func() {
			var filminvalidtypeJson, filminvalidtypeCsv media.Media
			var errorfilminvalidtypeJson, errorfilminvalidtypeCsv error

			BeforeEach(func() {
				filminvalidtypeJson, errorfilminvalidtypeJson = serviceScraperJson.ScrapeWorkPage(unknownMediaPage)
				filminvalidtypeCsv, errorfilminvalidtypeCsv = serviceScraperCsv.ScrapeWorkPage(unknownMediaPage)
			})

			It("Should return an empty media object", func() {
				Expect(filminvalidtypeJson.GetWork()).To(BeNil())
				Expect(filminvalidtypeJson.GetPage()).To(BeNil())
				Expect(filminvalidtypeJson.GetMediaType()).To(Equal(media.UnknownMediaType))

				Expect(filminvalidtypeCsv.GetWork()).To(BeNil())
				Expect(filminvalidtypeCsv.GetPage()).To(BeNil())
				Expect(filminvalidtypeCsv.GetMediaType()).To(Equal(media.UnknownMediaType))
			})

			It("Should return an appropriate error", func() {
				Expect(errorfilminvalidtypeJson).To(Equal(media.ErrUnknownMediaType))
				Expect(errorfilminvalidtypeCsv).To(Equal(media.ErrUnknownMediaType))
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
