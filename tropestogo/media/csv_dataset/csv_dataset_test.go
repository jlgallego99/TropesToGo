package csv_dataset_test

import (
	"encoding/csv"
	"errors"
	"fmt"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"github.com/jlgallego99/TropesToGo/media/csv_dataset"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	oldboyUrl = "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003"
)

var headers = []string{"title", "year", "lastupdated", "url", "mediatype", "tropes", "subtropes", "subtropes_namespaces"}

var repository *csv_dataset.CSVRepository
var errorRepository, errRemoveAll, errAddMedia, errPersist error
var mediaEntry media.Media
var reader *csv.Reader
var datasetFile *os.File
var tropes map[tropestogo.Trope]struct{}
var numTropes int

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var _ = BeforeSuite(func() {
	repository, errorRepository = csv_dataset.NewCSVRepository("dataset")

	const max = 10
	const min = 2
	numTropes = seededRand.Intn(max-min) + min
	tropes = createTropeSet(numTropes)
	subTropes := createSubTropeSet(numTropes)
	for subTrope := range subTropes {
		tropes[subTrope] = struct{}{}
	}

	tvTropesPage, _ := tropestogo.NewPage(oldboyUrl)
	mediaEntry, _ = media.NewMedia("Oldboy", "2003", time.Now(), tropes, tvTropesPage, media.Film)
})

var _ = Describe("CsvDataset", func() {
	BeforeEach(func() {
		reader, _ = repository.GetReader()
	})

	AfterEach(func() {
		// Reset file
		repository.RemoveAll()
	})

	Context("Create CSV Repository", func() {
		BeforeEach(func() {
			errPersist = repository.Persist()
		})

		It("Should have created a CSV file", func() {
			Expect("dataset.csv").To(BeAnExistingFile())
		})

		It("Shouldn't return an error", func() {
			Expect(errorRepository).To(BeNil())
		})

		It("Should only have the headers", func() {
			records, err := reader.ReadAll()

			Expect(err).To(BeNil())
			Expect(len(records)).To(Equal(1))
			Expect(records[0]).To(Equal(headers))
		})

		It("Shouldn't be able to persist anything", func() {
			Expect(errors.Is(errPersist, csv_dataset.ErrPersist)).To(BeTrue())
		})
	})

	Context("Add a Media to the CSV file", func() {
		BeforeEach(func() {
			errAddMedia = repository.AddMedia(mediaEntry)
			errPersist = repository.Persist()
		})

		It("Should have added the correct record to the CSV", func() {
			records, err := reader.ReadAll()

			Expect(records[0]).To(Equal(headers))
			Expect(len(records[1]) > 0).To(BeTrue())
			Expect(records[1][0]).To(Not(BeEmpty()))
			Expect(records[1][1]).To(Not(BeEmpty()))
			Expect(records[1][3]).To(Not(BeEmpty()))
			Expect(records[1][4]).To(Not(BeEmpty()))
			Expect(len(strings.Split(records[1][5], ";"))).To(Equal(numTropes))
			Expect(len(strings.Split(records[1][6], ";")) > 0).To(BeTrue())
			Expect(len(strings.Split(records[1][7], ";")) > 0).To(BeTrue())
			Expect(err).To(BeNil())
		})

		It("Shouldn't return an error", func() {
			Expect(errAddMedia).To(BeNil())
			Expect(errPersist).To(BeNil())
		})
	})

	Context("Add duplicated Media to the CSV file", func() {
		BeforeEach(func() {
			errAddMedia = repository.AddMedia(mediaEntry)
			errAddMedia = repository.AddMedia(mediaEntry)
			errPersist = repository.Persist()
		})

		It("Should only be one record on the CSV file", func() {
			records, err := reader.ReadAll()

			Expect(err).To(BeNil())
			Expect(len(records)).To(Equal(2))
		})

		It("Should return an error", func() {
			Expect(errors.Is(errAddMedia, csv_dataset.ErrDuplicatedMedia)).To(BeTrue())
			Expect(errPersist).To(BeNil())
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
			Expect(records[0]).To(Equal(headers))
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
			repository, errorRepository = csv_dataset.NewCSVRepository("dataset")
		})

		It("A CSV file shouldn't exist", func() {
			Expect("dataset.csv").To(Not(BeAnExistingFile()))
		})

		It("Should return an error", func() {
			Expect(errors.Is(errRemoveAll, csv_dataset.ErrFileNotExists)).To(BeTrue())
		})
	})

	Context("Update the Year, URL and tropes of a Film in the CSV file", func() {
		var errUpdate error

		BeforeEach(func() {
			errAddMedia = repository.AddMedia(mediaEntry)
			errPersist = repository.Persist()

			const max = 10
			const min = 2
			numTropes = seededRand.Intn(max-min) + min

			// Create the new Media to be updated
			newTropes := createTropeSet(numTropes)
			newSubTropes := createSubTropeSet(numTropes)
			for subTrope := range newSubTropes {
				newTropes[subTrope] = struct{}{}
			}
			tvTropesPage, _ := tropestogo.NewPage(oldboyUrl)
			updatedMediaEntry, _ := media.NewMedia("Oldboy", "2013", time.Now(), newTropes, tvTropesPage, media.Film)

			errUpdate = repository.UpdateMedia("Oldboy", "2003", updatedMediaEntry)
		})

		It("Should have the new record updated", func() {
			records, err := reader.ReadAll()

			Expect(err).To(BeNil())
			Expect(errPersist).To(BeNil())
			Expect(len(records)).To(Equal(2))
			Expect(records[0]).To(Equal(headers))
			Expect(records[1][0]).To(Not(BeEmpty()))
			Expect(records[1][1]).To(Not(BeEmpty()))
			Expect(records[1][3]).To(Not(BeEmpty()))
			Expect(records[1][4]).To(Not(BeEmpty()))
			Expect(records[1][1]).To(Not(Equal(mediaEntry.GetWork().Year)))
			Expect(records[1][2]).To(Not(Equal(mediaEntry.GetWork().LastUpdated)))
			Expect(records[1][3]).To(Not(Equal(mediaEntry.GetPage().GetUrl())))
			Expect(records[1][5]).To(Not(Equal(createTropesString(mediaEntry.GetWork().Tropes))))
			Expect(len(strings.Split(records[1][5], ";"))).To(Equal(numTropes))
			Expect(len(strings.Split(records[1][6], ";")) > 0).To(BeTrue())
			Expect(len(strings.Split(records[1][7], ";")) > 0).To(BeTrue())
		})

		It("Shouldn't return an error", func() {
			Expect(errUpdate).To(BeNil())
			Expect(errPersist).To(BeNil())
		})
	})

	Context("Persist an already persisted before record", func() {
		BeforeEach(func() {
			// Persist first
			errAddMedia = repository.AddMedia(mediaEntry)
			errPersist = repository.Persist()

			// Try to persist again the same Media
			errAddMedia = repository.AddMedia(mediaEntry)
			errPersist = repository.Persist()
		})

		It("Should only be one Media record on the CSV file", func() {
			records, err := reader.ReadAll()

			Expect(err).To(BeNil())
			Expect(len(records)).To(Equal(2))
		})
	})
})

