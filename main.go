package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry-incubator/diego-enabler/diego_support"
)

type DiegoEnabler struct{}

func (c *DiegoEnabler) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Diego-Enabler",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
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

	fmt.Printf("Setting %s Deigo support to %t\n", appName, on)
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if output, err := d.SetDiegoFlag(app.Guid, on); err != nil {
		fmt.Println("err 1", err, output)
		exitWithError(err, output)
	}
	sayOk()

	fmt.Printf("Verifying %s Deigo support is set to %t\n", appName, on)
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
