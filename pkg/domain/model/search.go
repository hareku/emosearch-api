package model

import "time"

// SearchID is the identifier of Search domain.
type SearchID string

// Search is the structure of a searching configuration.
type Search struct {
	SearchID           SearchID
	UserID             UserID
	Title              string
	Query              string
	NextSearchUpdateAt time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
