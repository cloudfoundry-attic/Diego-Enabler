package errorhelpers

import "errors"

var SpecifyOrgOrSpaceError = errors.New("Cannot specify org together with space.")

func ErrorIfOrgAndSpacesSet(orgName, spaceName string) error {
	if orgName != "" && spaceName != "" {
		return SpecifyOrgOrSpaceError
	}
	return nil
}
