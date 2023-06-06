package csv_dataset_test

import (
	"encoding/csv"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"github.com/jlgallego99/TropesToGo/media/csv_dataset"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/url"
	"os"
	"strings"
	"time"
)

var repository *csv_dataset.CSVRepository
var errorRepository, errRemoveAll, errAddMedia error
var mediaEntry media.Media
var reader *csv.Reader
var datasetFile *os.File

var _ = BeforeSuite(func() {
	repository, errorRepository = csv_dataset.NewCSVRepository("dataset", ',')

	tropes := make(map[tropestogo.Trope]struct{})
	trope1, _ := tropestogo.NewTrope("AdaptationalLocationChange", tropestogo.TropeIndex(1))
	trope2, _ := tropestogo.NewTrope("AdaptationNameChange", tropestogo.TropeIndex(1))
	trope3, _ := tropestogo.NewTrope("AgeGapRomance", tropestogo.TropeIndex(2))
	tropes[trope1] = struct{}{}
	tropes[trope2] = struct{}{}
	tropes[trope3] = struct{}{}
	tvTropesUrl, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003")
	tvTropesPage := &tropestogo.Page{
		URL:         tvTropesUrl,
		LastUpdated: time.Now(),
	}
	mediaEntry, _ = media.NewMedia("Oldboy", "2003", time.Now(), tropes, tvTropesPage, media.Film)
})

var _ = Describe("CsvDataset", func() {
	BeforeEach(func() {
		errAddMedia = repository.AddMedia(mediaEntry)
		datasetFile, _ = os.Open("dataset.csv")
		reader, _ = repository.GetReader()
	})

	AfterEach(func() {
		// Reset file
		repository.RemoveAll()
	})

	Context("Create CSV Repository", func() {
		It("Should have created a CSV file", func() {
			Expect("dataset.csv").To(BeAnExistingFile())
		})

		It("Should have a delimiter", func() {
			Expect(repository.GetDelimiter()).To(Equal(','))
		})

		It("Shouldn't return an error", func() {
			Expect(errorRepository).To(BeNil())
		})

		It("Should only have the headers", func() {
			records, err := reader.ReadAll()

			Expect(err).To(BeNil())
			Expect(len(records)).To(Equal(2))
			Expect(records[0]).To(Equal([]string{"title", "year", "lastupdated", "url", "mediatype", "tropes", "tropes_index"}))
		})
	})

	Context("Add a Media to the CSV file", func() {
		It("Should have added the correct record to the CSV", func() {
			records, err := reader.ReadAll()

			Expect(len(records[0])).To(Equal(7))
			Expect(len(records[1])).To(Equal(7))
			Expect(records[1][0]).To(Equal("Oldboy"))
			Expect(records[1][1]).To(Equal("2003"))
			Expect(records[1][3]).To(Equal("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003"))
			Expect(records[1][4]).To(Equal("Film"))
			Expect(len(strings.Split(records[1][5], ";"))).To(Equal(3))
			Expect(strings.Contains(records[1][5], "AdaptationalLocationChange")).To(BeTrue())
			Expect(strings.Contains(records[1][5], "AdaptationNameChange")).To(BeTrue())
			Expect(strings.Contains(records[1][5], "AgeGapRomance")).To(BeTrue())
			Expect(strings.Contains(records[1][6], "GenreTrope")).To(BeTrue())
			Expect(strings.Contains(records[1][6], "MediaTrope")).To(BeTrue())
			Expect(err).To(BeNil())
		})

		It("Shouldn't return an error", func() {
			Expect(errAddMedia).To(BeNil())
		})
	})

	Context("Remove CSV file contents", func() {
		BeforeEach(func() {
			errRemoveAll = repository.RemoveAll()
		})

		It("Should still exist a CSV file", func() {
			Expect("dataset.csv").To(BeAnExistingFile())
		})

		It("Should only have the headers", func() {
			records, err := reader.ReadAll()

			Expect(err).To(BeNil())
			Expect(len(records)).To(Equal(1))
			Expect(records[0]).To(Equal([]string{"title", "year", "lastupdated", "url", "mediatype", "tropes", "tropes_index"}))
		})

		It("Should have no errors", func() {
			Expect(errRemoveAll).To(BeNil())
		})
	})

	Context("Remove contents of CSV file that doesn't exist", func() {
		BeforeEach(func() {
			os.Remove("dataset.csv")
			errRemoveAll = repository.RemoveAll()
		})

		AfterEach(func() {
			repository, errorRepository = csv_dataset.NewCSVRepository("dataset", ',')
		})

		It("Shouldn't exist a CSV file", func() {
			Expect("dataset.csv").To(Not(BeAnExistingFile()))
		})

		It("Should return an error", func() {
			Expect(errRemoveAll).To(Equal(csv_dataset.ErrFileNotExists))
		})
	})

	Context("Update the Year, URL and tropes of a Film in the CSV file", func() {
		var errUpdate error

		BeforeEach(func() {
			// Create the new Media to be updated
			trope1, _ := tropestogo.NewTrope("AdaptationalComicRelief", tropestogo.TropeIndex(1))
			trope2, _ := tropestogo.NewTrope("AdaptationalHeroism", tropestogo.TropeIndex(1))
			tropes := make(map[tropestogo.Trope]struct{})
			tropes[trope1] = struct{}{}
			tropes[trope2] = struct{}{}

			updatedUrl, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2013")
			tvTropesPage := &tropestogo.Page{
				URL:         updatedUrl,
				LastUpdated: time.Now(),
			}

			updatedMediaEntry, _ := media.NewMedia("Oldboy", "2013", time.Now(), tropes, tvTropesPage, media.Film)

			errUpdate = repository.UpdateMedia("Oldboy", "2003", updatedMediaEntry)
		})

		It("Should have the new record updated", func() {
			records, err := reader.ReadAll()

			Expect(err).To(BeNil())
			Expect(len(records)).To(Equal(2))
			Expect(records[0]).To(Equal([]string{"title", "year", "lastupdated", "url", "mediatype", "tropes", "tropes_index"}))
			Expect(records[1][0]).To(Equal("Oldboy"))
			Expect(records[1][1]).To(Equal("2013"))
			Expect(records[1][3]).To(Equal("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2013"))
			Expect(records[1][4]).To(Equal("Film"))
			Expect(len(strings.Split(records[1][5], ";"))).To(Equal(2))
			Expect(strings.Contains(records[1][5], "AdaptationalComicRelief")).To(BeTrue())
			Expect(strings.Contains(records[1][5], "AdaptationalHeroism")).To(BeTrue())
			Expect(strings.Contains(records[1][6], "GenreTrope")).To(BeTrue())
		})

		It("Shouldn't return an error", func() {
			Expect(errUpdate).To(BeNil())
		})
	})
})

var _ = AfterSuite(func() {
	datasetFile.Close()
	os.Remove("dataset.csv")
})
