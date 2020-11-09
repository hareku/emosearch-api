package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/guregu/dynamo"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
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
		TweetID:  d.TweetID,
		SearchID: d.SearchID,
		AuthorID: d.AuthorID,
		Text:     d.Text,
		SentimentScore: &sentiment.Score{
			Negative: d.SentimentScoreNegative,
			Positive: d.SentimentScorePositive,
			Mixed:    d.SentimentScoreMixed,
			Neutral:  d.SentimentScoreNeutral,
		},
		TweetCreatedAt: d.TweetCreatedAt,
		CreatedAt:      time.Time(d.CreatedAt),
		UpdatedAt:      time.Time(d.UpdatedAt),
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

func (r *dynamoDBTweetRepository) List(ctx context.Context, input *repository.TweetRepositoryListInput) ([]model.Tweet, error) {
	var dTweets []dynamoDBTweet

	q := r.dynamoDB.
		Get("PK", fmt.Sprintf("SEARCH#%s", input.SearchID)).
		Order(false).
		Limit(input.Limit)
	if input.UntilID != 0 {
		q.Range("SK", dynamo.Less, fmt.Sprintf("TWEET#%d", input.UntilID))
	}

	err := q.AllWithContext(ctx, &dTweets)

	if errors.Is(err, dynamo.ErrNotFound) {
		return []model.Tweet{}, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("dynamo error: %w", err)
	}

	tweets := []model.Tweet{}

	for i := 0; i < len(dTweets); i++ {
		tweets = append(tweets, *dTweets[i].NewTweetModel())
	}

	return tweets, nil
}

func (r *dynamoDBTweetRepository) LatestTweetID(ctx context.Context, searchID model.SearchID) (model.TweetID, error) {
	var dynamoTweet dynamoDBTweet
	err := r.dynamoDB.
		Get("PK", fmt.Sprintf("SEARCH#%s", searchID)).
		Limit(1).
		Order(false).
		OneWithContext(ctx, &dynamoTweet)

	if errors.Is(err, dynamo.ErrNotFound) {
		return 0, repository.ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("dynamo error: %w", err)
	}

	return dynamoTweet.TweetID, nil
}
