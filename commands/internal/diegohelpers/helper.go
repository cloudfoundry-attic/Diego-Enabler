package diegohelpers

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry-incubator/diego-enabler/commands/internal/displayhelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/diego_support"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
	"github.com/cloudfoundry/cli/plugin"
)

func ToggleDiegoSupport(on bool, cliConnection plugin.CliConnection, appName string) error {
	d := diego_support.NewDiegoSupport(cliConnection)

	fmt.Printf("Setting %s Diego support to %t\n", appName, on)
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		return err
	}

	if output, err := d.SetDiegoFlag(app.Guid, on); err != nil {
		return fmt.Errorf("%s\n%s", err, strings.Join(output, "\n"))
	}
	displayhelpers.SayOK()

	fmt.Printf("Verifying %s Diego support is set to %t\n", appName, on)
	app, err = cliConnection.GetApp(appName)
	if err != nil {
		return err
	}

	if app.Diego == on {
		displayhelpers.SayOK()
	} else {
		displayhelpers.SayFailed()
		return fmt.Errorf("Diego support for %s is NOT set to %t\n\n", appName, on)
	}

	return nil
}

func IsDiegoEnabled(cliConnection plugin.CliConnection, appName string) error {
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		return err
	}

	if app.Guid == "" {
		displayhelpers.SayFailed()
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

func NewAppsGetterFunc(cliConnection plugin.CliConnection, orgName string, runtime ui.Runtime) (thingdoer.AppsGetterFunc, error) {
	diegoAppsCommand := thingdoer.AppsGetter{}
	if orgName != "" {
		org, err := cliConnection.GetOrg(orgName)
		if err != nil || org.Guid == "" {
			return nil, OrgNotFoundErr{OrganizationName: orgName}
		}
		diegoAppsCommand.OrganizationGuid = org.Guid
	}

	var appsGetterFunc = diegoAppsCommand.DiegoApps
	if runtime == ui.DEA {
		appsGetterFunc = diegoAppsCommand.DeaApps
	}

	return appsGetterFunc, nil
}
