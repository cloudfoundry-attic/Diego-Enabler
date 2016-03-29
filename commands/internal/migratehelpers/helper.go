package migratehelpers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"sync"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/internal/displayhelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/commands/internal/verificationhelpers"
	"github.com/cloudfoundry-incubator/diego-enabler/diego_support"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
	"github.com/cloudfoundry-incubator/diego-enabler/thingdoer"
	"github.com/cloudfoundry-incubator/diego-enabler/ui"
	"github.com/cloudfoundry/cli/plugin"
)

type MigrateAppsOpts struct {
	Organization string `short:"o"`
	MaxInFlight  string `short:"p"`
}

type MigrateApps struct {
	Opts               MigrateAppsOpts
	Runtime            ui.Runtime
	AppsGetterFunc     thingdoer.AppsGetterFunc
	MigrateAppsCommand *ui.MigrateAppsCommand
}

func (cmd *MigrateApps) Execute(cliConnection plugin.CliConnection) error {
	cmd.MigrateAppsCommand.BeforeAll()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	apiClient, err := newAPIClient(cliConnection)
	if err != nil {
		return err
	}

	appRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetAppsRequest),
	)

	pageParser := api.PageParser{}
	apps, err := cmd.AppsGetterFunc(
		models.ApplicationsParser{},
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
		models.SpacesParser{},
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

	maxInFlight := 1
	if cmd.Opts.MaxInFlight != "" {
		maxInFlight, err = strconv.Atoi(cmd.Opts.MaxInFlight)
		if err != nil || maxInFlight <= 0 || maxInFlight > 100 {
			return fmt.Errorf("Invalid maximum apps in flight: %s\nValue for MAX_IN_FLIGHT must be an integer between 1 and 100", cmd.Opts.MaxInFlight)
		}
	}

	warnings := cmd.migrateApps(cliConnection, apps, spaceMap, maxInFlight)
	cmd.MigrateAppsCommand.AfterAll(len(apps), warnings)

	return nil
}

func NewMigrateAppsCommand(cliConnection plugin.CliConnection, opts MigrateAppsOpts, runtime ui.Runtime) (ui.MigrateAppsCommand, error) {
	username, err := cliConnection.Username()
	if err != nil {
		return ui.MigrateAppsCommand{}, err
	}

	return ui.MigrateAppsCommand{
		Username:     username,
		Organization: opts.Organization,
		Runtime:      runtime,
	}, nil
}

func newAPIClient(cliConnection plugin.CliConnection) (*api.ApiClient, error) {
	if err := verificationhelpers.VerifyLoggedIn(cliConnection); err != nil {
		return nil, err
	}

	accessToken, err := cliConnection.AccessToken()
	if err != nil {
		return nil, err
	}

	apiEndpoint, err := cliConnection.ApiEndpoint()
	if err != nil {
		return nil, err
	}

	apiClient, err := api.NewApiClient(apiEndpoint, accessToken)
	if err != nil {
		return nil, err
	}

	return apiClient, nil
}

type migrateAppFunc func(appPrinter *displayhelpers.AppPrinter, diegoSupport *diego_support.DiegoSupport) bool

func (cmd *MigrateApps) migrateApp(appPrinter *displayhelpers.AppPrinter, diegoSupport *diego_support.DiegoSupport) bool {
	cmd.MigrateAppsCommand.BeforeEach(appPrinter)

	var waitTime time.Duration
	if appPrinter.App.State == models.Started {
		waitTime = 1 * time.Minute
		timeout := os.Getenv("CF_STARTUP_TIMEOUT")
		if timeout != "" {
			t, err := strconv.Atoi(timeout)

			if err == nil {
				waitTime = time.Duration(float32(t)/5.0*60.0) * time.Second
			}
		}
	}

	_, err := diegoSupport.SetDiegoFlag(appPrinter.App.Guid, cmd.Runtime == ui.Diego)
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Continuing...")
		// WARNING: No authorization to migrate app APP_NAME in org ORG_NAME / space SPACE_NAME to RUNTIME as PERSON...
		return false
	}

	printDot := time.NewTicker(5 * time.Second)
	go func() {
		for range printDot.C {
			cmd.MigrateAppsCommand.DuringEach(appPrinter)
		}
	}()

	time.Sleep(waitTime)
	printDot.Stop()

	cmd.MigrateAppsCommand.CompletedEach(appPrinter)

	return true
}

func (cmd *MigrateApps) migrateApps(cliConnection plugin.CliConnection, apps models.Applications, spaceMap map[string]models.Space, maxInFlight int) int {
	if len(apps) < maxInFlight {
		maxInFlight = len(apps)
	}

	runningAppsChan := generateAppsChan(apps)
	outputsChan, waitDone := processAppsChan(cliConnection, spaceMap, cmd.migrateApp, runningAppsChan, maxInFlight, len(apps))

	waitDone.Wait()
	close(outputsChan)

	return outputAppsChan(outputsChan)
}

func generateAppsChan(apps models.Applications) chan models.Application {
	runningAppsChan := make(chan models.Application)
	go func() {
		defer close(runningAppsChan)
		for _, app := range apps {
			runningAppsChan <- app
		}
	}()

	return runningAppsChan
}

func processAppsChan(
	cliConnection plugin.CliConnection,
	spaceMap map[string]models.Space,
	migrate migrateAppFunc,
	appsChan chan models.Application,
	maxInFlight int,
	outputSize int) (chan bool, *sync.WaitGroup) {
	var waitDone sync.WaitGroup

	output := make(chan bool, outputSize)

	diegoSupport := diego_support.NewDiegoSupport(cliConnection)

	for i := 0; i < maxInFlight; i++ {
		waitDone.Add(1)

		go func() {
			defer waitDone.Done()

			for app := range appsChan {
				a := &displayhelpers.AppPrinter{
					App:    app,
					Spaces: spaceMap,
				}
				output <- migrate(a, diegoSupport)
			}
		}()
	}
	return output, &waitDone
}

func outputAppsChan(outputsChan chan bool) int {
	warnings := 0

	for success := range outputsChan {
		if !success {
			warnings++
		}
	}

	return warnings
}
