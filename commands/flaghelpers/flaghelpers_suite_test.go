package flaghelpers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFlaghelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flaghelpers Suite")
}
