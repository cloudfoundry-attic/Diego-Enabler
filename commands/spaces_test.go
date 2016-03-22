package commands_test

import (
	"errors"

	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/commands"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/fakes"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spaces", func() {
	var fakeRequestFactory *fakes.FakeRequestFactory
	var fakeCloudControllerClient *fakes.FakeCloudControllerClient
	var fakeSpacesParser *fakes.FakeSpacesParser
	var fakePaginatedParser *fakes.FakePaginatedParser
	var spaces models.Spaces
	var testRequest *http.Request
	var testResponse *http.Response

	var err error

	generateApiResponse := func(body string) *http.Response {
		return &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(body)),
		}
	}

	BeforeEach(func() {
		fakeCloudControllerClient = new(fakes.FakeCloudControllerClient)
		fakeSpacesParser = new(fakes.FakeSpacesParser)
		fakePaginatedParser = new(fakes.FakePaginatedParser)
		fakeRequestFactory = new(fakes.FakeRequestFactory)

		testRequest, err = http.NewRequest("GET", "something", strings.NewReader(""))
		Expect(err).NotTo(HaveOccurred())
		testResponse = generateApiResponse("")

		fakeRequestFactory.Returns(testRequest, nil)
		fakeCloudControllerClient.DoReturns(testResponse, nil)
	})

	JustBeforeEach(func() {
		spaces, err = commands.Spaces(fakeRequestFactory.Spy, fakeCloudControllerClient, fakeSpacesParser, fakePaginatedParser)
	})

	It("should create a request", func() {
		expectedParams := map[string]interface{}{
			"inline-relations-depth": 1,
		}

		Expect(fakeRequestFactory.CallCount()).To(Equal(1))
		_, params := fakeRequestFactory.ArgsForCall(0)
		Expect(params).To(Equal(expectedParams))
	})

	Context("when creating the request fails", func() {
		var disaster = errors.New("OH NOOOOOOO")
		BeforeEach(func() {
			fakeRequestFactory.Returns(new(http.Request), disaster)
		})

		It("should return the error", func() {
			Expect(spaces).To(BeEmpty())
			Expect(err).To(Equal(disaster))
		})
	})

	Context("when creating the request succeeds", func() {
		BeforeEach(func() {
			fakeRequestFactory.Returns(testRequest, nil)
		})

		It("should make a request", func() {
			Expect(fakeCloudControllerClient.DoCallCount()).To(Equal(1))

			expectedRequest, err := http.NewRequest("GET", "something", strings.NewReader(""))
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeCloudControllerClient.DoArgsForCall(0)).To(Equal(expectedRequest))
		})

		Context("when making the request fails", func() {
			var requestError error

			BeforeEach(func() {
				requestError = errors.New("request execution failed")
				fakeCloudControllerClient.DoReturns(new(http.Response), requestError)
			})

			It("should return the request error", func() {
				Expect(spaces).To(BeEmpty())
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
					Expect(spaces).To(BeEmpty())
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
					})

					It("calls for more results", func() {
						Expect(fakeRequestFactory.CallCount()).To(Equal(2))

						_, params := fakeRequestFactory.ArgsForCall(1)
						Expect(params["page"]).To(Equal(2))
					})

					Context("when the parsing fails", func() {
						var spaces models.Spaces
						var parseError error

						BeforeEach(func() {
							parseError = errors.New("parsing json failed")
							fakeSpacesParser.ParseReturns(spaces, parseError)
						})

						It("returns the parse error", func() {
							Expect(spaces).To(BeEmpty())
							Expect(err).To(Equal(parseError))
						})
					})

					Context("when the parsing succeeds", func() {
						var parsedSpaces models.Spaces = models.Spaces{
							models.Space{
								models.SpaceEntity{
									Name: "lumpy space princess",
								},
								models.SpaceMetadata{
									Guid: "some-guid",
								},
							},
						}

						BeforeEach(func() {
							// for each call of Parse
							fakeSpacesParser.ParseReturns(parsedSpaces, nil)
						})

						It("returns a list of diego Spaces", func() {
							expectedSpaces := models.Spaces{
								models.Space{
									models.SpaceEntity{
										Name: "lumpy space princess",
									},
									models.SpaceMetadata{
										Guid: "some-guid",
									},
								},
								models.Space{
									models.SpaceEntity{
										Name: "lumpy space princess",
									},
									models.SpaceMetadata{
										Guid: "some-guid",
									},
								},
							}

							Expect(spaces).To(Equal(expectedSpaces))
							Expect(err).NotTo(HaveOccurred())
						})
					})
				})
			})
		})
	})
})
