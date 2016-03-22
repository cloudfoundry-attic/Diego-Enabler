package commands

import "github.com/cloudfoundry-incubator/diego-enabler/models"

type SpacesParser interface {
	Parse([]byte) (models.Spaces, error)
}

func Spaces(factory RequestFactory, client CloudControllerClient, spacesParser SpacesParser, pageParser PaginatedParser, spaceGuids []string) (models.Spaces, error) {
	var noSpaces models.Spaces

	return noSpaces, nil
}
