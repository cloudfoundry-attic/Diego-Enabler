package main

import (
	"errors"
	"fmt"
	"os"

	"crypto/tls"
	"net/http"

	"strings"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/commands"
	"github.com/cloudfoundry-incubator/diego-enabler/diego_support"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
)

type DiegoEnabler struct{}

func (c *DiegoEnabler) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Diego-Enabler",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     "enable-diego",
				HelpText: "enable Diego support for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf enable-diego APP_NAME",
				},
			},
			{
				Name:     "disable-diego",
				HelpText: "disable Diego support for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf disable-diego APP_NAME",
				},
			},
			{
				Name:     "has-diego-enabled",
				HelpText: "Check if Diego support is enabled for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf has-diego-enabled APP_NAME",
				},
			},
			{
				Name:     "diego-apps",
				HelpText: "Lists all apps running on the Diego runtime that are visible to the user",
				UsageDetails: plugin.Usage{
					Usage: "cf diego-apps",
				},
			},
			{
				Name:     "dea-apps",
				HelpText: "Lists all apps running on the DEA runtime that are visible to the user",
				UsageDetails: plugin.Usage{
					Usage: "cf dea-apps",
				},
			},
			{
				Name:     "migrate-apps",
				HelpText: "Migrate all apps to Diego/DEA",
				UsageDetails: plugin.Usage{
					Usage: `cf migrate-apps (diego | dea)

WARNING:
   Migration of a running app causes a restart. Stopped apps will be configured to run on the target runtime but are not started.`,
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(DiegoEnabler))
}

func (c *DiegoEnabler) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "enable-diego" && len(args) == 2 {
		c.toggleDiegoSupport(true, cliConnection, args[1])
	} else if args[0] == "disable-diego" && len(args) == 2 {
		c.toggleDiegoSupport(false, cliConnection, args[1])
	} else if args[0] == "has-diego-enabled" && len(args) == 2 {
		c.isDiegoEnabled(cliConnection, args[1])
	} else if args[0] == "diego-apps" && len(args) == 1 {
		c.showApps(cliConnection, commands.DiegoApps)
	} else if args[0] == "dea-apps" && len(args) == 1 {
		c.showApps(cliConnection, commands.DeaApps)
	} else if args[0] == "migrate-apps" && len(args) == 2 {
		runtime := strings.ToLower(args[1])

		if runtime == "diego" {
			c.migrateAppsToDiego(cliConnection)
		} else if runtime == "dea" {
			c.migrateAppsToDea(cliConnection)
		} else {
			c.showUsage(args)
		}
	} else {
		c.showUsage(args)
	}
}

func (c *DiegoEnabler) showApps(cliConnection plugin.CliConnection, appsGetter func(commands.ApplicationsParser, commands.PaginatedRequester) (models.Applications, error)) {
	username, err := cliConnection.Username()
	if err != nil {
		exitWithError(err, []string{})
	}

	if err := verifyLoggedIn(cliConnection); err != nil {
		exitWithError(err, []string{})
	}

	accessToken, err := cliConnection.AccessToken()
	if err != nil {
		exitWithError(err, []string{})
	}

	fmt.Printf("Getting apps on the Diego runtime as %s...\n", terminal.EntityNameColor(username))

	pageParser := api.PageParser{}
	appsParser := models.ApplicationsParser{}
	spacesParser := models.SpacesParser{}

	apiEndpoint, err := cliConnection.ApiEndpoint()
	if err != nil {
		exitWithError(err, []string{})
	}

	apiClient, err := api.NewApiClient(apiEndpoint, accessToken)
	if err != nil {
		exitWithError(err, []string{})
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	appRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetAppsRequest),
	)

	apps, err := appsGetter(
		appsParser,
		&api.PaginatedRequester{
			RequestFactory: appRequestFactory,
			Client:         httpClient,
			PageParser:     pageParser,
		},
	)
	if err != nil {
		exitWithError(err, []string{})
	}

	spaceRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetSpacesRequest),
	)

	spaces, err := commands.Spaces(
		spacesParser,
		&api.PaginatedRequester{
			RequestFactory: spaceRequestFactory,
			Client:         httpClient,
			PageParser:     pageParser,
		},
	)
	if err != nil {
		exitWithError(err, []string{})
	}

	spaceMap := make(map[string]models.Space)
	for _, space := range spaces {
		spaceMap[space.Guid] = space
	}

	sayOk()

	traceEnv := os.Getenv("CF_TRACE")
	traceLogger := trace.NewLogger(false, traceEnv, "")
	ui := terminal.NewUI(os.Stdin, terminal.NewTeePrinter(), traceLogger)

	headers := []string{
		"name",
		"space",
		"org",
	}
	t := terminal.NewTable(ui, headers)

	for _, app := range apps {
		t.Add(app.Name, spaceDisplayFor(app, spaceMap), orgDisplayFor(app, spaceMap))
	}

	t.Print()
}

func (c *DiegoEnabler) migrateAppsToDiego(cliConnection plugin.CliConnection) {
	c.migrateApps(cliConnection, commands.DeaApps, true)
}

func (c *DiegoEnabler) migrateAppsToDea(cliConnection plugin.CliConnection) {
	c.migrateApps(cliConnection, commands.DiegoApps, false)
}

