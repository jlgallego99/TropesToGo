package tropestogo

import "github.com/google/uuid"

// Work in TvTropes is any production which has a narrative that uses tropes
type Work struct {
	ID    uuid.UUID
	Title string
}
