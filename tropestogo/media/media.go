package media

import (
	"errors"
	tropestogo "github.com/jlgallego99/TropesToGo"
)

var (
	ErrEmptyTitle = errors.New("a media work must have a title")
)

// MediaType enumerates all supported Media types in TropesToGo
type MediaType int64

const (
	Film MediaType = iota
	Series
	Anime
	VideoGames
)

// Media holds the logic of a work with its tropes that exist within a particular medium in TvTropes
type Media struct {
	work      *tropestogo.Work // Root entity
	tropes    []tropestogo.Trope
	mediatype MediaType
}

// NewMedia is a factory that creates a Media aggregate with validations
func NewMedia(title string) (Media, error) {
	return Media{}, nil
}