func (c *DiegoEnabler) migrateApps(cliConnection plugin.CliConnection, appsGetter func(commands.ApplicationsParser, commands.PaginatedRequester) (models.Applications, error), enableDiego bool) {
	username, err := cliConnection.Username()
	if err != nil {
		exitWithError(err, []string{})
	}
	fmt.Printf("Migrating apps to Diego as %s...\n", terminal.EntityNameColor(username))

	if err := verifyLoggedIn(cliConnection); err != nil {
		exitWithError(err, []string{})
	}

	accessToken, err := cliConnection.AccessToken()
	if err != nil {
		exitWithError(err, []string{})
	}

	pageParser := api.PageParser{}
	appsParser := models.ApplicationsParser{}
	spacesParser := models.SpacesParser{}

	apiEndpoint, err := cliConnection.ApiEndpoint()
	if err != nil {
		exitWithError(err, []string{})
	}

	apiClient, err := api.NewApiClient(apiEndpoint, accessToken)
	if err != nil {
		exitWithError(err, []string{})
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	appRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetAppsRequest),
	)

	apps, err := appsGetter(
		appsParser,
		&api.PaginatedRequester{
			RequestFactory: appRequestFactory,
			Client:         httpClient,
			PageParser:     pageParser,
		},
	)
	if err != nil {
		exitWithError(err, []string{})
	}

	spaceRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetSpacesRequest),
	)

	spaces, err := commands.Spaces(
		spacesParser,
		&api.PaginatedRequester{
			RequestFactory: spaceRequestFactory,
			Client:         httpClient,
			PageParser:     pageParser,
		},
	)
	if err != nil {
		exitWithError(err, []string{})
	}

	spaceMap := make(map[string]models.Space)
	for _, space := range spaces {
		spaceMap[space.Guid] = space
	}

	diegoSupport := diego_support.NewDiegoSupport(cliConnection)

	var runtime string
	if enableDiego {
		runtime = "Diego"
	} else {
		runtime = "DEA"
	}

	warnings := 0
	for _, app := range apps {
		fmt.Println()
		orgName := orgDisplayFor(app, spaceMap)
		spaceName := spaceDisplayFor(app, spaceMap)

		fmt.Printf(
			"Migrating app %s in org %s / space %s to %s as %s...\n",
			terminal.EntityNameColor(app.Name),
			terminal.EntityNameColor(orgName),
			terminal.EntityNameColor(spaceName),
			terminal.EntityNameColor(runtime),
			terminal.EntityNameColor(username),
		)

		_, err := diegoSupport.SetDiegoFlag(app.Guid, enableDiego)
		if err != nil {
			warnings += 1
			fmt.Println("Error: ", err)
			fmt.Println("Continuing...")
			// WARNING: No authorization to migrate app APP_NAME in org ORG_NAME / space SPACE_NAME to RUNTIME as PERSON...
			continue
		}

		fmt.Printf(
			"Completed migrating app %s in org %s / space %s to %s as %s...\n",
			terminal.EntityNameColor(app.Name),
			terminal.EntityNameColor(orgName),
			terminal.EntityNameColor(spaceName),
			terminal.EntityNameColor(runtime),
			terminal.EntityNameColor(username),
		)
	}

	fmt.Println()
	fmt.Printf("Migration to %s completed: %d apps, %d warnings\n", terminal.EntityNameColor(runtime), len(apps), warnings)
}

func spaceDisplayFor(app models.Application, spaces map[string]models.Space) string {
	var display string

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

func orgDisplayFor(app models.Application, spaces map[string]models.Space) string {
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

	return space.OrganizationGuid
}

func (c *DiegoEnabler) showUsage(args []string) {
	for _, cmd := range c.GetMetadata().Commands {
		if cmd.Name == args[0] {
			fmt.Println("Invalid Usage: \n", cmd.UsageDetails.Usage)
		}
	}
}

func (c *DiegoEnabler) toggleDiegoSupport(on bool, cliConnection plugin.CliConnection, appName string) {
	d := diego_support.NewDiegoSupport(cliConnection)

	fmt.Printf("Setting %s Diego support to %t\n", appName, on)
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if output, err := d.SetDiegoFlag(app.Guid, on); err != nil {
		fmt.Println("err 1", err, output)
		exitWithError(err, output)
	}
	sayOk()

	fmt.Printf("Verifying %s Diego support is set to %t\n", appName, on)
	app, err = cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if app.Diego == on {
		sayOk()
	} else {
		sayFailed()
		fmt.Printf("Diego support for %s is NOT set to %t\n\n", appName, on)
		os.Exit(1)
	}
}

func (c *DiegoEnabler) isDiegoEnabled(cliConnection plugin.CliConnection, appName string) {
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if app.Guid == "" {
		sayFailed()
		fmt.Printf("App %s not found\n\n", appName)
		os.Exit(1)
	}

	fmt.Println(app.Diego)
}

func exitWithError(err error, output []string) {
	sayFailed()
	fmt.Println("Error: ", err)
	for _, str := range output {
		fmt.Println(str)
	}
	os.Exit(1)
}

func say(message string, color uint, bold int) string {
	return fmt.Sprintf("\033[%d;%dm%s\033[0m", bold, color, message)
}

func sayOk() {
	fmt.Println(say("Ok\n", 32, 1))
}

func sayFailed() {
	fmt.Println(say("FAILED", 31, 1))
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
