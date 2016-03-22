package commands

import (
	"github.com/cloudfoundry-incubator/diego-enabler/api"
	"github.com/cloudfoundry-incubator/diego-enabler/models"
)

//go:generate counterfeiter . SpacesParser
type SpacesParser interface {
	Parse([]byte) (models.Spaces, error)
}

func Spaces(requestFactory RequestFactory, client CloudControllerClient, spacesParser SpacesParser, pageParser PaginatedParser) (models.Spaces, error) {
	var noSpaces models.Spaces

	filter := api.Filters{}

	params := map[string]interface{}{
		"inline-relations-depth": 1,
	}

	responseBodies, err := paginatedRequester(requestFactory, filter, params, client, pageParser)
	if err != nil {
		return noSpaces, err
	}

	var spaces models.Spaces

	for _, nextBody := range responseBodies {
		apps, err := spacesParser.Parse(nextBody)
		if err != nil {
			return noSpaces, err
		}

		spaces = append(spaces, apps...)
	}

	return spaces, nil
}
