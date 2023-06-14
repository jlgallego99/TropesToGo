package media_test

import (
	"errors"
	"net/url"
	"time"

	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	avengersUrl = "https://tvtropes.org/pmwiki/pmwiki.php/Film/TheAvengers2012"
)

var _ = Describe("Media", func() {
	var tvTropesPage *tropestogo.Page
	var lastUpdated time.Time
	tropes := make(map[tropestogo.Trope]struct{})

	BeforeEach(func() {
		trope1, _ := tropestogo.NewTrope("AccentUponTheWrongSyllable", tropestogo.TropeIndex(0))
		trope2, _ := tropestogo.NewTrope("ChekhovsGun", tropestogo.TropeIndex(0))
		tropes[trope1] = struct{}{}
		tropes[trope2] = struct{}{}
		lastUpdated = time.Now()

		tvTropesUrl, _ := url.Parse(avengersUrl)
		tvTropesPage, _ = tropestogo.NewPage(tvTropesUrl)
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
				mediaNoPage, errMediaNoPage = media.NewMedia("TheAvengers", "2012", lastUpdated, tropes, nil, media.Film)
			})

			It("Should return an empty object", func() {
				Expect(mediaNoPage.GetWork()).To(BeNil())
				Expect(mediaNoPage.GetPage()).To(BeNil())
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
				mediaNoTitle, errMediaNoTitle = media.NewMedia("", "2012", lastUpdated, tropes, nil, media.Film)
			})

			It("Should return an empty object", func() {
				Expect(mediaNoTitle.GetWork()).To(BeNil())
				Expect(mediaNoTitle.GetPage()).To(BeNil())
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
				Expect(mediaNoType.GetPage()).To(BeNil())
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
				Expect(mediaWrongYear.GetPage()).To(BeNil())
				Expect(mediaWrongYear.GetMediaType()).To(Equal(media.MediaType(0)))
			})

			It("Should raise a proper error", func() {
				Expect(errors.Is(errMediaWrongYear, media.ErrInvalidYear)).To(BeTrue())
			})
		})
	})
})
