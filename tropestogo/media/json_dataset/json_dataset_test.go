package json_dataset_test

import (
	"encoding/json"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"github.com/jlgallego99/TropesToGo/media/csv_dataset"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/url"
	"os"
	"time"
)

var repository *json_dataset.JSONRepository
var errorRepository, errRemoveAll error
var mediaEntry media.Media
var datasetFile *os.File

var _ = BeforeSuite(func() {
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

var _ = Describe("JsonDataset", func() {
	BeforeEach(func() {
		repository, errorRepository = json_dataset.NewJSONRepository("dataset")
		datasetFile, _ = os.Open("dataset.json")
	})

	AfterEach(func() {
		// Reset file
		//repository.RemoveAll()
	})

	Context("Create JSON repository", func() {
		It("Should have created a JSON file", func() {
			Expect("dataset.json").To(BeAnExistingFile())
		})

		It("Shouldn't return an error", func() {
			Expect(errorRepository).To(BeNil())
		})
	})

	Context("Add a Media to the JSON file", func() {
		var errAddMedia error

		BeforeEach(func() {
			errAddMedia = repository.AddMedia(mediaEntry)
		})

		It("Should have all the correct fields", func() {
			var dataset json_dataset.JSONDataset
			fileContents, _ := os.ReadFile("dataset.json")
			err := json.Unmarshal(fileContents, &dataset)

			Expect(err).To(BeNil())
			Expect(dataset.Tropestogo[0].Title).To(Equal("Oldboy"))
			Expect(dataset.Tropestogo[0].Year).To(Equal("2003"))
			Expect(dataset.Tropestogo[0].URL).To(Equal("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003"))
			Expect(dataset.Tropestogo[0].MediaType).To(Equal("Film"))
			Expect(len(dataset.Tropestogo[0].Tropes)).To(Equal(3))
		})

		It("Shouldn't return an error", func() {
			Expect(errAddMedia).To(BeNil())
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

		It("Shouldn't exist a JSON file", func() {
			Expect("dataset.json").To(Not(BeAnExistingFile()))
		})

		It("Should return an error", func() {
			Expect(errRemoveAll).To(Equal(csv_dataset.ErrFileNotExists))
		})
	})
})

var _ = AfterSuite(func() {
	datasetFile.Close()
	//os.Remove("dataset.json")
})
