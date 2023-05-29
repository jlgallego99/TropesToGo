package csv_dataset_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCsvDataset(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CsvDataset Suite")
}
