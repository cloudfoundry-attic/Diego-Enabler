package diego_support_test

import (
	"errors"

	"github.com/cloudfoundry/cli/plugin/fakes"

	"github.com/simonleung8/diego-enabler/diego_support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DiegoSupport", func() {
	var (
		fakeCliConnection *fakes.FakeCliConnection
		diegoSupport      *diego_support.DiegoSupport
	)

	BeforeEach(func() {
		fakeCliConnection = &fakes.FakeCliConnection{}
		fakeCliConnection.CliCommandWithoutTerminalOutputReturns([]string{""}, nil)
		diegoSupport = diego_support.NewDiegoSupport(fakeCliConnection)
	})

	Describe("SetDiegoFlag", func() {
		It("invokes CliCommandWithoutTerminalOutput()", func() {
			diegoSupport.SetDiegoFlag("123", false)

			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
		})

		It("calls cli core command 'curl'", func() {
			diegoSupport.SetDiegoFlag("123", false)

			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[0]).To(Equal("curl"))
		})

		It("hits the /v2/apps endpoint", func() {
			diegoSupport.SetDiegoFlag("test-app-guid", false)

			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[1]).To(Equal("/v2/apps/test-app-guid"))
		})

		It("uses the 'PUT' method", func() {
			diegoSupport.SetDiegoFlag("test-app-guid", false)

			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[2]).To(Equal("-X"))
			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[3]).To(Equal("PUT"))
		})

		It("includes http data in the body to set diego flag", func() {
			diegoSupport.SetDiegoFlag("test-app-guid", true)

			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[4]).To(Equal("-d"))
			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[5]).To(Equal(`{"diego":true}`))
		})

		It("returns the output and the error from 'curl'", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns([]string{"This is the fake output from curl", "some other content"}, errors.New("error from curl"))

			output, err := diegoSupport.SetDiegoFlag("test-app-guid", false)
			Ω(output[0]).To(Equal("This is the fake output from curl"))
			Ω(output[1]).To(Equal("some other content"))
			Ω(err).To(HaveOccurred())
			Ω(err.Error()).To(Equal("error from curl"))
		})

		It("parses the output and returns any diego specific error", func() {
			response := []string{`{"code": 10000,	"description": "diego not supported",	"error_code": "12345"}`}

			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(response, nil)

			output, err := diegoSupport.SetDiegoFlag("test-app-guid", false)
			Ω(output[0]).To(ContainSubstring(`"code": 10000`))
			Ω(output[0]).To(ContainSubstring("diego not supported"))
			Ω(output[0]).To(ContainSubstring("12345"))
			Ω(err).To(HaveOccurred())
			Ω(err.Error()).To(Equal("12345 - diego not supported"))
		})
	})
})
