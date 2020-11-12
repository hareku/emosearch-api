package model

import (
	"time"

	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
	"github.com/hareku/emosearch-api/pkg/domain/twitter"
)

// TweetID is the identifier of Tweet domain.
type TweetID int64

// TwitterUser represents Twitter user object.
type TwitterUser struct {
	ID              int64 `json:",string"`
	Name            string
	ScreenName      string
	ProfileImageURL string
}

// Tweet is the structure of a tweet.
type Tweet struct {
	TweetID            TweetID `json:",string"`
	SearchID           SearchID
	AuthorID           int64 `json:",string"`
	User               *TwitterUser
	Text               string
	SentimentScore     *sentiment.Score `dynamo:",null"`
	SentimentLabel     sentiment.Label
	Entities           *twitter.Entities `dynamo:",null"`
	ExpirationUnixTime int64             `json:"-"`
	TweetCreatedAt     time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
