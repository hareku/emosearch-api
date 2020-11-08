package batch

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hareku/emosearch-api/pkg/registry"
)

type handler struct {
	registry registry.Registry
}

// Handler provides the gate of AWS Lambda.
type Handler interface {
	StartSearch()
}

// NewLambdaHandler returns an instance of Handler.
func NewLambdaHandler(registry registry.Registry) Handler {
	return &handler{registry}
}

// Start Lambda function.
func (h *handler) StartSearch() {
	lambda.Start(h.registry.NewBatchUsecase().UpdateAllUserSearches)
}
