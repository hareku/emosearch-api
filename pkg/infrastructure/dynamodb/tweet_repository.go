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
	PK                    string
	SK                    string
	TweetSentimentIndexPK string
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
		PK:                    fmt.Sprintf("SEARCH#%s", tweet.SearchID),
		SK:                    fmt.Sprintf("TWEET#%d", tweet.TweetID),
		TweetSentimentIndexPK: r.buildTweetSentimentIndexPK(tweet.SearchID, tweet.SentimentLabel),
		Tweet:                 tweet,
	}

	err := r.dynamoDB.Put(&dtweet).RunWithContext(ctx)

	if err != nil {
		return fmt.Errorf("dynamo error: %w", err)
	}

	return nil
}

func (r *dynamoDBTweetRepository) List(ctx context.Context, input *repository.TweetRepositoryListInput) ([]model.Tweet, error) {
	var dTweets []dynamoDBTweet

	q := r.buildListQuery(input)
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

func (r *dynamoDBTweetRepository) buildListQuery(input *repository.TweetRepositoryListInput) *dynamo.Query {
	var q *dynamo.Query

	if input.SentimentLabel == nil {
		q = r.dynamoDB.Get("PK", fmt.Sprintf("SEARCH#%s", input.SearchID))
	} else {
		q = r.dynamoDB.
			Get("TweetSentimentIndexPK", r.buildTweetSentimentIndexPK(input.SearchID, *input.SentimentLabel)).
			Index("TweetSentimentIndex")
	}

	q.Order(false).Limit(input.Limit)

	if input.UntilID != 0 {
		q.Range("SK", dynamo.Less, fmt.Sprintf("TWEET#%d", input.UntilID))
	}

	return q
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

func (r *dynamoDBTweetRepository) buildTweetSentimentIndexPK(ID model.SearchID, label sentiment.Label) string {
	labelKey := label
	if label == sentiment.LabelPositive || label == sentiment.LabelNegative {
		labelKey = "POS_OR_NEG__" + labelKey
	}
	return fmt.Sprintf("SEARCH#%s#%s", ID, labelKey)
}
