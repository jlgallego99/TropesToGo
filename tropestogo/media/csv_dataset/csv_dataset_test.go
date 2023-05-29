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
	"time"
)

var repository *csv_dataset.CSVRepository
var errorRepository error
var mediaEntry media.Media
var reader *csv.Reader
var datasetFile *os.File

var _ = BeforeSuite(func() {
	repository, errorRepository = csv_dataset.NewCSVRepository(',')

	datasetFile, _ = os.Open("dataset.csv")
	reader = csv.NewReader(datasetFile)

	tropes := make(map[tropestogo.Trope]struct{})
	trope1, _ := tropestogo.NewTrope("AccentUponTheWrongSyllable", tropestogo.TropeIndex(0))
	trope2, _ := tropestogo.NewTrope("ChekhovsGun", tropestogo.TropeIndex(0))
	tropes[trope1] = struct{}{}
	tropes[trope2] = struct{}{}
	tvTropesUrl, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003")
	tvTropesPage := &tropestogo.Page{
		URL:         tvTropesUrl,
		LastUpdated: time.Now(),
	}
	mediaEntry, _ = media.NewMedia("TheAvengers", "2012", time.Now(), tropes, tvTropesPage, media.Film)
})

var _ = Describe("CsvDataset", func() {
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
	})

	Context("Add a Media to the CSV file", func() {
		var errAddMedia error

		BeforeEach(func() {
			errAddMedia = repository.AddMedia(mediaEntry)
		})

		It("Should have added the correct record to the CSV", func() {
			record, err := reader.Read()

			Expect(len(record)).To(Equal(2))
			Expect(record[0]).To(Equal("TheAvengers"))
			Expect(record[1]).To(Equal("AccentUponTheWrongSyllable;ChekhovsGun"))
			Expect(err).To(BeNil())
		})

		It("Shouldn't return an error", func() {
			Expect(errAddMedia).To(BeNil())
		})

		AfterEach(func() {
			// Delete file contents
			os.Truncate("dataset.csv", 0)
		})
	})
})

var _ = AfterSuite(func() {
	datasetFile.Close()
})
