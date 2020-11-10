package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/validator"
	"github.com/hareku/emosearch-api/pkg/usecase"
)

func (h *handler) registerSearchRoutes() {
	h.router.Route("GET", "/searches", h.fetchSearches())
	h.router.Route("GET", "/searches/:id", h.fetchSearch())
	h.router.Route("DELETE", "/searches/:id", h.deleteSearch())
	h.router.Route("POST", "/searches", h.createSearch())
}

func (h *handler) fetchSearches() lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		u := h.registry.NewSearchUsecase()
		searches, err := u.ListUserSearches(ctx)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		return lmdrouter.MarshalResponse(http.StatusOK, nil, searches)
	}
}

type fetchSearchInput struct {
	SearchID model.SearchID `lambda:"path.id"`
}

func (h *handler) fetchSearch() lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		var input fetchSearchInput
		err = lmdrouter.UnmarshalRequest(req, false, &input)
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

		return lmdrouter.MarshalResponse(http.StatusOK, nil, search)
	}
}

type deleteSearchInput struct {
	SearchID model.SearchID `lambda:"path.id"`
}

func (h *handler) deleteSearch() lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		var input deleteSearchInput
		err = lmdrouter.UnmarshalRequest(req, false, &input)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		u := h.registry.NewSearchUsecase()
		err = u.DeleteUserSearch(ctx, input.SearchID)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		return lmdrouter.MarshalResponse(http.StatusNoContent, nil, nil)
	}
}

type createSearchInput struct {
	Query string `json:"Query"`
}

func (h *handler) createSearch() lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		var input createSearchInput
		err = lmdrouter.UnmarshalRequest(req, true, &input)
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		u := h.registry.NewSearchUsecase()
		search, err := u.Create(ctx, &usecase.SearchUsecaseCreateInput{
			Query: input.Query,
		})
		var errv validator.ErrValidation
		if errors.As(err, &errv) {
			return h.handleValidationErrors(errv)
		}
		if err != nil {
			return lmdrouter.HandleError(err)
		}

		return lmdrouter.MarshalResponse(http.StatusCreated, nil, search)
	}
}
