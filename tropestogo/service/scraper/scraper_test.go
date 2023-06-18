package scraper_test

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"github.com/jlgallego99/TropesToGo/media"
	"github.com/jlgallego99/TropesToGo/media/csv_dataset"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	"io"
	"net/url"
	"os"
	"strings"

	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/service/scraper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	oldboyUrl        = "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003"
	anewhopeUrl      = "https://tvtropes.org/pmwiki/pmwiki.php/Film/ANewHope"
	avengersUrl      = "https://tvtropes.org/pmwiki/pmwiki.php/Film/TheAvengers2012"
	mediaUrl         = "https://tvtropes.org/pmwiki/pmwiki.php/Main/Media"
	googleUrl        = "https://www.google.com/"
	attackontitanUrl = "https://tvtropes.org/pmwiki/pmwiki.php/Manga/AttackOnTitan"
)

var (
	avengersSubpageFiles = []string{"resources/theavengers_tropesAtoD.html",
		"resources/theavengers_tropesEtoL.html",
		"resources/theavengers_tropesMtoP.html",
		"resources/theavengers_tropesQtoZ.html"}
)

// A scraper service for test purposes
var serviceScraperJson, serviceScraperCsv, invalidScraper *scraper.ServiceScraper
var newScraperJsonErr, newScraperCsvErr, invalidScraperErr, errPersistJson, errPersistCsv error
var csvRepositoryErr, jsonRepositoryErr error
var csvRepository, jsonRepository media.RepositoryMedia
var pageReaderJson, pageReaderCsv *os.File

var _ = BeforeSuite(func() {
	// Create two scrapers, one for the JSON dataset and the other for the CSV dataset
	csvRepository, csvRepositoryErr = csv_dataset.NewCSVRepository("dataset")
	jsonRepository, jsonRepositoryErr = json_dataset.NewJSONRepository("dataset")
	serviceScraperJson, newScraperJsonErr = scraper.NewServiceScraper(scraper.ConfigMediaRepository(jsonRepository))
	serviceScraperCsv, newScraperCsvErr = scraper.NewServiceScraper(scraper.ConfigMediaRepository(csvRepository))

	// Create invalid scraper
	invalidScraper, invalidScraperErr = scraper.NewServiceScraper(scraper.ConfigIndexRepository(nil))
})

