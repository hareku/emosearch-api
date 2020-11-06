package model

// UserID is the identifier of User domain.
type UserID string

// User contains user's data and auth tokens.
type User struct {
	UserID                   UserID
	TwitterAccessToken       string
	TwitterAccessTokenSecret string
}
