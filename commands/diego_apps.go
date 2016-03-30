package commands

import (
	"github.com/cloudfoundry-incubator/diego-enabler/commands/internal/diegohelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/internal/listhelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
)

type DiegoAppsCommand struct {
	Organization string `short:"o" long:"organization" value-name:"ORG" description:"Organization to restrict the app migration to"`
	Space        string `short:"s" long:"space" value-name:"SPACE" description:"Space in the targeted organization to limit results to"`
}

func (command DiegoAppsCommand) Execute([]string) error {
	cliConnection := DiegoEnabler.CLIConnection
	runtime := ui.Diego

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
