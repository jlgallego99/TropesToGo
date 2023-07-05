package trope

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
	case UnknownTropeIndex, GenreTrope, MediaTrope, NarrativeTrope, TopicalTrope:
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
	// title is the trope immutable name and is recognised by it
	title string
	// index is the conceptual group of tropes to which this trope belongs
	index TropeIndex
	// isMain represents if the trope is on the main Work page
	isMain bool
	// subpage refers to the subpage within a Work this trope belongs
	subpage string
}

// NewTrope is a factory that creates a valid Trope value object by receiving its name and the index to which it belongs
// It checks if the index is valid and returns an ErrUnknownIndex if it's not
func NewTrope(title string, index TropeIndex, subpage string) (Trope, error) {
	if len(title) == 0 {
		return Trope{}, ErrMissingValues
	}

	if !index.IsValid() {
		return Trope{}, ErrUnknownIndex
	}

	return Trope{
		title:   title,
		index:   index,
		isMain:  subpage == "",
		subpage: subpage,
	}, nil
}

// GetTitle returns the main identifier of the trope: its title
func (trope Trope) GetTitle() string {
	return trope.title
}

// GetIndex returns the main category this trope belongs to in narratives
func (trope Trope) GetIndex() TropeIndex {
	return trope.index
}

// GetIsMain returns a boolean indicating whether it's a trope on the main Work page
func (trope Trope) GetIsMain() bool {
	return trope.isMain
}

// GetSubpage returns the Work subpage this trope belongs to
func (trope Trope) GetSubpage() string {
	return trope.subpage
}
