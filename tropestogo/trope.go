package tropestogo

import "errors"

var (
	ErrMissingValues = errors.New("one or more fields are missing")
	ErrUnknownIndex  = errors.New("unknown trope index")
)

// TropeIndex enumerates all index a trope can belong to in TvTropes
type TropeIndex int64

const (
	UnknownTropeIndex TropeIndex = iota
	GenreTrope
	MediaTrope
	NarrativeTrope
	TopicalTrope
)

func (index TropeIndex) IsValid() bool {
	switch index {
	case GenreTrope, MediaTrope, NarrativeTrope, TopicalTrope:
		return true
	}

	return false
}

// ToTropeIndex converts a string to a MediaType
func ToTropeIndex(tropeIndexString string) (TropeIndex, error) {
	for tropeindex := UnknownTropeIndex + 1; tropeindex <= GenreTrope; tropeindex++ {
		if tropeIndexString == tropeindex.String() {
			return tropeindex, nil
		}
	}

	return UnknownTropeIndex, ErrUnknownIndex
}

// Implement Stringer interface for comparing string media types and avoid using literals
func (index TropeIndex) String() string {
	switch index {
	case GenreTrope:
		return "GenreTrope"
	case MediaTrope:
		return "MediaTrope"
	case NarrativeTrope:
		return "NarrativeTrope"
	case TopicalTrope:
		return "TopicalTrope"
	default:
		return "UnknownTropeIndex"
	}
}

// Trope represents a reiterative resource that is collected in TvTropes
type Trope struct {
	// A trope has an immutable name and is recognised by it
	title string
	// index is the conceptual group of tropes to which this trope belongs
	index TropeIndex
}

// NewTrope is a factory that creates a valid Trope value object
func NewTrope(title string, index TropeIndex) (Trope, error) {
	if len(title) == 0 {
		return Trope{}, ErrMissingValues
	}

	if !index.IsValid() {
		return Trope{}, ErrUnknownIndex
	}

	return Trope{
		title: title,
		index: index,
	}, nil
}

func (trope Trope) GetTitle() string {
	return trope.title
}

func (trope Trope) GetIndex() TropeIndex {
	return trope.index
}
