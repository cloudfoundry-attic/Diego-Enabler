package commands_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/commands"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/fakes"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeaApps", func() {
	var (
		fakePaginatedRequester *fakes.FakePaginatedRequester
		fakeApplicationsParser *fakes.FakeApplicationsParser
		apps                   models.Applications
		err                    error
	)

	BeforeEach(func() {
		fakePaginatedRequester = new(fakes.FakePaginatedRequester)
		fakeApplicationsParser = new(fakes.FakeApplicationsParser)
	})

	JustBeforeEach(func() {
		apps, err = commands.DeaApps(fakeApplicationsParser, fakePaginatedRequester)
	})

	It("should create a request with diego filter set to false", func() {
		expectedFilters := api.EqualFilter{
			Name:  "diego",
			Value: false,
		}

		Expect(fakePaginatedRequester.DoCallCount()).To(Equal(1))
		filters, _ := fakePaginatedRequester.DoArgsForCall(0)
		Expect(filters).To(Equal(expectedFilters))
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

			It("returns a list of diego applications", func() {
				expectedApps := models.Applications{
					models.Application{
						models.ApplicationEntity{
							Diego: false,
						},
						models.ApplicationMetadata{
							Guid: "some-guid",
						},
					},
					models.Application{
						models.ApplicationEntity{
							Diego: false,
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
