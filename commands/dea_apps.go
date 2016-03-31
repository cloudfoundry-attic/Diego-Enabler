package commands

import (
	"github.com/cloudfoundry-incubator/diego-enabler/commands/diegohelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/listhelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
)

type DeaAppsCommand struct {
	Organization string `short:"o" value-name:"ORG" description:"Organization to restrict the app migration to"`
	Space        string `short:"s" value-name:"SPACE" description:"Space in the targeted organization to limit results to"`
}

func (command DeaAppsCommand) Execute([]string) error {
	cliConnection := DiegoEnabler.CLIConnection
	runtime := ui.DEA

	appsGetter, err := diegohelpers.NewAppsGetterFunc(cliConnection, command.Organization, command.Space, runtime)
	if err != nil {
		return err
	}

	listAppsCommand, err := listhelpers.NewListAppsCommand(cliConnection, command.Organization, command.Space, runtime)
	if err != nil {
		return err
	}

	err = listhelpers.ListApps(cliConnection, appsGetter, &listAppsCommand)
	if err != nil {
		return err
	}
	return nil
}
