package index

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
)

var (
	ErrURLNotFound = errors.New("the URL does not exist")
)

// Index holds the logic of all discovered pages by the crawler
type Index struct {
	pages *[]tropestogo.Page
}

func NewIndex() (Index, error) {
	return Index{}, nil
}
