package diegosupport_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/diego-enabler/diegosupport"
	"github.com/cloudfoundry-incubator/diego-enabler/diegosupport/diegosupportfakes"
	"github.com/cloudfoundry/cli/plugin/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("DiegoSupport", func() {
	var (
		fakeCliConnection *diegosupportfakes.FakeCliConnection
		diegoSupport      *diegosupport.DiegoSupport
	)

	BeforeEach(func() {
		fakeCliConnection = &diegosupportfakes.FakeCliConnection{}
		fakeCliConnection.CliCommandWithoutTerminalOutputReturns([]string{""}, nil)
		diegoSupport = diegosupport.NewDiegoSupport(fakeCliConnection)
	})

	Describe("SetDiegoFlag", func() {
		Context("when constructing the api call", func() {
			It("invokes CliCommandWithoutTerminalOutput()", func() {
				diegoSupport.SetDiegoFlag("123", false)

				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			})

			It("calls cli core command 'curl'", func() {
				diegoSupport.SetDiegoFlag("123", false)

				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[0]).To(Equal("curl"))
			})

			It("hits the /v2/apps endpoint", func() {
				diegoSupport.SetDiegoFlag("test-app-guid", false)

				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[1]).To(Equal("/v2/apps/test-app-guid"))
			})

			It("uses the 'PUT' method", func() {
				diegoSupport.SetDiegoFlag("test-app-guid", false)

				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[2]).To(Equal("-X"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[3]).To(Equal("PUT"))
			})

			It("includes http data in the body to set diego flag", func() {
				diegoSupport.SetDiegoFlag("test-app-guid", true)

				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[4]).To(Equal("-d"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[5]).To(Equal(`{"diego":true}`))
			})
		})

		Context("when we do not encounter an error from the API call", func() {
			It("is able to set the diego flag", func() {
				apiOutput := []string{"{", `"key":"value"`, "}"}
				fakeCliConnection.CliCommandWithoutTerminalOutputReturns(apiOutput, nil)

				output, err := diegoSupport.SetDiegoFlag("test-app-guid", false)
				Expect(output).To(Equal(apiOutput))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when we encounter an error from the API call", func() {
			It("returns the output and the error from 'curl'", func() {
				fakeCliConnection.CliCommandWithoutTerminalOutputReturns([]string{"This is the fake output from curl", "some other content"}, errors.New("error from curl"))

				output, err := diegoSupport.SetDiegoFlag("test-app-guid", false)
				Expect(output[0]).To(Equal("This is the fake output from curl"))
				Expect(output[1]).To(Equal("some other content"))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error from curl"))
			})

			It("parses the output and returns any diego specific error", func() {
				response := []string{`{"code": 10000,	"description": "diego not supported",	"error_code": "12345"}`}

				fakeCliConnection.CliCommandWithoutTerminalOutputReturns(response, nil)

				output, err := diegoSupport.SetDiegoFlag("test-app-guid", false)
				Expect(output[0]).To(ContainSubstring(`"code": 10000`))
				Expect(output[0]).To(ContainSubstring("diego not supported"))
				Expect(output[0]).To(ContainSubstring("12345"))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("12345 - diego not supported"))
			})
		})
	})

	Describe("HasRoutes", func() {
		Context("when the app has no routes", func() {
			BeforeEach(func() {
				fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{}, nil)
			})

			It("returns false and no error", func() {
				hasRoutes, err := diegoSupport.HasRoutes("some-app")
				Expect(hasRoutes).To(BeFalse())
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeCliConnection.GetAppCallCount()).To(Equal(1))
				Expect(fakeCliConnection.GetAppArgsForCall(0)).To(Equal("some-app"))
			})
		})

		Context("when the app has routes", func() {
			BeforeEach(func() {
				fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{
					Routes: []plugin_models.GetApp_RouteSummary{
						{},
					},
				}, nil)
			})

			It("returns true and no error", func() {
				hasRoutes, err := diegoSupport.HasRoutes("some-app")
				Expect(hasRoutes).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when an error is encountered", func() {
			var expectedError error

			BeforeEach(func() {
				expectedError = errors.New("some-error")
				fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{}, expectedError)
			})

			It("returns the error", func() {
				_, err := diegoSupport.HasRoutes("some-app")
				Expect(err).To(MatchError(expectedError))
			})
		})
	})

	Describe("WarnNoRoutes", func() {
		var (
			output *gbytes.Buffer
		)

		BeforeEach(func() {
			output = gbytes.NewBuffer()
		})

		Context("when the app has no routes", func() {
			BeforeEach(func() {
				fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{}, nil)
				fakeCliConnection.GetCurrentSpaceReturns(plugin_models.Space{
					SpaceFields: plugin_models.SpaceFields{
						Name: "some-space",
					},
				}, nil)
				fakeCliConnection.GetSpaceReturns(plugin_models.GetSpace_Model{
					GetSpaces_Model: plugin_models.GetSpaces_Model{
						Name: "some-space",
					},
					Organization: plugin_models.GetSpace_Orgs{
						Name: "some-org",
					},
				}, nil)

				fakeCliConnection.UsernameReturns("some-user", nil)
			})

			It("writes a warning to the output", func() {
				err := diegoSupport.WarnNoRoutes("some-app", output)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeCliConnection.GetAppCallCount()).To(Equal(1))
				Expect(fakeCliConnection.GetAppArgsForCall(0)).To(Equal("some-app"))

				Expect(fakeCliConnection.GetCurrentSpaceCallCount()).To(Equal(1))

				Expect(fakeCliConnection.GetSpaceCallCount()).To(Equal(1))
				Expect(fakeCliConnection.GetSpaceArgsForCall(0)).To(Equal("some-space"))

				Expect(fakeCliConnection.UsernameCallCount()).To(Equal(1))

				Expect(output).To(gbytes.Say("WARNING: Assuming health check of type process \\('none'\\) for app with no mapped routes\\. Use 'cf set-health-check' to change this\\. App .+some-app.+ to .+Diego.+ in space .+some-space.+ / org .+some-org.+ as .+some-user.+"))
			})
		})

		Context("when the app has routes", func() {
			BeforeEach(func() {
				fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{
					Routes: []plugin_models.GetApp_RouteSummary{
						{},
					},
				}, nil)
			})

			It("does not write a warning to the output", func() {
				err := diegoSupport.WarnNoRoutes("some-app", output)
				Expect(err).ToNot(HaveOccurred())

				Expect(output.Contents()).To(BeEmpty())
			})
		})

		Context("when getting the app returns an error", func() {
			var expectedError error

			BeforeEach(func() {
				expectedError = errors.New("some-error")
				fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{}, expectedError)
			})

			It("returns the error", func() {
				err := diegoSupport.WarnNoRoutes("some-app", output)
				Expect(err).To(MatchError(expectedError))
			})
		})

		Context("when getting the current space returns an error", func() {
			var expectedError error

			BeforeEach(func() {
				expectedError = errors.New("some-error")
				fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{}, nil)
				fakeCliConnection.GetCurrentSpaceReturns(plugin_models.Space{}, expectedError)
			})

			It("returns the error", func() {
				err := diegoSupport.WarnNoRoutes("some-app", output)
				Expect(err).To(MatchError(expectedError))
			})
		})

		Context("when getting the space returns an error", func() {
			var expectedError error

			BeforeEach(func() {
				expectedError = errors.New("some-error")
				fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{}, nil)
				fakeCliConnection.GetCurrentSpaceReturns(plugin_models.Space{
					SpaceFields: plugin_models.SpaceFields{
						Name: "some-space",
					},
				}, nil)
				fakeCliConnection.GetSpaceReturns(plugin_models.GetSpace_Model{}, expectedError)
			})

			It("returns the error", func() {
				err := diegoSupport.WarnNoRoutes("some-app", output)
				Expect(err).To(MatchError(expectedError))
			})
		})

		Context("when getting the username returns an error", func() {
			var expectedError error
			BeforeEach(func() {
				expectedError = errors.New("some-error")
				fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{}, nil)
				fakeCliConnection.GetCurrentSpaceReturns(plugin_models.Space{
					SpaceFields: plugin_models.SpaceFields{
						Name: "some-space",
					},
				}, nil)
				fakeCliConnection.GetSpaceReturns(plugin_models.GetSpace_Model{
					GetSpaces_Model: plugin_models.GetSpaces_Model{
						Name: "some-space",
					},
					Organization: plugin_models.GetSpace_Orgs{
						Name: "some-org",
					},
				}, nil)

				fakeCliConnection.UsernameReturns("some-user", expectedError)
			})

			It("returns the error", func() {
				err := diegoSupport.WarnNoRoutes("some-app", output)
				Expect(err).To(MatchError(expectedError))
			})
		})
	})
})
