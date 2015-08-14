package diego_support_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDiegoSupport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DiegoSupport Suite")
}
