package twitter

import (
	"context"
	"fmt"
	"time"

	_twitter "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	domain_twitter "github.com/hareku/emosearch-api/pkg/domain/twitter"
)

type twitterOauth1Client struct {
	config *oauth1.Config
}

// NewTwitterOauth1Client creates Client of domain Twitter.
func NewTwitterOauth1Client(config *oauth1.Config) domain_twitter.Client {
	return &twitterOauth1Client{config}
}

func (c *twitterOauth1Client) Search(ctx context.Context, input *domain_twitter.SearchInput) ([]domain_twitter.Tweet, error) {
	client := c.makeTwitterClient(ctx, input.TwitterAccessToken, input.TwitterAccessTokenSecret)
	search, _, err := client.Search.Tweets(&_twitter.SearchTweetParams{
		Query:     input.Query,
		MaxID:     input.MaxID,
		SinceID:   input.SinceID,
		TweetMode: "extended",
		Count:     100,
	})
	if err != nil {
		return nil, fmt.Errorf("twitter error: %w", err)
	}

	var tweets []domain_twitter.Tweet
	statuses := search.Statuses

	for i := 0; i < len(statuses); i++ {
		createdAt, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", statuses[i].CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("tweet created_at parse error: %w", err)
		}
		tweets = append(tweets, domain_twitter.Tweet{
			TweetID:   statuses[i].ID,
			AuthorID:  statuses[i].User.ID,
			Text:      statuses[i].FullText,
			CreatedAt: createdAt,
		})
	}

	return tweets, nil
}

func (c *twitterOauth1Client) makeTwitterClient(ctx context.Context, accessToken string, accessTokenSecret string) *_twitter.Client {
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := c.config.Client(ctx, token)
	return _twitter.NewClient(httpClient)
}
