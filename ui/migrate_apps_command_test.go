package ui_test

import (
	. "github.com/cloudfoundry-incubator/diego-enabler/ui"
	"github.com/cloudfoundry-incubator/diego-enabler/ui/uifakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("MigrateAppsCommand", func() {
	var (
		command    MigrateAppsCommand
		appPrinter *uifakes.FakeApplicationPrinter
		output     *Buffer
	)

	BeforeEach(func() {
		command = MigrateAppsCommand{
			Username: "some-user",
		}

		appPrinter = new(uifakes.FakeApplicationPrinter)
		appPrinter.OrganizationReturns("some-org")
		appPrinter.SpaceReturns("some-space")
		appPrinter.NameReturns("some-app")

		output = NewBuffer()
	})

	Describe("HealthCheckNoneWarning", func() {
		It("writes a warning to the output", func() {
			command.HealthCheckNoneWarning(appPrinter, output)

			Eventually(output).Should(Say("WARNING: Assuming health check of type process \\('none'\\) for app with no mapped routes\\. Use 'cf set-health-check' to change this\\. App .+some-app.+ to .+Diego.+ in space .+some-space.+ / org .+some-org.+ as .+some-user.+"))
		})
	})
})
