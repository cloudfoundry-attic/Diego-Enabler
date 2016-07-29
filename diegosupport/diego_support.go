package diegosupport

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

//go:generate counterfeiter . CliConnection
type CliConnection interface {
	CliCommandWithoutTerminalOutput(args ...string) ([]string, error)
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
