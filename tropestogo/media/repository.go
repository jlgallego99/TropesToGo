package media

// RepositoryMedia defines an interface for all kinds of repositories of media tropes in TvTropes
// The interface allows us to implement multiple structs that handle different data formats like CSV or JSON
// sharing common methods
type RepositoryMedia interface {
	// AddMedia adds a new Media (Work with its Tropes) to the dataset
	AddMedia(Media) error

	// UpdateMedia updates a Media (Work with its Tropes) within the dataset
	// It distinguishes between works with the same name by both its title and year
	UpdateMedia(string, string, Media) error

	// RemoveAll delete all Media entries on the repository
	RemoveAll() error

	// Persist adds all repository Media objects to the proper dataset
	Persist() error
}
