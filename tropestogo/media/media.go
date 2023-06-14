package media

import (
	"encoding/json"
	"errors"
	"fmt"
	tropestogo "github.com/jlgallego99/TropesToGo"
	"regexp"
	"time"
)

var (
	ErrMissingValues    = errors.New("one or more fields are missing")
	ErrInvalidYear      = errors.New("year is invalid")
	ErrUnknownMediaType = errors.New("unknown media type")
)

// MediaType enumerates all supported Media types in TropesToGo
type MediaType int64

const (
	UnknownMediaType MediaType = iota
	Film
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
		return "UnknownMediaType"
	}
}

func (mediatype MediaType) IsValid() bool {
	switch mediatype {
	case Film, Series, Anime, VideoGames:
		return true
	}

	return false
}

// ToMediaType converts a string to a MediaType
func ToMediaType(mediaTypeString string) (MediaType, error) {
	for mediatype := UnknownMediaType + 1; mediatype <= VideoGames; mediatype++ {
		if mediaTypeString == mediatype.String() {
			return mediatype, nil
		}
	}

	return UnknownMediaType, fmt.Errorf("%w: "+mediaTypeString, ErrUnknownMediaType)
}

// Media holds the logic of all Works with its tropes that exist within a particular medium in TvTropes
type Media struct {
	// work is the root entity, holds the work information and its tropes
	work *tropestogo.Work

	// page is the TvTropes webpage from where the Work information is extracted
	page tropestogo.Page

	// MediaType is the media index that this work belongs to
	mediaType MediaType
}

type JsonResponse struct {
	Title       string      `json:"title"`
	Year        string      `json:"year"`
	MediaType   string      `json:"media_type"`
	LastUpdated string      `json:"last_updated"`
	URL         string      `json:"url"`
	Tropes      []JsonTrope `json:"tropes"`
}

type JsonTrope struct {
	Title string `json:"title"`
	Index string `json:"index"`
}

func (media Media) MarshalJSON() ([]byte, error) {
	tropes := GetJsonTropes(media)

	return json.Marshal(&JsonResponse{
		Title:       media.work.Title,
		Year:        media.work.Year,
		MediaType:   media.mediaType.String(),
		LastUpdated: media.work.LastUpdated.Format("2006-01-02 15:04:05"),
		URL:         media.page.GetUrl().String(),
		Tropes:      tropes,
	})
}

func GetJsonTropes(media Media) []JsonTrope {
	var tropes []JsonTrope
	for trope := range media.GetWork().Tropes {
		title := trope.GetTitle()
		index := trope.GetIndex().String()

		if title != "" && index != "" /*&& index != "UnknownTropeIndex"*/ {
			tropes = append(tropes, JsonTrope{
				Title: trope.GetTitle(),
				Index: trope.GetIndex().String(),
			})
		}
	}

	return tropes
}

// NewMedia is a factory that creates a Media aggregate with validations
func NewMedia(title, year string, lastUpdated time.Time, tropes map[tropestogo.Trope]struct{}, page tropestogo.Page, mediaType MediaType) (Media, error) {
	if page.GetUrl() == nil {
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
			return Media{}, fmt.Errorf("%w: "+year, ErrInvalidYear)
		}
	}

	if !mediaType.IsValid() {
		return Media{}, fmt.Errorf("%w: "+mediaType.String(), ErrUnknownMediaType)
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

func (media Media) GetPage() tropestogo.Page {
	return media.page
}

func (media Media) GetMediaType() MediaType {
	return media.mediaType
}
