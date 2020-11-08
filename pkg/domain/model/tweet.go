package model

import (
	"time"

	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
)

// TweetID is the identifier of Tweet domain.
type TweetID string

// Tweet is the structure of a tweet.
type Tweet struct {
	TweetID        TweetID
	SearchID       SearchID
	AuthorID       string
	Text           string
	SentimentScore *sentiment.Score
	TweetCreatedAt time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
