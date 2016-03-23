package ui

import (
	"fmt"

	"github.com/cloudfoundry/cli/cf/terminal"
)

type MigrateAppsCommand struct {
	Username     string
	Runtime      Runtime
	Organization string
}

func (c *MigrateAppsCommand) BeforeAll() {
	if c.Organization == "" {
		fmt.Printf(
			"Migrating apps to %s as %s...\n",
			terminal.EntityNameColor(c.Runtime.String()),
			terminal.EntityNameColor(c.Username),
		)
	} else {
		fmt.Printf(
			"Migrating apps to %s in org %s as %s...\n",
			terminal.EntityNameColor(c.Runtime.String()),
			terminal.EntityNameColor(c.Organization),
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
		"Completed migrating app %s in org %s / space %s to %s as %s...\n",
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

func (c *MigrateAppsCommand) AfterAll(attempts, warnings int) {
	fmt.Println()
	fmt.Printf("Migration to %s completed: %d apps, %d warnings\n", terminal.EntityNameColor(c.Runtime.String()), attempts, warnings)
}
