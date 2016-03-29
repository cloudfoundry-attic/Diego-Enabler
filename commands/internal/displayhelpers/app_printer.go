package displayhelpers

import "github.com/cloudfoundry-incubator/diego-enabler/models"

type AppPrinter struct {
	App    models.Application
	Spaces map[string]models.Space
}

func (a *AppPrinter) Name() string {
	return a.App.Name
}

func (a *AppPrinter) Organization() string {
	spaces := a.Spaces
	app := a.App

	if len(spaces) == 0 {
		return ""
	}

	space, ok := spaces[app.SpaceGuid]
	if !ok {
		return ""
	}

	if space.Organization.Name != "" {
		return space.Organization.Name
	}

	return space.Organization.Guid
}

func (a *AppPrinter) Space() string {
	var display string
	spaces := a.Spaces
	app := a.App

	if len(spaces) == 0 {
		display = app.SpaceGuid
	} else {
		space, ok := spaces[app.SpaceGuid]
		if ok {
			display = space.Name
		} else {
			display = app.SpaceGuid
		}
	}

	return display
}
