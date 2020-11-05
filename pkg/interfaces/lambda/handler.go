package lambda

import (
	"context"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hareku/emosearch-api/pkg/registry"
)

type handler struct {
	registry registry.Registry
	router   *lmdrouter.Router
}

// Handler provides the gate of AWS Lambda.
type Handler interface {
	Start()
}

// NewLambdaHandler returns an instance of LambdaHandler.
func NewLambdaHandler(registry registry.Registry) Handler {
	h := &handler{
		registry,
		lmdrouter.NewRouter("/v1", authMiddleware(registry.NewAuthenticator()), corsMiddleware()),
	}

	h.registerRoutes()

	return h
}

// Start Lambda function.
func (h *handler) Start() {
	lambda.Start(h.router.Handler)
}

func (h *handler) registerRoutes() {
	h.router.Route("GET", "/searches", h.fetchSearches())
}

func returnInternalServerError() (events.APIGatewayProxyResponse, error) {
	return lmdrouter.HandleError(lmdrouter.HTTPError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error.",
	})
}

func (h *handler) fetchSearches() lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		u := h.registry.NewSearchUsecase()
		userID, err := h.registry.NewAuthenticator().UserID(ctx)
		if err != nil {
			return returnInternalServerError()
		}

		searches, err := u.ListByUserID(ctx, userID)
		if err != nil {
			return returnInternalServerError()
		}

		return lmdrouter.MarshalResponse(http.StatusOK, nil, searches)
	}
}
