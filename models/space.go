package models

import "encoding/json"

type Spaces []Space

type SpaceEntity struct {
	Name string `json:"name"`
}

type SpaceMetadata struct {
	Guid string `json:"guid"`
}

type SpacesResponse struct {
	Resources Spaces `json:"resources"`
}

type Space struct {
	SpaceEntity `json:"entity"`
	SpaceMetadata `json:"metadata"`
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
