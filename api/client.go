package api

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cloudfoundry/cli/plugin/models"
)

type Client struct {
	BaseUrl   *url.URL
	AuthToken string
}

//go:generate counterfeiter . Connection

type Connection interface {
	IsLoggedIn() (bool, error)
	IsSSLDisabled() (bool, error)
	ApiEndpoint() (string, error)
	AccessToken() (string, error)

	Username() (string, error)

	CliCommandWithoutTerminalOutput(args ...string) ([]string, error)
	GetApp(string) (plugin_models.GetAppModel, error)
	GetOrg(string) (plugin_models.GetOrg_Model, error)
	GetSpace(string) (plugin_models.GetSpace_Model, error)
}

var NotLoggedInError = errors.New("You must be logged in")

func NewClient(connection Connection) (*Client, error) {
	if connected, err := connection.IsLoggedIn(); !connected {
		if err != nil {
			return nil, err
		}
		return nil, NotLoggedInError
	}

	rawURL, err := connection.ApiEndpoint()
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	authToken, err := connection.AccessToken()
	if err != nil {
		return nil, err
	}

	client := &Client{
		BaseUrl:   u,
		AuthToken: authToken,
	}

	return client, nil
}

func (c *Client) NewGetAppsRequest() (*http.Request, error) {
	req := &http.Request{
		Method: "GET",
		URL:    c.BaseUrl,
	}
	req.URL.Path = "/v2/apps"

	return req, nil
}

func (c *Client) NewGetSpacesRequest() (*http.Request, error) {
	req := &http.Request{
		Method: "GET",
		URL:    c.BaseUrl,
	}
	req.URL.Path = "/v2/spaces"

	return req, nil
}

func (c *Client) HandleFiltersAndParameters(next func() (*http.Request, error)) func(filter Filter, params map[string]interface{}) (*http.Request, error) {
	return func(filter Filter, params map[string]interface{}) (*http.Request, error) {
		req, err := next()
		if err != nil {
			return new(http.Request), err
		}

		req.URL.RawQuery = generateParams(filter, params).Encode()
		return req, nil
	}
}

func (c *Client) Authorize(next func() (*http.Request, error)) func() (*http.Request, error) {
	return func() (*http.Request, error) {
		req, err := next()
		if err != nil {
			return new(http.Request), err
		}

		header := http.Header{}
		header.Set("Authorization", c.AuthToken)

		req.Header = header
		return req, nil
	}
}

func generateParams(filter Filter, params map[string]interface{}) url.Values {
	values := url.Values{}
	q := filter.ToFilterQueryParam()
	if q != "" {
		values.Set("q", q)
	}

	for k, v := range params {
		values.Set(k, fmt.Sprint(v))
	}
	return values
}

func NewHttpClient(cliConnection Connection) (*http.Client, error) {
	skipVerify, err := cliConnection.IsSSLDisabled()
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
			Proxy:           http.ProxyFromEnvironment,
		},
	}
	return httpClient, nil
}
