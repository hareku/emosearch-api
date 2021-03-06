package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
)

func (h *handler) registerTweetRoutes() {
	h.router.Route("GET", "/searches/:search_id/tweets", h.fetchTweets())
}

type fetchTweetsInput struct {
	SearchID       model.SearchID `lambda:"path.search_id"`
	UntilID        model.TweetID  `lambda:"query.until_id"`
	Limit          int64          `lambda:"query.limit"`
	SentimentLabel string         `lambda:"query.sentiment_label"`
}

type fetchTweetsRes struct {
	Tweets  []model.Tweet
	HasMore bool
}

func (h *handler) fetchTweets() lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		var input fetchTweetsInput
		err = lmdrouter.UnmarshalRequest(req, false, &input)
		if err != nil {
			return lmdrouter.HandleError(fmt.Errorf("failed to parse input: %w", err))
		}

		u := h.registry.NewSearchUsecase()
		search, err := u.GetUserSearch(ctx, input.SearchID)
		if err != nil {
			return lmdrouter.HandleError(fmt.Errorf("failed to fetch user search: %w", err))
		}
		if search == nil {
			return lmdrouter.HandleError(lmdrouter.HTTPError{
				Code:    http.StatusNotFound,
				Message: "specified search was not found",
			})
		}

		r := h.registry.NewTweetRepository()
		listInput := &repository.TweetRepositoryListInput{
			SearchID:       input.SearchID,
			UntilID:        input.UntilID,
			Limit:          input.Limit + 1,
			SentimentLabel: nil,
		}
		if input.SentimentLabel != "" {
			label := sentiment.Label(input.SentimentLabel)
			listInput.SentimentLabel = &label
		}

		tweets, err := r.List(ctx, listInput)
		if err != nil {
			return lmdrouter.HandleError(fmt.Errorf("failed to fetch tweets: %w", err))
		}

		pagination := fetchTweetsRes{
			Tweets:  tweets,
			HasMore: false,
		}
		if len(pagination.Tweets) > int(input.Limit) {
			pagination.Tweets = tweets[:len(tweets)-1]
			pagination.HasMore = true
		}

		return lmdrouter.MarshalResponse(http.StatusOK, nil, pagination)
	}
}
