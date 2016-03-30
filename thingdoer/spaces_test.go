package thingdoer_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/diego-enabler/models"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer/thingdoerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spaces", func() {
	var (
		fakePaginatedRequester *thingdoerfakes.FakePaginatedRequester
		fakeSpacesParser       *thingdoerfakes.FakeSpacesParser
		spaces                 models.Spaces
		err                    error
	)

	BeforeEach(func() {
		fakePaginatedRequester = new(thingdoerfakes.FakePaginatedRequester)
		fakeSpacesParser = new(thingdoerfakes.FakeSpacesParser)
	})

	JustBeforeEach(func() {
		spaces, err = thingdoer.Spaces(fakeSpacesParser, fakePaginatedRequester)
	})

	It("should create a request inline-relations-depth of 1", func() {
		expectedParams := map[string]interface{}{
			"inline-relations-depth": 1,
		}

		Expect(fakePaginatedRequester.DoCallCount()).To(Equal(1))
		_, params := fakePaginatedRequester.DoArgsForCall(0)
		Expect(params).To(Equal(expectedParams))
	})

	Context("when the paginated requester fails", func() {
		var requestError error

		BeforeEach(func() {
			requestError = errors.New("making API requests failed")
			fakePaginatedRequester.DoReturns([][]byte{}, requestError)
		})

		It("returns the requester error", func() {
			Expect(spaces).To(BeEmpty())
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
			var parsedApps models.Spaces = models.Spaces{
				models.Space{
					models.SpaceEntity{
						Name: "space-foo",
					},
					models.SpaceMetadata{
						Guid: "some-guid",
					},
				},
			}

			BeforeEach(func() {
				// for each call of Parse
				fakeSpacesParser.ParseReturns(parsedApps, nil)
			})

			It("returns a list of diego applications", func() {
				expectedApps := models.Spaces{
					models.Space{
						models.SpaceEntity{
							Name: "space-foo",
						},
						models.SpaceMetadata{
							Guid: "some-guid",
						},
					},
					models.Space{
						models.SpaceEntity{
							Name: "space-foo",
						},
						models.SpaceMetadata{
							Guid: "some-guid",
						},
					},
				}

				Expect(spaces).To(Equal(expectedApps))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
