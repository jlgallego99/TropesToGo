package json_dataset

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/jlgallego99/TropesToGo/media"
)

var (
	ErrFileNotExists   = errors.New("CSV dataset file does not exist")
	ErrDuplicatedMedia = errors.New("duplicated media: the record already exists on the dataset")
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
		LastUpdated: med.GetWork().LastUpdated.Format("2006-01-02 15:04:05"),
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

	// Append the Media only if it doesn't exist yet on the dataset
	for _, datasetMedia := range dataset.Tropestogo {
		if datasetMedia.Title == med.GetWork().Title && datasetMedia.Year == med.GetWork().Year {
			return ErrDuplicatedMedia
		}
	}

	dataset.Tropestogo = append(dataset.Tropestogo, record)
	jsonBytes, err := json.Marshal(dataset)
	if err != nil {
		return err
	}

	repository.Lock()
	errWriteFile := os.WriteFile("dataset.json", jsonBytes, 0644)
	repository.Unlock()

	return errWriteFile
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
			dataset.Tropestogo[pos].LastUpdated = med.GetWork().LastUpdated.Format("2006-01-02 15:04:05")
			dataset.Tropestogo[pos].URL = med.GetPage().URL.String()
			dataset.Tropestogo[pos].Tropes = tropes

			break
		}
	}

	// Update the record and marshal to the file
	repository.Lock()
	jsonBytes, err := json.Marshal(dataset)
	if err != nil {
		return err
	}

	errWriteFile := os.WriteFile("dataset.json", jsonBytes, 0644)
	repository.Unlock()

	return errWriteFile
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