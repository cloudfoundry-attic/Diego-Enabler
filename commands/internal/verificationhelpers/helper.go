package verificationhelpers

import (
	"errors"

	"github.com/cloudfoundry/cli/plugin"
)

func VerifyLoggedIn(cliCon plugin.CliConnection) error {
	var result error

	if connected, err := cliCon.IsLoggedIn(); !connected {
		result = NotLoggedInError

		if err != nil {
			result = err
		}
	}

	return result
}

var NotLoggedInError = errors.New("You must be logged in")
