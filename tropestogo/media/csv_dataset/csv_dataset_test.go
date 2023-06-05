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
	BeforeEach(func() {
		errAddMedia = repository.AddMedia(mediaEntry)
	})

	AfterEach(func() {
		// Reset file
		os.Create("dataset.csv")
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
	})

	Context("Add a Media to the CSV file", func() {
		It("Should have added the correct record to the CSV", func() {
			record, err := reader.Read()

			Expect(len(record)).To(Equal(2))
			Expect(strings.Trim(record[0], "\x00")).To(Equal("TheAvengers"))
			Expect(record[1]).To(Equal("AccentUponTheWrongSyllable;ChekhovsGun"))
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
			_, err := os.Stat("dataset.csv")
			Expect(err).To(BeNil())
		})

		It("Should be empty", func() {
			file, _ := os.Stat("dataset.csv")
			Expect(file.Size()).To(BeZero())
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

		It("Shouldn't exist a CSV file", func() {
			_, err := os.Stat("dataset.csv")
			Expect(err).To(Not(BeNil()))
		})

		It("Should return an error", func() {
			Expect(errRemoveAll).To(Equal(csv_dataset.ErrFileNotExists))
		})
	})
})

var _ = AfterSuite(func() {
	datasetFile.Close()
	os.Remove("dataset.csv")
})
