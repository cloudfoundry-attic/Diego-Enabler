package migratehelpers_test

import (
	"errors"
	"io"
	"os"

	"github.com/cloudfoundry-incubator/diego-enabler/commands/displayhelpers"
	. "github.com/cloudfoundry-incubator/diego-enabler/commands/migratehelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/migratehelpers/migratehelpersfakes"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("MigrateApps", func() {
	var (
		command MigrateApps
	)

	Describe("MigrateApp", func() {
		Context("when migrating the app fails", func() {
			var (
				success      bool
				diegoSupport *migratehelpersfakes.FakeDiegoFlagSetter
				appPrinter   *displayhelpers.AppPrinter
				buf          *gbytes.Buffer
				stdout       *os.File
			)

			BeforeEach(func() {
				buf = gbytes.NewBuffer()
				stdout = captureStdout(buf)

				diegoSupport = new(migratehelpersfakes.FakeDiegoFlagSetter)
				appPrinter = &displayhelpers.AppPrinter{
					App: models.Application{
						ApplicationEntity: models.ApplicationEntity{
							Name:      "some-app",
							Diego:     true,
							State:     "STARTED",
							SpaceGuid: "some-space-guid",
						},
						ApplicationMetadata: models.ApplicationMetadata{
							Guid: "some-app-guid",
						},
					},
					Spaces: map[string]models.Space{},
				}
				command = MigrateApps{
					MaxInFlight:    1,
					Runtime:        ui.Diego,
					AppsGetterFunc: nil,
					MigrateAppsCommand: &ui.MigrateAppsCommand{
						Username:     "some-username",
						Runtime:      ui.Diego,
						Organization: "some-organization",
						Space:        "some-space",
					},
				}
			})

			AfterEach(func() {
				os.Stdout.Close()
				os.Stdout = stdout
			})

			JustBeforeEach(func() {
				success = command.MigrateApp(appPrinter, diegoSupport)
			})

			Context("when the user does not have permissions to migrate apps", func() {
				BeforeEach(func() {
					diegoSupport.SetDiegoFlagReturns(nil, errors.New("CF-NotAuthorized - You are not authorized to perform the requested action"))
				})

				It("returns a warning", func() {
					Expect(success).To(BeFalse())
					Eventually(buf).Should(gbytes.Say("WARNING"))
				})
			})

			Context("for any other reason", func() {
				BeforeEach(func() {
					diegoSupport.SetDiegoFlagReturns(nil, errors.New("disaster"))
				})

				It("returns an error", func() {
					Expect(success).To(BeFalse())
					Eventually(buf).Should(gbytes.Say("Error"))
				})
			})
		})
	})
})

func captureStdout(buf *gbytes.Buffer) *os.File {
	stdout := os.Stdout
	r, w, err := os.Pipe()
	Expect(err).NotTo(HaveOccurred())
	os.Stdout = w
	go func() {
		_, err = io.Copy(buf, r)
		buf.Close()
		r.Close()
	}()
	Expect(err).NotTo(HaveOccurred())
	return stdout
}
