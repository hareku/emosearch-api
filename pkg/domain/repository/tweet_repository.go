package repository

import (
	"context"

	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
)

// TweetRepository provides CRUD methods for Tweet domain.
type TweetRepository interface {
	Store(ctx context.Context, tweet *model.Tweet) error
	BatchStore(ctx context.Context, tweets []*model.Tweet) error
	LatestTweetID(ctx context.Context, searchID model.SearchID) (model.TweetID, error)
	List(ctx context.Context, input *TweetRepositoryListInput) ([]model.Tweet, error)
}

// TweetRepositoryListInput is used for List method of Tweet repository.
type TweetRepositoryListInput struct {
	SearchID       model.SearchID
	Limit          int64
	UntilID        model.TweetID
	SentimentLabel *sentiment.Label
}
