package tvtropespages_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTvTropesPages(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TvTropesPages Suite")
}
