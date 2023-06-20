package json_dataset_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	oldboyUrl = "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003"
)

var repository *json_dataset.JSONRepository
var errorRepository, errRemoveAll, errAddMedia, errPersist error
var mediaEntry media.Media
var datasetFile *os.File
var tropes map[tropestogo.Trope]struct{}
var numTropes int

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var _ = BeforeSuite(func() {
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

var _ = Describe("JsonDataset", func() {
	BeforeEach(func() {
		repository, errorRepository = json_dataset.NewJSONRepository("dataset")
		datasetFile, _ = os.Open("dataset.json")
	})

	AfterEach(func() {
		// Reset file
		repository.RemoveAll()
	})

	Context("Create JSON repository", func() {
		BeforeEach(func() {
			errPersist = repository.Persist()
		})

		It("Should have created a JSON file", func() {
			Expect("dataset.json").To(BeAnExistingFile())
		})

		It("Shouldn't return an error", func() {
			Expect(errorRepository).To(BeNil())
		})

		It("Shouldn't be able to persist anything", func() {
			Expect(errors.Is(errPersist, json_dataset.ErrPersist)).To(BeTrue())
		})
	})

	Context("Add a Media to the JSON file", func() {
		BeforeEach(func() {
			errAddMedia = repository.AddMedia(mediaEntry)
			errPersist = repository.Persist()
		})

		It("Should have all the correct fields", func() {
			var dataset json_dataset.JSONDataset
			fileContents, _ := os.ReadFile("dataset.json")
			err := json.Unmarshal(fileContents, &dataset)

			Expect(err).To(BeNil())
			Expect(dataset.Tropestogo[0].Title).To(Not(BeEmpty()))
			Expect(dataset.Tropestogo[0].Year).To(Not(BeEmpty()))
			Expect(dataset.Tropestogo[0].URL).To(Not(BeEmpty()))
			Expect(dataset.Tropestogo[0].MediaType).To(Not(BeEmpty()))
			Expect(len(dataset.Tropestogo[0].Tropes) > 0).To(BeTrue())
			Expect(len(dataset.Tropestogo[0].SubTropes) > 0).To(BeTrue())
		})

		It("Shouldn't return an error", func() {
			Expect(errAddMedia).To(BeNil())
			Expect(errPersist).To(BeNil())
		})
	})

	Context("Add duplicated Media to the JSON file", func() {
		BeforeEach(func() {
			errAddMedia = repository.AddMedia(mediaEntry)
			errAddMedia = repository.AddMedia(mediaEntry)
			errPersist = repository.Persist()
		})

		It("Should only be one record on the JSON file", func() {
			var dataset json_dataset.JSONDataset
			fileContents, _ := os.ReadFile("dataset.json")
			err := json.Unmarshal(fileContents, &dataset)

			Expect(err).To(BeNil())
			Expect(len(dataset.Tropestogo)).To(Equal(1))
		})

		It("Should return an error", func() {
			Expect(errors.Is(errAddMedia, json_dataset.ErrDuplicatedMedia)).To(BeTrue())
			Expect(errPersist).To(BeNil())
		})
	})

	Context("Remove JSON file contents", func() {
		BeforeEach(func() {
			errRemoveAll = repository.RemoveAll()
		})

		It("Should still exist a JSON file", func() {
			Expect("dataset.json").To(BeAnExistingFile())
		})

		It("Should have an empty key of the main array", func() {
			var dataset json_dataset.JSONDataset
			fileContents, _ := os.ReadFile("dataset.json")
			err := json.Unmarshal(fileContents, &dataset)

			Expect(err).To(BeNil())
			Expect(dataset.Tropestogo).To(BeEmpty())
		})

		It("Should have no errors", func() {
			Expect(errRemoveAll).To(BeNil())
		})
	})

	Context("Remove contents of JSON file that doesn't exist", func() {
		BeforeEach(func() {
			os.Remove("dataset.json")
			errRemoveAll = repository.RemoveAll()
		})

		AfterEach(func() {
			repository, errorRepository = json_dataset.NewJSONRepository("dataset")
		})

		It("A JSON file shouldn't exist", func() {
			Expect("dataset.json").To(Not(BeAnExistingFile()))
		})

		It("Should return an error", func() {
			Expect(errors.Is(errRemoveAll, json_dataset.ErrFileNotExists)).To(BeTrue())
		})
	})

	Context("Update the Year, URL and tropes of a Film in the JSON file", func() {
		var errUpdate error

		BeforeEach(func() {
			errAddMedia = repository.AddMedia(mediaEntry)
			errPersist = repository.Persist()

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
			var dataset json_dataset.JSONDataset
			fileContents, _ := os.ReadFile("dataset.json")
			err := json.Unmarshal(fileContents, &dataset)

			Expect(err).To(BeNil())
			Expect(errPersist).To(BeNil())
			Expect(len(dataset.Tropestogo)).To(Equal(1))
			Expect(dataset.Tropestogo[0].Title).To(Not(BeEmpty()))
			Expect(dataset.Tropestogo[0].Year).To(Not(BeEmpty()))
			Expect(dataset.Tropestogo[0].URL).To(Not(BeEmpty()))
			Expect(dataset.Tropestogo[0].MediaType).To(Not(BeEmpty()))
			Expect(len(dataset.Tropestogo[0].Tropes) > 0).To(BeTrue())
			Expect(len(dataset.Tropestogo[0].SubTropes) > 0).To(BeTrue())
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

		It("Should only be one Media record on the JSON file", func() {
			var dataset json_dataset.JSONDataset
			fileContents, _ := os.ReadFile("dataset.json")
			err := json.Unmarshal(fileContents, &dataset)

			Expect(err).To(BeNil())
			Expect(len(dataset.Tropestogo)).To(Equal(1))
		})
	})
})

var _ = AfterSuite(func() {
	datasetFile.Close()
	os.Remove("dataset.json")
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
