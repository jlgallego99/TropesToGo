package json_dataset

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jlgallego99/TropesToGo/media"
	"os"
	"time"
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
	ErrParseTime       = errors.New("error parsing the timestamp string from the dataset")
)

const timeLayout = "2006-01-02 15:04:05"

// JSONDataset is an intermediate structure for marshaling/unmarshalling data from the JSON dataset
type JSONDataset struct {
	Tropestogo []media.JsonResponse `json:"tropestogo"`
}

// JSONRepository implements the RepositoryMedia for creating and handling JSON datasets of all the scraped data on TvTropes
// It has an internal data structure of Media objects for better performance that can be persisted into a file all in one go
type JSONRepository struct {
	// name of the file dataset
	name string

	// data is the intermediate dataset added here before persisting it all at once
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

// NewJSONRepository is the constructor for JSONRepository objects that handle JSON datasets
// It receives the name that the JSON dataset file will have and creates the file with a "tropestogo" key with an empty array
// It will return an ErrCreateJson error if the file couldn't be created
func NewJSONRepository(name string) (*JSONRepository, error) {
	// If the file doesn't exist, create it
	if _, errStat := os.Stat(name + ".json"); errStat != nil {
		f, errCreate := os.Create(name + ".json")
		if errCreate != nil {
			return nil, Error(name, ErrCreateJson, errCreate)
		}

		f.WriteString("{\"tropestogo\": []}")
	}

	repository := &JSONRepository{
		name: name + ".json",
	}

	return repository, nil
}

// AddMedia adds a newMedia Media object to the in-memory dataset, so it can be later persisted
// There can only be unique objects on the dataset, so it will return an ErrDuplicatedMedia error if the Media object already exists
func (repository *JSONRepository) AddMedia(newMedia media.Media) error {
	for _, mediaData := range repository.data {
		if mediaData.GetWork().Title == newMedia.GetWork().Title && mediaData.GetWork().Year == newMedia.GetWork().Year {
			return Error("Title: "+newMedia.GetWork().Title, ErrDuplicatedMedia, nil)
		}
	}

	repository.data = append(repository.data, newMedia)

	return nil
}

// UpdateMedia updates a record already written on the dataset by checking if it has the same title and year, because that differentiates a record
// It returns an ErrReadJson, ErrWriteJson or an ErrUnmarshalJson error if the dataset couldn't be read, written or unmarshalled into a internal structure
func (repository *JSONRepository) UpdateMedia(title string, year string, updateMedia media.Media) error {
	var dataset JSONDataset

	fileContents, errReadDataset := os.ReadFile(repository.name)
	if errReadDataset != nil {
		return Error(repository.name, ErrReadJson, errReadDataset)
	}

	errUnmarshal := json.Unmarshal(fileContents, &dataset)
	if errUnmarshal != nil {
		return Error(repository.name, ErrUnmarshalJson, errUnmarshal)
	}

	for pos, record := range dataset.Tropestogo {
		if record.Title == title && record.Year == year {
			tropes, subTropes := media.GetJsonTropes(updateMedia)
			dataset.Tropestogo[pos].Title = updateMedia.GetWork().Title
			dataset.Tropestogo[pos].Year = updateMedia.GetWork().Year
			dataset.Tropestogo[pos].MediaType = updateMedia.GetMediaType().String()
			dataset.Tropestogo[pos].LastUpdated = formatDate(updateMedia.GetWork().LastUpdated)
			dataset.Tropestogo[pos].URL = updateMedia.GetPage().GetUrl().String()
			dataset.Tropestogo[pos].Tropes = tropes
			dataset.Tropestogo[pos].SubTropes = subTropes

			break
		}
	}

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

// RemoveAll deletes all data on both the in-memory intermediate data and on the dataset file
// It tries to recreate the dataset, so it will return an ErrCreateJson error if that wasn't possible
// If the dataset file doesn't exist, it returns an ErrFileNotExists error
func (repository *JSONRepository) RemoveAll() error {
	var err error
	var f *os.File

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

// Persist writes all intermediate Media data into the proper dataset file and empties the structure, because it has already been persisted
// It checks whether the new records are already on the dataset file, but doesn't return an error, but simply skips it
// If the internal data structure is empty, it will do nothing and return an ErrPersist error
// It returns an ErrReadJson, ErrWriteJson or an ErrUnmarshalJson error if the dataset couldn't be read, written or unmarshalled into a internal structure
func (repository *JSONRepository) Persist() error {
	if len(repository.data) == 0 {
		return Error(repository.name, ErrPersist, nil)
	}

	var dataset JSONDataset

	fileContents, errReadDataset := os.ReadFile(repository.name)
	if errReadDataset != nil {
		return Error(repository.name, ErrReadJson, errReadDataset)
	}

	errUnmarshal := json.Unmarshal(fileContents, &dataset)
	if errUnmarshal != nil {
		return Error(repository.name, ErrUnmarshalJson, errUnmarshal)
	}

	for _, mediaData := range repository.data {
		exists := false
		for _, datasetMedia := range dataset.Tropestogo {
			if datasetMedia.Title == mediaData.GetWork().Title && datasetMedia.Year == mediaData.GetWork().Year {
				exists = true
				break
			}
		}

		if !exists {
			tropes, subTropes := media.GetJsonTropes(mediaData)
			record := media.JsonResponse{
				Title:       mediaData.GetWork().Title,
				Year:        mediaData.GetWork().Year,
				MediaType:   mediaData.GetMediaType().String(),
				LastUpdated: formatDate(mediaData.GetWork().LastUpdated),
				URL:         mediaData.GetPage().GetUrl().String(),
				Tropes:      tropes,
				SubTropes:   subTropes,
			}

			dataset.Tropestogo = append(dataset.Tropestogo, record)
		}
	}

	repository.data = []media.Media{}

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

// GetWorkPages retrieves all persisted Work urls on the JSON dataset and the last time they were updated
// Returns a map relating page URLs to the last time they were updated
func (repository *JSONRepository) GetWorkPages() (map[string]time.Time, error) {
	var dataset JSONDataset
	datasetPages := make(map[string]time.Time, 0)

	fileContents, errReadDataset := os.ReadFile(repository.name)
	if errReadDataset != nil {
		return nil, Error(repository.name, ErrReadJson, errReadDataset)
	}

	errUnmarshal := json.Unmarshal(fileContents, &dataset)
	if errUnmarshal != nil {
		return nil, Error(repository.name, ErrUnmarshalJson, errUnmarshal)
	}

	for _, record := range dataset.Tropestogo {
		lastUpdated, errLastUpdated := time.Parse(timeLayout, record.LastUpdated)
		if errLastUpdated != nil {
			return nil, Error(repository.name, ErrParseTime, errLastUpdated)
		}

		datasetPages[record.URL] = lastUpdated
	}

	return datasetPages, nil
}

// formatDate transforms a date to a unified string format across all datasets
func formatDate(date time.Time) string {
	return date.Format(timeLayout)
}
