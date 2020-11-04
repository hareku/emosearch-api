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
	reg    registry.Registry
	router *lmdrouter.Router
}

// Handler provides the gate of AWS Lambda.
type Handler interface {
	Start()
}

// NewLambdaHandler returns an instance of LambdaHandler.
func NewLambdaHandler(reg registry.Registry) Handler {
	h := &handler{
		reg,
		lmdrouter.NewRouter("/v1/", authMiddleware(reg.NewAuthenticator())),
	}

	h.registerRoutes()

	return h
}

// Start Lambda function.
func (h *handler) Start() {
	lambda.Start(h.router.Handler)
}

func (h *handler) registerRoutes() {
	h.router.Route("GET", "/@me", fetchMe)
}

type fetchMeInput struct {
	Authorization string `lambda:"header.Authorization"`
}

func fetchMe(ctx context.Context, req events.APIGatewayProxyRequest) (
	res events.APIGatewayProxyResponse,
	err error,
) {
	var input fetchMeInput
	err = lmdrouter.UnmarshalRequest(req, false, &input)
	if err != nil {
		return lmdrouter.HandleError(err)
	}

	var data interface{}
	return lmdrouter.MarshalResponse(http.StatusOK, nil, data)
}