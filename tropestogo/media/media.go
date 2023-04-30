package media

import tropestogo "github.com/jlgallego99/TropesToGo"

// MediaType enumerates all supported Media types in TropesToGo
type MediaType int64

const (
	Film MediaType = iota
	Series
	Anime
	VideoGames
)

// Media holds the logic of a work with it tropes that exist within a particular medium in TvTropes
type Media struct {
	work      *tropestogo.Work // Root entity
	tropes    []*tropestogo.Trope
	mediatype MediaType
}
