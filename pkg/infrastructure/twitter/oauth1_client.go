package twitter

import (
	"context"
	"fmt"
	"strings"

	_twitter "github.com/dghubble/go-twitter/twitter"
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
	search, _, err := client.Search.Tweets(&_twitter.SearchTweetParams{
		Query:           addExcludeRetweetOption(input.Query),
		MaxID:           input.MaxID,
		SinceID:         input.SinceID,
		IncludeEntities: _twitter.Bool(true),
		TweetMode:       "extended",
		Count:           100,
	})
	if err != nil {
		return nil, fmt.Errorf("twitter error: %w", err)
	}

	tweets := []dtwitter.Tweet{}

	for _, st := range search.Statuses {
		createdAt, err := st.CreatedAtTime()
		if err != nil {
			return nil, fmt.Errorf("tweet created_at parse error: %w", err)
		}
		tweets = append(tweets, dtwitter.Tweet{
			TweetID:  st.ID,
			AuthorID: st.User.ID,
			User: &dtwitter.User{
				ID:              st.User.ID,
				Name:            st.User.Name,
				ScreenName:      st.User.ScreenName,
				ProfileImageURL: st.User.ProfileImageURLHttps,
			},
			Entities:  makeEntities(&st),
			Text:      st.FullText,
			CreatedAt: createdAt,
		})
	}

	return tweets, nil
}

func makeEntities(tweet *_twitter.Tweet) *dtwitter.Entities {
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

func (c *twitterOauth1Client) makeTwitterClient(ctx context.Context, accessToken string, accessTokenSecret string) *_twitter.Client {
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := c.config.Client(ctx, token)
	return _twitter.NewClient(httpClient)
}

func addExcludeRetweetOption(query string) string {
	if !strings.Contains(query, "-filter:retweets") {
		query += " -filter:retweets"
	}

	return query
}
