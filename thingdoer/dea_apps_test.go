package thingdoer_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/api/apifakes"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer/thingdoerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeaApps", func() {
	var (
		command                thingdoer.AppsGetter
		fakeConnection         *apifakes.FakeConnection
		fakeApplicationsParser *thingdoerfakes.FakeApplicationsParser
		fakePaginatedRequester *thingdoerfakes.FakePaginatedRequester
		apps                   models.Applications
		err                    error
	)

	BeforeEach(func() {
		fakeConnection = new(apifakes.FakeConnection)
		command = thingdoer.AppsGetter{CliConnection: fakeConnection}
		fakePaginatedRequester = new(thingdoerfakes.FakePaginatedRequester)
		fakeApplicationsParser = new(thingdoerfakes.FakeApplicationsParser)
	})

	Context("DeaApps", func() {
		JustBeforeEach(func() {
			apps, err = command.DeaApps(fakeApplicationsParser, fakePaginatedRequester)
		})

		It("should create a request with diego filter set to false", func() {
			expectedFilters := api.Filters{
				api.EqualFilter{
					Name:  "diego",
					Value: false,
				},
			}

			Expect(fakePaginatedRequester.DoCallCount()).To(Equal(1))
			filters, _ := fakePaginatedRequester.DoArgsForCall(0)
			Expect(filters).To(Equal(expectedFilters))
		})

		Context("when an organization name is specified", func() {
			BeforeEach(func() {
				command.OrganizationGuid = "some-organization-guid"
			})

			It("should create a request with organization guid set", func() {
				expectedFilters := api.Filters{
					api.EqualFilter{
						Name:  "diego",
						Value: false,
					},
					api.EqualFilter{
						Name:  "organization_guid",
						Value: "some-organization-guid",
					},
				}

				Expect(fakePaginatedRequester.DoCallCount()).To(Equal(1))
				filters, _ := fakePaginatedRequester.DoArgsForCall(0)
				Expect(filters).To(Equal(expectedFilters))
			})
		})

		Context("when an space name is specified", func() {
			BeforeEach(func() {
				command.SpaceGuid = "some-space-guid"
			})

			It("should create a request with space guid set", func() {
				expectedFilters := api.Filters{
					api.EqualFilter{
						Name:  "diego",
						Value: false,
					},
					api.EqualFilter{
						Name:  "space_guid",
						Value: "some-space-guid",
					},
				}

				Expect(fakePaginatedRequester.DoCallCount()).To(Equal(1))
				filters, _ := fakePaginatedRequester.DoArgsForCall(0)
				Expect(filters).To(Equal(expectedFilters))
			})
		})

		Context("when the paginated requester fails", func() {
			var requestError error

			BeforeEach(func() {
				requestError = errors.New("making API requests failed")
				fakePaginatedRequester.DoReturns([][]byte{}, requestError)
			})

			It("returns the requester error", func() {
				Expect(apps).To(BeEmpty())
				Expect(err).To(Equal(requestError))
			})
		})

		Context("When the paginated requester succeeds", func() {
			BeforeEach(func() {
				responseBodies := [][]byte{
					[]byte("some-json"),
					[]byte("some-other-json"),
				}
				fakePaginatedRequester.DoReturns(responseBodies, nil)
			})

			Context("when the parsing fails", func() {
				var apps models.Applications
				var parseError error

				BeforeEach(func() {
					parseError = errors.New("parsing json failed")
					fakeApplicationsParser.ParseReturns(apps, parseError)
				})

				It("returns the parse error", func() {
					Expect(apps).To(BeEmpty())
					Expect(err).To(Equal(parseError))
				})
			})

			Context("when the parsing succeeds", func() {
				var parsedApps models.Applications = models.Applications{
					models.Application{
						models.ApplicationEntity{
							Name:  "app-1",
							Diego: false,
						},
						models.ApplicationMetadata{
							Guid: "some-guid",
						},
					},
				}

				BeforeEach(func() {
					// for each call of Parse
					fakeApplicationsParser.ParseReturns(parsedApps, nil)
				})

				Context("when no errors are encountered getting routes for apps", func() {
					BeforeEach(func() {
						fakeConnection.CliCommandWithoutTerminalOutputReturns(
							[]string{
								"{",
								`"total_results": 15`,
								"}",
							},
							nil,
						)
					})

					It("returns a list of diego applications", func() {
						expectedApps := models.Applications{
							models.Application{
								models.ApplicationEntity{
									Name:      "app-1",
									Diego:     false,
									HasRoutes: true,
								},
								models.ApplicationMetadata{
									Guid: "some-guid",
								},
							},
							models.Application{
								models.ApplicationEntity{
									Name:      "app-1",
									Diego:     false,
									HasRoutes: true,
								},
								models.ApplicationMetadata{
									Guid: "some-guid",
								},
							},
						}

						Expect(err).NotTo(HaveOccurred())
						Expect(apps).To(Equal(expectedApps))
					})
				})

				Context("when getting routes for an application fails", func() {
					BeforeEach(func() {
						fakeConnection.CliCommandWithoutTerminalOutputReturns(
							nil,
							errors.New("getting routes error"),
						)
					})

					It("returns a getting routes error", func() {
						Expect(err).To(MatchError("Unable to get routes for app 'app-1'\ngetting routes error"))
					})
				})
			})
		})
	})

	Describe("ApplicationHasRoutes", func() {
		Context("when the application has routes", func() {
			BeforeEach(func() {
				fakeConnection.CliCommandWithoutTerminalOutputReturns(
					[]string{
						"{",
						`"total_results": 15,`,
						`"total_pages": 1`,
						"}",
					},
					nil,
				)
			})

			It("returns true", func() {
				hasRoutes, err := command.ApplicationHasRoutes("some-app-guid")
				Expect(err).ToNot(HaveOccurred())
				Expect(hasRoutes).To(BeTrue())

				Expect(fakeConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
				Expect(fakeConnection.CliCommandWithoutTerminalOutputArgsForCall(0)).To(Equal(
					[]string{
						"curl",
						"/v2/apps/some-app-guid/routes",
					}))
			})
		})

		Context("when the application does not have routes", func() {
			BeforeEach(func() {
				fakeConnection.CliCommandWithoutTerminalOutputReturns(
					[]string{
						"{",
						`   "total_results": 0,`,
						`   "total_pages": 1`,
						"}",
					},
					nil,
				)
			})

			It("returns false", func() {
				hasRoutes, err := command.ApplicationHasRoutes("some-app-guid")
				Expect(err).ToNot(HaveOccurred())
				Expect(hasRoutes).To(BeFalse())
			})
		})

		Context("when a rpc server error is encountered", func() {
			var returnedErr error

			BeforeEach(func() {
				returnedErr = errors.New("rpc error")
				fakeConnection.CliCommandWithoutTerminalOutputReturns(
					nil,
					returnedErr,
				)
			})

			It("returns the error", func() {
				_, err := command.ApplicationHasRoutes("some-app-guid")
				Expect(err).To(MatchError(returnedErr))
			})
		})

		Context("when a http error is encountered", func() {
			BeforeEach(func() {
				fakeConnection.CliCommandWithoutTerminalOutputReturns(
					[]string{
						"502 Bad Gateway: Registered endpoint failed to handle the request.",
					},
					nil,
				)
			})

			It("returns the error", func() {
				_, err := command.ApplicationHasRoutes("some-app-guid")
				Expect(err).To(MatchError("502 Bad Gateway: Registered endpoint failed to handle the request."))
			})
		})

		Context("when a cloud controller error is encountered", func() {
			BeforeEach(func() {
				fakeConnection.CliCommandWithoutTerminalOutputReturns(
					[]string{
						"{",
						`   "code": 10000,`,
						`   "description": "Unknown request",`,
						`   "error_code": "CF-NotFound"`,
						"}",
					},
					nil,
				)
			})

			It("returns a cloud controller error", func() {
				_, err := command.ApplicationHasRoutes("some-app-guid")
				Expect(err).To(MatchError(`CC code:       10000
CC error code: CF-NotFound
Description:   Unknown request`))
			})
		})
	})
})
