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
	*model.Tweet
}

func (d *dynamoDBTweet) NewTweetModel() *model.Tweet {
	return d.Tweet
}

func (r *dynamoDBTweetRepository) Store(ctx context.Context, tweet *model.Tweet) error {
	createdAt := time.Now()
	tweet.CreatedAt = createdAt
	tweet.UpdatedAt = createdAt

	dtweet := dynamoDBTweet{
		PK:    fmt.Sprintf("SEARCH#%s", tweet.SearchID),
		SK:    fmt.Sprintf("TWEET#%d", tweet.TweetID),
		Tweet: tweet,
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
