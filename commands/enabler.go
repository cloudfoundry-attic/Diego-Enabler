package commands

import "github.com/cloudfoundry/cli/plugin"

type Enabler struct {
	CLIConnection plugin.CliConnection

	EnableDiego     EnableDiegoCommand     `command:"enable-diego" description:"enable Diego support for an app"`
	DisableDiego    DisableDiegoCommand    `command:"disable-diego" description:"disable Diego support for an app"`
	HasDiegoEnabled HasDiegoEnabledCommand `command:"has-diego-enabled" description:"Check if Diego support is enabled for an app"`
	DiegoApps       DiegoAppsCommand       `command:"diego-apps" description:"Lists all apps running on the Diego runtime that are visible to the user"`
	DeaApps         DeaAppsCommand         `command:"dea-apps" description:"Lists all apps running on the DEA runtime that are visible to the user"`
	MigrateApps     MigrateAppsCommand     `command:"migrate-apps" description:"Migrate all apps to Diego/DEA"`
	UninstallPlugin UninstallHook          `command:"CLI-MESSAGE-UNINSTALL"`
}

var DiegoEnabler Enabler
