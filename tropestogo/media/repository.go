package media

import (
	"github.com/google/uuid"
)

// RepositoryMedia defines an interface for all kinds of repositories of media tropes in TvTropes
// The interface allows us to implement multiple structs that handle different data formats like CSV or JSON
// sharing common methods
type RepositoryMedia interface {
	GetAll() ([]Media, error)
	Get(uuid.UUID) (Media, error)
	Add(Media) error
	Update(Media) error
}

type CSVRepository struct {
}

type JSONRepository struct {
}
