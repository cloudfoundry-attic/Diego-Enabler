package thingdoer

import (
	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
)

func (c AppsGetter) DeaApps(appsParser ApplicationsParser, paginatedRequester PaginatedRequester) (models.Applications, error) {
	var noApps models.Applications

	filter := api.Filters{
		api.EqualFilter{
			Name:  "diego",
			Value: false,
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
	} else if c.SpaceGuid != "" {
		filter = append(
			filter,
			api.EqualFilter{
				Name:  "space_guid",
				Value: c.SpaceGuid,
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
