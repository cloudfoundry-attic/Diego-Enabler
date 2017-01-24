package thingdoer

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
)

func handleCCError(ccErrResponse string) error {
	var ccErr struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		ErrorCode   string `json:"error_code"`
	}

	err := json.Unmarshal([]byte(ccErrResponse), &ccErr)
	if err != nil {
		return fmt.Errorf("Unexpected response:\n%s", ccErrResponse)
	}

	return fmt.Errorf("Cloud controller error:\nCode:          %d\nDescription:   %s\nError Code:    %s", ccErr.Code, ccErr.Description, ccErr.ErrorCode)
}

func (c AppsGetter) ApplicationHasRoutes(appGUID string) (bool, error) {
	response, err := c.CliConnection.CliCommandWithoutTerminalOutput("curl", fmt.Sprintf("/v2/apps/%s/routes", appGUID))
	if err != nil {
		return false, err
	}

	strResponse := strings.Join(response, "")

	var ccMetadata struct {
		TotalResults int `json:"total_results"`
	}
	err = json.Unmarshal([]byte(strResponse), &ccMetadata)
	if err != nil {
		return false, errors.New(strResponse)
	}

	if strings.Contains(strResponse, `"error_code":`) {
		err = handleCCError(strResponse)
	}

	return ccMetadata.TotalResults > 0, err
}

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

	for i, app := range applications {
		hasRoutes, err := c.ApplicationHasRoutes(app.Guid)
		if err != nil {
			return noApps, fmt.Errorf("Error getting routes for app '%s': %s", app.Name, err.Error())
		}

		applications[i].HasRoutes = hasRoutes
	}

	return applications, nil
}
