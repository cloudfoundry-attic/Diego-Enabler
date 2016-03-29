package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/diego-enabler/commands/internal/diegohelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/internal/migratehelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
)

type MigrateAppsPositionalArgs struct {
	Runtime string `positional-arg-name:"runtime" required:"true" description:"dea or diego"`
}

type MigrateAppsCommand struct {
	RequiredOptions MigrateAppsPositionalArgs `positional-args:"yes"`
	Organization    string                    `short:"o" long:"organization" value-name:"ORG" description:"Organization to restrict the app migration to"`
	MaxInFlight     int                       `short:"p" long:"parallel" value-name:"MAX_IN_FLIGHT" default:"1" description:"Maximum number of apps to migrate in parallel (maximum: 100)"`
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

	opts := migratehelpers.MigrateAppsOpts{
		Organization: command.Organization,
		//TODO: do not pass in a string, pass in an int
		MaxInFlight: fmt.Sprintf("%d", command.MaxInFlight),
	}

	appsGetter, err := diegohelpers.NewAppsGetterFunc(cliConnection, opts.Organization, runtime.Flip())
	if err != nil {
		return err
	}

	migrateAppsCommand, err := migratehelpers.NewMigrateAppsCommand(cliConnection, opts, runtime)
	if err != nil {
		return err
	}

	cmd := migratehelpers.MigrateApps{
		Opts:               opts,
		Runtime:            runtime,
		AppsGetterFunc:     appsGetter,
		MigrateAppsCommand: &migrateAppsCommand,
	}

	return cmd.Execute(cliConnection)
}
