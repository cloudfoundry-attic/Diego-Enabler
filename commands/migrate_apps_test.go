package commands_test

import (
	. "github.com/cloudfoundry-incubator/diego-enabler/commands"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/errorhelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MigrateAppsCommand", func() {
	var (
		command         MigrateAppsCommand
		requiredOptions MigrateAppsPositionalArgs

		err error
	)

	JustBeforeEach(func() {
		err = command.Execute([]string{})
	})

	Context("when both organization and space are passed", func() {
		BeforeEach(func() {
			requiredOptions = MigrateAppsPositionalArgs{
				Runtime: string(ui.DEA),
			}

			command = MigrateAppsCommand{
				RequiredOptions: requiredOptions,
				Space:           "some-space",
				Organization:    "some-organization",
			}
		})

		It("returns an error", func() {
			Expect(err).To(Equal(errorhelpers.SpecifyOrgOrSpaceError))
		})
	})
})
