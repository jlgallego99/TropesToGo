package tropestogo_test

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
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
	// Do not log during testing
	zerolog.SetGlobalLevel(zerolog.Disabled)

	tvtropespages = tropestogo.NewTvTropesPages()
})

var _ = Describe("Tvtropespages", func() {
	AfterEach(func() {
		tvtropespages.Pages = make(map[tropestogo.Page]*tropestogo.TvTropesSubpages, 0)
	})

	Context("Add a Film page with subpages to the entity", func() {
		var errAddValid, errAddSubpage error
		var mainPage tropestogo.Page

		BeforeEach(func() {
			mainPage, errAddValid = tvtropespages.AddTvTropesPage(oldboyUrl, false, nil)
			errAddSubpage = tvtropespages.AddSubpages(oldboyUrl, oldboySuburls, false, nil)
		})

		It("Shouldn't return an error", func() {
			Expect(errAddValid).To(BeNil())
			Expect(errAddSubpage).To(BeNil())
		})

		It("Should have added the Pages", func() {
			Expect(isPageEmpty(mainPage)).To(Not(BeTrue()))
			Expect(len(tvtropespages.Pages) > 0).To(BeTrue())

			for page := range tvtropespages.Pages {
				Expect(page.GetDocument()).To(BeNil())
			}
		})

		It("Should have added the Subpages", func() {
			for mainpage, subpages := range tvtropespages.Pages {
				Expect(mainpage.GetPageType()).To(Not(BeZero()))
				Expect(mainpage.GetUrl()).To(Not(BeNil()))
				Expect(len(subpages.Subpages) > 0).To(BeTrue())

				for subpage := range subpages.Subpages {
					Expect(subpage.GetUrl()).To(Not(BeNil()))
					Expect(subpage.GetDocument()).To(BeNil())
				}
			}
		})
	})

	Context("Add Trope and Index Pages without subpages to the entity", func() {
		var errAddValid2, errAddValid3 error
		var errAddSubpage2, errAddSubpage3 error
		var mainPage2, mainPage3 tropestogo.Page

		BeforeEach(func() {
			mainPage2, errAddValid2 = tvtropespages.AddTvTropesPage(tropeUrl, false, nil)
			errAddSubpage2 = tvtropespages.AddSubpages(tropeUrl, []string{}, false, nil)

			mainPage3, errAddValid3 = tvtropespages.AddTvTropesPage(indexUrl, false, nil)
			errAddSubpage3 = tvtropespages.AddSubpages(indexUrl, []string{}, false, nil)
		})

		It("Shouldn't return an error", func() {
			Expect(errAddValid2).To(BeNil())
			Expect(errAddSubpage2).To(BeNil())

			Expect(errAddValid3).To(BeNil())
			Expect(errAddSubpage3).To(BeNil())
		})

		It("Should have added the Pages", func() {
			Expect(len(tvtropespages.Pages) > 0).To(BeTrue())
			Expect(isPageEmpty(mainPage2)).To(Not(BeTrue()))
			Expect(isPageEmpty(mainPage3)).To(Not(BeTrue()))
		})

		It("Should have no Subpages", func() {
			for _, subpage := range tvtropespages.Pages {
				Expect(len(subpage.Subpages)).To(BeZero())
			}
		})
	})

	Context("Add subpages to an unknown Page", func() {
		var errAddSubpage error

		BeforeEach(func() {
			errAddSubpage = tvtropespages.AddSubpages("NotAnUrl", oldboySuburls, false, nil)
		})

		It("Should return an error", func() {
			Expect(errors.Is(errAddSubpage, tropestogo.ErrNotFound)).To(BeTrue())
		})
	})

	Context("Add a duplicated Page", func() {
		var errAddDuplicated error
		var duplicatedPage tropestogo.Page

		BeforeEach(func() {
			duplicatedPage, errAddDuplicated = tvtropespages.AddTvTropesPage(oldboyUrl, false, nil)
			duplicatedPage, errAddDuplicated = tvtropespages.AddTvTropesPage(oldboyUrl, false, nil)
		})

		It("Should return an error", func() {
			Expect(errors.Is(errAddDuplicated, tropestogo.ErrDuplicatedPage)).To(BeTrue())
		})

		It("Shouldn't have added the second Page", func() {
			Expect(len(tvtropespages.Pages)).To(Equal(1))
		})

		It("Should have returned an empty main Page", func() {
			Expect(isPageEmpty(duplicatedPage)).To(BeTrue())
		})
	})

	Context("Add and empty string url", func() {
		var errAddEmpty error
		var emptyPage tropestogo.Page

		BeforeEach(func() {
			emptyPage, errAddEmpty = tvtropespages.AddTvTropesPage("", false, nil)
		})

		It("Should return an error", func() {
			Expect(errors.Is(errAddEmpty, tropestogo.ErrEmptyUrl)).To(BeTrue())
		})

		It("Shouldn't have added anything", func() {
			Expect(len(tvtropespages.Pages)).To(Equal(0))
		})

		It("Should have returned an empty Page", func() {
			Expect(isPageEmpty(emptyPage)).To(BeTrue())
		})
	})

	Context("Add a badly formated url", func() {
		var errAddEmpty error
		var badPage tropestogo.Page

		BeforeEach(func() {
			badPage, errAddEmpty = tvtropespages.AddTvTropesPage("htp$p%^^^&&***!!!!!", false, nil)
		})

		It("Should return an error", func() {
			Expect(errors.Is(errAddEmpty, tropestogo.ErrBadUrl)).To(BeTrue())
		})

		It("Shouldn't have added anything", func() {
			Expect(len(tvtropespages.Pages)).To(Equal(0))
		})

		It("Should have returned an empty Page", func() {
			Expect(isPageEmpty(badPage)).To(BeTrue())
		})
	})
})

func isPageEmpty(page tropestogo.Page) bool {
	return page.GetPageType() == 0 && page.GetUrl() == nil && page.GetDocument() == nil
}
