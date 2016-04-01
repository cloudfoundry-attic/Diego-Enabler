package listhelpers

import (
	"os"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/displayhelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
)

func ListApps(cliConnection api.Connection, appsGetterFunc thingdoer.AppsGetterFunc, listAppsCommand *ui.ListAppsCommand) error {
	listAppsCommand.BeforeAll()

	appsParser := models.ApplicationsParser{}
	spacesParser := models.SpacesParser{}

	apiClient, err := api.NewClient(cliConnection)
	if err != nil {
		return err
	}

	appRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetAppsRequest),
	)

	appPaginatedRequester, err := api.NewPaginatedRequester(cliConnection, appRequestFactory)
	if err != nil {
		return err
	}

	apps, err := appsGetterFunc(
		appsParser,
		appPaginatedRequester,
	)
	if err != nil {
		return err
	}

	spaceRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetSpacesRequest),
	)
	spacesPaginatedRequester, err := api.NewPaginatedRequester(cliConnection, spaceRequestFactory)
	if err != nil {
		return err
	}

	spaces, err := thingdoer.Spaces(
		spacesParser,
		spacesPaginatedRequester,
	)
	if err != nil {
		return err
	}

	spaceMap := make(map[string]models.Space)
	for _, space := range spaces {
		spaceMap[space.Guid] = space
	}

	var appPrinters []ui.ApplicationPrinter
	for _, a := range apps {
		appPrinters = append(appPrinters, &displayhelpers.AppPrinter{
			App:    a,
			Spaces: spaceMap,
		})
	}

	listAppsCommand.AfterAll(appPrinters)

	return nil
}

func NewListAppsCommand(cliConnection api.Connection, orgName string, spaceName string, runtime ui.Runtime) (ui.ListAppsCommand, error) {
	username, err := cliConnection.Username()
	if err != nil {
		return ui.ListAppsCommand{}, err
	}

	if spaceName != "" {
		space, err := cliConnection.GetSpace(spaceName)
		if err != nil || space.Guid == "" {
			return ui.ListAppsCommand{}, err
		}
		orgName = space.Organization.Name
	}

	traceEnv := os.Getenv("CF_TRACE")
	traceLogger := trace.NewLogger(false, traceEnv, "")
	tUI := terminal.NewUI(os.Stdin, terminal.NewTeePrinter(), traceLogger)

	cmd := ui.ListAppsCommand{
		Username:     username,
		Organization: orgName,
		Space:        spaceName,
		UI:           tUI,
		Runtime:      runtime,
	}
	return cmd, nil
}
