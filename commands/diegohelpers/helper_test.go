package diegohelpers_test

import (
	"github.com/cloudfoundry-incubator/diego-enabler/api/apifakes"
	. "github.com/cloudfoundry-incubator/diego-enabler/commands/diegohelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DiegoHelpers", func() {
	var (
		fakeApi *apifakes.FakeConnection
	)

	BeforeEach(func() {
		fakeApi = new(apifakes.FakeConnection)
	})

	Describe("ToggleDiegoSupport", func() {
		Context("when disabling diego", func() {
			BeforeEach(func() {
				ToggleDiegoSupport(false, fakeApi, "some-app")
			})

			It("should not check that there are no routes", func() {
				Expect(fakeApi.GetCurrentSpaceCallCount()).To(Equal(0))
				Expect(fakeApi.GetSpaceCallCount()).To(Equal(0))
				Expect(fakeApi.UsernameCallCount()).To(Equal(0))
			})
		})

		Context("when enabling diego", func() {
			BeforeEach(func() {
				ToggleDiegoSupport(true, fakeApi, "some-app")
			})

			It("should not check that there are no routes", func() {
				Expect(fakeApi.GetCurrentSpaceCallCount()).To(Equal(1))
				Expect(fakeApi.GetSpaceCallCount()).To(Equal(1))
				Expect(fakeApi.UsernameCallCount()).To(Equal(1))
			})
		})
	})
})
