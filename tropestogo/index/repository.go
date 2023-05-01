package index

import (
	"errors"
	"github.com/google/uuid"
	tropestogo "github.com/jlgallego99/TropesToGo"
)

var (
	ErrPageNotFound = errors.New("the page was not found in the index repository")
)

// RepositoryIndex defines an interface for operations within the crawler indexing of TvTropes
type RepositoryIndex interface {
	GetAll() ([]tropestogo.Page, error)
	Get(uuid.UUID) (tropestogo.Page, error)
	Add(tropestogo.Page) error
	Update(tropestogo.Page) error
}
