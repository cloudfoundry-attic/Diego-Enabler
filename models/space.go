package models

import "encoding/json"

type Spaces []Space

type SpaceEntity struct {
	Name             string       `json:"name"`
	OrganizationGuid string       `json:"organization_guid"`
	Organization     Organization `json:"organization"`
}

type SpaceMetadata struct {
	Guid string `json:"guid"`
}

type SpacesResponse struct {
	Resources Spaces `json:"resources"`
}

type Space struct {
	SpaceEntity   `json:"entity"`
	SpaceMetadata `json:"metadata"`
}

type Organization struct {
	OrganizationEntity   `json:"entity"`
	OrganizationMetadata `json:"metadata"`
}

type OrganizationEntity struct {
	Name string `json:"name"`
}

type OrganizationMetadata struct {
	Guid string `json:"guid"`
}

type SpacesParser struct{}

func (a SpacesParser) Parse(body []byte) (Spaces, error) {
	var response SpacesResponse
	var emptySpaces Spaces

	err := json.Unmarshal(body, &response)
	if err != nil {
		return emptySpaces, err
	}

	return response.Resources, nil
}
