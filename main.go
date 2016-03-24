package main

import (
	"errors"
	"fmt"
	"os"

	"crypto/tls"
	"net/http"

	"strings"

	"strconv"
	"time"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/commands"
	"github.com/cloudfoundry-incubator/diego-enabler/diego_support"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/jessevdk/go-flags"
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
					Usage: `cf diego-apps [-o ORG]

OPTIONS:
   -o      Organization to restrict the app migration to`,
				},
			},
			{
				Name:     "dea-apps",
				HelpText: "Lists all apps running on the DEA runtime that are visible to the user",
				UsageDetails: plugin.Usage{
					Usage: `cf dea-apps [-o ORG]

OPTIONS:
   -o      Organization to restrict the app migration to`,
				},
			},
			{
				Name:     "migrate-apps",
				HelpText: "Migrate all apps to Diego/DEA",
				UsageDetails: plugin.Usage{
					Usage: `cf migrate-apps (diego | dea) [-o ORG]

WARNING:
   Migration of a running app causes a restart. Stopped apps will be configured to run on the target runtime but are not started.

OPTIONS:
   -o      Organization to restrict the app migration to`,
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
	} else if args[0] == "diego-apps" {
		opts := parseArgs(args)
		diegoAppsCommand := newDiegoAppsCommand(cliConnection, opts)
		listAppsCommand := newListAppsCommand(cliConnection, opts)
		listAppsCommand.Runtime = ui.Runtime("diego")

		c.showApps(cliConnection, diegoAppsCommand.DiegoApps, listAppsCommand)
	} else if args[0] == "dea-apps" {
		opts := parseArgs(args)
		diegoAppsCommand := newDiegoAppsCommand(cliConnection, opts)
		listAppsCommand := newListAppsCommand(cliConnection, opts)
		listAppsCommand.Runtime = ui.Runtime("dea")

		c.showApps(cliConnection, diegoAppsCommand.DeaApps, listAppsCommand)
	} else if args[0] == "migrate-apps" && len(args) >= 2 {
		opts := parseArgs(args)
		diegoAppsCommand := newDiegoAppsCommand(cliConnection, opts)

		runtime := strings.ToLower(args[1])
		migrateAppsCommand := newMigrateAppsCommand(cliConnection, opts, runtime)

		if runtime == "diego" {
			c.migrateApps(cliConnection, diegoAppsCommand.DeaApps, true, migrateAppsCommand)
		} else if runtime == "dea" {
			c.migrateApps(cliConnection, diegoAppsCommand.DiegoApps, false, migrateAppsCommand)
		} else {
			c.showUsage(args)
		}
	} else {
		c.showUsage(args)
	}
}

func newListAppsCommand(cliConnection plugin.CliConnection, opts Opts) *ui.ListAppsCommand {
	username, err := cliConnection.Username()
	if err != nil {
		exitWithError(err, []string{})
	}

	traceEnv := os.Getenv("CF_TRACE")
	traceLogger := trace.NewLogger(false, traceEnv, "")
	tUI := terminal.NewUI(os.Stdin, terminal.NewTeePrinter(), traceLogger)

	cmd := &ui.ListAppsCommand{
		Username:     username,
		Organization: opts.Organization,
		UI:           tUI,
	}
	return cmd
}

func newMigrateAppsCommand(cliConnection plugin.CliConnection, opts Opts, runtime string) *ui.MigrateAppsCommand {
	username, err := cliConnection.Username()
	if err != nil {
		exitWithError(err, []string{})
	}

	cmd := &ui.MigrateAppsCommand{
		Username:     username,
		Organization: opts.Organization,
		Runtime:      ui.Runtime(runtime),
	}

	return cmd
}

type Opts struct {
	Organization string `short:"o"`
}

func parseArgs(args []string) Opts {
	var opts Opts

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		exitWithError(err, []string{})
	}

	return opts
}

func newDiegoAppsCommand(cliConnection plugin.CliConnection, opts Opts) commands.DiegoAppsCommand {
	diegoAppsCommand := commands.DiegoAppsCommand{}
	if opts.Organization != "" {
		org, err := cliConnection.GetOrg(opts.Organization)
		if err != nil {
			exitWithError(err, []string{})
		}
		diegoAppsCommand.OrganizationGuid = org.Guid
	}
	return diegoAppsCommand
}

func (c *DiegoEnabler) showApps(cliConnection plugin.CliConnection, appsGetter func(commands.ApplicationsParser, commands.PaginatedRequester) (models.Applications, error), p *ui.ListAppsCommand) {
	if err := verifyLoggedIn(cliConnection); err != nil {
		exitWithError(err, []string{})
	}

	accessToken, err := cliConnection.AccessToken()
	if err != nil {
		exitWithError(err, []string{})
	}

	p.BeforeAll()

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

	var appPrinters []ui.ApplicationPrinter
	for _, a := range apps {
		appPrinters = append(appPrinters, &appPrinter{
			app:    a,
			spaces: spaceMap,
		})
	}

	p.AfterAll(appPrinters)
}

func (c *DiegoEnabler) migrateApps(cliConnection plugin.CliConnection, appsGetter func(commands.ApplicationsParser, commands.PaginatedRequester) (models.Applications, error), enableDiego bool, p *ui.MigrateAppsCommand) {
	p.BeforeAll()

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

	warnings := 0
	for _, app := range apps {
		a := &appPrinter{
			app:    app,
			spaces: spaceMap,
		}

		p.BeforeEach(a)

		var waitTime time.Duration
		if app.State == models.Started {
			waitTime = 1 * time.Minute
			timeout := os.Getenv("CF_STARTUP_TIMEOUT")
			if timeout != "" {
				t, err := strconv.Atoi(timeout)

				if err == nil {
					waitTime = time.Duration(float32(t)/5.0*60.0) * time.Second
				}
			}
		}

		_, err := diegoSupport.SetDiegoFlag(app.Guid, enableDiego)
		if err != nil {
			warnings += 1
			fmt.Println("Error: ", err)
			fmt.Println("Continuing...")
			// WARNING: No authorization to migrate app APP_NAME in org ORG_NAME / space SPACE_NAME to RUNTIME as PERSON...
			continue
		}

		printDot := time.NewTicker(5 * time.Second)
		endDot := time.After(waitTime)
	dance:
		for {
			select {
			case <-endDot:
				printDot.Stop()
				break dance
			case <-printDot.C:
				p.DuringEach(a)
			}
		}

		p.CompletedEach(a)
	}

	p.AfterAll(len(apps), warnings)
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

	return space.OrganizationGuid
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
