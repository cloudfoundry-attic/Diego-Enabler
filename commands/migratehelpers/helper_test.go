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
	. "github.com/onsi/gomega/gbytes"
)

func captureStdout(buf *Buffer) *os.File {
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

var _ = Describe("MigrateApps", func() {
	var (
		command MigrateApps
	)

	Describe("MigrateApp", func() {
		Context("when migrating the app fails", func() {
			var (
				success      int
				diegoSupport *migratehelpersfakes.FakeDiegoFlagSetter
				appPrinter   *displayhelpers.AppPrinter
				buf          *Buffer
				stdout       *os.File
			)

			BeforeEach(func() {
				buf = NewBuffer()
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
					Spaces: map[string]models.Space{
						"some-space-guid": models.Space{
							SpaceEntity: models.SpaceEntity{
								Name: "some-space",
								Organization: models.Organization{
									OrganizationEntity: models.OrganizationEntity{
										Name: "some-org",
									},
								},
							},
						},
					},
				}
				command = MigrateApps{
					MaxInFlight:    1,
					Runtime:        ui.Diego,
					AppsGetterFunc: nil,
					MigrateAppsCommand: &ui.MigrateAppsCommand{
						Username:     "some-user",
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
					Expect(success).To(Equal(Warning))
					Eventually(buf).Should(Say("WARNING"))
				})
			})

			Context("for any other reason", func() {
				BeforeEach(func() {
					diegoSupport.SetDiegoFlagReturns(nil, errors.New("disaster"))
				})

				It("returns an error", func() {
					Expect(success).To(Equal(Err))
					Eventually(buf).Should(Say("Error: Failed to migrate app"))
				})
			})

			Context("when the runtime is Diego", func() {
				BeforeEach(func() {
					// error so the test doesn't wait
					diegoSupport.SetDiegoFlagReturns(nil, errors.New("disaster"))
					command.Runtime = ui.Diego
				})

				Context("when the app has zero routes", func() {
					BeforeEach(func() {
						diegoSupport.HasRoutesReturns(false, nil)
					})

					It("displays a no route warning", func() {
						Eventually(buf).Should(Say("WARNING: Assuming health check of type process \\('none'\\) for app with no mapped routes\\. Use 'cf set-health-check' to change this\\. App .+some-app.+ to .+Diego.+ in space .+some-space.+ / org .+some-org.+ as .+some-user.+"))

						Expect(diegoSupport.HasRoutesCallCount()).To(Equal(1))
						Expect(diegoSupport.HasRoutesArgsForCall(0)).To(Equal("some-app"))
					})
				})

				Context("when the app has routes", func() {
					BeforeEach(func() {
						diegoSupport.HasRoutesReturns(true, nil)
					})

					It("does not display a no route warning", func() {
						Consistently(buf).ShouldNot(Say("WARNING: Assuming health check"))
					})
				})

				Context("when checking the app routes returns an error", func() {
					BeforeEach(func() {
						expectedErr := errors.New("warn no routes error")

						diegoSupport.HasRoutesReturns(false, expectedErr)
					})

					It("returns the error", func() {
						Expect(success).To(Equal(Err))
					})
				})
			})

			Context("when the runtime is DEA", func() {
				BeforeEach(func() {
					// error so the test doesn't wait
					diegoSupport.SetDiegoFlagReturns(nil, errors.New("disaster"))
					command.Runtime = ui.DEA
				})

				It("does not check if the app has zero routes", func() {
					Expect(diegoSupport.HasRoutesCallCount()).To(Equal(0))
				})
			})
		})
	})
})
