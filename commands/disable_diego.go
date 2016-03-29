package commands

import "github.com/cloudfoundry-incubator/diego-enabler/commands/internal/diegohelpers"

type DisableDiegoCommand struct {
	RequiredOptions DisableDiegoPositionalArgs `positional-args:"yes"`
}

type DisableDiegoPositionalArgs struct {
	AppName string `positional-arg-name:"APP_NAME" required:"true" description:"The app name"`
}

func (command DisableDiegoCommand) Execute([]string) error {
	return diegohelpers.ToggleDiegoSupport(false, DiegoEnabler.CLIConnection, command.RequiredOptions.AppName)
}
