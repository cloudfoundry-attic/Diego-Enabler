package ui

import (
	"fmt"

	"github.com/cloudfoundry/cli/cf/terminal"
)

type ListAppsCommand struct {
	Username     string
	Runtime      Runtime
	Organization string
	UI           terminal.UI
}

func (c *ListAppsCommand) BeforeAll() {
	if c.Organization == "" {
		fmt.Printf(
			"Getting apps on the %s runtime as %s...\n",
			terminal.EntityNameColor(c.Runtime.String()),
			terminal.EntityNameColor(c.Username),
		)
	} else {
		fmt.Printf(
			"Getting apps on the %s runtime in org %s as %s...\n",
			terminal.EntityNameColor(c.Runtime.String()),
			terminal.EntityNameColor(c.Organization),
			terminal.EntityNameColor(c.Username),
		)
	}
}

func (c *ListAppsCommand) AfterAll(apps []ApplicationPrinter) {
	sayOk()

	headers := []string{
		"name",
		"space",
		"org",
	}
	t := terminal.NewTable(c.UI, headers)

	for _, app := range apps {
		t.Add(app.Name(), app.Space(), app.Organization())
	}

	t.Print()
}
