package json_dataset

import (
	"encoding/json"
	"errors"
	"github.com/jlgallego99/TropesToGo/media"
	"os"
	"sync"
	"time"
)

var (
	ErrFileNotExists = errors.New("CSV dataset file does not exist")
)

type JSONDataset struct {
	Tropestogo []media.JsonResponse `json:"tropestogo"`
}

type JSONRepository struct {
	sync.Mutex
	name string
}

func NewJSONRepository(name string) (*JSONRepository, error) {
	f, err := os.Create(name + ".json")
	f.WriteString("{\"tropestogo\": []}")

	repository := &JSONRepository{
		name: name + ".json",
	}

	return repository, err
}

func (repository *JSONRepository) AddMedia(med media.Media) error {
	var dataset JSONDataset

	tropes := media.GetJsonTropes(med)
	record := media.JsonResponse{
		Title:       med.GetWork().Title,
		Year:        med.GetWork().Year,
		MediaType:   med.GetMediaType().String(),
		LastUpdated: med.GetWork().LastUpdated.Format(time.DateTime),
		URL:         med.GetPage().URL.String(),
		Tropes:      tropes,
	}

	fileContents, errReadDataset := os.ReadFile("dataset.json")
	if errReadDataset != nil {
		return errReadDataset
	}

	// Get the JSON array and append the new Media object
	errUnmarshal := json.Unmarshal(fileContents, &dataset)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	dataset.Tropestogo = append(dataset.Tropestogo, record)

	jsonBytes, err := json.Marshal(dataset)
	if err != nil {
		return err
	}

	return os.WriteFile("dataset.json", jsonBytes, 0644)
}

func (repository *JSONRepository) UpdateMedia(s string, s2 string, media media.Media) error {
	//TODO implement me
	panic("implement me")
}

func (repository *JSONRepository) RemoveAll() error {
	var err error
	var f *os.File
	if _, err = os.Stat(repository.name); err == nil {
		f, err = os.Create(repository.name)
		if err != nil {
			return err
		}

		f.WriteString("{\"tropestogo\": []}")
		return nil
	} else {
		return ErrFileNotExists
	}
}
