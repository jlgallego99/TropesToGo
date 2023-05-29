package csv_dataset

import "sync"

type CSVRepository struct {
	sync.Mutex
}
