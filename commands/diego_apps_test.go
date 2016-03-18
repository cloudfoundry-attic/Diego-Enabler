package commands_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/diego-enabler/commands"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/fakes"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"strings"
)

var _ = Describe("DiegoApps", func() {
	var fakeCliCon *fakes.FakeCliConnection
	var fakeRequestFactory *fakes.FakeRequestFactory
	var fakeCloudControllerClient *fakes.FakeCloudControllerClient
	var fakeResponseParser *fakes.FakeResponseParser
	var apps models.Applications

	var err error

	BeforeEach(func() {
		fakeCliCon = new(fakes.FakeCliConnection)
		fakeRequestFactory = new(fakes.FakeRequestFactory)
		fakeCloudControllerClient = new(fakes.FakeCloudControllerClient)
		fakeResponseParser = new(fakes.FakeResponseParser)
	})

	JustBeforeEach(func() {
		apps, err = commands.DiegoApps(fakeCliCon, fakeRequestFactory, fakeCloudControllerClient, fakeResponseParser)
	})

	Context("when logged in", func() {
		var testRequest *http.Request

		BeforeEach(func() {
			testRequest, err = http.NewRequest("GET", "something", strings.NewReader(""))
			Expect(err).NotTo(HaveOccurred())

			fakeCliCon.IsLoggedInReturns(true, nil)
			fakeCliCon.AccessTokenReturns("", nil)
			fakeRequestFactory.NewGetAppsRequestReturns(testRequest, nil)
		})

		It("tries to get an access token", func() {
			Expect(fakeCliCon.AccessTokenCallCount()).To(Equal(1))
		})

		Context("when getting the access token fails", func() {
			accessTokenErr := errors.New("failed to get access token")

			BeforeEach(func() {
				fakeCliCon.AccessTokenReturns("", accessTokenErr)
			})

			It("should return the error", func() {
				Expect(apps).To(BeEmpty())
				Expect(err).To(Equal(accessTokenErr))
			})
		})

		Context("when getting the access token succeeds", func() {
			testAccessToken := "testAuthToken"

			BeforeEach(func() {
				fakeCliCon.AccessTokenReturns(testAccessToken, nil)
			})

			It("should create a request", func() {
				expectedParams := map[string]interface{}{
					"diego": true,
				}

				Expect(fakeRequestFactory.NewGetAppsRequestCallCount()).To(Equal(1))
				Expect(fakeRequestFactory.NewGetAppsRequestArgsForCall(0)).To(Equal(expectedParams))
			})

			Context("when creating the request fails", func() {
				var disaster = errors.New("OH NOOOOOOO")
				BeforeEach(func() {
					fakeRequestFactory.NewGetAppsRequestReturns(new(http.Request), disaster)
				})

				It("should return the error", func() {
					Expect(apps).To(BeEmpty())
					Expect(err).To(Equal(disaster))
				})
			})

			Context("when creating the request succeeds", func() {
				BeforeEach(func() {
					fakeRequestFactory.NewGetAppsRequestReturns(testRequest, nil)
				})

				It("should make a request with the auth token as an Authorization header", func() {
					Expect(fakeCloudControllerClient.DoCallCount()).To(Equal(1))

					expectedRequest, err := http.NewRequest("GET", "something", strings.NewReader(""))
					Expect(err).NotTo(HaveOccurred())
					expectedRequest.Header.Set("Authorization", testAccessToken)

					Expect(fakeCloudControllerClient.DoArgsForCall(0)).To(Equal(expectedRequest))
				})

				Context("when making the request fails", func() {
					var requestError error

					BeforeEach(func() {
						requestError = errors.New("request execution failed")
						fakeCloudControllerClient.DoReturns(new(http.Response), requestError)
					})

					It("should return the request error", func() {
						Expect(apps).To(BeEmpty())
						Expect(err).To(Equal(requestError))
					})
				})

				Context("when making the request succeeds", func() {
					response := new(http.Response)

					BeforeEach(func() {
						fakeCloudControllerClient.DoReturns(response, nil)
					})

					It("should parse the response", func() {
						Expect(fakeResponseParser.ParseCallCount()).To(Equal(1))
						Expect(fakeResponseParser.ParseArgsForCall(0)).To(Equal(response))
					})

					Context("when the parsing fails", func() {
						var apps models.Applications
						var parseError error

						BeforeEach(func() {
							parseError = errors.New("parsing json failed")
							fakeResponseParser.ParseReturns(apps, parseError)
						})

						It("returns the parse error", func() {
							Expect(apps).To(BeEmpty())
							Expect(err).To(Equal(parseError))
						})
					})

					Context("when the parsing succeeds", func() {
						var parsedApps models.Applications = []models.Application{
							models.Application{Diego: true},
						}

						BeforeEach(func() {
							fakeResponseParser.ParseReturns(parsedApps, nil)
						})

						It("returns a list of diego applications", func() {
							Expect(apps).To(Equal(parsedApps))
							Expect(err).NotTo(HaveOccurred())
						})
					})
				})
			})
		})
	})

	Context("when not logged in", func() {
		BeforeEach(func() {
			fakeCliCon.IsLoggedInReturns(false, nil)
		})

		It("should error", func() {
			Expect(err).To(Equal(commands.NotLoggedInError))
		})
	})

	Context("when the CliConnection fails", func() {
		BeforeEach(func() {
			fakeCliCon.IsLoggedInReturns(false, errors.New("horrible things"))
		})

		It("returns the error", func() {
			Expect(err).To(Equal(errors.New("horrible things")))
		})
	})
})
