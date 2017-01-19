package diegohelpers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDiegohelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Diegohelpers Suite")
}
