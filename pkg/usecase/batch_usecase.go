package usecase

import (
	"context"
	"strconv"

	"github.com/hareku/emosearch-api/pkg/domain/auth"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
	"github.com/hareku/emosearch-api/pkg/domain/twitter"
)

// BatchUsecase provides usecases of Batch domain.
type BatchUsecase interface {
	UpdateAllUserSearches(ctx context.Context) error
}

type batchUsecase struct {
	authenticator     auth.Authenticator
	userUsecase       UserUsecase
	searchUsecase     SearchUsecase
	twitterClient     twitter.Client
	sentimentDetector sentiment.Detector
}

// NewBatchUsecaseInput is the input of NewBatchUsecase.
type NewBatchUsecaseInput struct {
	Authenticator     auth.Authenticator
	UserUsecase       UserUsecase
	SearchUsecase     SearchUsecase
	TwitterClient     twitter.Client
	SentimentDetector sentiment.Detector
}

// NewBatchUsecase creates BatchUsecase.
func NewBatchUsecase(input *NewBatchUsecaseInput) BatchUsecase {
	return &batchUsecase{
		authenticator:     input.Authenticator,
		userUsecase:       input.UserUsecase,
		searchUsecase:     input.SearchUsecase,
		twitterClient:     input.TwitterClient,
		sentimentDetector: input.SentimentDetector,
	}
}

func (u *batchUsecase) UpdateAllUserSearches(ctx context.Context) error {
	ids, nextToken, err := u.authenticator.ListUserID(ctx, "")
	for {
		if len(ids) == 0 {
			break
		}
		if err != nil {
			return err
		}

		for i := 0; i < len(ids); i++ {
			err = u.updateUserSearches(ctx, ids[i])
			if err != nil {
				return err
			}
		}

		ids, nextToken, err = u.authenticator.ListUserID(ctx, nextToken)
	}

	return nil
}

func (u *batchUsecase) updateUserSearches(ctx context.Context, userID model.UserID) error {
	searches, err := u.searchUsecase.ListByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for i := 0; i < len(searches); i++ {
		err = u.collectTweets(ctx, searches[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *batchUsecase) collectTweets(ctx context.Context, search *model.Search) error {
	user, err := u.userUsecase.FindByID(ctx, search.UserID)
	if err != nil {
		return err
	}

	input := twitter.SearchInput{
		Query:                    search.Query,
		TwitterAccessToken:       user.TwitterAccessToken,
		TwitterAccessTokenSecret: user.TwitterAccessTokenSecret,
		// TODO: set SinceID by stored latest tweet.
	}

	tweets, err := u.twitterClient.Search(ctx, &input)
	for {
		if len(tweets) == 0 {
			break
		}
		if err != nil {
			return err
		}

		for i := 0; i < len(tweets); i++ {
			err = u.storeTweet(ctx, &tweets[i])
			if err != nil {
				return err
			}
		}

		maxID, err := strconv.ParseInt(tweets[len(tweets)-1].TweetID, 10, 64)
		if err != nil {
			return err
		}
		input.MaxID = maxID

		tweets, err = u.twitterClient.Search(ctx, &input)
	}

	return nil
}

func (u *batchUsecase) storeTweet(ctx context.Context, tweet *twitter.Tweet) error {
	// score, err := u.sentimentDetector.Detect(ctx, tweet.Text)
	// if err != nil {
	// 	return err
	// }
	return nil
}