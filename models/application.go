package models

import "encoding/json"

type Applications []Application

type ApplicationMetadata struct {
	Guid string `json:"guid"`
}

const (
	Started = "STARTED"
	Stopped = "STOPPED"
)

type ApplicationEntity struct {
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
	State     string `json:"state"`
	SpaceGuid string `json:"space_guid"`
	//PackageUpdatedAt     *time.Time
	//PackageState         string
	//StagingFailedReason  string
	//AppPorts             []int
	//Stack                *GetApp_Stack
	//Instances            []GetApp_AppInstanceFields
	//Routes               []GetApp_RouteSummary
	//Services             []GetApp_ServiceSummary
	HasRoutes bool
}

type ApplicationsResponse struct {
	Resources Applications `json:"resources"`
}

type Application struct {
	ApplicationEntity   `json:"entity"`
	ApplicationMetadata `json:"metadata"`
}

type ApplicationsParser struct{}

func (a ApplicationsParser) Parse(body []byte) (Applications, error) {
	var response ApplicationsResponse
	var emptyApplications Applications

	err := json.Unmarshal(body, &response)
	if err != nil {
		return emptyApplications, err
	}

	return response.Resources, nil
}
