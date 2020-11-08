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

// Tweet represents Twitter Tweet.
type Tweet struct {
	TweetID   int64
	AuthorID  int64
	Text      string
	CreatedAt time.Time
}
