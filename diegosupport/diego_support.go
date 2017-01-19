package diegosupport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin/models"
)

//go:generate counterfeiter . CliConnection
type CliConnection interface {
	CliCommandWithoutTerminalOutput(args ...string) ([]string, error)
	GetApp(string) (plugin_models.GetAppModel, error)
	GetCurrentSpace() (plugin_models.Space, error)
	GetSpace(string) (plugin_models.GetSpace_Model, error)
	Username() (string, error)
}

type DiegoSupport struct {
	cli CliConnection
}

type diegoError struct {
	Code        int64  `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	ErrorCode   string `json:"error_code,omitempty"`
}

func NewDiegoSupport(cli CliConnection) *DiegoSupport {
	return &DiegoSupport{
		cli: cli,
	}
}

func (d *DiegoSupport) SetDiegoFlag(appGuid string, enable bool) ([]string, error) {
	output, err := d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid, "-X", "PUT", "-d", `{"diego":`+strconv.FormatBool(enable)+`}`)
	if err != nil {
		return output, err
	}

	if err = checkDiegoError(strings.Join(output, "")); err != nil {
		return output, err
	}

	return output, nil
}

func checkDiegoError(jsonRsp string) error {
	b := []byte(jsonRsp)
	diegoErr := diegoError{}
	err := json.Unmarshal(b, &diegoErr)
	if err != nil {
		return err
	}

	if diegoErr.ErrorCode != "" || diegoErr.Code != 0 {
		return errors.New(diegoErr.ErrorCode + " - " + diegoErr.Description)
	}

	return nil
}

func (d *DiegoSupport) WarnNoRoutes(appName string, output io.Writer) error {
	hasRoutes, err := d.HasRoutes(appName)
	if err != nil {
		return err
	}

	if hasRoutes {
		return nil
	}

	// Couldn't find a better way to get from the app to the space.
	// The app doesn't known the space name and there is no
	// cliConnection.GetSpace that accepts a space GUID.
	currentSpace, err := d.cli.GetCurrentSpace()
	if err != nil {
		return err
	}

	space, err := d.cli.GetSpace(currentSpace.Name)
	if err != nil {
		return err
	}

	username, err := d.cli.Username()
	if err != nil {
		return err
	}

	fmt.Fprintf(
		output,
		"WARNING: Assuming health check of type process ('none') for app with no mapped routes. Use 'cf set-health-check' to change this. App %s to %s in space %s / org %s as %s\n",
		terminal.EntityNameColor(appName),
		terminal.EntityNameColor("Diego"),
		terminal.EntityNameColor(space.Name),
		terminal.EntityNameColor(space.Organization.Name),
		terminal.EntityNameColor(username),
	)

	return nil
}

func (d *DiegoSupport) HasRoutes(appName string) (bool, error) {
	app, err := d.cli.GetApp(appName)
	if err != nil {
		return false, err
	}

	return len(app.Routes) > 0, nil
}
