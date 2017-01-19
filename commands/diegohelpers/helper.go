package diegohelpers

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/diegosupport"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
)

func ToggleDiegoSupport(on bool, cliConnection api.Connection, appName string) error {
	d := diegosupport.NewDiegoSupport(cliConnection)

	fmt.Printf("Setting %s Diego support to %t\n", appName, on)
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		return err
	}

	err = d.WarnNoRoutes(appName, os.Stderr)
	if err != nil {
		return err
	}

	if output, err := d.SetDiegoFlag(app.Guid, on); err != nil {
		return fmt.Errorf("%s\n%s", err, strings.Join(output, "\n"))
	}
	ui.SayOK()

	fmt.Printf("Verifying %s Diego support is set to %t\n", appName, on)
	app, err = cliConnection.GetApp(appName)
	if err != nil {
		return err
	}

	if app.Diego == on {
		ui.SayOK()
	} else {
		return fmt.Errorf("Diego support for %s is NOT set to %t\n\n", appName, on)
	}

	return nil
}

func IsDiegoEnabled(cliConnection api.Connection, appName string) error {
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		return err
	}

	if app.Guid == "" {
		return fmt.Errorf("App %s not found\n\n", appName)
	}

	fmt.Println(app.Diego)

	return nil
}

type OrgNotFoundErr struct {
	OrganizationName string
}

func (e OrgNotFoundErr) Error() string {
	return fmt.Sprintf("Organization not found: %s", e.OrganizationName)
}

type SpaceNotFoundErr struct {
	SpaceName string
}

func (e SpaceNotFoundErr) Error() string {
	return fmt.Sprintf("Space not found: %s", e.SpaceName)
}

func NewAppsGetterFunc(
	cliConnection api.Connection,
	orgName string,
	spaceName string,
	runtime ui.Runtime,
) (thingdoer.AppsGetterFunc, error) {
	diegoAppsCommand := thingdoer.AppsGetter{}

	if orgName != "" {
		org, err := cliConnection.GetOrg(orgName)
		if err != nil || org.Guid == "" {
			return nil, OrgNotFoundErr{OrganizationName: orgName}
		}
		diegoAppsCommand.OrganizationGuid = org.Guid
	} else if spaceName != "" {
		space, err := cliConnection.GetSpace(spaceName)
		if err != nil || space.Guid == "" {
			return nil, SpaceNotFoundErr{SpaceName: spaceName}
		}
		diegoAppsCommand.SpaceGuid = space.Guid
	}

	var appsGetterFunc = diegoAppsCommand.DiegoApps
	if runtime == ui.DEA {
		appsGetterFunc = diegoAppsCommand.DeaApps
	}

	return appsGetterFunc, nil
}
