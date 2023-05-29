package json_dataset

import "sync"

type JSONRepository struct {
	sync.Mutex
}
