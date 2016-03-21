package models

import "encoding/json"

type Applications []Application

type Application struct {
	//Guid                 string
	Name string `json:"name"`
	//BuildpackUrl         string
	//Command              string
	Diego bool
	//DetectedStartCommand string
	//DiskQuota            int64 // in Megabytes
	//EnvironmentVars      map[string]interface{}
	//InstanceCount        int
	//Memory               int64 // in Megabytes
	//RunningInstances     int
	//HealthCheckTimeout   int
	//State                string
	//SpaceGuid            string
	//PackageUpdatedAt     *time.Time
	//PackageState         string
	//StagingFailedReason  string
	//AppPorts             []int
	//Stack                *GetApp_Stack
	//Instances            []GetApp_AppInstanceFields
	//Routes               []GetApp_RouteSummary
	//Services             []GetApp_ServiceSummary
}

type ApplicationsResponse struct {
	Resources []ApplicationResource `json:"resources"`
}

type ApplicationResource struct {
	Application `json:"entity"`
}

type ApplicationsParser struct{}

func (a ApplicationsParser) Parse(body []byte) (Applications, error) {
	var response ApplicationsResponse

	var applications Applications
	var emptyApplications Applications

	err := json.Unmarshal(body, &response)
	if err != nil {
		return emptyApplications, err
	}

	for _, app := range response.Resources {
		applications = append(applications, app.Application)
	}

	return applications, nil
}
