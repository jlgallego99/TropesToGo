package tropestogo_test

import (
	tropestogo "github.com/jlgallego99/TropesToGo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Trope", func() {
	var validTrope, tropeNoTitle, tropeNoIndex tropestogo.Trope
	var errValidTrope, errNoTitle, errNoIndex error

	BeforeEach(func() {
		validTrope, errValidTrope = tropestogo.NewTrope("ChekhovsGun", tropestogo.NarrativeTrope)
		tropeNoTitle, errNoTitle = tropestogo.NewTrope("", tropestogo.TopicalTrope)
		tropeNoIndex, errNoIndex = tropestogo.NewTrope("ChekhovsGun", tropestogo.TropeIndex(100))
	})

	Describe("Create a Trope", func() {
		Context("The Trope is created correctly", func() {
			It("Shouldn't return an empty object", func() {
				Expect(validTrope.GetTitle()).To(Equal("ChekhovsGun"))
				Expect(validTrope.GetIndex()).To(Equal(tropestogo.NarrativeTrope))
			})

			It("Shouldn't raise an error", func() {
				Expect(errValidTrope).To(BeNil())
			})
		})

		Context("The Trope doesn't have a title", func() {
			It("Should return an empty object", func() {
				Expect(tropeNoTitle.GetTitle()).To(BeEmpty())
				Expect(tropeNoTitle.GetIndex()).To(BeZero())
			})

			It("Should raise a proper error", func() {
				Expect(errNoTitle).To(Equal(tropestogo.ErrMissingValues))
			})
		})

		Context("The Trope doesn't have a valid index", func() {
			It("Should return an empty object", func() {
				Expect(tropeNoIndex.GetTitle()).To(BeEmpty())
				Expect(tropeNoIndex.GetIndex()).To(BeZero())
			})

			It("Should raise a proper error", func() {
				Expect(errNoIndex).To(Equal(tropestogo.ErrUnknownIndex))
			})
		})
	})
})
