package media

import (
	tropestogo "github.com/jlgallego99/TropesToGo"
)

// RepositoryMedia defines an interface for all kinds of repositories of media tropes in TvTropes
// The interface allows us to implement multiple structs that handle different data formats like CSV or JSON
// sharing common methods
type RepositoryMedia interface {
	// AddMedia adds a new Media (Work with its Tropes) to the dataset
	AddMedia(Media) error

	// UpdateMedia updates a Media (Work with its Tropes) within the dataset
	UpdateMedia(Media) error

	// GetTvTropes returns all Media found in TvTropes
	GetTvTropes() ([]Media, error)

	// GetMedia returns a Work with its Tropes
	GetMedia(tropestogo.Work) ([]Media, error)

	// GetMedia returns all Media within a MediaType (for example, all films)
	GetMediaType(MediaType) ([]Media, error)
}
