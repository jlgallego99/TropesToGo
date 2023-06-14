package tropestogo_test

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("Page", func() {
	var validPage, invalidPage, nullPage *tropestogo.Page
	var validUrl, invalidUrl *url.URL
	var errValidPage, errInvalidPage, errNullPage error

	Context("Create a TvTropes Page object", func() {
		BeforeEach(func() {
			validUrl, _ = url.Parse("https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003")
			validPage, errValidPage = tropestogo.NewPage(validUrl)
		})

		It("Shouldn't return an error", func() {
			Expect(errValidPage).To(BeNil())
		})

		It("Should return a correct Page object", func() {
			Expect(validPage.URL).To(Equal(validUrl))
			Expect(validPage.LastUpdated).To(Not(BeZero()))
		})
	})

	Context("Create a Page object of a web that isn't TvTropes", func() {
		BeforeEach(func() {
			invalidUrl, _ = url.Parse("htttps://google.com")
			invalidPage, errInvalidPage = tropestogo.NewPage(invalidUrl)
		})

		It("Should return an error", func() {
			Expect(errors.Is(errInvalidPage, tropestogo.ErrNotTvTropes)).To(BeTrue())
		})

		It("Shouldn't return a Page object", func() {
			Expect(invalidPage).To(BeNil())
		})
	})

	Context("Create a Page with a null URL", func() {
		BeforeEach(func() {
			nullPage, errNullPage = tropestogo.NewPage(nil)
		})

		It("Should return an error", func() {
			Expect(errors.Is(errNullPage, tropestogo.ErrBadUrl)).To(BeTrue())
		})

		It("Shouldn't return a Page object", func() {
			Expect(nullPage).To(BeNil())
		})
	})
})
