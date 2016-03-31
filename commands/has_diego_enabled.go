package commands

import "github.com/cloudfoundry-incubator/diego-enabler/commands/diegohelpers"

type HasDiegoEnabledCommand struct {
	RequiredOptions HasDiegoEnabledPositionalArgs `positional-args:"yes"`
}

type HasDiegoEnabledPositionalArgs struct {
	AppName string `positional-arg-name:"APP_NAME" required:"true" description:"The app name"`
}

func (command HasDiegoEnabledCommand) Execute([]string) error {
	return diegohelpers.IsDiegoEnabled(DiegoEnabler.CLIConnection, command.RequiredOptions.AppName)
}
