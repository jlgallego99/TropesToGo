package media

// RepositoryMedia defines an interface for all kinds of repositories of media tropes in TvTropes
// The interface allows us to implement multiple structs that handle different data formats like CSV or JSON
// sharing common methods
type RepositoryMedia interface {
	// AddMedia adds a new Media (Work with its Tropes) to the dataset
	AddMedia(Media) error

	// UpdateMedia updates a Media (Work with its Tropes) within the dataset
	UpdateMedia(Media) error
}

type CSVRepository struct {
}

type JSONRepository struct {
}
