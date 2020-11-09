package api

import (
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hareku/emosearch-api/pkg/domain/validator"
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
		lmdrouter.NewRouter("/v1", loggerMiddleware, authMiddleware(registry.NewAuthenticator()), corsMiddleware()),
	}

	h.registerRoutes()

	return h
}

// Start Lambda function.
func (h *handler) Start() {
	lambda.Start(h.router.Handler)
}

func (h *handler) registerRoutes() {
	h.registerSearchRoutes()
	h.registerUserRoutes()
	h.registerTweetRoutes()
}

func (h *handler) handleValidationErrors(verr validator.ErrValidation) (events.APIGatewayProxyResponse, error) {
	body := verr.ToMap()
	return lmdrouter.MarshalResponse(http.StatusUnprocessableEntity, nil, body)
}
