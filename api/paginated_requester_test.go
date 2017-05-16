package api_test

import (
	"errors"

	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/api/apifakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PaginatedRequester", func() {
	var fakeRequestFactory *apifakes.FakeRequestFactory
	var fakeCloudControllerClient *apifakes.FakeCloudControllerClient
	var fakePaginatedParser *apifakes.FakePaginatedParser
	var fakeFilter *apifakes.FakeFilter
	var params map[string]interface{}
	var testRequest *http.Request
	var testResponse *http.Response

	var paginatedRequester *api.PaginatedRequester

	var err error
	var responseBodies [][]byte

	generateApiResponse := func(body string) *http.Response {
		return &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(body)),
		}
	}

	BeforeEach(func() {
		fakeCloudControllerClient = new(apifakes.FakeCloudControllerClient)
		fakePaginatedParser = new(apifakes.FakePaginatedParser)
		fakeRequestFactory = new(apifakes.FakeRequestFactory)
		fakeFilter = new(apifakes.FakeFilter)
		params = make(map[string]interface{})

		testRequest, err = http.NewRequest("GET", "something", strings.NewReader(""))
		Expect(err).NotTo(HaveOccurred())
		testResponse = generateApiResponse("")

		fakeRequestFactory.Returns(testRequest, nil)
		fakeCloudControllerClient.DoReturns(testResponse, nil)

		paginatedRequester = &api.PaginatedRequester{
			RequestFactory: fakeRequestFactory.Spy,
			Client:         fakeCloudControllerClient,
			PageParser:     fakePaginatedParser,
		}
	})

	JustBeforeEach(func() {
		responseBodies, err = paginatedRequester.Do(fakeFilter, params)
	})

	It("should create a request", func() {
		Expect(fakeRequestFactory.CallCount()).To(Equal(1))
		filters, _ := fakeRequestFactory.ArgsForCall(0)
		Expect(filters).To(Equal(fakeFilter))
	})

	Context("when creating the request fails", func() {
		var disaster = errors.New("OH NOOOOOOO")
		BeforeEach(func() {
			fakeRequestFactory.Returns(new(http.Request), disaster)
		})

		It("should return the error", func() {
			Expect(responseBodies).To(BeEmpty())
			Expect(err).To(Equal(disaster))
		})
	})

	Context("when creating the request succeeds", func() {
		BeforeEach(func() {
			fakeRequestFactory.Returns(testRequest, nil)
		})

		It("should make a request", func() {
			Expect(fakeCloudControllerClient.DoCallCount()).To(Equal(1))

			_, err := http.NewRequest("GET", "something", strings.NewReader(""))
			Expect(err).NotTo(HaveOccurred())

			actualRequest := fakeCloudControllerClient.DoArgsForCall(0)
			Expect(actualRequest.Method).To(Equal("GET"))
			Expect(actualRequest.URL.Path).To(Equal("something"))
		})

		Context("when making the request fails", func() {
			var requestError error

			BeforeEach(func() {
				requestError = errors.New("request execution failed")
				fakeCloudControllerClient.DoReturns(new(http.Response), requestError)
			})

			It("should return the request error", func() {
				Expect(responseBodies).To(BeEmpty())
				Expect(err).To(Equal(requestError))
			})
		})

		Context("when making the request succeeds", func() {
			response := generateApiResponse("")

			BeforeEach(func() {
				fakeCloudControllerClient.DoReturns(response, nil)
			})

			It("parses it to find out the number of pages", func() {
				Expect(fakePaginatedParser.ParseCallCount()).To(Equal(1))
			})

			Context("when parsing for the number of pages fails", func() {
				var parseErr error

				BeforeEach(func() {
					parseErr = errors.New("some err")
					fakePaginatedParser.ParseReturns(api.PaginatedResponse{}, parseErr)
				})

				It("returns the error", func() {
					Expect(responseBodies).To(BeEmpty())
					Expect(err).To(Equal(parseErr))
				})
			})

			Context("when parsing for the number of pages succeeds", func() {
				Context("when there's only one page", func() {
					BeforeEach(func() {
						fakePaginatedParser.ParseReturns(api.PaginatedResponse{
							TotalPages: 1,
						}, nil)
					})

					It("does not make more API calls", func() {
						Expect(fakeRequestFactory.CallCount()).To(Equal(1))

						_, params := fakeRequestFactory.ArgsForCall(0)
						Expect(params["page"]).To(BeNil())
					})
				})

				Context("when there's more than one page", func() {
					BeforeEach(func() {
						fakePaginatedParser.ParseReturns(api.PaginatedResponse{
							TotalPages: 2,
						}, nil)

						testResponse = generateApiResponse("some-body")
						testResponse2 := generateApiResponse("some-second-body")

						var i int
						fakeCloudControllerClient.DoStub = func(*http.Request) (*http.Response, error) {
							i += 1

							if i == 1 {
								return testResponse, nil
							}

							return testResponse2, nil
						}
					})

					It("calls for more results", func() {
						Expect(fakeRequestFactory.CallCount()).To(Equal(2))

						_, params := fakeRequestFactory.ArgsForCall(1)
						Expect(params["page"]).To(Equal(2))
					})

					It("contains the list of all byte slices of response bodies", func() {
						Expect(responseBodies).To(Equal([][]byte{
							[]byte("some-body"),
							[]byte("some-second-body"),
						}))
					})
				})
			})
		})
	})
})
