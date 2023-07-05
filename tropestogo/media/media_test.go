package media_test

import (
	"errors"
	"fmt"
	trope "github.com/jlgallego99/TropesToGo/trope"
	"github.com/jlgallego99/TropesToGo/tvtropespages"
	"math/rand"
	"time"

	"github.com/jlgallego99/TropesToGo/media"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	avengersUrl = "https://tvtropes.org/pmwiki/pmwiki.php/Film/TheAvengers2012"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var _ = Describe("Media", func() {
	var tvTropesPage tvtropespages.Page
	var lastUpdated time.Time
	tropes := make(map[trope.Trope]struct{})

	BeforeEach(func() {
		trope1, _ := trope.NewTrope("AccentUponTheWrongSyllable", trope.TropeIndex(0), "")
		trope2, _ := trope.NewTrope("ChekhovsGun", trope.TropeIndex(0), "")
		tropes[trope1] = struct{}{}
		tropes[trope2] = struct{}{}
		lastUpdated = time.Now()

		tvTropesPage, _ = tvtropespages.NewPage(avengersUrl, false, nil)
	})

	AfterEach(func() {
		tropes = make(map[trope.Trope]struct{})
	})

	Describe("Create Media", func() {
		Context("The Media is created correctly", func() {
			var validMedia media.Media
			var errValidMedia error

			BeforeEach(func() {
				validMedia, errValidMedia = media.NewMedia("TheAvengers", "2012", lastUpdated, tropes, tvTropesPage, media.Film)
			})

			It("Should return a valid object", func() {
				Expect(validMedia.GetWork()).To(Not(BeNil()))
				Expect(validMedia.GetPage()).To(Not(BeNil()))
				Expect(validMedia.GetMediaType()).To(Equal(media.Film))
				Expect(len(validMedia.GetWork().Tropes)).To(Equal(2))
				Expect(validMedia.GetWork().Title).To(Equal("TheAvengers"))
				Expect(validMedia.GetWork().Year).To(Equal("2012"))
				Expect(validMedia.GetWork().LastUpdated).To(Equal(lastUpdated))
				Expect(validMedia.GetPage()).To(Equal(tvTropesPage))
			})

			It("Shouldn't raise an error", func() {
				Expect(errValidMedia).To(BeNil())
			})
		})

		Context("The Media has no page", func() {
			var mediaNoPage media.Media
			var errMediaNoPage error

			BeforeEach(func() {
				mediaNoPage, errMediaNoPage = media.NewMedia("TheAvengers", "2012", lastUpdated, tropes, tvtropespages.Page{}, media.Film)
			})

			It("Should return an empty object", func() {
				Expect(mediaNoPage.GetWork()).To(BeNil())
				Expect(mediaNoPage.GetPage().GetUrl()).To(BeNil())
				Expect(mediaNoPage.GetPage().GetPageType()).To(BeZero())
				Expect(mediaNoPage.GetMediaType()).To(Equal(media.MediaType(0)))
			})

			It("Should raise a proper error", func() {
				Expect(errMediaNoPage).To(Equal(media.ErrMissingValues))
			})
		})

		Context("The Media has no title", func() {
			var mediaNoTitle media.Media
			var errMediaNoTitle error

			BeforeEach(func() {
				mediaNoTitle, errMediaNoTitle = media.NewMedia("", "2012", lastUpdated, tropes, tvtropespages.Page{}, media.Film)
			})

			It("Should return an empty object", func() {
				Expect(mediaNoTitle.GetWork()).To(BeNil())
				Expect(mediaNoTitle.GetPage().GetUrl()).To(BeNil())
				Expect(mediaNoTitle.GetPage().GetPageType()).To(BeZero())
				Expect(mediaNoTitle.GetMediaType()).To(Equal(media.MediaType(0)))
			})

			It("Should raise a proper error", func() {
				Expect(errMediaNoTitle).To(Equal(media.ErrMissingValues))
			})
		})

		Context("The Media has no media type", func() {
			var mediaNoType media.Media
			var errMediaNoType error

			BeforeEach(func() {
				mediaNoType, errMediaNoType = media.NewMedia("TheAvengers", "2012", lastUpdated, tropes, tvTropesPage, media.MediaType(100))
			})

			It("Should return an empty object", func() {
				Expect(mediaNoType.GetWork()).To(BeNil())
				Expect(mediaNoType.GetPage().GetUrl()).To(BeNil())
				Expect(mediaNoType.GetPage().GetPageType()).To(BeZero())
				Expect(mediaNoType.GetMediaType()).To(Equal(media.MediaType(0)))
			})

			It("Should raise a proper error", func() {
				Expect(errors.Is(errMediaNoType, media.ErrUnknownMediaType)).To(BeTrue())
			})
		})

		Context("The Media year is not a valid year number", func() {
			var mediaWrongYear media.Media
			var errMediaWrongYear error

			BeforeEach(func() {
				mediaWrongYear, errMediaWrongYear = media.NewMedia("TheAvengers", "2012aaaaa", lastUpdated, tropes, tvTropesPage, media.Film)
			})

			It("Should return an empty object", func() {
				Expect(mediaWrongYear.GetWork()).To(BeNil())
				Expect(mediaWrongYear.GetPage().GetUrl()).To(BeNil())
				Expect(mediaWrongYear.GetPage().GetPageType()).To(BeZero())
				Expect(mediaWrongYear.GetMediaType()).To(Equal(media.MediaType(0)))
			})

			It("Should raise a proper error", func() {
				Expect(errors.Is(errMediaWrongYear, media.ErrInvalidYear)).To(BeTrue())
			})
		})

		Context("The Media has same SubTropes on different SubWikis", func() {
			const max = 10
			const min = 2
			numTropes := seededRand.Intn(max-min) + min
			var mediaAllTropes media.Media
			var errMediaAllTropes error

			BeforeEach(func() {
				tropes = createTropes(numTropes, randomTrope)
				subTropes := createTropes(numTropes, randomSubTrope)
				for subTrope := range subTropes {
					tropes[subTrope] = struct{}{}
				}

				mediaAllTropes, errMediaAllTropes = media.NewMedia("TheAvengers", "2012", lastUpdated, tropes, tvTropesPage, media.Film)
			})

			It("Shouldn't return an error", func() {
				Expect(errMediaAllTropes).To(BeNil())
			})

			It("Should have main tropes and sub tropes", func() {
				Expect(mediaAllTropes.GetWork().Tropes).To(Not(BeEmpty()))
				Expect(mediaAllTropes.GetWork().SubTropes).To(Not(BeEmpty()))
			})

			It("Shouldn't have repeated tropes or sub tropes", func() {
				Expect(areTropesUnique(mediaAllTropes.GetWork().Tropes)).To(BeTrue())
				Expect(areTropesUnique(mediaAllTropes.GetWork().SubTropes)).To(BeTrue())
			})

			It("Should have added all tropes and SubTropes because they are from different SubWikis", func() {
				Expect(len(mediaAllTropes.GetWork().Tropes) + len(mediaAllTropes.GetWork().SubTropes)).To(Equal(numTropes * 2))
			})
		})
	})
})

func areTropesUnique(tropes map[trope.Trope]struct{}) bool {
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

// createTropes generates a map of numTropes size applying a callback function to all elements
func createTropes(numTropes int, callback func() trope.Trope) map[trope.Trope]struct{} {
	tropeset := make(map[trope.Trope]struct{}, numTropes)

	for i := 0; i < numTropes; i++ {
		tropeset[callback()] = struct{}{}
	}

	return tropeset
}

var randomTrope = func() trope.Trope {
	trope, _ := trope.NewTrope("Trope"+fmt.Sprint(seededRand.Int()), 1, "")
	return trope
}

var randomSubTrope = func() trope.Trope {
	subWikis := []string{"SubWiki1", "SubWiki2"}
	trope, _ := trope.NewTrope("Trope"+fmt.Sprint(seededRand.Int()), 1, subWikis[seededRand.Intn(1)])

	return trope
}
