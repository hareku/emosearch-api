package api

import (
	"context"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

func (h *handler) registerTweetRoutes() {
	h.router.Route("GET", "/searches/:search_id/tweets", h.fetchTweetes())
}

type fetchTweetsInput struct {
	SearchID model.SearchID `lambda:"path.search_id"`
	UntilID  model.TweetID  `lambda:"query.until_id"`
	Limit    int64          `lambda:"query.limit"`
}

func (h *handler) fetchTweetes() lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		var input fetchTweetsInput
		err = lmdrouter.UnmarshalRequest(req, true, &input)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		u := h.registry.NewSearchUsecase()
		search, err := u.GetUserSearch(ctx, input.SearchID)
		if err != nil {
			return lmdrouter.HandleError(err)
		}
		if search == nil {
			return lmdrouter.HandleError(lmdrouter.HTTPError{
				Code:    http.StatusNotFound,
				Message: "specified search was not found",
			})
		}

		r := h.registry.NewTweetRepository()
		tweets, err := r.List(ctx, &repository.TweetRepositoryListInput{
			SearchID: input.SearchID,
			UntilID:  input.UntilID,
			Limit:    input.Limit,
		})
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		return lmdrouter.MarshalResponse(http.StatusOK, nil, tweets)
	}
}
