package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
	"github.com/hareku/emosearch-api/pkg/domain/twitter"
)

// BatchUsecase provides usecases of Batch domain.
type BatchUsecase interface {
	CollectTweets(ctx context.Context, searchID model.SearchID) error
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

func (u *batchUsecase) CollectTweets(ctx context.Context, searchID model.SearchID) error {
	// search, err := u.searchUsecase.

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

	return nil
}

func (u *batchUsecase) storeTweet(ctx context.Context, search *model.Search, tweet *twitter.Tweet) error {
	score, err := u.sentimentDetector.Detect(ctx, tweet.Text)
	if err != nil {
		return err
	}

	dtweet := model.Tweet{
		TweetID:        model.TweetID(tweet.TweetID),
		SearchID:       search.SearchID,
		AuthorID:       tweet.AuthorID,
		Text:           tweet.Text,
		SentimentScore: score,
		TweetCreatedAt: tweet.CreatedAt,
	}

	err = u.tweetRepository.Store(ctx, &dtweet)
	if err != nil {
		return fmt.Errorf("tweet storing error: %w", err)
	}

	return nil
}
