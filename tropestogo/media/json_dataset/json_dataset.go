package json_dataset

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jlgallego99/TropesToGo/media"
	"os"
)

var (
	ErrFileNotExists   = errors.New("CSV dataset file does not exist")
	ErrDuplicatedMedia = errors.New("duplicated media, the record already exists on the dataset")
	ErrReadJson        = errors.New("error reading JSON file")
	ErrCreateJson      = errors.New("error creating JSON file")
	ErrOpenJson        = errors.New("error opening JSON file")
	ErrWriteJson       = errors.New("error writing on the JSON file")
	ErrUnmarshalJson   = errors.New("error unmarshalling JSON file")
	ErrMarshalJson     = errors.New("error marshalling JSON")
)

type JSONDataset struct {
	Tropestogo []media.JsonResponse `json:"tropestogo"`
}

type JSONRepository struct {
	name string
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

	fileContents, errReadDataset := os.ReadFile(repository.name)
	if errReadDataset != nil {
		return Error(repository.name, ErrReadJson, errReadDataset)
	}

	// Get the JSON array and append the new Media object
	errUnmarshal := json.Unmarshal(fileContents, &dataset)
	if errUnmarshal != nil {
		return Error(repository.name, ErrUnmarshalJson, errUnmarshal)
	}

	// Append the Media only if it doesn't exist yet on the dataset
	for _, datasetMedia := range dataset.Tropestogo {
		if datasetMedia.Title == med.GetWork().Title && datasetMedia.Year == med.GetWork().Year {
			return Error("%w Title: "+med.GetWork().Title, ErrDuplicatedMedia, nil)
		}
	}

	dataset.Tropestogo = append(dataset.Tropestogo, record)
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

func (repository *JSONRepository) UpdateMedia(title string, year string, med media.Media) error {
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