var _ = Describe("Scraper", func() {
	AfterEach(func() {
		pageReaderCsv.Close()
		pageReaderJson.Close()
	})

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
		Context("URL belongs to a TvTropes Work page with tropes on a list", func() {
			var validTvTropesPageCsv, validTvTropesPageJson bool
			var errTvTropesCsv, errTvTropesJson error

			BeforeEach(func() {
				tvTropesUrl, _ := url.Parse(oldboyUrl)
				pageReaderJson, _ = os.Open("resources/oldboy2003.html")
				pageReaderCsv, _ = os.Open("resources/oldboy2003.html")

				validTvTropesPageJson, errTvTropesJson = serviceScraperJson.CheckValidWorkPage(pageReaderJson, tvTropesUrl)
				validTvTropesPageCsv, errTvTropesCsv = serviceScraperCsv.CheckValidWorkPage(pageReaderCsv, tvTropesUrl)
			})

			It("Should mark the page as valid", func() {
				Expect(validTvTropesPageJson).To(BeTrue())
				Expect(validTvTropesPageCsv).To(BeTrue())
			})

			It("Shouldn't return an error", func() {
				Expect(errTvTropesJson).To(BeNil())
				Expect(errTvTropesCsv).To(BeNil())
			})
		})

		Context("URL belongs to a TvTropes Work page with tropes on subpages", func() {
			var validTvTropesPage2Csv, validTvTropesPage2Json bool
			var errTvTropes2Csv, errTvTropes2Json error

			BeforeEach(func() {
				tvTropesUrl2, _ := url.Parse(avengersUrl)
				pageReaderJson, _ = os.Open("resources/theavengers2012.html")
				pageReaderCsv, _ = os.Open("resources/theavengers2012.html")

				validTvTropesPage2Json, errTvTropes2Json = serviceScraperJson.CheckValidWorkPage(pageReaderJson, tvTropesUrl2)
				validTvTropesPage2Csv, errTvTropes2Csv = serviceScraperCsv.CheckValidWorkPage(pageReaderCsv, tvTropesUrl2)
			})

			It("Should mark the page as valid", func() {
				Expect(validTvTropesPage2Json).To(BeTrue())
				Expect(validTvTropesPage2Csv).To(BeTrue())
			})

			It("Shouldn't return an error", func() {
				Expect(errTvTropes2Json).To(BeNil())
				Expect(errTvTropes2Csv).To(BeNil())
			})
		})

		Context("URL belongs to a TvTropes Work page with tropes on folders", func() {
			var validTvTropesPage3Csv, validTvTropesPage3Json bool
			var errTvTropes3Csv, errTvTropes3Json error

			BeforeEach(func() {
				tvTropesUrl3, _ := url.Parse(anewhopeUrl)
				pageReaderCsv, _ = os.Open("resources/anewhope.html")
				pageReaderJson, _ = os.Open("resources/anewhope.html")

				validTvTropesPage3Json, errTvTropes3Json = serviceScraperJson.CheckValidWorkPage(pageReaderJson, tvTropesUrl3)
				validTvTropesPage3Csv, errTvTropes3Json = serviceScraperCsv.CheckValidWorkPage(pageReaderCsv, tvTropesUrl3)
			})

			It("Should mark the page as valid", func() {
				Expect(validTvTropesPage3Json).To(BeTrue())
				Expect(validTvTropesPage3Csv).To(BeTrue())
			})

			It("Shouldn't return an error", func() {
				Expect(errTvTropes3Json).To(BeNil())
				Expect(errTvTropes3Csv).To(BeNil())
			})
		})

		Context("URL belongs to TvTropes but isn't from a Work page", func() {
			var validNotWorkPageCsv, validNotWorkPageJson bool
			var errNotWorkPageCsv, errNotWorkPageJson error

			BeforeEach(func() {
				notWorkUrl, _ := url.Parse(mediaUrl)

				validNotWorkPageJson, errNotWorkPageJson = serviceScraperJson.CheckIsWorkPage(notWorkUrl)
				validNotWorkPageCsv, errNotWorkPageCsv = serviceScraperCsv.CheckIsWorkPage(notWorkUrl)
			})

			It("Should mark the page as invalid", func() {
				Expect(validNotWorkPageJson).To(BeFalse())
				Expect(validNotWorkPageCsv).To(BeFalse())
			})

			It("Should return an appropriate error", func() {
				Expect(errors.Is(errNotWorkPageJson, scraper.ErrNotWorkPage)).To(BeTrue())
				Expect(errors.Is(errNotWorkPageCsv, scraper.ErrNotWorkPage)).To(BeTrue())
			})
		})

		Context("URL isn't from TvTropes", func() {
			var validDifferentPageCsv, validDifferentPageJson bool
			var errDifferentCsv, errDifferentJson error

			BeforeEach(func() {
				differentUrl, _ := url.Parse(googleUrl)
				validDifferentPageJson, errDifferentJson = serviceScraperJson.CheckIsWorkPage(differentUrl)
				validDifferentPageCsv, errDifferentCsv = serviceScraperCsv.CheckIsWorkPage(differentUrl)
			})

			It("Should mark the page as invalid", func() {
				Expect(validDifferentPageJson).To(BeFalse())
				Expect(validDifferentPageCsv).To(BeFalse())
			})

			It("Should return an appropriate error", func() {
				Expect(errors.Is(errDifferentJson, scraper.ErrNotTvTropes)).To(BeTrue())
				Expect(errors.Is(errDifferentCsv, scraper.ErrNotTvTropes)).To(BeTrue())
			})
		})
	})

	Describe("Scrape Film Page", func() {
		Context("Valid Film Page with tropes on a simple list", func() {
			var validfilm1Csv, validfilm1Json media.Media
			var errorfilm1Csv, errorfilm1Json error

			BeforeEach(func() {
				tvTropesUrl, _ := url.Parse(oldboyUrl)
				pageReaderCsv, _ = os.Open("resources/oldboy2003.html")
				pageReaderJson, _ = os.Open("resources/oldboy2003.html")

				validfilm1Json, errorfilm1Json = serviceScraperJson.ScrapeWorkPage(pageReaderJson, []io.Reader{}, tvTropesUrl)
				validfilm1Csv, errorfilm1Csv = serviceScraperCsv.ScrapeWorkPage(pageReaderCsv, []io.Reader{}, tvTropesUrl)
				errPersistJson = serviceScraperJson.Persist()
				errPersistCsv = serviceScraperCsv.Persist()
			})

			It("Shouldn't return an error", func() {
				Expect(errorfilm1Json).To(BeNil())
				Expect(errorfilm1Csv).To(BeNil())
				Expect(errPersistJson).To(BeNil())
				Expect(errPersistCsv).To(BeNil())
			})

			It("Should have correct fields", func() {
				testValidScrapedMedia(validfilm1Csv, "Oldboy", "2003", media.Film)
				testValidScrapedMedia(validfilm1Json, "Oldboy", "2003", media.Film)
			})

			It("Shouldn't have repeated tropes", func() {
				uniqueCsv := areTropesUnique(validfilm1Csv.GetWork().Tropes)
				uniqueJson := areTropesUnique(validfilm1Json.GetWork().Tropes)

				Expect(uniqueCsv).To(BeTrue())
				Expect(uniqueJson).To(BeTrue())
			})

			It("Should have added a correct record on the JSON repository", func() {
				testJsonRepository("Oldboy", "2003", "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003", "Film")
			})

			It("Should have added a correct record on the CSV repository", func() {
				testCsvRepository("Oldboy", "2003", "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003", "Film")
			})
		})

		Context("Valid Film Page with tropes distributed on main sub pages", func() {
			var validfilm2Csv, validfilm2Json media.Media
			var errorfilm2Csv, errorfilm2Json error

			BeforeEach(func() {
				tvTropesUrl2, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/TheAvengers2012")
				pageReaderCsv, _ = os.Open("resources/theavengers2012.html")
				pageReaderJson, _ = os.Open("resources/theavengers2012.html")

				var subpageReadersCsv []io.Reader
				var subpageReadersJson []io.Reader
				for _, subpageFile := range avengersSubpageFiles {
					subpageReaderCsv, _ := os.Open(subpageFile)
					subpageReaderJson, _ := os.Open(subpageFile)

					subpageReadersCsv = append(subpageReadersCsv, subpageReaderCsv)
					subpageReadersJson = append(subpageReadersJson, subpageReaderJson)
				}

				validfilm2Csv, errorfilm2Json = serviceScraperJson.ScrapeWorkPage(pageReaderJson, subpageReadersCsv, tvTropesUrl2)
				validfilm2Json, errorfilm2Csv = serviceScraperCsv.ScrapeWorkPage(pageReaderCsv, subpageReadersJson, tvTropesUrl2)
				errPersistJson = serviceScraperJson.Persist()
				errPersistCsv = serviceScraperCsv.Persist()
			})

			It("Shouldn't return an error", func() {
				Expect(errorfilm2Csv).To(BeNil())
				Expect(errorfilm2Json).To(BeNil())
			})

			It("Should have correct fields", func() {
				testValidScrapedMedia(validfilm2Csv, "The Avengers", "2012", media.Film)
				testValidScrapedMedia(validfilm2Json, "The Avengers", "2012", media.Film)
			})

			It("Shouldn't have repeated tropes", func() {
				uniqueCsv := areTropesUnique(validfilm2Csv.GetWork().Tropes)
				uniqueJson := areTropesUnique(validfilm2Json.GetWork().Tropes)

				Expect(uniqueCsv).To(BeTrue())
				Expect(uniqueJson).To(BeTrue())
			})

			It("Should have added a correct record on the JSON repository", func() {
				testJsonRepository("Oldboy", "2003", "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003", "Film")
			})

			It("Should have added a correct record on the CSV repository", func() {
				testCsvRepository("Oldboy", "2003", "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003", "Film")
			})
		})

		Context("Valid Film Page with tropes on folders", func() {
			var validfilm3Csv, validfilm3Json media.Media
			var errorfilm3Csv, errorfilm3Json error

			BeforeEach(func() {
				tvTropesUrl3, _ := url.Parse(anewhopeUrl)
				pageReaderCsv, _ = os.Open("resources/anewhope.html")
				pageReaderJson, _ = os.Open("resources/anewhope.html")

				validfilm3Json, errorfilm3Json = serviceScraperJson.ScrapeWorkPage(pageReaderJson, []io.Reader{}, tvTropesUrl3)
				validfilm3Csv, errorfilm3Csv = serviceScraperCsv.ScrapeWorkPage(pageReaderCsv, []io.Reader{}, tvTropesUrl3)
				errPersistJson = serviceScraperJson.Persist()
				errPersistCsv = serviceScraperCsv.Persist()
			})

			It("Shouldn't return an error", func() {
				Expect(errorfilm3Csv).To(BeNil())
				Expect(errorfilm3Json).To(BeNil())
				Expect(errPersistCsv).To(BeNil())
				Expect(errPersistJson).To(BeNil())
			})

			It("Should have correct fields", func() {
				testValidScrapedMedia(validfilm3Csv, "A New Hope", "", media.Film)
				testValidScrapedMedia(validfilm3Json, "A New Hope", "", media.Film)
			})

			It("Shouldn't have repeated tropes", func() {
				uniqueCsv := areTropesUnique(validfilm3Csv.GetWork().Tropes)
				uniqueJson := areTropesUnique(validfilm3Json.GetWork().Tropes)

				Expect(uniqueCsv).To(BeTrue())
				Expect(uniqueJson).To(BeTrue())
			})

			It("Should have added a correct record on the JSON repository", func() {
				testJsonRepository("Oldboy", "2003", "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003", "Film")
			})

			It("Should have added a correct record on the CSV repository", func() {
				testCsvRepository("Oldboy", "2003", "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003", "Film")
			})
		})

		Context("Invalid Film because the media type isn't supported", func() {
			var filminvalidtypeJson, filminvalidtypeCsv media.Media
			var errorfilminvalidtypeJson, errorfilminvalidtypeCsv error

			BeforeEach(func() {
				tvTropesUrlUnknown, _ := url.Parse(attackontitanUrl)
				pageReaderCsv, _ = os.Open("resources/attackontitan.html")
				pageReaderJson, _ = os.Open("resources/attackontitan.html")

				filminvalidtypeJson, errorfilminvalidtypeJson = serviceScraperJson.ScrapeWorkPage(pageReaderCsv, []io.Reader{}, tvTropesUrlUnknown)
				filminvalidtypeCsv, errorfilminvalidtypeCsv = serviceScraperCsv.ScrapeWorkPage(pageReaderJson, []io.Reader{}, tvTropesUrlUnknown)
				errPersistCsv = serviceScraperCsv.Persist()
				errPersistJson = serviceScraperJson.Persist()
			})

			It("Should return an empty media object", func() {
				Expect(filminvalidtypeJson.GetWork()).To(BeNil())
				Expect(filminvalidtypeJson.GetPage().GetUrl()).To(BeNil())
				Expect(filminvalidtypeJson.GetPage().GetPageType()).To(BeZero())
				Expect(filminvalidtypeJson.GetMediaType()).To(Equal(media.UnknownMediaType))

				Expect(filminvalidtypeCsv.GetWork()).To(BeNil())
				Expect(filminvalidtypeCsv.GetPage().GetUrl()).To(BeNil())
				Expect(filminvalidtypeCsv.GetPage().GetPageType()).To(BeZero())
				Expect(filminvalidtypeCsv.GetMediaType()).To(Equal(media.UnknownMediaType))
			})

			It("Should return an appropriate error", func() {
				Expect(errors.Unwrap(errorfilminvalidtypeJson)).To(Equal(media.ErrUnknownMediaType))
				Expect(errors.Unwrap(errorfilminvalidtypeCsv)).To(Equal(media.ErrUnknownMediaType))
				Expect(errors.Unwrap(errPersistCsv)).To(Equal(csv_dataset.ErrPersist))
				Expect(errors.Unwrap(errPersistJson)).To(Equal(json_dataset.ErrPersist))
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

func testCsvRepository(title, year, URL, mediaType string) {
	f, errOpen := os.Open("dataset.csv")
	reader := csv.NewReader(f)
	records, errReadCSV := reader.ReadAll()

	Expect(errOpen).To(BeNil())
	Expect(errReadCSV).To(BeNil())

	Expect(len(records[0])).To(Equal(7))
	Expect(len(records[1])).To(Equal(7))
	Expect(records[1][0]).To(Equal(title))
	Expect(records[1][1]).To(Equal(year))
	Expect(records[1][3]).To(Equal(URL))
	Expect(records[1][4]).To(Equal(mediaType))
	Expect(len(strings.Split(records[1][5], ";")) > 0).To(BeTrue())
}

func testJsonRepository(title, year, URL, mediaType string) {
	var dataset json_dataset.JSONDataset
	fileContents, _ := os.ReadFile("dataset.json")
	err := json.Unmarshal(fileContents, &dataset)

	Expect(err).To(BeNil())
	Expect(dataset.Tropestogo[0].Title).To(Equal(title))
	Expect(dataset.Tropestogo[0].Year).To(Equal(year))
	Expect(dataset.Tropestogo[0].URL).To(Equal(URL))
	Expect(dataset.Tropestogo[0].MediaType).To(Equal(mediaType))
	Expect(len(dataset.Tropestogo[0].Tropes) > 0).To(BeTrue())
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
