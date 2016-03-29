package migratehelpers

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/diego-enabler/models"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
	"github.com/cloudfoundry/cli/plugin"
)

//Originally in listhelpers/helper.go
func OrgNotFound(org string) error {
	return fmt.Errorf("Organization not found: %s", org)
}

func NewAppsGetterFunc(cliConnection plugin.CliConnection, orgName string, runtime ui.Runtime) (thingdoer.AppsGetterFunc, error) {
	diegoAppsCommand := thingdoer.AppsGetter{}
	if orgName != "" {
		org, err := cliConnection.GetOrg(orgName)
		if err != nil || org.Guid == "" {
			return nil, OrgNotFound(orgName)
		}
		diegoAppsCommand.OrganizationGuid = org.Guid
	}

	var appsGetterFunc = diegoAppsCommand.DiegoApps
	if runtime == ui.DEA {
		appsGetterFunc = diegoAppsCommand.DeaApps
	}

	return appsGetterFunc, nil
}

type appPrinter struct {
	app    models.Application
	spaces map[string]models.Space
}

func (a *appPrinter) Name() string {
	return a.app.Name
}

func (a *appPrinter) Organization() string {
	spaces := a.spaces
	app := a.app

	if len(spaces) == 0 {
		return ""
	}

	space, ok := spaces[app.SpaceGuid]
	if !ok {
		return ""
	}

	if space.Organization.Name != "" {
		return space.Organization.Name
	}

	return space.Organization.Guid
}

func (a *appPrinter) Space() string {
	var display string
	spaces := a.spaces
	app := a.app

	if len(spaces) == 0 {
		display = app.SpaceGuid
	} else {
		space, ok := spaces[app.SpaceGuid]
		if ok {
			display = space.Name
		} else {
			display = app.SpaceGuid
		}
	}

	return display
}

func verifyLoggedIn(cliCon plugin.CliConnection) error {
	var result error

	if connected, err := cliCon.IsLoggedIn(); !connected {
		result = NotLoggedInError

		if err != nil {
			result = err
		}
	}

	return result
}

var NotLoggedInError = errors.New("You must be logged in")
