package commands

import "github.com/cloudfoundry-incubator/diego-enabler/commands/internal/diegohelpers"

type EnableDiegoCommand struct {
	RequiredOptions EnableDiegoPositionalArgs `positional-args:"yes"`
}

type EnableDiegoPositionalArgs struct {
	AppName string `positional-arg-name:"APP_NAME" required:"true" description:"The app name"`
}

func (command EnableDiegoCommand) Execute([]string) error {
	return diegohelpers.ToggleDiegoSupport(true, DiegoEnabler.CLIConnection, command.RequiredOptions.AppName)
}
