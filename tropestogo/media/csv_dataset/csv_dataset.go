package csv_dataset

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/jlgallego99/TropesToGo/media"
	"os"
	"strings"
	"time"
)

var (
	ErrFileNotExists   = errors.New("CSV dataset file does not exist")
	ErrDuplicatedMedia = errors.New("duplicated media, the record already exists on the dataset")
	ErrReadCsv         = errors.New("error reading CSV file")
	ErrCreateCsv       = errors.New("error creating CSV file")
	ErrOpenCsv         = errors.New("error opening CSV file")
	ErrWriteCsv        = errors.New("error writing on the CSV file")
	ErrPersist         = errors.New("can't persist data on the CSV file because there's none")
	ErrParseTime       = errors.New("error parsing the timestamp string from the dataset")
)

var Headers = []string{"title", "year", "lastupdated", "url", "mediatype", "tropes", "subtropes", "subtropes_namespaces"}

const timeLayout = "2006-01-02 15:04:05"

// CSVRepository implements the RepositoryMedia for creating and handling CSV datasets of all the scraped data on TvTropes
type CSVRepository struct {
	// name of the file dataset
	name string

	// writer for modifying the loaded CSV dataset that this repository manages
	writer *csv.Writer

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

// NewCSVRepository is the constructor for CSVRepository objects that handle CSV datasets
// It receives the name that the CSV dataset file will have and creates and empty file with only the column headers
// It will return an ErrCreateCsv error if the file couldn't be created
func NewCSVRepository(name string) (*CSVRepository, error) {
	// If the file doesn't exist, create it
	var csvFile *os.File
	var writer *csv.Writer
	if _, errStat := os.Stat(name + ".csv"); errStat != nil {
		var errCreate error
		csvFile, errCreate = os.Create(name + ".csv")
		if errCreate != nil {
			return nil, Error(name, ErrCreateCsv, errCreate)
		}

		writer = csv.NewWriter(csvFile)
		writer.Write(Headers)
		writer.Flush()
	} else {
		csvFile, _ = os.Open(name + ".csv")
		writer = csv.NewWriter(csvFile)
	}

	repository := &CSVRepository{
		name:   name + ".csv",
		writer: writer,
	}

	return repository, nil
}

// GetReader returns a new CSV reader object starting from the top of the file
// If the dataset file doesn't exist, it returns an ErrOpenCsv error
func (repository *CSVRepository) GetReader() (*csv.Reader, error) {
	dataset, err := os.Open(repository.name)
	if err != nil {
		return nil, Error(repository.name, ErrOpenCsv, err)
	}

	reader := csv.NewReader(dataset)
	return reader, nil
}

// AddMedia adds a newMedia Media object to the in-memory dataset, so it can be later persisted
// There can only be unique objects on the dataset, so it will return an ErrDuplicatedMedia error if the Media object already exists
func (repository *CSVRepository) AddMedia(newMedia media.Media) error {
	for _, mediaData := range repository.data {
		if mediaData.GetWork().Title == newMedia.GetWork().Title && mediaData.GetWork().Year == newMedia.GetWork().Year {
			return Error("Title: "+newMedia.GetWork().Title, ErrDuplicatedMedia, nil)
		}
	}

	repository.data = append(repository.data, newMedia)

	return nil
}

// UpdateMedia updates a media record already written on the dataset by checking if it has the same title and year, because that differentiates a record
// It returns an ErrReadCsv or ErrWriteCsv error if the dataset file couldn't be read or written
func (repository *CSVRepository) UpdateMedia(title string, year string, media media.Media) error {
	reader, errReader := repository.GetReader()
	if errReader != nil {
		return Error(repository.name, ErrReadCsv, errReader)
	}

	records, errReadAll := reader.ReadAll()
	if errReadAll != nil {
		return Error(repository.name, ErrReadCsv, errReadAll)
	}

	updateLine := -1
	for pos, record := range records {
		if record[0] == title && record[1] == year {
			updateLine = pos
			break
		}
	}

	input, _ := os.ReadFile(repository.name)
	lines := strings.Split(string(input), "\n")
	for linePos := range lines {
		if linePos == updateLine {
			updatedRecord := CreateMediaRecord(media)
			lines[linePos] = strings.Join(updatedRecord, ",")

			break
		}
	}
	output := strings.Join(lines, "\n")
	errWrite := os.WriteFile(repository.name, []byte(output), 0644)
	if errWrite != nil {
		return Error(repository.name, errWrite, nil)
	}

	return nil
}

// RemoveAll deletes all data on both the in-memory intermediate data and on the dataset file
// It tries to recreate the dataset, so it will return an ErrCreateCsv error if that wasn't possible
// If the dataset file doesn't exist, it returns an ErrFileNotExists error
func (repository *CSVRepository) RemoveAll() error {
	repository.data = []media.Media{}

	if _, err := os.Stat(repository.name); err == nil {
		csvFile, errRemove := os.Create(repository.name)
		if errRemove != nil {
			return Error(repository.name, errRemove, nil)
		}

		repository.writer = csv.NewWriter(csvFile)

		repository.writer.Write(Headers)
		repository.writer.Flush()

		return nil
	} else {
		pwd, _ := os.Getwd()

		return Error("at "+pwd+"/"+repository.name, ErrFileNotExists, nil)
	}
}

// Persist writes all intermediate Media data into the proper dataset file and empties the structure, because it has already been persisted
// It checks whether the new records are already on the dataset file, but doesn't return an error, but simply skips it
// If the internal data structure is empty, it will do nothing and return an ErrPersist error
// It returns an ErrReadCsv or ErrWriteCsv error if the dataset file couldn't be read or written
func (repository *CSVRepository) Persist() error {
	if len(repository.data) == 0 {
		return Error(repository.name, ErrPersist, nil)
	}

	reader, errReader := repository.GetReader()
	if errReader != nil {
		return Error(repository.name, ErrReadCsv, errReader)
	}

	records, errReadAll := reader.ReadAll()
	if errReadAll != nil {
		return Error(repository.name, ErrReadCsv, errReadAll)
	}

	for _, mediaData := range repository.data {
		exists := false
		for _, record := range records {
			if record[0] == mediaData.GetWork().Title && record[1] == mediaData.GetWork().Year {
				exists = true
				break
			}
		}

		if !exists {
			record := CreateMediaRecord(mediaData)

			err := repository.writer.Write(record)
			if err != nil {
				return Error(repository.name, ErrWriteCsv, err)
			}
			repository.writer.Flush()
		}
	}

	repository.data = []media.Media{}

	return nil
}

// CreateMediaRecord forms a proper string record from a Media object for inserting in a CSV file
// Each value on the returned array is a column value for the CSV file
func CreateMediaRecord(media media.Media) []string {
	var tropes []string
	var subTropes []string
	var subTropesNamespaces []string
	//var indexes []string

	for trope := range media.GetWork().Tropes {
		title := trope.GetTitle()
		index := trope.GetIndex().String()

		if title != "" && index != "" /*&& index != "UnknownTropeIndex"*/ {
			tropes = append(tropes, title)
			//indexes = append(indexes, index)
		}
	}

	for subTrope := range media.GetWork().SubTropes {
		title := subTrope.GetTitle()
		namespace := subTrope.GetSubpage()

		if title != "" && namespace != "" {
			subTropes = append(subTropes, title)
			subTropesNamespaces = append(subTropesNamespaces, namespace)
		}
	}

	record := []string{media.GetWork().Title, media.GetWork().Year, media.GetWork().LastUpdated.Format(timeLayout),
		media.GetPage().GetUrl().String(), media.GetMediaType().String(), strings.Join(tropes, ";"),
		strings.Join(subTropes, ";"), strings.Join(subTropesNamespaces, ";")}

	return record
}

// GetWorkPages retrieves all persisted Work urls on the CSV dataset and the last time they were updated
// Returns a map relating page URLs to the last time they were updated
func (repository *CSVRepository) GetWorkPages() (map[string]time.Time, error) {
	datasetPages := make(map[string]time.Time, 0)

	reader, errReader := repository.GetReader()
	if errReader != nil {
		return nil, Error(repository.name, ErrReadCsv, errReader)
	}

	records, errReadAll := reader.ReadAll()
	if errReadAll != nil {
		return nil, Error(repository.name, ErrReadCsv, errReadAll)
	}

	// Only iterate from the second row onwards (ignoring the first row, the headers)
	records = append(records[:0], records[1:]...)
	for _, record := range records {
		lastUpdated, errLastUpdated := time.Parse(timeLayout, record[2])
		if errLastUpdated != nil {
			return nil, Error(repository.name, ErrParseTime, errLastUpdated)
		}

		datasetPages[record[3]] = lastUpdated
	}

	return datasetPages, nil
}
