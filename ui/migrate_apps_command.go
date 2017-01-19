package ui

import (
	"fmt"
	"io"

	"github.com/cloudfoundry/cli/cf/terminal"
)

type MigrateAppsCommand struct {
	Username     string
	Runtime      Runtime
	Organization string
	Space        string
}

func (c *MigrateAppsCommand) BeforeAll() {
	switch {
	case c.Organization != "" && c.Space != "":
		fmt.Printf(
			"Migrating apps to %s in org %s / %s as %s...\n",
			terminal.EntityNameColor(c.Runtime.String()),
			terminal.EntityNameColor(c.Organization),
			terminal.EntityNameColor(c.Space),
			terminal.EntityNameColor(c.Username),
		)
	case c.Organization != "":
		fmt.Printf(
			"Migrating apps to %s in org %s as %s...\n",
			terminal.EntityNameColor(c.Runtime.String()),
			terminal.EntityNameColor(c.Organization),
			terminal.EntityNameColor(c.Username),
		)
	default:
		fmt.Printf(
			"Migrating apps to %s as %s...\n",
			terminal.EntityNameColor(c.Runtime.String()),
			terminal.EntityNameColor(c.Username),
		)
	}
}

func (c *MigrateAppsCommand) BeforeEach(app ApplicationPrinter) {
	fmt.Println()
	fmt.Printf(
		"Migrating app %s in org %s / space %s to %s as %s...\n",
		terminal.EntityNameColor(app.Name()),
		terminal.EntityNameColor(app.Organization()),
		terminal.EntityNameColor(app.Space()),
		terminal.EntityNameColor(c.Runtime.String()),
		terminal.EntityNameColor(c.Username),
	)
}

func (c *MigrateAppsCommand) CompletedEach(app ApplicationPrinter) {
	fmt.Println()
	fmt.Printf(
		"Completed migrating app %s in org %s / space %s to %s as %s\n",
		terminal.EntityNameColor(app.Name()),
		terminal.EntityNameColor(app.Organization()),
		terminal.EntityNameColor(app.Space()),
		terminal.EntityNameColor(c.Runtime.String()),
		terminal.EntityNameColor(c.Username),
	)
}

func (c *MigrateAppsCommand) DuringEach(app ApplicationPrinter) {
	fmt.Print(".")
}

func (c *MigrateAppsCommand) AfterAll(attempts, warnings int, errors int) {
	successes := attempts - warnings - errors
	fmt.Println()
	fmt.Printf("Migration to %s completed: %d apps, %d errors, %d warnings\n", terminal.EntityNameColor(c.Runtime.String()), successes, errors, warnings)
}

func (c *MigrateAppsCommand) UserWarning(app ApplicationPrinter) {
	fmt.Printf(
		"WARNING: No authorization to migrate app %s to %s in space %s / org %s as %s\n",
		terminal.EntityNameColor(app.Name()),
		terminal.EntityNameColor(c.Runtime.String()),
		terminal.EntityNameColor(app.Space()),
		terminal.EntityNameColor(app.Organization()),
		terminal.EntityNameColor(c.Username),
	)
}

func (c *MigrateAppsCommand) FailMigrate(app ApplicationPrinter, err error) {
	fmt.Printf(
		"Error: Failed to migrate app %s to %s in space %s / org %s as %s: %s",
		terminal.EntityNameColor(app.Name()),
		terminal.EntityNameColor(c.Runtime.String()),
		terminal.EntityNameColor(app.Space()),
		terminal.EntityNameColor(app.Organization()),
		terminal.EntityNameColor(c.Username),
		terminal.EntityNameColor(err.Error()),
	)
}

func (c *MigrateAppsCommand) HealthCheckNoneWarning(app ApplicationPrinter, writer io.Writer) {
	fmt.Fprintf(
		writer,
		"WARNING: Assuming health check of type process ('none') for app with no mapped routes. Use 'cf set-health-check' to change this. App %s to %s in space %s / org %s as %s\n",
		terminal.EntityNameColor(app.Name()),
		terminal.EntityNameColor("Diego"),
		terminal.EntityNameColor(app.Space()),
		terminal.EntityNameColor(app.Organization()),
		terminal.EntityNameColor(c.Username),
	)
}
