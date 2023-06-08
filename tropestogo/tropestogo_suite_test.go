package tropestogo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTropestogo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tropestogo Suite")
}
