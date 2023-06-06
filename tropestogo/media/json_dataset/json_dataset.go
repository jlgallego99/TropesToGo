package json_dataset

import (
	"github.com/jlgallego99/TropesToGo/media"
	"os"
	"sync"
)

type JSONRepository struct {
	sync.Mutex
	name string
}

func NewJSONRepository(name string) (*JSONRepository, error) {
	_, err := os.Create(name + ".json")

	repository := &JSONRepository{
		name: name + ".json",
	}

	return repository, err
}

func (repository *JSONRepository) AddMedia(media media.Media) error {
	//TODO implement me
	panic("implement me")
}

func (repository *JSONRepository) UpdateMedia(s string, s2 string, media media.Media) error {
	//TODO implement me
	panic("implement me")
}

func (repository *JSONRepository) RemoveAll() error {
	//TODO implement me
	panic("implement me")
}
