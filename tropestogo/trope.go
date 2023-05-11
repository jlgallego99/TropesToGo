package tropestogo

// TropeIndex enumerates all index a trope can belong to in TvTropes
type TropeIndex int64

const (
	GenreTrope TropeIndex = iota
	MediaTrope
	NarrativeTrope
	TopicalTrope
)

// Trope represents a reiterative resource that is collected in TvTropes
type Trope struct {
	// A trope has an immutable name and is recognised by it
	title string
	// index is the conceptual group of tropes to which this trope belongs
	index TropeIndex
}
