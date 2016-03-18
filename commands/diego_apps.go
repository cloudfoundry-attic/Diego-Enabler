package commands

import (
	"errors"
	"net/http"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
)

var NotLoggedInError = errors.New("You must be logged in")

//go:generate counterfeiter . RequestFactory
type RequestFactory interface {
	NewGetAppsRequest(api.Filter, map[string]interface{}) (*http.Request, error)
}

//go:generate counterfeiter . CloudControllerClient
type CloudControllerClient interface {
	Do(*http.Request) (*http.Response, error)
}

//go:generate counterfeiter . ResponseParser
type ResponseParser interface {
	Parse(*http.Response) (models.Applications, error)
}

//go:generate counterfeiter . PaginatedParser
type PaginatedParser interface {
	Parse(*http.Response) (api.PaginatedResponse, error)
}

//go:generate counterfeiter . CliConnection
type CliConnection interface {
	IsLoggedIn() (bool, error)
	AccessToken() (string, error)
}

func DiegoApps(cliCon CliConnection, factory RequestFactory, client CloudControllerClient, appsParser ResponseParser, pageParser PaginatedParser) (models.Applications, error) {
	var noApps models.Applications

	if err := verifyLoggedIn(cliCon); err != nil {
		return noApps, err
	}

	filter := api.EqualFilter{
		Name:  "diego",
		Value: true,
	}

	params := map[string]interface{}{}

	req, err := factory.NewGetAppsRequest(filter, params)
	if err != nil {
		return noApps, err
	}

	accessToken, err := cliCon.AccessToken()
	if err != nil {
		return noApps, err
	}

	req.Header.Set("Authorization", accessToken)

	var responses []*http.Response

	res, err := client.Do(req)
	if err != nil {
		return noApps, err
	}
	responses = append(responses, res)

	paginatedRes, err := pageParser.Parse(res)
	if err != nil {
		return noApps, err
	}
	for page := 2; page <= paginatedRes.TotalPages; page++ {
		// construct a new request with the current page
		params["page"] = page
		req, err := factory.NewGetAppsRequest(filter, params)
		if err != nil {
			return noApps, err
		}

		req.Header.Set("Authorization", accessToken)

		// perform the request
		res, err := client.Do(req)
		if err != nil {
			return noApps, err
		}

		responses = append(responses, res)
	}

	var applications models.Applications

	for _, nextRes := range responses {
		apps, err := appsParser.Parse(nextRes)
		if err != nil {
			return noApps, err
		}

		applications = append(applications, apps...)
	}

	return applications, nil
}

func verifyLoggedIn(cliCon CliConnection) error {
	var result error

	if connected, err := cliCon.IsLoggedIn(); !connected {
		result = NotLoggedInError

		if err != nil {
			result = err
		}
	}

	return result
}
