package commands

import (
	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
)

func DeaApps(requestFactory RequestFactory, client CloudControllerClient, appsParser ApplicationsParser, pageParser PaginatedParser) (models.Applications, error) {
	var noApps models.Applications

	filter := api.EqualFilter{
		Name:  "diego",
		Value: false,
	}

	params := map[string]interface{}{}

	responseBodies, err := paginatedRequester(requestFactory, filter, params, client, pageParser)
	if err != nil {
		return noApps, err
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
