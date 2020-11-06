package lambda

import (
	"context"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/hareku/emosearch-api/pkg/domain/model"
)

func (h *handler) registerUserRoutes() {
	h.router.Route("GET", "/users/@me", h.fetchMe())
	h.router.Route("POST", "/users/@me", h.registerUser())
}

func (h *handler) fetchMe() lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		u := h.registry.NewUserRepository()
		userID, err := h.registry.NewAuthenticator().UserID(ctx)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		user, err := u.FindByID(ctx, userID)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		return lmdrouter.MarshalResponse(http.StatusCreated, nil, user)
	}
}

type registerUserInput struct {
	TwitterAccessToken       string `json:"twitter_access_token"`
	TwitterAccessTokenSecret string `json:"twitter_access_token_secret"`
}

func (h *handler) registerUser() lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		var input registerUserInput
		err = lmdrouter.UnmarshalRequest(req, true, &input)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		u := h.registry.NewUserRepository()
		userID, err := h.registry.NewAuthenticator().UserID(ctx)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		user, err := u.FindByID(ctx, userID)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		if user == nil {
			user = &model.User{
				UserID:                   userID,
				TwitterAccessToken:       input.TwitterAccessToken,
				TwitterAccessTokenSecret: input.TwitterAccessTokenSecret,
			}

			err = u.Create(ctx, user)
			if err != nil {
				return lmdrouter.HandleError(err)
			}
		}

		return lmdrouter.MarshalResponse(http.StatusOK, nil, user)
	}
}