var _ = AfterSuite(func() {
	datasetFile.Close()
	os.Remove("dataset.csv")
})

// createTropeSet generates a generic set of N correct tropes
func createTropeSet(numTropes int) map[tropestogo.Trope]struct{} {
	tropeset := make(map[tropestogo.Trope]struct{})
	for i := 0; i < numTropes; i++ {
		trope, _ := tropestogo.NewTrope("Trope"+fmt.Sprint(i), 1, "")
		tropeset[trope] = struct{}{}
	}

	return tropeset
}

// createSubTropeSet generates a generic set of N correct SubTropes of different SubWikis at random
func createSubTropeSet(numTropes int) map[tropestogo.Trope]struct{} {
	subWikis := []string{"SubWiki1", "SubWiki2"}

	tropeset := make(map[tropestogo.Trope]struct{})
	for i := 0; i < numTropes; i++ {
		trope, _ := tropestogo.NewTrope("Trope"+fmt.Sprint(i), 1, subWikis[0])
		tropeset[trope] = struct{}{}

		trope, _ = tropestogo.NewTrope("Trope"+fmt.Sprint(i), 1, subWikis[1])
		tropeset[trope] = struct{}{}
	}

	return tropeset
}

// createTropesString generates a string of all tropes titles joined by a semicolon
func createTropesString(tropes map[tropestogo.Trope]struct{}) string {
	var tropeTitles []string
	for trope := range tropes {
		tropeTitles = append(tropeTitles, trope.GetTitle())
	}

	return strings.Join(tropeTitles, ";")
}
