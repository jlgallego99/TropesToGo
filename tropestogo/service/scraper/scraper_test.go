package scraper_test

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"github.com/jlgallego99/TropesToGo/media/csv_dataset"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	"github.com/jlgallego99/TropesToGo/service/scraper"
	"net/url"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	oldboyUrl             = "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003"
	oldboyResource        = "resources/oldboy2003.html"
	anewhopeUrl           = "https://tvtropes.org/pmwiki/pmwiki.php/Film/ANewHope"
	anewhopeResource      = "resources/anewhope.html"
	avengersUrl           = "https://tvtropes.org/pmwiki/pmwiki.php/Film/TheAvengers2012"
	avengersResource      = "resources/theavengers2012.html"
	mediaUrl              = "https://tvtropes.org/pmwiki/pmwiki.php/Main/Media"
	googleUrl             = "https://www.google.com/"
	attackontitanUrl      = "https://tvtropes.org/pmwiki/pmwiki.php/Manga/AttackOnTitan"
	attackontitanResource = "resources/attackontitan.html"
	emptyResource         = "resources/empty.html"
)

var (
	avengersSubpageFiles = []string{"resources/theavengers_tropesAtoD.html",
		"resources/theavengers_tropesEtoL.html",
		"resources/theavengers_tropesMtoP.html",
		"resources/theavengers_tropesQtoZ.html"}
	oldboySubpageFiles = []string{"resources/oldboy_awesome.html", "resources/oldboy_fridge.html",
		"resources/oldboy_laconic.html", "resources/oldboy_trivia.html", "resources/oldboy_ymmv.html",
		"resources/oldboy_videoexamples.html"}
	headers = []string{"title", "year", "lastupdated", "url", "mediatype", "tropes", "subtropes", "subtropes_namespaces"}
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
				pageReaderJson, _ = os.Open(oldboyResource)
				pageReaderCsv, _ = os.Open(oldboyResource)

				docJson, _ := goquery.NewDocumentFromReader(pageReaderJson)
				docCsv, _ := goquery.NewDocumentFromReader(pageReaderCsv)

				validTvTropesPageJson, errTvTropesJson = serviceScraperJson.CheckValidWorkPage(docJson, tvTropesUrl)
				validTvTropesPageCsv, errTvTropesCsv = serviceScraperCsv.CheckValidWorkPage(docCsv, tvTropesUrl)
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
				pageReaderJson, _ = os.Open(avengersResource)
				pageReaderCsv, _ = os.Open(avengersResource)

				docJson, _ := goquery.NewDocumentFromReader(pageReaderJson)
				docCsv, _ := goquery.NewDocumentFromReader(pageReaderCsv)

				validTvTropesPage2Json, errTvTropes2Json = serviceScraperJson.CheckValidWorkPage(docJson, tvTropesUrl2)
				validTvTropesPage2Csv, errTvTropes2Csv = serviceScraperCsv.CheckValidWorkPage(docCsv, tvTropesUrl2)
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
				pageReaderCsv, _ = os.Open(anewhopeResource)
				pageReaderJson, _ = os.Open(anewhopeResource)

				docJson, _ := goquery.NewDocumentFromReader(pageReaderJson)
				docCsv, _ := goquery.NewDocumentFromReader(pageReaderCsv)

				validTvTropesPage3Json, errTvTropes3Json = serviceScraperJson.CheckValidWorkPage(docJson, tvTropesUrl3)
				validTvTropesPage3Csv, errTvTropes3Json = serviceScraperCsv.CheckValidWorkPage(docCsv, tvTropesUrl3)
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
				pageReaderJson, _ = os.Open(emptyResource)
				pageReaderCsv, _ = os.Open(emptyResource)

				docJson, _ := goquery.NewDocumentFromReader(pageReaderJson)
				docCsv, _ := goquery.NewDocumentFromReader(pageReaderCsv)

				validNotWorkPageJson, errNotWorkPageJson = serviceScraperJson.CheckValidWorkPage(docJson, notWorkUrl)
				validNotWorkPageCsv, errNotWorkPageCsv = serviceScraperCsv.CheckValidWorkPage(docCsv, notWorkUrl)
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
				pageReaderJson, _ = os.Open(emptyResource)
				pageReaderCsv, _ = os.Open(emptyResource)

				docJson, _ := goquery.NewDocumentFromReader(pageReaderJson)
				docCsv, _ := goquery.NewDocumentFromReader(pageReaderCsv)

				validDifferentPageJson, errDifferentJson = serviceScraperJson.CheckValidWorkPage(docJson, differentUrl)
				validDifferentPageCsv, errDifferentCsv = serviceScraperCsv.CheckValidWorkPage(docCsv, differentUrl)
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

	Describe("Scrape and persist an invalid Film because the media type isn't supported", func() {
		var filminvalidtypeJson, filminvalidtypeCsv media.Media
		var errorfilminvalidtypeJson, errorfilminvalidtypeCsv error

		BeforeEach(func() {
			tvTropesUrlUnknown, _ := url.Parse(attackontitanUrl)
			pageReaderCsv, _ = os.Open(attackontitanResource)
			pageReaderJson, _ = os.Open(attackontitanResource)

			docJson, _ := goquery.NewDocumentFromReader(pageReaderJson)
			docCsv, _ := goquery.NewDocumentFromReader(pageReaderCsv)
			subDocs := make([]*goquery.Document, 0)

			filminvalidtypeJson, errorfilminvalidtypeJson = serviceScraperJson.ScrapeFromDocuments(docJson, subDocs, tvTropesUrlUnknown)
			filminvalidtypeCsv, errorfilminvalidtypeCsv = serviceScraperCsv.ScrapeFromDocuments(docCsv, subDocs, tvTropesUrlUnknown)
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

		It("Should have empty datasets", func() {
			// Check empty JSON dataset
			var dataset json_dataset.JSONDataset
			fileContents, _ := os.ReadFile("dataset.json")
			err := json.Unmarshal(fileContents, &dataset)

			Expect(err).To(BeNil())
			Expect(dataset.Tropestogo).To(BeEmpty())

			// Check empty CSV dataset
			datasetFile, errReader := os.Open("dataset.csv")
			Expect(errReader).To(BeNil())
			reader := csv.NewReader(datasetFile)
			records, errReadAll := reader.ReadAll()
			Expect(errReadAll).To(BeNil())

			Expect(err).To(BeNil())
			Expect(len(records)).To(Equal(1))
			Expect(records[0]).To(Equal(headers))
		})
	})

	Describe("Scrape different Film Pages and persist on the dataset", func() {
		var validfilm1Csv, validfilm2Csv, validfilm3Csv, validfilm1Json, validfilm2Json, validfilm3Json media.Media
		var errorfilm1Csv, errorfilm2Csv, errorfilm3Csv, errorfilm1Json, errorfilm2Json, errorfilm3Json error

		BeforeEach(func() {
			tvTropesUrl, _ := url.Parse(oldboyUrl)
			tvTropesUrl2, _ := url.Parse(avengersUrl)
			tvTropesUrl3, _ := url.Parse(anewhopeUrl)

			var subpageDocsCsv []*goquery.Document
			var subpageDocsJson []*goquery.Document

			// Scrape Oldboy
			subpageDocsJson, subpageDocsCsv = loadSubpageFiles(oldboySubpageFiles)
			pageReaderCsv, _ = os.Open(oldboyResource)
			pageReaderJson, _ = os.Open(oldboyResource)
			docJson, _ := goquery.NewDocumentFromReader(pageReaderJson)
			docCsv, _ := goquery.NewDocumentFromReader(pageReaderCsv)
			validfilm1Json, errorfilm1Json = serviceScraperJson.ScrapeFromDocuments(docJson, subpageDocsJson, tvTropesUrl)
			validfilm1Csv, errorfilm1Csv = serviceScraperCsv.ScrapeFromDocuments(docCsv, subpageDocsCsv, tvTropesUrl)

			// Scrape The Avengers
			subpageDocsJson, subpageDocsCsv = loadSubpageFiles(avengersSubpageFiles)
			pageReaderCsv, _ = os.Open(avengersResource)
			pageReaderJson, _ = os.Open(avengersResource)
			docJson, _ = goquery.NewDocumentFromReader(pageReaderJson)
			docCsv, _ = goquery.NewDocumentFromReader(pageReaderCsv)
			validfilm2Csv, errorfilm2Json = serviceScraperJson.ScrapeFromDocuments(docCsv, subpageDocsCsv, tvTropesUrl2)
			validfilm2Json, errorfilm2Csv = serviceScraperCsv.ScrapeFromDocuments(docJson, subpageDocsJson, tvTropesUrl2)

			// Scrape A New Hope
			pageReaderCsv, _ = os.Open(anewhopeResource)
			pageReaderJson, _ = os.Open(anewhopeResource)
			docJson, _ = goquery.NewDocumentFromReader(pageReaderJson)
			docCsv, _ = goquery.NewDocumentFromReader(pageReaderCsv)
			subDocs := make([]*goquery.Document, 0)
			validfilm3Json, errorfilm3Json = serviceScraperJson.ScrapeFromDocuments(docJson, subDocs, tvTropesUrl3)
			validfilm3Csv, errorfilm3Csv = serviceScraperCsv.ScrapeFromDocuments(docCsv, subDocs, tvTropesUrl3)

			// Persist all data
			errPersistJson = serviceScraperJson.Persist()
			errPersistCsv = serviceScraperCsv.Persist()
		})

		It("Shouldn't return any errors on scraping the Film", func() {
			Expect(errorfilm1Json).To(BeNil())
			Expect(errorfilm1Csv).To(BeNil())

			Expect(errorfilm2Json).To(BeNil())
			Expect(errorfilm2Csv).To(BeNil())

			Expect(errorfilm3Json).To(BeNil())
			Expect(errorfilm3Csv).To(BeNil())
		})

		It("Shouldn't return any persisting errors", func() {
			Expect(errPersistJson).To(BeNil())
			Expect(errPersistCsv).To(BeNil())
		})

		It("Should have no empty or null fields", func() {
			testValidScrapedMedia(validfilm1Csv)
			testValidScrapedMedia(validfilm1Json)

			testValidScrapedMedia(validfilm2Csv)
			testValidScrapedMedia(validfilm2Json)

			testValidScrapedMedia(validfilm3Csv)
			testValidScrapedMedia(validfilm3Json)
		})

		It("Shouldn't have repeated tropes", func() {
			Expect(areTropesUnique(validfilm1Csv.GetWork().Tropes)).To(BeTrue())
			Expect(areSubTropesUnique(validfilm1Csv.GetWork().SubTropes)).To(BeTrue())
			Expect(areTropesUnique(validfilm1Json.GetWork().Tropes)).To(BeTrue())
			Expect(areSubTropesUnique(validfilm1Json.GetWork().SubTropes)).To(BeTrue())

			Expect(areTropesUnique(validfilm2Csv.GetWork().Tropes)).To(BeTrue())
			Expect(areTropesUnique(validfilm2Json.GetWork().Tropes)).To(BeTrue())

			Expect(areTropesUnique(validfilm3Csv.GetWork().Tropes)).To(BeTrue())
			Expect(areTropesUnique(validfilm3Json.GetWork().Tropes)).To(BeTrue())
		})

		It("Each record on the repository should have the correct columns/keys and they mustn't be empty", func() {
			// Check JSON dataset
			var dataset json_dataset.JSONDataset
			fileContents, _ := os.ReadFile("dataset.json")
			err := json.Unmarshal(fileContents, &dataset)

			Expect(err).To(BeNil())
			Expect(dataset.Tropestogo).To(Not(BeEmpty()))
			for _, record := range dataset.Tropestogo {
				jsonStringTropes := make([]string, 0)
				for _, datasetTrope := range record.Tropes {
					jsonStringTropes = append(jsonStringTropes, datasetTrope.Title)
				}

				Expect(record.Title).To(Not(BeEmpty()))
				Expect(record.URL).To(Not(BeEmpty()))
				Expect(record.LastUpdated).To(Not(BeEmpty()))
				Expect(record.MediaType).To(Equal(media.Film.String()))

				areRepositoryTropesUnique(jsonStringTropes)
			}

			// Check CSV dataset
			datasetFile, errReader := os.Open("dataset.csv")
			Expect(errReader).To(BeNil())
			reader := csv.NewReader(datasetFile)
			records, errReadAll := reader.ReadAll()
			Expect(errReadAll).To(BeNil())

			Expect(err).To(BeNil())
			Expect(records[0]).To(Equal(headers))
			Expect(len(records) > 1).To(BeTrue())
			for _, record := range records {
				Expect(record[0]).To(Not(BeEmpty()))
				Expect(record[2]).To(Not(BeEmpty()))
				Expect(record[3]).To(Not(BeEmpty()))
				Expect(record[4]).To(Not(BeEmpty()))
				areRepositoryTropesUnique(strings.Split(record[5], ";"))
				areRepositorySubTropesUnique(strings.Split(record[6], ";"), strings.Split(record[7], ";"))
			}
		})
	})
})

func testValidScrapedMedia(validMedia media.Media) {
	Expect(validMedia.GetWork().Title).To(Not(BeEmpty()))
	Expect(validMedia.GetMediaType()).To(Equal(media.Film))
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

func areSubTropesUnique(tropes map[tropestogo.Trope]struct{}) bool {
	visited := make(map[string]bool, 0)
	for trope := range tropes {
		if visited[trope.GetTitle()+trope.GetSubpage()] {
			return false
		} else {
			visited[trope.GetTitle()+trope.GetSubpage()] = true
		}
	}

	return true
}

func areRepositoryTropesUnique(tropes []string) bool {
	visited := make(map[string]bool, 0)
	for _, trope := range tropes {
		if visited[trope] == true {
			return false
		} else {
			visited[trope] = true
		}
	}

	return true
}

func areRepositorySubTropesUnique(subTropes []string, subPages []string) bool {
	visited := make(map[string]bool, 0)
	for i, trope := range subTropes {
		if visited[trope+subPages[i]] == true {
			return false
		} else {
			visited[trope+subPages[i]] = true
		}
	}

	return true
}

func loadSubpageFiles(fileNames []string) ([]*goquery.Document, []*goquery.Document) {
	var subpageDocsCsv []*goquery.Document
	var subpageDocsJson []*goquery.Document

	for _, subpageFile := range fileNames {
		subpageReaderCsv, _ := os.Open(subpageFile)
		subpageReaderJson, _ := os.Open(subpageFile)

		docCsv, _ := goquery.NewDocumentFromReader(subpageReaderCsv)
		docJson, _ := goquery.NewDocumentFromReader(subpageReaderJson)

		subpageDocsCsv = append(subpageDocsCsv, docCsv)
		subpageDocsJson = append(subpageDocsJson, docJson)
	}

	return subpageDocsJson, subpageDocsCsv
}
