package csv_dataset

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/jlgallego99/TropesToGo/media"
	"os"
	"strings"
)

var (
	ErrFileNotExists   = errors.New("CSV dataset file does not exist")
	ErrDuplicatedMedia = errors.New("duplicated media, the record already exists on the dataset")
	ErrReadCsv         = errors.New("error reading CSV file")
	ErrCreateCsv       = errors.New("error creating CSV file")
	ErrOpenCsv         = errors.New("error opening CSV file")
	ErrWriteCsv        = errors.New("error writing on the CSV file")
	ErrPersist         = errors.New("can't persist data on the CSV file because there's none")
)

type CSVRepository struct {
	name   string
	writer *csv.Writer
	data   []media.Media
}

// Error formats a generic error
func Error(message string, err error, subErr error) error {
	if subErr != nil {
		return fmt.Errorf("%w: "+message+"\n%w", err, subErr)
	} else {
		return fmt.Errorf("%w: "+message+"", err)
	}
}

func NewCSVRepository(name string) (*CSVRepository, error) {
	csvFile, err := os.Create(name + ".csv")
	if err != nil {
		return nil, Error(name, ErrCreateCsv, err)
	}

	writer := csv.NewWriter(csvFile)

	repository := &CSVRepository{
		name:   name + ".csv",
		writer: writer,
	}

	repository.writer.Write([]string{"title", "year", "lastupdated", "url", "mediatype", "tropes", "tropes_index"})
	repository.writer.Flush()

	return repository, nil
}

func (repository *CSVRepository) GetReader() (*csv.Reader, error) {
	dataset, err := os.Open(repository.name)
	if err != nil {
		return nil, Error(repository.name, ErrOpenCsv, err)
	}

	reader := csv.NewReader(dataset)
	return reader, nil
}

func (repository *CSVRepository) AddMedia(newMedia media.Media) error {
	for _, mediaData := range repository.data {
		if mediaData.GetWork().Title == newMedia.GetWork().Title && mediaData.GetWork().Year == newMedia.GetWork().Year {
			return Error("Title: "+newMedia.GetWork().Title, ErrDuplicatedMedia, nil)
		}
	}

	repository.data = append(repository.data, newMedia)

	return nil
}

// UpdateMedia updates a record in the CSV files by searching the title and year
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

func (repository *CSVRepository) RemoveAll() error {
	repository.data = []media.Media{}

	if _, err := os.Stat(repository.name); err == nil {
		csvFile, errRemove := os.Create(repository.name)
		if errRemove != nil {
			return Error(repository.name, errRemove, nil)
		}

		repository.writer = csv.NewWriter(csvFile)

		repository.writer.Write([]string{"title", "year", "lastupdated", "url", "mediatype", "tropes", "tropes_index"})
		repository.writer.Flush()

		return nil
	} else {
		pwd, _ := os.Getwd()

		return Error("at "+pwd+"/"+repository.name, ErrFileNotExists, nil)
	}
}

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

// CreateMediaRecord forms a string properly separated for inserting in a CSV file
func CreateMediaRecord(media media.Media) []string {
	var tropes []string
	var indexes []string
	for trope := range media.GetWork().Tropes {
		title := trope.GetTitle()
		index := trope.GetIndex().String()

		if title != "" && index != "" /*&& index != "UnknownTropeIndex"*/ {
			tropes = append(tropes, title)
			indexes = append(indexes, index)
		}
	}

	record := []string{media.GetWork().Title, media.GetWork().Year, media.GetWork().LastUpdated.Format("2006-01-02 15:04:05"),
		media.GetPage().GetUrl().String(), media.GetMediaType().String(), strings.Join(tropes, ";"), strings.Join(indexes, ";")}

	return record
}
