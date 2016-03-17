package main_test

import (
	"errors"
	"os/exec"

	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/cloudfoundry/cli/testhelpers/rpc_server"
	fake_rpc_handlers "github.com/cloudfoundry/cli/testhelpers/rpc_server/fakes"

	. "github.com/cloudfoundry-incubator/diego-enabler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("DiegoEnabler", func() {
	var (
		validPluginPath = "./main.exe"
	)

	Describe("GetMetadata()", func() {
		var (
			cmd DiegoEnabler
		)

		BeforeEach(func() {
			cmd = DiegoEnabler{}
		})

		It("contains 3 commands", func() {})
	})

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
				Expect(session).To(gbytes.Say("Invalid Usage"))
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
					Expect(rpcHandlers.GetAppCallCount()).To(Equal(2))
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
					Expect(session).To(gbytes.Say("Ok"))
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
				Expect(session).To(gbytes.Say("Invalid Usage"))
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
					Expect(session).To(gbytes.Say("Invalid Usage"))
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
