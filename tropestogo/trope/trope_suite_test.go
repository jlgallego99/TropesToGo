package trope_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTrope(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trope Suite")
}
