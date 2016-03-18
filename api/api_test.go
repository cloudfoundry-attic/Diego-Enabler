package api_test

import (
	. "github.com/cloudfoundry-incubator/diego-enabler/api"

	fakes "github.com/cloudfoundry-incubator/diego-enabler/api/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Api", func() {
	Describe("NewGetAppsRequest", func() {
		var (
			apiClient  *ApiClient
			fakeFilter *fakes.FakeFilter
			params     map[string]interface{}
			baseUrl    string

			request *http.Request
			err     error
		)

		BeforeEach(func() {
			baseUrl = "https://api.my-crazy-domain.com"
			fakeFilter = new(fakes.FakeFilter)
			params = map[string]interface{}{}
		})

		JustBeforeEach(func() {
			apiClient, err = NewApiClient(baseUrl)
			Expect(err).NotTo(HaveOccurred())

			request, err = apiClient.NewGetAppsRequest(fakeFilter, params)
		})

		It("works", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when given filters", func() {
			BeforeEach(func() {
				fakeFilter.ToFilterQueryParamReturns("something")
			})

			It("puts the filter into `q`", func() {
				Expect(request.URL.Query().Get("q")).To(Equal("something"))
			})
		})

		Context("when given params", func() {
			BeforeEach(func() {
				params = map[string]interface{}{"param1": "paramValue", "param2": "paramValue2"}
			})

			It("adds the params to the request", func() {
				Expect(request.URL.Query().Get("param1")).To(Equal("paramValue"))
				Expect(request.URL.Query().Get("param2")).To(Equal("paramValue2"))
			})
		})

		It("hits the appropriate API URL", func() {
			Expect(request.Method).To(Equal("GET"))
			Expect(request.URL.String()).To(Equal("https://api.my-crazy-domain.com/v2/apps"))
		})
	})

	Describe("EqualFilter", func() {
		It("serializes to name:val", func() {
			filter := EqualFilter{
				Name:  "foo",
				Value: true,
			}

			Expect(filter.ToFilterQueryParam()).To(Equal("foo:true"))

			filter = EqualFilter{
				Name:  "something",
				Value: 2,
			}

			Expect(filter.ToFilterQueryParam()).To(Equal("something:2"))

			filter = EqualFilter{
				Name:  "quux",
				Value: "bar",
			}

			Expect(filter.ToFilterQueryParam()).To(Equal("quux:bar"))
		})
	})

	Describe("Filters", func() {
		It("combines its filters together with semicolons", func() {
			filter1 := new(fakes.FakeFilter)
			filter1.ToFilterQueryParamReturns("something>2")

			filter2 := new(fakes.FakeFilter)
			filter2.ToFilterQueryParamReturns("bar::baaz")

			filters := Filters{
				filter1,
				filter2,
			}

			Expect(filters.ToFilterQueryParam()).To(Equal("something>2;bar::baaz"))
		})
	})
})
