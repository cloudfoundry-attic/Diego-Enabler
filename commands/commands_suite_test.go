package commands_test

import (
	"github.com/cloudfoundry-incubator/diego-enabler/api/apifakes"
	"github.com/cloudfoundry-incubator/diego-enabler/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var (
	fakeConnection *apifakes.FakeConnection
)

var _ = BeforeEach(func() {
	fakeConnection = new(apifakes.FakeConnection)
	commands.DiegoEnabler.CLIConnection = fakeConnection
})
