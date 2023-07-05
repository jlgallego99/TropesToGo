package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jlgallego99/TropesToGo/trope"
	"github.com/jlgallego99/TropesToGo/tvtropespages"
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

// String is an implementation of the Stringer interface for comparing string media types and avoid using literals
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

// IsValid checks whether a MediaType is known or not
func (mediatype MediaType) IsValid() bool {
	switch mediatype {
	case Film, Series, Anime, VideoGames:
		return true
	}

	return false
}

// ToMediaType converts a string to a MediaType
// It returns an ErrUnknownMediaType if the MediaType isn't recognized
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
	work *trope.Work

	// page is the TvTropes webpage from where the Work information is extracted
	page tvtropespages.Page

	// MediaType is the media index that this work belongs to
	mediaType MediaType
}

// JsonResponse is an object for marshaling/unmarshalling a single Media object in Json
type JsonResponse struct {
	Title       string      `json:"title"`
	Year        string      `json:"year"`
	MediaType   string      `json:"media_type"`
	LastUpdated string      `json:"last_updated"`
	URL         string      `json:"url"`
	Tropes      []JsonTrope `json:"tropes"`
	SubTropes   []JsonTrope `json:"sub_tropes"`
}

// JsonTrope is part of JsonResponse, and represent a trope with the index to which it belongs
type JsonTrope struct {
	Title     string `json:"title"`
	Namespace string `json:"namespace"`
}

// MarshalJSON implements Marshaller interface for custom marshalling of Media objects
// Returns a byte array that can be marshalled into a JSON file
func (media Media) MarshalJSON() ([]byte, error) {
	tropes, subTropes := GetJsonTropes(media)

	return json.Marshal(&JsonResponse{
		Title:       media.work.Title,
		Year:        media.work.Year,
		MediaType:   media.mediaType.String(),
		LastUpdated: media.work.LastUpdated.Format("2006-01-02 15:04:05"),
		URL:         media.page.GetUrl().String(),
		Tropes:      tropes,
		SubTropes:   subTropes,
	})
}

// GetJsonTropes receives a media object and transforms it into a JsonTrope array with all its tropes for correct marshalling
// Return two JsonTrope arrays, the first for the main tropes and the second for the sub tropes
func GetJsonTropes(media Media) ([]JsonTrope, []JsonTrope) {
	var tropes, subTropes []JsonTrope
	mediaType := media.GetMediaType().String()
	for trope := range media.GetWork().Tropes {
		title := trope.GetTitle()
		namespace := trope.GetSubpage()

		if title != "" && namespace == "" /*&& index != "UnknownTropeIndex"*/ {
			tropes = append(tropes, JsonTrope{
				Title:     title,
				Namespace: mediaType,
			})
		}
	}

	for subTrope := range media.GetWork().SubTropes {
		title := subTrope.GetTitle()
		namespace := subTrope.GetSubpage()

		if title != "" && namespace != "" /*&& index != "UnknownTropeIndex"*/ {
			subTropes = append(subTropes, JsonTrope{
				Title:     title,
				Namespace: namespace,
			})
		}
	}

	return tropes, subTropes
}

// NewMedia is a factory that creates a Media aggregate with validations from a title, year, a set of all tropes, a page object and a media type object
// It divides the tropes between main and secondary
// It returns a correctly formed Media object and an error of type ErrMissingValues if the title or page are empty
// an ErrInvalidYear if the year isn't real or an ErrUnknownMediaType if the received media type isn't known
func NewMedia(title, year string, lastUpdated time.Time, tropes map[trope.Trope]struct{}, page tvtropespages.Page, mediaType MediaType) (Media, error) {
	if page.GetUrl() == nil {
		return Media{}, ErrMissingValues
	}

	if len(title) == 0 {
		return Media{}, ErrMissingValues
	}

	if len(year) > 0 {
		r, _ := regexp.Compile("^[0-9]{4}$")

		if !r.MatchString(year) {
			return Media{}, fmt.Errorf("%w: "+year, ErrInvalidYear)
		}
	}

	if !mediaType.IsValid() {
		return Media{}, fmt.Errorf("%w: "+mediaType.String(), ErrUnknownMediaType)
	}

	mainTropes := make(map[trope.Trope]struct{})
	subTropes := make(map[trope.Trope]struct{})
	for trope := range tropes {
		if trope.GetIsMain() {
			mainTropes[trope] = struct{}{}
		} else {
			subTropes[trope] = struct{}{}
		}
	}

	work := &trope.Work{
		Title:       title,
		Year:        year,
		LastUpdated: lastUpdated,
		Tropes:      mainTropes,
		SubTropes:   subTropes,
	}

	return Media{
		work:      work,
		page:      page,
		mediaType: mediaType,
	}, nil
}

// GetWork returns the Work object that this media object manages
func (media Media) GetWork() *trope.Work {
	return media.work
}

// GetPage returns the Page object that this media object manages
func (media Media) GetPage() tvtropespages.Page {
	return media.page
}

// GetMediaType return the type this media belongs to
func (media Media) GetMediaType() MediaType {
	return media.mediaType
}
