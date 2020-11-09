package model

import (
	"time"

	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
)

// TweetID is the identifier of Tweet domain.
type TweetID int64

// Tweet is the structure of a tweet.
type Tweet struct {
	TweetID        TweetID `json:",string"`
	SearchID       SearchID
	AuthorID       int64 `json:",string"`
	Text           string
	SentimentScore *sentiment.Score
	TweetCreatedAt time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
