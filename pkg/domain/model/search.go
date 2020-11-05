package model

// SearchID is the identifier of Search domain.
type SearchID string

// Search is the structure of a searching configuration.
type Search struct {
	UserID    string
	SearchID  string
	Title     string
	Query     string
	CreatedAt string
	UpdatedAt string
}
