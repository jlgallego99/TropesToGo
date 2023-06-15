package json_dataset_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJsonDataset(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JsonDataset Suite")
}
