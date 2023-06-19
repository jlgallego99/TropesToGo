package tropestogo_test

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	oldboySuburls = []string{"https://tvtropes.org/pmwiki/pmwiki.php/Awesome/Oldboy2003",
		"https://tvtropes.org/pmwiki/pmwiki.php/Fridge/Oldboy2003",
		"https://tvtropes.org/pmwiki/pmwiki.php/Laconic/Oldboy2003",
		"https://tvtropes.org/pmwiki/pmwiki.php/Trivia/Oldboy2003",
		"https://tvtropes.org/pmwiki/pmwiki.php/YMMV/Oldboy2003"}
)

var tvtropespages *tropestogo.TvTropesPages

var _ = BeforeSuite(func() {
	tvtropespages = tropestogo.NewTvTropesPages()
})

var _ = Describe("Tvtropespages", func() {
	AfterEach(func() {
		tvtropespages.Pages = make(map[tropestogo.Page]*tropestogo.TvTropesSubpages, 0)
	})

	Context("Add a Film page with subpages to the entity", func() {
		var errAddValid error

		BeforeEach(func() {
			errAddValid = tvtropespages.AddTvTropesPage(oldboyUrl, oldboySuburls)
		})

		It("Shouldn't return an error", func() {
			Expect(errAddValid).To(BeNil())
		})

		It("Should have added the Pages", func() {
			Expect(len(tvtropespages.Pages) > 0).To(BeTrue())
		})

		It("Should have added the Subpages", func() {
			for mainpage, subpages := range tvtropespages.Pages {
				Expect(mainpage.GetPageType()).To(Not(BeZero()))
				Expect(mainpage.GetUrl()).To(Not(BeNil()))
				Expect(len(subpages.Subpages) > 0).To(BeTrue())

				for subpage := range subpages.Subpages {
					Expect(subpage.GetUrl()).To(Not(BeNil()))
				}
			}
		})
	})

	Context("Add Trope and Index Pages without subpages to the entity", func() {
		var errAddValid2, errAddValid3 error

		BeforeEach(func() {
			errAddValid2 = tvtropespages.AddTvTropesPage(tropeUrl, []string{})
			errAddValid3 = tvtropespages.AddTvTropesPage(indexUrl, []string{})
		})

		It("Shouldn't return an error", func() {
			Expect(errAddValid2).To(BeNil())
			Expect(errAddValid3).To(BeNil())
		})

		It("Should have added the Pages", func() {
			Expect(len(tvtropespages.Pages) > 0).To(BeTrue())
		})

		It("Should have no Subpages", func() {
			for _, subpage := range tvtropespages.Pages {
				Expect(len(subpage.Subpages)).To(BeZero())
			}
		})
	})

	Context("Add a duplicated Page", func() {
		var errAddDuplicated error

		BeforeEach(func() {
			errAddDuplicated = tvtropespages.AddTvTropesPage(oldboyUrl, oldboySuburls)
			errAddDuplicated = tvtropespages.AddTvTropesPage(oldboyUrl, oldboySuburls)
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
			errAddEmpty = tvtropespages.AddTvTropesPage("", []string{})
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
			errAddEmpty = tvtropespages.AddTvTropesPage("htp$p%^^^&&***!!!!!", []string{})
		})

		It("Should return an error", func() {
			Expect(errors.Is(errAddEmpty, tropestogo.ErrBadUrl)).To(BeTrue())
		})

		It("Shouldn't have added anything", func() {
			Expect(len(tvtropespages.Pages)).To(Equal(0))
		})
	})
})
