package commands

import (
	"github.com/cloudfoundry-incubator/diego-enabler/commands/diegohelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/errorhelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/flaghelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/migratehelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
)

type MigrateAppsPositionalArgs struct {
	Runtime string `positional-arg-name:"runtime" required:"true" description:"dea or diego"`
}

type MigrateAppsCommand struct {
	RequiredOptions MigrateAppsPositionalArgs `positional-args:"yes"`
	Organization    string                    `short:"o" value-name:"ORG" description:"Organization to restrict the app migration to"`
	Space           string                    `short:"s" value-name:"SPACE" description:"Space in the targeted organization to restrict the app migration to"`
	MaxInFlight     flaghelpers.ParallelFlag  `short:"p" value-name:"MAX_IN_FLIGHT" default:"1" description:"Maximum number of apps to migrate in parallel (maximum: 100)"`
}

//TODO: Figure out how to output this warning in the help
//WARNING:
//   Migration of a running app causes a restart. Stopped apps will be configured to run on the target runtime but are not started.

func (command MigrateAppsCommand) Execute([]string) error {
	cliConnection := DiegoEnabler.CLIConnection

	runtime, err := ui.ParseRuntime(command.RequiredOptions.Runtime)
	if err != nil {
		return err
	}

	err = errorhelpers.ErrorIfOrgAndSpacesSet(command.Organization, command.Space)
	if err != nil {
		return err
	}

	appsGetter, err := diegohelpers.NewAppsGetterFunc(cliConnection, command.Organization, command.Space, runtime.Flip())
	if err != nil {
		return err
	}

	migrateAppsCommand, err := migratehelpers.NewMigrateAppsCommand(cliConnection, command.Organization, command.Space, runtime)
	if err != nil {
		return err
	}

	cmd := migratehelpers.MigrateApps{
		MaxInFlight:        command.MaxInFlight.Value,
		Runtime:            runtime,
		AppsGetterFunc:     appsGetter,
		MigrateAppsCommand: &migrateAppsCommand,
	}

	return cmd.Execute(cliConnection)
}
