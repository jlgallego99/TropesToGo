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

func (repository *JSONRepository) UpdateMedia(title string, year string, med media.Media) error {
	var dataset JSONDataset

	fileContents, errReadDataset := os.ReadFile("dataset.json")
	if errReadDataset != nil {
		return errReadDataset
	}

	// Get the dataset on structs
	errUnmarshal := json.Unmarshal(fileContents, &dataset)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	// Look for the record that needs to be updated
	for pos, record := range dataset.Tropestogo {
		if record.Title == title && record.Year == year {
			tropes := media.GetJsonTropes(med)
			dataset.Tropestogo[pos].Title = med.GetWork().Title
			dataset.Tropestogo[pos].Year = med.GetWork().Year
			dataset.Tropestogo[pos].MediaType = med.GetMediaType().String()
			dataset.Tropestogo[pos].LastUpdated = med.GetWork().LastUpdated.Format(time.DateTime)
			dataset.Tropestogo[pos].URL = med.GetPage().URL.String()
			dataset.Tropestogo[pos].Tropes = tropes

			break
		}
	}

	// Update the record and marshal to the file
	jsonBytes, err := json.Marshal(dataset)
	if err != nil {
		return err
	}

	return os.WriteFile("dataset.json", jsonBytes, 0644)
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
