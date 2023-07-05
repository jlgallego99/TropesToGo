package trope

import (
	"errors"
	"time"
)

var (
	ErrEmptyTitle = errors.New("a media work must have a title")
)

// Work in TvTropes is any production with a story that has tropes. It's a mutable entity because its information can be updated
type Work struct {
	// Title of the work in TvTropes
	Title string
	// Year is the release year date of the Work, which serves to differentiate it with other Works which may have the same name
	Year string
	// LastUpdated is the last time the Work information was updated
	LastUpdated time.Time
	// Tropes that define the Work. Is a set, which means that all tropes are unique
	Tropes map[Trope]struct{}
	// SubTropes that belong to any of the SubWikis of the Work. Is a set, which means that all SubTropes are unique
	// There can't be two Tropes on the same SubWiki, but the same Trope can be in different SubWikis
	SubTropes map[Trope]struct{}
}
