package listhelpers

import (
	"errors"

	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
	"github.com/cloudfoundry/cli/plugin"
)

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
