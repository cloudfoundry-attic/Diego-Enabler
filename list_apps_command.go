package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/jessevdk/go-flags"
)

func OrgNotFound(org string, err error) error {
	return fmt.Errorf("Organization not found: %s, %+v", org, err)
}

type ListApps struct {
	opts            Opts
	Runtime         ui.Runtime
	appsGetterFunc  thingdoer.AppsGetterFunc
	listAppsCommand *ui.ListAppsCommand
}

type Opts struct {
	Organization string `short:"o"`
}

func PrepareListApps(args []string, cliConnection plugin.CliConnection) (ListApps, error) {
	empty := ListApps{}

	runtime, err := ui.ParseRuntime(args[0])
	if err != nil {
		return empty, err
	}

	var opts Opts
	_, err = flags.ParseArgs(&opts, args[1:])
	if err != nil {
		exitWithError(err, []string{})
	}

	appsGetter, err := NewAppsGetterFunc(cliConnection, opts.Organization, runtime)
	if err != nil {
		return empty, err
	}

	listAppsCommand, err := newListAppsCommand(cliConnection, opts.Organization, runtime)
	if err != nil {
		return empty, err
	}

	return ListApps{
		Runtime:         runtime,
		opts:            opts,
		appsGetterFunc:  appsGetter,
		listAppsCommand: &listAppsCommand,
	}, nil
}

func (cmd *ListApps) Execute(cliConnection plugin.CliConnection) error {

	if err := verifyLoggedIn(cliConnection); err != nil {
		return err
	}

	accessToken, err := cliConnection.AccessToken()
	if err != nil {
		return err
	}

	cmd.listAppsCommand.BeforeAll()

	pageParser := api.PageParser{}
	appsParser := models.ApplicationsParser{}
	spacesParser := models.SpacesParser{}

	apiEndpoint, err := cliConnection.ApiEndpoint()
	if err != nil {
		return err
	}

	apiClient, err := api.NewApiClient(apiEndpoint, accessToken)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	appRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetAppsRequest),
	)

	apps, err := cmd.appsGetterFunc(
		appsParser,
		&api.PaginatedRequester{
			RequestFactory: appRequestFactory,
			Client:         httpClient,
			PageParser:     pageParser,
		},
	)
	if err != nil {
		return err
	}

	spaceRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetSpacesRequest),
	)

	spaces, err := thingdoer.Spaces(
		spacesParser,
		&api.PaginatedRequester{
			RequestFactory: spaceRequestFactory,
			Client:         httpClient,
			PageParser:     pageParser,
		},
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
		appPrinters = append(appPrinters, &appPrinter{
			app:    a,
			spaces: spaceMap,
		})
	}

	cmd.listAppsCommand.AfterAll(appPrinters)

	return nil
}

func newListAppsCommand(cliConnection plugin.CliConnection, orgName string, runtime ui.Runtime) (ui.ListAppsCommand, error) {
	username, err := cliConnection.Username()
	if err != nil {
		return ui.ListAppsCommand{}, err
	}

	traceEnv := os.Getenv("CF_TRACE")
	traceLogger := trace.NewLogger(false, traceEnv, "")
	tUI := terminal.NewUI(os.Stdin, terminal.NewTeePrinter(), traceLogger)

	cmd := ui.ListAppsCommand{
		Username:     username,
		Organization: orgName,
		UI:           tUI,
		Runtime:      runtime,
	}
	return cmd, nil
}
