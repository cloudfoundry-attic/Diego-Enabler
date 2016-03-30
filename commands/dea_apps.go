package commands

import (
	"github.com/cloudfoundry-incubator/diego-enabler/commands/internal/diegohelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/internal/listhelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
)

type DeaAppsCommand struct {
	Organization string `short:"o" long:"organization" value-name:"ORG" description:"Organization to restrict the app migration to"`
}

func (command DeaAppsCommand) Execute([]string) error {
	cliConnection := DiegoEnabler.CLIConnection
	runtime := ui.DEA

	appsGetter, err := diegohelpers.NewAppsGetterFunc(cliConnection, command.Organization, runtime)
	if err != nil {
		return err
	}

	listAppsCommand, err := listhelpers.NewListAppsCommand(cliConnection, command.Organization, runtime)
	if err != nil {
		return err
	}

	err = listhelpers.ListApps(cliConnection, appsGetter, &listAppsCommand)
	if err != nil {
		return err
	}
	return nil
}
