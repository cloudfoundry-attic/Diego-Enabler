package flaghelpers_test

import (
	. "github.com/cloudfoundry-incubator/diego-enabler/commands/internal/flaghelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParallelFlag", func() {
	var parallelFlag ParallelFlag
	BeforeEach(func() {
		parallelFlag = ParallelFlag{}
	})

	Describe("positive values", func() {
		Context("value is less than or equal to 100", func() {
			It("does not error", func() {
				Expect(parallelFlag.UnmarshalFlag("1")).ToNot(HaveOccurred())
				Expect(parallelFlag.Value).To(Equal(1))
			})
		})
	})

	Describe("non-positive values", func() {
		Describe("zero values", func() {
			It("returns an error", func() {
				err := parallelFlag.UnmarshalFlag("0")
				_, ok := err.(InvalidParallelValueError)
				Expect(ok).To(BeTrue())
			})
		})

		Describe("negative values", func() {
			It("returns an error", func() {
				err := parallelFlag.UnmarshalFlag("-1")
				_, ok := err.(InvalidParallelValueError)
				Expect(ok).To(BeTrue())
			})
		})

		Describe("non-number values", func() {
			It("returns an error", func() {
				err := parallelFlag.UnmarshalFlag("banana")
				_, ok := err.(InvalidParallelValueError)
				Expect(ok).To(BeTrue())
			})
		})
	})
})
