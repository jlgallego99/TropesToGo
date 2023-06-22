package crawler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCrawler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crawler Suite")
}
