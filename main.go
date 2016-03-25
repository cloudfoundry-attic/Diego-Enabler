package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/diego-enabler/commands"
	"github.com/cloudfoundry-incubator/diego-enabler/diego_support"
	"github.com/cloudfoundry/cli/plugin"
)

type DiegoEnabler struct{}

func (c *DiegoEnabler) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Diego-Enabler",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     "enable-diego",
				HelpText: "enable Diego support for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf enable-diego APP_NAME",
				},
			},
			{
				Name:     "disable-diego",
				HelpText: "disable Diego support for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf disable-diego APP_NAME",
				},
			},
			{
				Name:     "has-diego-enabled",
				HelpText: "Check if Diego support is enabled for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf has-diego-enabled APP_NAME",
				},
			},
			{
				Name:     "diego-apps",
				HelpText: "Lists all apps running on the Diego runtime that are visible to the user",
				UsageDetails: plugin.Usage{
					Usage: `cf diego-apps [-o ORG]

OPTIONS:
   -o      Organization to restrict the app migration to`,
				},
			},
			{
				Name:     "dea-apps",
				HelpText: "Lists all apps running on the DEA runtime that are visible to the user",
				UsageDetails: plugin.Usage{
					Usage: `cf dea-apps [-o ORG]

OPTIONS:
   -o      Organization to restrict the app migration to`,
				},
			},
			{
				Name:     "migrate-apps",
				HelpText: "Migrate all apps to Diego/DEA",
				UsageDetails: plugin.Usage{
					Usage: `cf migrate-apps (diego | dea) [-o ORG]

WARNING:
   Migration of a running app causes a restart. Stopped apps will be configured to run on the target runtime but are not started.

OPTIONS:
   -o      Organization to restrict the app migration to
   -p      Maximum number of apps to migrate in parallel (Default: 1, maximum: 100)`,
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(DiegoEnabler))
}

func (c *DiegoEnabler) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "enable-diego" && len(args) == 2 {
		c.toggleDiegoSupport(true, cliConnection, args[1])
	} else if args[0] == "disable-diego" && len(args) == 2 {
		c.toggleDiegoSupport(false, cliConnection, args[1])
	} else if args[0] == "has-diego-enabled" && len(args) == 2 {
		c.isDiegoEnabled(cliConnection, args[1])
	} else if args[0] == "diego-apps" {
		opts := []string{"diego"}
		cmd, err := commands.PrepareListApps(append(opts, args[1:]...), cliConnection)
		if err != nil {
			exitWithError(err, []string{})
		}

		err = cmd.Execute(cliConnection)
		if err != nil {
			exitWithError(err, []string{})
		}
	} else if args[0] == "dea-apps" {
		opts := []string{"dea"}
		cmd, err := commands.PrepareListApps(append(opts, args[1:]...), cliConnection)
		if err != nil {
			exitWithError(err, []string{})
		}

		err = cmd.Execute(cliConnection)
		if err != nil {
			exitWithError(err, []string{})
		}
	} else if args[0] == "migrate-apps" {
		cmd, err := commands.PrepareMigrateApps(args[1:], cliConnection)
		if err != nil {
			exitWithError(err, []string{})
		}
		err = cmd.Execute(cliConnection)
		if err != nil {
			exitWithError(err, []string{})
		}
	} else {
		c.showUsage(args)
	}
}

func (c *DiegoEnabler) showUsage(args []string) {
	for _, cmd := range c.GetMetadata().Commands {
		if cmd.Name == args[0] {
			fmt.Println("Invalid Usage: \n", cmd.UsageDetails.Usage)
		}
	}
}

func (c *DiegoEnabler) toggleDiegoSupport(on bool, cliConnection plugin.CliConnection, appName string) {
	d := diego_support.NewDiegoSupport(cliConnection)

	fmt.Printf("Setting %s Diego support to %t\n", appName, on)
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if output, err := d.SetDiegoFlag(app.Guid, on); err != nil {
		fmt.Println("err 1", err, output)
		exitWithError(err, output)
	}
	sayOk()

	fmt.Printf("Verifying %s Diego support is set to %t\n", appName, on)
	app, err = cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if app.Diego == on {
		sayOk()
	} else {
		sayFailed()
		fmt.Printf("Diego support for %s is NOT set to %t\n\n", appName, on)
		os.Exit(1)
	}
}

func (c *DiegoEnabler) isDiegoEnabled(cliConnection plugin.CliConnection, appName string) {
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if app.Guid == "" {
		sayFailed()
		fmt.Printf("App %s not found\n\n", appName)
		os.Exit(1)
	}

	fmt.Println(app.Diego)
}

func exitWithError(err error, output []string) {
	sayFailed()
	fmt.Println("Error: ", err)
	for _, str := range output {
		fmt.Println(str)
	}
	os.Exit(1)
}

func say(message string, color uint, bold int) string {
	return fmt.Sprintf("\033[%d;%dm%s\033[0m", bold, color, message)
}

func sayOk() {
	fmt.Println(say("Ok\n", 32, 1))
}

func sayFailed() {
	fmt.Println(say("FAILED", 31, 1))
}
