package models

type Applications []Application

type Application struct {
	//Guid                 string
	//Name                 string
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
