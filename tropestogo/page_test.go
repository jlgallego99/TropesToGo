package tropestogo_test

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/url"
)

const (
	oldboyUrl = "https://tvtropes.org/pmwiki/pmwiki.php/Film/Oldboy2003"
	tropeUrl  = "https://tvtropes.org/pmwiki/pmwiki.php/Main/AboveGoodAndEvil"
	googleUrl = "https://www.google.com/"
	indexUrl  = "https://tvtropes.org/pmwiki/pagelist_having_pagetype_in_namespace.php?n=Film&t=work"
)

var _ = Describe("Page", func() {
	var validPage, invalidPage, nullPage tropestogo.Page
	var validUrl, invalidUrl *url.URL
	var errValidPage, errInvalidPage, errNullPage error

	Context("Create a TvTropes Page object of a Film page", func() {
		BeforeEach(func() {
			validUrl, _ = url.Parse(oldboyUrl)
			validPage, errValidPage = tropestogo.NewPage(validUrl)
		})

		It("Shouldn't return an error", func() {
			Expect(errValidPage).To(BeNil())
		})

		It("Should return a correct Page object", func() {
			Expect(validPage.GetUrl()).To(Not(BeNil()))
			Expect(validPage.GetPageType()).To(Not(BeZero()))
			Expect(validPage.GetUrl()).To(Equal(validUrl))
		})

		It("Should be of type WorkPage", func() {
			Expect(validPage.GetPageType()).To(Equal(tropestogo.WorkPage))
		})
	})

	Context("Create a TvTropes Page object of a Trope page", func() {
		BeforeEach(func() {
			validUrl, _ = url.Parse(tropeUrl)
			validPage, errValidPage = tropestogo.NewPage(validUrl)
		})

		It("Shouldn't return an error", func() {
			Expect(errValidPage).To(BeNil())
		})

		It("Should return a correct Page object", func() {
			Expect(validPage.GetUrl()).To(Not(BeNil()))
			Expect(validPage.GetPageType()).To(Not(BeZero()))
			Expect(validPage.GetUrl()).To(Equal(validUrl))
		})

		It("Should be of type MainPage", func() {
			Expect(validPage.GetPageType()).To(Equal(tropestogo.MainPage))
		})
	})

	Context("Create a TvTropes Page object of a Index page", func() {
		BeforeEach(func() {
			validUrl, _ = url.Parse(indexUrl)
			validPage, errValidPage = tropestogo.NewPage(validUrl)
		})

		It("Shouldn't return an error", func() {
			Expect(errValidPage).To(BeNil())
		})

		It("Should return a correct Page object", func() {
			Expect(validPage.GetUrl()).To(Not(BeNil()))
			Expect(validPage.GetPageType()).To(Not(BeZero()))
			Expect(validPage.GetUrl()).To(Equal(validUrl))
		})

		It("Should be of type MainPage", func() {
			Expect(validPage.GetPageType()).To(Equal(tropestogo.IndexPage))
		})
	})

	Context("Create a Page object of a web that isn't TvTropes", func() {
		BeforeEach(func() {
			invalidUrl, _ = url.Parse(googleUrl)
			invalidPage, errInvalidPage = tropestogo.NewPage(invalidUrl)
		})

		It("Should return an error", func() {
			Expect(errors.Is(errInvalidPage, tropestogo.ErrNotTvTropes)).To(BeTrue())
		})

		It("Should return an empty Page object", func() {
			Expect(invalidPage.GetUrl()).To(BeNil())
			Expect(invalidPage.GetPageType()).To(BeZero())
		})
	})

	Context("Create a Page with a null URL", func() {
		BeforeEach(func() {
			nullPage, errNullPage = tropestogo.NewPage(nil)
		})

		It("Should return an error", func() {
			Expect(errors.Is(errNullPage, tropestogo.ErrBadUrl)).To(BeTrue())
		})

		It("Should return an empty Page object", func() {
			Expect(nullPage.GetUrl()).To(BeNil())
			Expect(nullPage.GetPageType()).To(BeZero())
		})
	})
})
