package index

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
)

var (
	ErrPageNotFound = errors.New("the page was not found in the index repository")
)

// RepositoryIndex defines an interface for operations related to the Crawler
type RepositoryIndex interface {
	// Add a new crawled page to the Index
	AddPage(tropestogo.Page) error

	// UpdatePage updates a page in the Index based on if it's been updated or not
	UpdatePage(tropestogo.Page) error
}
