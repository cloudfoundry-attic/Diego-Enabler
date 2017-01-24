package thingdoer

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
)

type ccError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	ErrorCode   string `json:"error_code"`
}

func (e ccError) Error() string {
	return fmt.Sprintf("CC code:       %d\nCC error code: %s\nDescription:   %s",
		e.Code, e.ErrorCode, e.Description)
}

func handleCCError(ccErrResponse string) error {
	ccErr := ccError{}
	err := json.Unmarshal([]byte(ccErrResponse), &ccErr)
	if err != nil {
		return fmt.Errorf("Unexpected response:\n%s", ccErrResponse)
	}

	return ccErr
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
			return noApps, fmt.Errorf("Unable to get routes for app '%s'\n%s", app.Name, err.Error())
		}

		applications[i].HasRoutes = hasRoutes
	}

	return applications, nil
}
