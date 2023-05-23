package media_test

import (
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/url"
	"time"
)

var _ = Describe("Media", func() {
	var validMedia, mediaNoPage, mediaNoTitle, mediaNoType, mediaWrongYear media.Media
	var errValidMedia, errMediaNoPage, errMediaNoTitle, errMediaNoType, errMediaWrongYear error
	var tvTropesPage *tropestogo.Page
	var lastUpdated time.Time

	BeforeEach(func() {
		tropes := make(map[tropestogo.Trope]struct{})
		trope1, _ := tropestogo.NewTrope("AccentUponTheWrongSyllable", tropestogo.TropeIndex(0))
		trope2, _ := tropestogo.NewTrope("ChekhovsGun", tropestogo.TropeIndex(0))
		tropes[trope1] = struct{}{}
		tropes[trope2] = struct{}{}
		lastUpdated = time.Now()

		tvTropesUrl, _ := url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/TheAvengers2012")
		tvTropesPage = &tropestogo.Page{
			URL:         tvTropesUrl,
			LastUpdated: time.Now(),
		}

		validMedia, errValidMedia = media.NewMedia("TheAvengers", "2012", lastUpdated, tropes, tvTropesPage, media.Film)
		mediaNoPage, errMediaNoPage = media.NewMedia("TheAvengers", "2012", lastUpdated, tropes, nil, media.Film)
		mediaNoTitle, errMediaNoTitle = media.NewMedia("", "2012", lastUpdated, tropes, nil, media.Film)
		mediaNoType, errMediaNoType = media.NewMedia("TheAvengers", "2012", lastUpdated, tropes, tvTropesPage, media.MediaType(100))
		mediaWrongYear, errMediaWrongYear = media.NewMedia("TheAvengers", "2012aaaaa", lastUpdated, tropes, tvTropesPage, media.Film)
	})

	Describe("Create Media", func() {
		Context("The Media is created correctly", func() {
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
			It("Should return an empty object", func() {
				Expect(mediaNoType.GetWork()).To(BeNil())
				Expect(mediaNoType.GetPage()).To(BeNil())
				Expect(mediaNoType.GetMediaType()).To(Equal(media.MediaType(0)))
			})

			It("Should raise a proper error", func() {
				Expect(errMediaNoType).To(Equal(media.ErrUnsupportedMediaType))
			})
		})

		Context("The Media year is not a valid year number", func() {
			It("Should return an empty object", func() {
				Expect(mediaWrongYear.GetWork()).To(BeNil())
				Expect(mediaWrongYear.GetPage()).To(BeNil())
				Expect(mediaWrongYear.GetMediaType()).To(Equal(media.MediaType(0)))
			})

			It("Should raise a proper error", func() {
				Expect(errMediaWrongYear).To(Equal(media.ErrInvalidYear))
			})
		})
	})
})
