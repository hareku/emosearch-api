package usecase

import (
	"context"
	"errors"
	"fmt"
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
	search, err := u.searchUsecase.Find(ctx, searchID, userID)
	if err != nil {
		return err
	}
	if search == nil {
		return nil
	}

	user, err := u.userUsecase.FindByID(ctx, search.UserID)
	if err != nil {
		return err
	}

	input := twitter.SearchInput{
		Query:                    search.Query,
		TwitterAccessToken:       user.TwitterAccessToken,
		TwitterAccessTokenSecret: user.TwitterAccessTokenSecret,
	}

	latestTweetID, err := u.tweetRepository.LatestTweetID(ctx, search.SearchID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}
	input.SinceID = int64(latestTweetID)

	tweets, err := u.twitterClient.Search(ctx, &input)
	if err != nil {
		return err
	}
	if len(tweets) == 0 {
		return nil
	}

	for {
		if err != nil {
			return err
		}
		// break because Input.MaxID returns results with an ID less than (that is, older than) or equal to the specified ID.
		if len(tweets) == 1 {
			break
		}

		for i := 0; i < len(tweets); i++ {
			err = u.storeTweet(ctx, search, &tweets[i])
			if err != nil {
				return err
			}
		}

		input.MaxID = tweets[len(tweets)-1].TweetID
		tweets, err = u.twitterClient.Search(ctx, &input)
	}

	err = u.searchUsecase.UpdateNextUpdateAt(ctx, search)
	if err != nil {
		return err
	}
	return nil
}

func (u *batchUsecase) storeTweet(ctx context.Context, search *model.Search, tweet *twitter.Tweet) error {
	dtweet := model.Tweet{
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
		SentimentScore:     nil,
		SentimentLabel:     sentiment.LabelUndetected,
		ExpirationUnixTime: time.Now().AddDate(0, 0, 14).Unix(),
		TweetCreatedAt:     tweet.CreatedAt,
	}

	if shouldDetectScore(tweet) {
		output, err := u.sentimentDetector.Detect(ctx, tweet.Text)
		if err != nil {
			return fmt.Errorf("failed to detect sentiment score of a tweet: %w", err)
		}
		dtweet.SentimentScore = &output.Score
		dtweet.SentimentLabel = output.Label
		dtweet.ExpirationUnixTime = time.Now().AddDate(0, 3, 0).Unix()
	}

	err := u.tweetRepository.Store(ctx, &dtweet)
	if err != nil {
		return fmt.Errorf("failed to store a tweet: %w", err)
	}

	return nil
}

func shouldDetectScore(tweet *twitter.Tweet) bool {
	ngURLs := []string{"youtu.be", "youtube.com", "nicovideo", "peing.net"}

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
	for _, med := range tweet.Entities.Media {
		textLen -= len(med.URL)
	}
	for _, men := range tweet.Entities.Mentions {
		textLen -= len(men.Tag)
	}
	for _, hash := range tweet.Entities.HashTags {
		textLen -= len(hash.Tag)
	}

	return textLen >= 160
}
