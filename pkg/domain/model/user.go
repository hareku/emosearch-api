package model

// UserID is the identifier of User domain.
type UserID string

// User contains user's data and auth tokens.
type User struct {
	ID                       UserID
	Name                     string
	Email                    string
	TwitterAccessToken       string
	TwitterAccessTokenSecret string
	CreatedAt                string
	UpdatedAt                string
}
