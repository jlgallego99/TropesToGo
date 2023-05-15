package media

import (
	"fmt"
	tropestogo "github.com/jlgallego99/TropesToGo"
)

// MediaType enumerates all supported Media types in TropesToGo
type MediaType int64

const (
	Film MediaType = iota
	Series
	Anime
	VideoGames
)

// Implement Stringer interface for comparing string media types and avoid using literals
func (mediatype MediaType) String() string {
	switch mediatype {
	case Film:
		return "Film"
	case Series:
		return "Series"
	case Anime:
		return "Anime"
	case VideoGames:
		return "VideoGames"
	default:
		return fmt.Sprintf("%d", int(mediatype))
	}
}

// Media holds the logic of all Works with its tropes that exist within a particular medium in TvTropes
type Media struct {
	// work is the root entity, holds the work information and its tropes
	work *tropestogo.Work

	// page is the TvTropes webpage from where the Work information is extracted
	page *tropestogo.Page

	// MediaType is the media index that this work belongs to
	MediaType MediaType
}

// NewMedia is a factory that creates a Media aggregate with validations
func NewMedia(title string) (Media, error) {
	return Media{}, nil
}
