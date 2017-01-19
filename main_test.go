package main_test

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/cloudfoundry/cli/testhelpers/rpc_server"
	fake_rpc_handlers "github.com/cloudfoundry/cli/testhelpers/rpc_server/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("DiegoEnabler", func() {
	Describe("Commands", func() {
		var (
			rpcHandlers *fake_rpc_handlers.FakeHandlers
			ts          *test_rpc_server.TestServer
			err         error
		)

		BeforeEach(func() {
			rpcHandlers = &fake_rpc_handlers.FakeHandlers{}
		})

		JustBeforeEach(func() {
			//set rpc.CallCoreCommand to a successful call
			rpcHandlers.CallCoreCommandStub = func(_ []string, retVal *bool) error {
				*retVal = true
				return nil
			}

			//set rpc.GetOutputAndReset to return empty string; this is used by CliCommand()/CliWithoutTerminalOutput()
			rpcHandlers.GetOutputAndResetStub = func(_ bool, retVal *[]string) error {
				*retVal = []string{"{}"}
				return nil
			}

			ts, err = test_rpc_server.NewTestRpcServer(rpcHandlers)
			Expect(err).NotTo(HaveOccurred())

			err = ts.Start()
			Expect(err).NotTo(HaveOccurred())

		})

		AfterEach(func() {
			ts.Stop()
		})

		Context("enable-diego", func() {
			var args []string

			JustBeforeEach(func() {
				args = []string{ts.Port(), "enable-diego", "test-app"}
			})

			It("needs APP_NAME as argument", func() {
				args = []string{ts.Port(), "enable-diego"}
				session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				session.Wait()
				Expect(session).To(gbytes.Say("required argument `APP_NAME` was not provided"))
			})

			Context("when the args are properly provided", func() {
				BeforeEach(func() {
					rpcHandlers.GetAppStub = func(_ string, retVal *plugin_models.GetAppModel) error {
						*retVal = plugin_models.GetAppModel{}
						return nil
					}
				})

				It("calls GetApp() twice, one to get app guid, another to verify flag is set", func() {
					session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					session.Wait()
					Expect(rpcHandlers.GetAppCallCount()).To(Equal(3))
				})
			})

			Context("when the app is found", func() {
				BeforeEach(func() {
					rpcHandlers.GetAppStub = func(_ string, retVal *plugin_models.GetAppModel) error {
						*retVal = plugin_models.GetAppModel{Guid: "test-app-guid"}
						return nil
					}
				})

				It("sets diego flag with /v2/apps endpoint", func() {
					session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					session.Wait()
					Expect(rpcHandlers.CallCoreCommandCallCount()).To(Equal(1))

					output, _ := rpcHandlers.CallCoreCommandArgsForCall(0)
					Expect(output[1]).To(ContainSubstring("v2/apps/test-app-guid"))
					Expect(output[5]).To(ContainSubstring(`"diego":true`))
				})
			})

			Context("when the app is not found", func() {
				BeforeEach(func() {
					rpcHandlers.GetAppStub = func(_ string, retVal *plugin_models.GetAppModel) error {
						*retVal = plugin_models.GetAppModel{Guid: "test-app-guid"}
						return errors.New("error in GetApp")
					}
				})

				It("exits with error when GetApp() returns error", func() {
					session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					session.Wait()

					Expect(session).To(gbytes.Say("error in GetApp"))
					Expect(session.ExitCode()).To(Equal(1))
				})
			})

			Context("when the app was successfully changed to Diego", func() {
				BeforeEach(func() {
					rpcHandlers.GetAppStub = func(_ string, retVal *plugin_models.GetAppModel) error {
						*retVal = plugin_models.GetAppModel{Guid: "test-app-guid", Diego: true}
						return nil
					}
				})

				It("exit 0 after veriftying the flag is correct set", func() {
					session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					session.Wait()

					Expect(session).To(gbytes.Say("Verifying test-app Diego support is set to true"))
					Expect(session).To(gbytes.Say("OK"))
					Expect(session.ExitCode()).To(Equal(0))
				})
			})

			Context("when the change to Diego failed", func() {
				BeforeEach(func() {
					rpcHandlers.GetAppStub = func(_ string, retVal *plugin_models.GetAppModel) error {
						*retVal = plugin_models.GetAppModel{Guid: "test-app-guid", Diego: false}
						return nil
					}
				})

				It("exit 1 after veriftying the flag is not correct set", func() {
					session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					session.Wait()

					Expect(session).To(gbytes.Say("Verifying test-app Diego support is set to true"))
					Expect(session).To(gbytes.Say("FAILED"))
					Expect(session.ExitCode()).To(Equal(1))
				})
			})

			Context("when the app has no routes", func() {
				BeforeEach(func() {
					rpcHandlers.GetAppStub = func(appName string, retVal *plugin_models.GetAppModel) error {
						Expect(appName).To(Equal("test-app"))

						*retVal = plugin_models.GetAppModel{
							Guid: "test-app-guid",
							Name: "test-app",
						}
						return nil
					}

					rpcHandlers.GetCurrentSpaceStub = func(_ string, retVal *plugin_models.Space) error {
						*retVal = plugin_models.Space{
							SpaceFields: plugin_models.SpaceFields{
								Name: "some-space",
							},
						}
						return nil
					}

					rpcHandlers.GetSpaceStub = func(spaceName string, retVal *plugin_models.GetSpace_Model) error {
						if spaceName != "some-space" {
							return fmt.Errorf("expected spaceName argument to be %s, was %s", "some-space", spaceName)
						}
						*retVal = plugin_models.GetSpace_Model{
							GetSpaces_Model: plugin_models.GetSpaces_Model{
								Name: "some-space",
							},
							Organization: plugin_models.GetSpace_Orgs{Name: "some-org"},
						}
						return nil
					}

					rpcHandlers.UsernameStub = func(_ string, retVal *string) error {
						*retVal = "some-user"
						return nil
					}
				})

				It("warns the user that the health check type will be set to none", func() {
					session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					session.Wait()
					// This output is colored but this assertion does not test for color. This is because the terminal package disables color when the binary is run in a non-tty session.
					Expect(session.Out).To(gbytes.Say("WARNING: Assuming health check of type process \\('none'\\) for app with no mapped routes\\. Use 'cf set-health-check' to change this\\. App test-app to Diego in space some-space / org some-org as some-user"))
				})
			})
		})

		Context("disable-diego", func() {
			var args []string

			JustBeforeEach(func() {
				args = []string{ts.Port(), "disable-diego", "test-app"}
			})

			It("needs APP_NAME as argument", func() {
				args = []string{ts.Port(), "disable-diego"}
				session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				session.Wait()
				Expect(session).To(gbytes.Say("required argument `APP_NAME` was not provided"))
			})

			Context("when the app is found", func() {
				BeforeEach(func() {
					rpcHandlers.GetAppStub = func(_ string, retVal *plugin_models.GetAppModel) error {
						*retVal = plugin_models.GetAppModel{Guid: "test-app-guid"}
						return nil
					}
				})

				It("sets diego flag with /v2/apps endpoint", func() {
					session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					session.Wait()
					Expect(rpcHandlers.CallCoreCommandCallCount()).To(Equal(1))

					output, _ := rpcHandlers.CallCoreCommandArgsForCall(0)
					Expect(output[1]).To(ContainSubstring("v2/apps/test-app-guid"))
					Expect(output[5]).To(ContainSubstring(`"diego":false`))
				})
			})

			Context("has-diego-enabled", func() {
				var args []string

				JustBeforeEach(func() {
					args = []string{ts.Port(), "has-diego-enabled", "test-app"}
				})

				It("needs APP_NAME as argument", func() {
					args = []string{ts.Port(), "has-diego-enabled"}
					session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					session.Wait()
					Expect(session).To(gbytes.Say("required argument `APP_NAME` was not provided"))
				})

				Context("when the params are properly provided", func() {
					BeforeEach(func() {
						rpcHandlers.GetAppStub = func(_ string, retVal *plugin_models.GetAppModel) error {
							*retVal = plugin_models.GetAppModel{}
							return nil
						}
					})

					It("calls GetApp() to get app model", func() {

						session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
						Expect(err).NotTo(HaveOccurred())

						session.Wait()
						Expect(rpcHandlers.GetAppCallCount()).To(Equal(1))
					})
				})

				Context("when the app does not exist", func() {
					BeforeEach(func() {
						rpcHandlers.GetAppStub = func(_ string, retVal *plugin_models.GetAppModel) error {
							*retVal = plugin_models.GetAppModel{Guid: ""}
							return nil
						}
					})

					It("notifies user app is not found", func() {
						session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
						Expect(err).NotTo(HaveOccurred())

						session.Wait()
						Expect(session).To(gbytes.Say("App test-app not found"))
						Expect(session.ExitCode()).To(Equal(1))
					})
				})

				Context("when the app is on Diego", func() {
					BeforeEach(func() {
						rpcHandlers.GetAppStub = func(_ string, retVal *plugin_models.GetAppModel) error {
							*retVal = plugin_models.GetAppModel{Guid: "test-app-guid", Diego: true}
							return nil
						}
					})

					It("outputs the app's Diego flag value", func() {
						session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
						Expect(err).NotTo(HaveOccurred())

						session.Wait()
						Expect(session).To(gbytes.Say("true"))
						Expect(session.ExitCode()).To(Equal(0))
					})
				})

			})
		})
	})
})
