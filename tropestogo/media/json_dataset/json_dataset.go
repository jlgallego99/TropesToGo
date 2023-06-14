package json_dataset

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jlgallego99/TropesToGo/media"
	"os"
)

var (
	ErrFileNotExists   = errors.New("JSON dataset file does not exist")
	ErrDuplicatedMedia = errors.New("duplicated media, the record already exists on the dataset")
	ErrReadJson        = errors.New("error reading JSON file")
	ErrCreateJson      = errors.New("error creating JSON file")
	ErrOpenJson        = errors.New("error opening JSON file")
	ErrWriteJson       = errors.New("error writing on the JSON file")
	ErrUnmarshalJson   = errors.New("error unmarshalling JSON file")
	ErrMarshalJson     = errors.New("error marshalling JSON")
	ErrPersist         = errors.New("can't persist data on the JSON file because there's none")
)

type JSONDataset struct {
	Tropestogo []media.JsonResponse `json:"tropestogo"`
}

type JSONRepository struct {
	name string
	data []media.Media
}

// Error formats a generic error
func Error(message string, err error, subErr error) error {
	if subErr != nil {
		return fmt.Errorf("%w: "+message+"\n%w", err, subErr)
	} else {
		return fmt.Errorf("%w: "+message+"", err)
	}
}

func NewJSONRepository(name string) (*JSONRepository, error) {
	f, err := os.Create(name + ".json")
	if err != nil {
		return nil, Error(name, ErrCreateJson, err)
	}

	f.WriteString("{\"tropestogo\": []}")

	repository := &JSONRepository{
		name: name + ".json",
	}

	return repository, nil
}

func (repository *JSONRepository) AddMedia(newMedia media.Media) error {
	// Check if the new Media is a duplicate or not by checking its title and year
	for _, mediaData := range repository.data {
		if mediaData.GetWork().Title == newMedia.GetWork().Title && mediaData.GetWork().Year == newMedia.GetWork().Year {
			return Error("Title: "+newMedia.GetWork().Title, ErrDuplicatedMedia, nil)
		}
	}

	// Add Media to the repository in memory
	repository.data = append(repository.data, newMedia)

	return nil
}

func (repository *JSONRepository) UpdateMedia(title string, year string, updateMedia media.Media) error {
	var dataset JSONDataset

	fileContents, errReadDataset := os.ReadFile(repository.name)
	if errReadDataset != nil {
		return Error(repository.name, ErrReadJson, errReadDataset)
	}

	// Get the dataset on structs
	errUnmarshal := json.Unmarshal(fileContents, &dataset)
	if errUnmarshal != nil {
		return Error(repository.name, ErrUnmarshalJson, errUnmarshal)
	}

	// Look for the record that needs to be updated
	for pos, record := range dataset.Tropestogo {
		if record.Title == title && record.Year == year {
			tropes := media.GetJsonTropes(updateMedia)
			dataset.Tropestogo[pos].Title = updateMedia.GetWork().Title
			dataset.Tropestogo[pos].Year = updateMedia.GetWork().Year
			dataset.Tropestogo[pos].MediaType = updateMedia.GetMediaType().String()
			dataset.Tropestogo[pos].LastUpdated = updateMedia.GetWork().LastUpdated.Format("2006-01-02 15:04:05")
			dataset.Tropestogo[pos].URL = updateMedia.GetPage().GetUrl().String()
			dataset.Tropestogo[pos].Tropes = tropes

			break
		}
	}

	// Update the record and marshal to the file
	jsonBytes, err := json.Marshal(dataset)
	if err != nil {
		return Error("", ErrMarshalJson, err)
	}

	errWriteFile := os.WriteFile(repository.name, jsonBytes, 0644)
	if errWriteFile != nil {
		return Error(repository.name, ErrWriteJson, errWriteFile)
	}

	return nil
}

func (repository *JSONRepository) RemoveAll() error {
	var err error
	var f *os.File

	// Empty the in-memory data
	repository.data = []media.Media{}

	if _, err = os.Stat(repository.name); err == nil {
		f, err = os.Create(repository.name)
		if err != nil {
			return Error(repository.name, ErrCreateJson, err)
		}

		f.WriteString("{\"tropestogo\": []}")
		return nil
	} else {
		pwd, _ := os.Getwd()

		return Error("at "+pwd+"/"+repository.name, ErrFileNotExists, nil)
	}
}

func (repository *JSONRepository) Persist() error {
	if len(repository.data) == 0 {
		return Error(repository.name, ErrPersist, nil)
	}

	var dataset JSONDataset

	fileContents, errReadDataset := os.ReadFile(repository.name)
	if errReadDataset != nil {
		return Error(repository.name, ErrReadJson, errReadDataset)
	}

	// Get the JSON array and append the new Media object
	errUnmarshal := json.Unmarshal(fileContents, &dataset)
	if errUnmarshal != nil {
		return Error(repository.name, ErrUnmarshalJson, errUnmarshal)
	}

	for _, mediaData := range repository.data {
		// Search if the value already exists on the dataset
		exists := false
		for _, datasetMedia := range dataset.Tropestogo {
			if datasetMedia.Title == mediaData.GetWork().Title && datasetMedia.Year == mediaData.GetWork().Year {
				exists = true
				break
			}
		}

		// Append without repeating
		if !exists {
			tropes := media.GetJsonTropes(mediaData)
			record := media.JsonResponse{
				Title:       mediaData.GetWork().Title,
				Year:        mediaData.GetWork().Year,
				MediaType:   mediaData.GetMediaType().String(),
				LastUpdated: mediaData.GetWork().LastUpdated.Format("2006-01-02 15:04:05"),
				URL:         mediaData.GetPage().GetUrl().String(),
				Tropes:      tropes,
			}

			dataset.Tropestogo = append(dataset.Tropestogo, record)
		}
	}

	// Empty in-memory data, because it has been persisted on the dataset file
	repository.data = []media.Media{}

	// Persist the old data + the new data in the JSON file
	jsonBytes, err := json.Marshal(dataset)
	if err != nil {
		return Error("", ErrMarshalJson, err)
	}

	errWriteFile := os.WriteFile(repository.name, jsonBytes, 0644)
	if errWriteFile != nil {
		return Error(repository.name, ErrWriteJson, errWriteFile)
	}

	return nil
}
