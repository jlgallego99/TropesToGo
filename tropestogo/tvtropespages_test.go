package tropestogo_test

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/url"
	"time"
)

var tvtropespages *tropestogo.TvTropesPages

var _ = BeforeSuite(func() {
	tvtropespages = tropestogo.NewTvTropesPages()
})

var _ = Describe("Tvtropespages", func() {
	AfterEach(func() {
		tvtropespages.Pages = make(map[tropestogo.Page]time.Time, 0)
	})

	Context("Add a Film, Trope and Index Pages to the entity", func() {
		var errAddValid, errAddValid2, errAddValid3 error

		BeforeEach(func() {
			errAddValid = tvtropespages.AddTvTropesPage(oldboyUrl)
			errAddValid2 = tvtropespages.AddTvTropesPage(tropeUrl)
			errAddValid3 = tvtropespages.AddTvTropesPage(indexUrl)
		})

		It("Shouldn't return an error", func() {
			Expect(errAddValid).To(BeNil())
			Expect(errAddValid2).To(BeNil())
			Expect(errAddValid3).To(BeNil())
		})

		It("Should have three Pages", func() {
			Expect(len(tvtropespages.Pages)).To(Equal(3))
		})

		It("Should have added the correct Pages", func() {
			var urls []*url.URL
			var pageTypes []tropestogo.PageType
			for page, lastUpdated := range tvtropespages.Pages {
				urls = append(urls, page.GetUrl())
				pageTypes = append(pageTypes, page.GetPageType())

				Expect(lastUpdated).To(Not(BeNil()))
			}

			Expect(urls[0].String()).To(Equal(oldboyUrl))
			Expect(pageTypes[0]).To(Equal(tropestogo.WorkPage))
			Expect(urls[1].String()).To(Equal(tropeUrl))
			Expect(pageTypes[1]).To(Equal(tropestogo.MainPage))
			Expect(urls[2].String()).To(Equal(indexUrl))
			Expect(pageTypes[2]).To(Equal(tropestogo.IndexPage))
		})
	})

	Context("Add a duplicated Page", func() {
		var errAddDuplicated error

		BeforeEach(func() {
			errAddDuplicated = tvtropespages.AddTvTropesPage(oldboyUrl)
			errAddDuplicated = tvtropespages.AddTvTropesPage(oldboyUrl)
		})

		It("Should return an error", func() {
			Expect(errors.Is(errAddDuplicated, tropestogo.ErrDuplicatedPage)).To(BeTrue())
		})

		It("Shouldn't have added the second Page", func() {
			Expect(len(tvtropespages.Pages)).To(Equal(1))
		})
	})

	Context("Add and empty string url", func() {
		var errAddEmpty error

		BeforeEach(func() {
			errAddEmpty = tvtropespages.AddTvTropesPage("")
		})

		It("Should return an error", func() {
			Expect(errors.Is(errAddEmpty, tropestogo.ErrEmptyUrl)).To(BeTrue())
		})

		It("Shouldn't have added anything", func() {
			Expect(len(tvtropespages.Pages)).To(Equal(0))
		})
	})

	Context("Add a badly formated url", func() {
		var errAddEmpty error

		BeforeEach(func() {
			errAddEmpty = tvtropespages.AddTvTropesPage("htp$p%^^^&&***!!!!!")
		})

		It("Should return an error", func() {
			Expect(errors.Is(errAddEmpty, tropestogo.ErrBadUrl)).To(BeTrue())
		})

		It("Shouldn't have added anything", func() {
			Expect(len(tvtropespages.Pages)).To(Equal(0))
		})
	})
})
