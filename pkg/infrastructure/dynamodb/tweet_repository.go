package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/guregu/dynamo"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

type dynamoDBTweetRepository struct {
	dynamoDB dynamo.Table
}

// NewDynamoDBTweetRepository creates TweetRepository which is implemented by DynamoDB.
func NewDynamoDBTweetRepository(dynamoDB dynamo.Table) repository.TweetRepository {
	return &dynamoDBTweetRepository{dynamoDB}
}

type dynamoDBTweet struct {
	PK string
	SK string

	TweetID                model.TweetID
	SearchID               model.SearchID
	AuthorID               int64
	Text                   string
	SentimentScoreMixed    *float64
	SentimentScoreNegative *float64
	SentimentScoreNeutral  *float64
	SentimentScorePositive *float64
	TweetCreatedAt         time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (d *dynamoDBTweet) NewTweetModel() *model.Tweet {
	return &model.Tweet{
		CreatedAt: time.Time(d.CreatedAt),
		UpdatedAt: time.Time(d.UpdatedAt),
	}
}

func (r *dynamoDBTweetRepository) Store(ctx context.Context, tweet *model.Tweet) error {
	createdAt := time.Now()
	tweet.CreatedAt = createdAt
	tweet.UpdatedAt = createdAt

	dtweet := dynamoDBTweet{
		PK: fmt.Sprintf("SEARCH#%s", tweet.SearchID),
		SK: fmt.Sprintf("TWEET#%d", tweet.TweetID),

		TweetID:                tweet.TweetID,
		AuthorID:               tweet.AuthorID,
		SearchID:               tweet.SearchID,
		Text:                   tweet.Text,
		SentimentScoreMixed:    tweet.SentimentScore.Mixed,
		SentimentScoreNegative: tweet.SentimentScore.Negative,
		SentimentScoreNeutral:  tweet.SentimentScore.Negative,
		SentimentScorePositive: tweet.SentimentScore.Neutral,
		TweetCreatedAt:         tweet.TweetCreatedAt,
		CreatedAt:              tweet.CreatedAt,
		UpdatedAt:              tweet.CreatedAt,
	}

	err := r.dynamoDB.Put(&dtweet).RunWithContext(ctx)

	if err != nil {
		return fmt.Errorf("dynamo error: %w", err)
	}

	return nil
}

func (r *dynamoDBTweetRepository) LatestTweetID(ctx context.Context, searchID model.SearchID) (model.TweetID, error) {
	var dynamoTweet dynamoDBTweet
	err := r.dynamoDB.
		Get("PK", fmt.Sprintf("SEARCH#%s", searchID)).
		Limit(1).
		Order(true).
		OneWithContext(ctx, &dynamoTweet)

	if errors.Is(err, dynamo.ErrNotFound) {
		return 0, repository.ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("dynamo error: %w", err)
	}

	return dynamoTweet.TweetID, nil
}
