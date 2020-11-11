package twitter

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	dtwitter "github.com/hareku/emosearch-api/pkg/domain/twitter"
)

type twitterOauth1Client struct {
	config *oauth1.Config
}

// NewTwitterOauth1Client creates Client of domain Twitter.
func NewTwitterOauth1Client(config *oauth1.Config) dtwitter.Client {
	return &twitterOauth1Client{config}
}

func (c *twitterOauth1Client) Search(ctx context.Context, input *dtwitter.SearchInput) ([]dtwitter.Tweet, error) {
	client := c.makeTwitterClient(ctx, input.TwitterAccessToken, input.TwitterAccessTokenSecret)
	search, _, err := client.Search.Tweets(&sdk.SearchTweetParams{
		Query:           addExcludeRetweetOption(input.Query),
		MaxID:           input.MaxID,
		SinceID:         input.SinceID,
		IncludeEntities: sdk.Bool(true),
		TweetMode:       "extended",
		Count:           100,
	})
	if err != nil {
		return nil, fmt.Errorf("twitter error: %w", err)
	}

	tweets := []dtwitter.Tweet{}

	for _, tweet := range search.Statuses {
		createdAt, err := tweet.CreatedAtTime()
		if err != nil {
			return nil, fmt.Errorf("tweet created_at parse error: %w", err)
		}
		tweets = append(tweets, dtwitter.Tweet{
			TweetID:  tweet.ID,
			AuthorID: tweet.User.ID,
			User: &dtwitter.User{
				ID:              tweet.User.ID,
				Name:            tweet.User.Name,
				ScreenName:      tweet.User.ScreenName,
				ProfileImageURL: tweet.User.ProfileImageURLHttps,
			},
			Entities:  makeEntities(&tweet),
			Text:      tweet.FullText,
			CreatedAt: createdAt,
		})
	}

	return tweets, nil
}

func makeEntities(tweet *sdk.Tweet) *dtwitter.Entities {
	entities := dtwitter.Entities{}

	for _, hashtag := range tweet.Entities.Hashtags {
		entities.HashTags = append(entities.HashTags, dtwitter.HashTag{
			Start: hashtag.Indices.Start(),
			End:   hashtag.Indices.End(),
			Tag:   hashtag.Text,
		})
	}

	for _, url := range tweet.Entities.Urls {
		entities.URLs = append(entities.URLs, dtwitter.URL{
			Start:       url.Indices.Start(),
			End:         url.Indices.End(),
			URL:         url.URL,
			DisplayURL:  url.DisplayURL,
			ExpandedURL: url.ExpandedURL,
		})
	}

	for _, mention := range tweet.Entities.UserMentions {
		entities.Mentions = append(entities.Mentions, dtwitter.Mention{
			Start: mention.Indices.Start(),
			End:   mention.Indices.End(),
			Tag:   mention.ScreenName,
		})
	}

	for _, medium := range tweet.Entities.Media {
		entities.Media = append(entities.Media, dtwitter.Medium{
			Start:    medium.Indices.Start(),
			End:      medium.Indices.End(),
			URL:      medium.URL,
			MediaURL: medium.MediaURLHttps,
		})
	}

	return &entities
}

func (c *twitterOauth1Client) makeTwitterClient(ctx context.Context, accessToken string, accessTokenSecret string) *sdk.Client {
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := c.config.Client(ctx, token)
	return sdk.NewClient(httpClient)
}

func addExcludeRetweetOption(query string) string {
	if !strings.Contains(query, "-filter:retweets") {
		query += " -filter:retweets"
	}

	return query
}
