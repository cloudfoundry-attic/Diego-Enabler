package commands

import (
	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
)

//go:generate counterfeiter . ApplicationsParser
type ApplicationsParser interface {
	Parse([]byte) (models.Applications, error)
}

type DiegoAppsCommand struct {
	OrganizationGuid string
}

func (c DiegoAppsCommand) DiegoApps(
	appsParser ApplicationsParser,
	paginatedRequester PaginatedRequester,
) (models.Applications, error) {
	var noApps models.Applications

	filter := api.Filters{
		api.EqualFilter{
			Name:  "diego",
			Value: true,
		},
	}

	if c.OrganizationGuid != "" {
		filter = append(
			filter,
			api.EqualFilter{
				Name:  "organization_guid",
				Value: c.OrganizationGuid,
			},
		)
	}
	params := map[string]interface{}{}

	responseBodies, err := paginatedRequester.Do(filter, params)
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
