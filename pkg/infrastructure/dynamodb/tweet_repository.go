package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
	AuthorID               string
	Text                   string
	SentimentScoreMixed    *float64
	SentimentScoreNegative *float64
	SentimentScoreNeutral  *float64
	SentimentScorePositive *float64
	TweetCreatedAt         dynamodbattribute.UnixTime
	CreatedAt              dynamodbattribute.UnixTime
	UpdatedAt              dynamodbattribute.UnixTime
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
		SK: fmt.Sprintf("TWEET#%s", tweet.TweedID),

		AuthorID:               tweet.AuthorID,
		SearchID:               tweet.SearchID,
		Text:                   tweet.Text,
		SentimentScoreMixed:    tweet.SentimentScore.Mixed,
		SentimentScoreNegative: tweet.SentimentScore.Negative,
		SentimentScoreNeutral:  tweet.SentimentScore.Negative,
		SentimentScorePositive: tweet.SentimentScore.Neutral,
		TweetCreatedAt:         dynamodbattribute.UnixTime(tweet.TweetCreatedAt),
		CreatedAt:              dynamodbattribute.UnixTime(tweet.CreatedAt),
		UpdatedAt:              dynamodbattribute.UnixTime(tweet.UpdatedAt),
	}

	err := r.dynamoDB.Put(&dtweet).RunWithContext(ctx)

	if err != nil {
		return fmt.Errorf("dynamo error: %w", err)
	}

	return nil
}