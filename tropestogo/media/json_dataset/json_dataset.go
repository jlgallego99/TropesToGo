package json_dataset

import (
	"encoding/json"
	"errors"
	"github.com/jlgallego99/TropesToGo/media"
	"os"
	"sort"
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
	_, err := os.Create(name + ".json")

	repository := &JSONRepository{
		name: name + ".json",
	}

	return repository, err
}

func (repository *JSONRepository) AddMedia(med media.Media) error {
	var dataset []media.JsonResponse

	var tropes []string
	for trope := range med.GetWork().Tropes {
		tropes = append(tropes, trope.GetTitle())
	}
	sort.Strings(tropes)
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
	json.Unmarshal(fileContents, &dataset)
	dataset = append(dataset, record)

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
	if _, err = os.Stat(repository.name); err == nil {
		_, err = os.Create(repository.name)

		return err
	} else {
		return ErrFileNotExists
	}
}
