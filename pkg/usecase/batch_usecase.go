package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
	"github.com/hareku/emosearch-api/pkg/domain/twitter"
)

// BatchUsecase provides usecases of Batch domain.
type BatchUsecase interface {
	CollectTweets(ctx context.Context, searchID model.SearchID, userID model.UserID) error
}

type batchUsecase struct {
	userUsecase       UserUsecase
	searchUsecase     SearchUsecase
	tweetRepository   repository.TweetRepository
	twitterClient     twitter.Client
	sentimentDetector sentiment.Detector
}

// NewBatchUsecaseInput is the input of NewBatchUsecase.
type NewBatchUsecaseInput struct {
	UserUsecase       UserUsecase
	SearchUsecase     SearchUsecase
	TweetRepository   repository.TweetRepository
	TwitterClient     twitter.Client
	SentimentDetector sentiment.Detector
}

// NewBatchUsecase creates BatchUsecase.
func NewBatchUsecase(input *NewBatchUsecaseInput) BatchUsecase {
	return &batchUsecase{
		userUsecase:       input.UserUsecase,
		searchUsecase:     input.SearchUsecase,
		tweetRepository:   input.TweetRepository,
		twitterClient:     input.TwitterClient,
		sentimentDetector: input.SentimentDetector,
	}
}

func (u *batchUsecase) CollectTweets(ctx context.Context, searchID model.SearchID, userID model.UserID) error {
	search, input, err := u.prepareSearch(ctx, searchID, userID)
	if err != nil {
		return fmt.Errorf("collect tweets preparation error: %w", err)
	}

	err = u.searchUsecase.UpdateNextUpdateAt(ctx, search)
	if err != nil {
		return fmt.Errorf("failed to save next search update at: %w", err)
	}

	err = u.runCollection(ctx, search, input)
	if err != nil {
		return fmt.Errorf("failed to collect tweets: %w", err)
	}

	return nil
}

func (u *batchUsecase) prepareSearch(ctx context.Context, searchID model.SearchID, userID model.UserID) (*model.Search, *twitter.SearchInput, error) {
	search, err := u.searchUsecase.Find(ctx, searchID, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch search: %w", err)
	}
	if search == nil {
		return nil, nil, fmt.Errorf("specified search (id: %s) not found: %w", searchID, err)
	}

	user, err := u.userUsecase.FindByID(ctx, search.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	latestTweetID, err := u.tweetRepository.LatestTweetID(ctx, search.SearchID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, nil, fmt.Errorf("failed to get latest collected tweet id: %w", err)
	}
	input := &twitter.SearchInput{
		Query:                    search.Query,
		TwitterAccessToken:       user.TwitterAccessToken,
		TwitterAccessTokenSecret: user.TwitterAccessTokenSecret,
		SinceID:                  int64(latestTweetID),
	}

	return search, input, nil
}

func (u *batchUsecase) runCollection(ctx context.Context, search *model.Search, input *twitter.SearchInput) error {
	tweetsBuf := []*twitter.Tweet{}

	for {
		tweets, err := u.twitterClient.Search(ctx, input)
		if err != nil {
			return fmt.Errorf("twitter search error: %w", err)
		}
		// MaxID option includes itself
		if input.MaxID != 0 && len(tweets) > 0 {
			tweets = tweets[1:]
		}

		if len(tweets) == 0 {
			break
		}

		for i := 0; i < len(tweets); i++ {
			tweet := &tweets[i]

			if shouldDetectScore(tweet) {
				tweetsBuf = append(tweetsBuf, tweet)
				if len(tweetsBuf) == 25 {
					err = u.batchStoreTweetsWithDetection(ctx, search, tweetsBuf)
					if err != nil {
						return fmt.Errorf("failed to batch store tweets with sentiment detection: %w", err)
					}
					tweetsBuf = tweetsBuf[:0]
				}
			}
		}

		input.MaxID = tweets[len(tweets)-1].TweetID
	}

	if len(tweetsBuf) > 0 {
		err := u.batchStoreTweetsWithDetection(ctx, search, tweetsBuf)
		if err != nil {
			return fmt.Errorf("failed to batch store tweets with sentiment detection: %w", err)
		}
	}

	return nil
}

func (u *batchUsecase) batchStoreTweetsWithDetection(ctx context.Context, search *model.Search, tweets []*twitter.Tweet) error {
	log.Printf("Writing %d tweets with sentiment detection.\n", len(tweets))

	textList := []*string{}
	for _, tweet := range tweets {
		textList = append(textList, &tweet.Text)
	}

	detectOutputs, err := u.sentimentDetector.BatchDetect(ctx, textList)
	if err != nil {
		return fmt.Errorf("failed to batch detect sentiment score %w", err)
	}

	modelTweets := []*model.Tweet{}

	for i, tweet := range tweets {
		detectOutput := detectOutputs[i]

		modelTweets = append(modelTweets, &model.Tweet{
			TweetID:  model.TweetID(tweet.TweetID),
			SearchID: search.SearchID,
			AuthorID: tweet.AuthorID,
			User: &model.TwitterUser{
				ID:              tweet.User.ID,
				Name:            tweet.User.Name,
				ScreenName:      tweet.User.ScreenName,
				ProfileImageURL: tweet.User.ProfileImageURL,
			},
			Entities:           tweet.Entities,
			Text:               tweet.Text,
			SentimentScore:     &detectOutput.Score,
			SentimentLabel:     detectOutput.Label,
			ExpirationUnixTime: time.Now().AddDate(0, 6, 0).Unix(),
			TweetCreatedAt:     tweet.CreatedAt,
		})
	}

	err = u.tweetRepository.BatchStore(ctx, modelTweets)
	if err != nil {
		return fmt.Errorf("failed to batch store tweets: %w", err)
	}

	return nil
}

func shouldDetectScore(tweet *twitter.Tweet) bool {
	ngURLs := []string{"youtu.be", "youtube.com", "nicovideo", "nico.ms", "peing.net"}

	for _, url := range tweet.Entities.URLs {
		for _, ngURL := range ngURLs {
			if strings.Contains(url.ExpandedURL, ngURL) {
				return false
			}
		}
	}

	textLen := len(tweet.Text)

	for _, url := range tweet.Entities.URLs {
		textLen -= len(url.URL)
	}
	for _, men := range tweet.Entities.Mentions {
		textLen -= len(men.Tag)
	}
	for _, hash := range tweet.Entities.HashTags {
		textLen -= len(hash.Tag)
	}

	return textLen >= 160
}
