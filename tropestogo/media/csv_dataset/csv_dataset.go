package csv_dataset

import (
	"encoding/csv"
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ErrFileNotExists   = errors.New("CSV dataset file does not exist")
	ErrRecordNotExists = errors.New("there isn't a record with that media title in the CSV dataset")
)

type CSVRepository struct {
	sync.Mutex
	name      string
	delimiter rune
	writer    *csv.Writer
}

func NewCSVRepository(name string, delimiter rune) (*CSVRepository, error) {
	csvFile, err := os.Create(name + ".csv")
	writer := csv.NewWriter(csvFile)
	writer.Comma = delimiter

	repository := &CSVRepository{
		name:      name + ".csv",
		delimiter: delimiter,
		writer:    writer,
	}

	// Add headers to the CSV file
	repository.Lock()
	repository.writer.Write([]string{"title", "year", "lastupdated", "url", "mediatype", "tropes"})
	repository.writer.Flush()
	repository.Unlock()

	return repository, err
}

func (repository *CSVRepository) GetDelimiter() rune {
	return repository.delimiter
}

func (repository *CSVRepository) GetReader() (*csv.Reader, error) {
	dataset, err := os.Open(repository.name)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(dataset)
	reader.Comma = repository.delimiter
	return reader, nil
}

func (repository *CSVRepository) AddMedia(media media.Media) error {
	record := CreateMediaRecord(media)

	// Add record to the CSV file
	// Mutual exclusion access to the repository
	repository.Lock()
	err := repository.writer.Write(record)
	repository.writer.Flush()
	repository.Unlock()

	return err
}

// UpdateMedia updates a record in the CSV files by searching the title and year
func (repository *CSVRepository) UpdateMedia(title string, year string, media media.Media) error {
	reader, errReader := repository.GetReader()
	if errReader != nil {
		return errReader
	}

	records, errReadCSV := reader.ReadAll()
	if errReadCSV != nil {
		return errReadCSV
	}

	// Look for the line that holds the record that needs to be updated
	updateLine := -1
	for pos, record := range records {
		if record[0] == title && record[1] == year {
			updateLine = pos
			break
		}
	}

	// Update record
	input, _ := os.ReadFile(repository.name)
	lines := strings.Split(string(input), "\n")
	for linePos, _ := range lines {
		if linePos == updateLine {
			updatedRecord := CreateMediaRecord(media)
			lines[linePos] = strings.Join(updatedRecord, ",")

			break
		}
	}
	output := strings.Join(lines, "\n")
	errWrite := os.WriteFile("dataset.csv", []byte(output), 0644)
	if errWrite != nil {
		return errWrite
	}

	return nil
}

func (repository *CSVRepository) GetMedia(work tropestogo.Work) ([]media.Media, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *CSVRepository) GetMediaOfType(mediaType media.MediaType) ([]media.Media, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *CSVRepository) RemoveAll() error {
	if _, err := os.Stat("dataset.csv"); err == nil {
		csvFile, errRemove := os.Create(repository.name)
		repository.writer = csv.NewWriter(csvFile)
		repository.writer.Comma = repository.delimiter

		// Add headers to the CSV file
		repository.Lock()
		repository.writer.Write([]string{"title", "year", "lastupdated", "url", "mediatype", "tropes"})
		repository.writer.Flush()
		repository.Unlock()

		return errRemove
	} else {
		return ErrFileNotExists
	}
}

// CreateMediaRecord forms a string properly separated for inserting in a CSV file
func CreateMediaRecord(media media.Media) []string {
	var tropes []string
	for trope := range media.GetWork().Tropes {
		tropes = append(tropes, trope.GetTitle())
	}

	// A record consists of the following fields: title,year,lastupdated,url,mediatype,tropes
	record := []string{media.GetWork().Title, media.GetWork().Year, media.GetWork().LastUpdated.Format(time.DateTime),
		media.GetPage().URL.String(), media.GetMediaType().String(), strings.Join(tropes, ";")}

	return record
}
