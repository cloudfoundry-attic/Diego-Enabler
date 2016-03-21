package commands

import (
	"net/http"

	"io/ioutil"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
)

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
	Parse([]byte) (models.Applications, error)
}

//go:generate counterfeiter . PaginatedParser
type PaginatedParser interface {
	Parse([]byte) (api.PaginatedResponse, error)
}

func DiegoApps(factory RequestFactory, client CloudControllerClient, appsParser ResponseParser, pageParser PaginatedParser) (models.Applications, error) {
	var noApps models.Applications

	filter := api.EqualFilter{
		Name:  "diego",
		Value: true,
	}

	params := map[string]interface{}{}

	req, err := factory.NewGetAppsRequest(filter, params)
	if err != nil {
		return noApps, err
	}

	var responseBodies [][]byte

	res, err := client.Do(req)
	if err != nil {
		return noApps, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return noApps, err
	}

	responseBodies = append(responseBodies, body)

	paginatedRes, err := pageParser.Parse(body)
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

		// perform the request
		res, err := client.Do(req)
		if err != nil {
			return noApps, err
		}

		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return noApps, err
		}

		responseBodies = append(responseBodies, body)
	}

	var applications models.Applications

	for _, nextBody := range responseBodies {
		apps, err := appsParser.Parse(nextBody)
		if err != nil {
			return noApps, err
		}

		applications = append(applications, apps...)
	}

	return applications, nil
}
