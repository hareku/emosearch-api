package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/hareku/emosearch-api/pkg/domain/validator"
	"github.com/hareku/emosearch-api/pkg/usecase"
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
		usecase := h.registry.NewUserUsecase()

		user, err := usecase.FetchAuthUser(ctx)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		return lmdrouter.MarshalResponse(http.StatusCreated, nil, user)
	}
}

type registerUserInput struct {
	TwitterAccessToken       string `json:"TwitterAccessToken" validate:"required,gt=1,lte=500"`
	TwitterAccessTokenSecret string `json:"TwitterAccessTokenSecret" validate:"required,gt=1,lte=500"`
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

		v := h.registry.NewValidator()
		var verr validator.ErrValidation
		if errors.As(v.StructCtx(ctx, input), &verr) {
			return h.handleValidationErrors(verr)
		}

		u := h.registry.NewUserUsecase()
		user, err := u.Register(ctx, usecase.UserUsecaseRegisterInput{
			TwitterAccessToken:       input.TwitterAccessToken,
			TwitterAccessTokenSecret: input.TwitterAccessTokenSecret,
		})
		if errors.Is(err, usecase.ErrUserAlreadyExist) {
			return lmdrouter.MarshalResponse(http.StatusOK, nil, user)
		}
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		return lmdrouter.MarshalResponse(http.StatusCreated, nil, user)
	}
}
