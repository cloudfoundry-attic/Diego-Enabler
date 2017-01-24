package thingdoer_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer/thingdoerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DiegoApps", func() {
	var (
		command                thingdoer.AppsGetter
		fakePaginatedRequester *thingdoerfakes.FakePaginatedRequester
		fakeApplicationsParser *thingdoerfakes.FakeApplicationsParser
		apps                   models.Applications
		err                    error
	)

	BeforeEach(func() {
		command = thingdoer.AppsGetter{}
		fakePaginatedRequester = new(thingdoerfakes.FakePaginatedRequester)
		fakeApplicationsParser = new(thingdoerfakes.FakeApplicationsParser)
	})

	JustBeforeEach(func() {
		apps, err = command.DiegoApps(fakeApplicationsParser, fakePaginatedRequester)
	})

	It("should create a request with diego filter set to true", func() {
		expectedFilters := api.Filters{
			api.EqualFilter{
				Name:  "diego",
				Value: true,
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
					Value: true,
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
					Value: true,
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
						Diego: true,
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

			It("returns a list of diego applications", func() {
				expectedApps := models.Applications{
					models.Application{
						models.ApplicationEntity{
							Diego: true,
						},
						models.ApplicationMetadata{
							Guid: "some-guid",
						},
					},
					models.Application{
						models.ApplicationEntity{
							Diego: true,
						},
						models.ApplicationMetadata{
							Guid: "some-guid",
						},
					},
				}

				Expect(apps).To(Equal(expectedApps))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
