package domain

// SearchID is the identifier of Search domain.
type SearchID string

// Search is the structure of a searching configuration.
type Search struct {
	ID        SearchID
	Query     string
	Title     string
	CreatedAt string
	UpdatedAt string
}
