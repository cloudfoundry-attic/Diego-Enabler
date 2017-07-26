package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/diego-enabler/commands"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/jessevdk/go-flags"
)

type DiegoEnabler struct{}

func (c *DiegoEnabler) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Diego-Enabler",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 2,
			Build: 4,
		},
		Commands: []plugin.Command{
			{
				Name:     "enable-diego",
				HelpText: "Migrate app to the Diego runtime",
				UsageDetails: plugin.Usage{
					Usage: `cf enable-diego APP_NAME

WARNING:
   Migration of a running app causes a restart. Stopped apps will be configured to run on the target runtime but are not started.`,
				},
			},
			{
				Name:     "disable-diego",
				HelpText: "Migrate app to the DEA runtime",
				UsageDetails: plugin.Usage{
					Usage: `cf disable-diego APP_NAME

WARNING:
   Migration of a running app causes a restart. Stopped apps will be configured to run on the target runtime but are not started.`,
				},
			},
			{
				Name:     "has-diego-enabled",
				HelpText: "Report whether an app is configured to run on the Diego runtime",
				UsageDetails: plugin.Usage{
					Usage: "cf has-diego-enabled APP_NAME",
				},
			},
			{
				Name:     "diego-apps",
				HelpText: "Lists all apps running on the Diego runtime that are visible to the user",
				UsageDetails: plugin.Usage{
					Usage: `cf diego-apps [-o ORG | -s SPACE]

OPTIONS:
   -o      Organization to restrict the app migration to,
   -s      Space in the targeted organization to limit results to`,
				},
			},
			{
				Name:     "dea-apps",
				HelpText: "Lists all apps running on the DEA runtime that are visible to the user",
				UsageDetails: plugin.Usage{
					Usage: `cf dea-apps [-o ORG | -s SPACE]

OPTIONS:
   -o      Organization to restrict the app migration to,
   -s      Space in the targeted organization to limit results to`,
				},
			},
			{
				Name:     "migrate-apps",
				HelpText: "Migrate all apps to Diego/DEA",
				UsageDetails: plugin.Usage{
					Usage: `cf migrate-apps (diego | dea) [-o ORG | -s SPACE] [-p MAX_IN_FLIGHT]

WARNING:
   Migration of a running app causes a restart. Stopped apps will be configured to run on the target runtime but are not started.

OPTIONS:
   -o      Organization to restrict the app migration to
   -s      Space in the targeted organization to restrict the app migration to
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
	commands.DiegoEnabler.CLIConnection = cliConnection
	parser := flags.NewParser(&commands.DiegoEnabler, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	_, err := parser.ParseArgs(args)
	if err != nil {
		ui.SayFailed()
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
