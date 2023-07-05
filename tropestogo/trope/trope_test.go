package trope_test

import (
	"github.com/jlgallego99/TropesToGo/trope"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Trope", func() {
	var validTrope, tropeNoTitle, tropeNoIndex trope.Trope
	var errValidTrope, errNoTitle, errNoIndex error

	BeforeEach(func() {
		validTrope, errValidTrope = trope.NewTrope("ChekhovsGun", trope.NarrativeTrope, "")
		tropeNoTitle, errNoTitle = trope.NewTrope("", trope.TopicalTrope, "")
		tropeNoIndex, errNoIndex = trope.NewTrope("ChekhovsGun", trope.TropeIndex(100), "")
	})

	Describe("Create a Trope", func() {
		Context("The Trope is created correctly", func() {
			It("Shouldn't return an empty object", func() {
				Expect(validTrope.GetTitle()).To(Equal("ChekhovsGun"))
				Expect(validTrope.GetIndex()).To(Equal(trope.NarrativeTrope))
				Expect(validTrope.GetSubpage()).To(BeEmpty())
				Expect(validTrope.GetIsMain()).To(BeTrue())
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
				Expect(errNoTitle).To(Equal(trope.ErrMissingValues))
			})
		})

		Context("The Trope doesn't have a valid index", func() {
			It("Should return an empty object", func() {
				Expect(tropeNoIndex.GetTitle()).To(BeEmpty())
				Expect(tropeNoIndex.GetIndex()).To(BeZero())
			})

			It("Should raise a proper error", func() {
				Expect(errNoIndex).To(Equal(trope.ErrUnknownIndex))
			})
		})
	})
})
