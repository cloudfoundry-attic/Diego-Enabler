package commands

import (
	"errors"
	"net/http"

	"github.com/cloudfoundry-incubator/diego-enabler/models"
)

var NotLoggedInError = errors.New("You must be logged in")

//go:generate counterfeiter . RequestFactory
type RequestFactory interface {
	NewGetAppsRequest(map[string]interface{}) (*http.Request, error)
}

//go:generate counterfeiter . CloudControllerClient
type CloudControllerClient interface {
	Do(*http.Request) (*http.Response, error)
}

//go:generate counterfeiter . ResponseParser
type ResponseParser interface {
	Parse(*http.Response) (models.Applications, error)
}

//go:generate counterfeiter . CliConnection
type CliConnection interface {
	IsLoggedIn() (bool, error)
	AccessToken() (string, error)
}

func DiegoApps(cliCon CliConnection, factory RequestFactory, client CloudControllerClient, resParser ResponseParser) (models.Applications, error) {
	var noApps models.Applications

	if err := verifyLoggedIn(cliCon); err != nil {
		return noApps, err
	}

	params := map[string]interface{}{
		"diego": true,
	}

	req, err := factory.NewGetAppsRequest(params)
	if err != nil {
		return noApps, err
	}

	accessToken, err := cliCon.AccessToken()
	if err != nil {
		return noApps, err
	}

	req.Header.Set("Authorization", accessToken)

	res, err := client.Do(req)
	if err != nil {
		return noApps, err
	}

	return resParser.Parse(res)
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
