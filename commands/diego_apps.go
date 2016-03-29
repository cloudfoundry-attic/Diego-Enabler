package commands

import (
	"github.com/cloudfoundry-incubator/diego-enabler/commands/internal/listhelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
)

type DiegoAppsCommand struct {
	Organization string `short:"o" long:"organization" value-name:"ORG" description:"Organization to restrict the app migration to"`
}

func (command DiegoAppsCommand) Execute([]string) error {
	cliConnection := DiegoEnabler.CLIConnection
	runtime := ui.Runtime(ui.Diego)

	opts := listhelpers.ListAppOpts{
		Organization: command.Organization,
	}

	appsGetter, err := listhelpers.NewAppsGetterFunc(cliConnection, opts.Organization, runtime)
	if err != nil {
		return err
	}

	listAppsCommand, err := listhelpers.NewListAppsCommand(cliConnection, opts.Organization, runtime)
	if err != nil {
		return err
	}

	cmd := listhelpers.ListApps{
		Opts:            opts,
		Runtime:         runtime,
		AppsGetterFunc:  appsGetter,
		ListAppsCommand: &listAppsCommand,
	}

	err = cmd.Execute(cliConnection)
	if err != nil {
		return err
	}
	return nil
}
