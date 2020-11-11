package twitter

import (
	"context"
	"time"
)

// Client provides twitter actions.
type Client interface {
	Search(ctx context.Context, input *SearchInput) ([]Tweet, error)
}

// SearchInput is the input for Search method.
type SearchInput struct {
	Query                    string
	TwitterAccessToken       string
	TwitterAccessTokenSecret string
	SinceID                  int64
	MaxID                    int64
}

// User represents Twitter user object.
type User struct {
	ID              int64 `json:",string"`
	Name            string
	ScreenName      string
	ProfileImageURL string
}

// Tweet represents Twitter Tweet.
type Tweet struct {
	TweetID   int64
	AuthorID  int64
	User      *User
	Text      string
	Entities  *Entities
	CreatedAt time.Time
}

// Entities represents Twitter Entities.
type Entities struct {
	HashTags []HashTag
	Mentions []Mention
	URLs     []URL
	Media    []Medium
}

// HashTag represents Twitter HashTag.
type HashTag struct {
	Start int
	End   int
	Tag   string
}

// Mention represents Twitter Mention.
type Mention struct {
	Start int
	End   int
	Tag   string
}

// URL represents Twitter URL.
type URL struct {
	Start       int
	End         int
	URL         string
	ExpandedURL string
	DisplayURL  string
}

// Medium represents Twitter Medium.
type Medium struct {
	Start       int
	End         int
	URL         string
	ExpandedURL string
}
