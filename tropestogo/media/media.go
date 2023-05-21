package media

import (
	"errors"
	"fmt"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"regexp"
	"time"
)

var (
	ErrMissingValues        = errors.New("one or more fields are missing")
	ErrInvalidYear          = errors.New("year is invalid")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
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

func (mediatype MediaType) IsValid() bool {
	switch mediatype {
	case Film, Series, Anime, VideoGames:
		return true
	}

	return false
}

// Media holds the logic of all Works with its tropes that exist within a particular medium in TvTropes
type Media struct {
	// work is the root entity, holds the work information and its tropes
	work *tropestogo.Work

	// page is the TvTropes webpage from where the Work information is extracted
	page *tropestogo.Page

	// MediaType is the media index that this work belongs to
	mediaType MediaType
}

// NewMedia is a factory that creates a Media aggregate with validations
func NewMedia(title, year string, lastUpdated time.Time, tropes []tropestogo.Trope, page *tropestogo.Page, mediaType MediaType) (Media, error) {
	if page == nil {
		return Media{}, ErrMissingValues
	}

	// Check if there's a title. A year can be empty, because not all media will have it extracted
	if len(title) == 0 {
		return Media{}, ErrMissingValues
	}

	// Check if the Year string represents a valid year number (4 digits between 0 and 9)
	if len(year) > 0 {
		r, _ := regexp.Compile("^[0-9]{4}$")

		if !r.MatchString(year) {
			return Media{}, ErrInvalidYear
		}
	}

	if !mediaType.IsValid() {
		return Media{}, ErrUnsupportedMediaType
	}

	work := &tropestogo.Work{
		Title:       title,
		Year:        year,
		LastUpdated: lastUpdated,
		Tropes:      tropes,
	}

	return Media{
		work:      work,
		page:      page,
		mediaType: mediaType,
	}, nil
}

func (media Media) GetWork() *tropestogo.Work {
	return media.work
}

func (media Media) GetPage() *tropestogo.Page {
	return media.page
}

func (media Media) GetMediaType() MediaType {
	return media.mediaType
}
