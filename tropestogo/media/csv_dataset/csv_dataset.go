package csv_dataset

import (
	"encoding/csv"
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"github.com/jlgallego99/TropesToGo/media"
	"os"
	"strings"
	"sync"
)

var (
	ErrFileNotExists = errors.New("CSV dataset file does not exist")
)

type CSVRepository struct {
	sync.Mutex
	name      string
	delimiter rune
	reader    *csv.Reader
	writer    *csv.Writer
}

func NewCSVRepository(name string, delimiter rune) (*CSVRepository, error) {
	csvFile, err := os.Create(name + ".csv")
	reader := csv.NewReader(csvFile)
	writer := csv.NewWriter(csvFile)
	reader.Comma = delimiter
	writer.Comma = delimiter

	repository := &CSVRepository{
		name:      name + ".csv",
		delimiter: delimiter,
		reader:    reader,
		writer:    writer,
	}

	return repository, err
}

func (repository *CSVRepository) GetDelimiter() rune {
	return repository.delimiter
}

func (repository *CSVRepository) AddMedia(media media.Media) error {
	var tropes []string
	for trope := range media.GetWork().Tropes {
		tropes = append(tropes, trope.GetTitle())
	}

	repository.Lock()
	err := repository.writer.Write([]string{media.GetWork().Title, strings.Join(tropes, ";")})
	repository.writer.Flush()
	repository.Unlock()

	return err
}

func (repository *CSVRepository) UpdateMedia(media media.Media) error {
	//TODO implement me
	panic("implement me")
}

func (repository *CSVRepository) GetMedia(work tropestogo.Work) ([]media.Media, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *CSVRepository) GetMediaOfType(mediaType media.MediaType) ([]media.Media, error) {
	//TODO implement me
	panic("implement me")
}
