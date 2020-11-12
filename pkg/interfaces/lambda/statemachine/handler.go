package statemachine

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
	"github.com/hareku/emosearch-api/pkg/registry"
)

type handler struct {
	registry registry.Registry
}

// Handler provides the gate of AWS Lambda.
type Handler interface {
	StartListSearches()
	StartCollectTweets()
}

// New returns an instance of Handler.
func New(registry registry.Registry) Handler {
	return &handler{registry}
}

func (h *handler) StartListSearches() {
	lambda.Start(h.listSearchesHandler)
}

// StartListSearchesRes is the response of ListSearches lambda function.
type StartListSearchesRes struct {
	Events []StartListSearchesResEvent `json:"events"`
}

// StartListSearchesResEvent is events of StartListSearchesRes.
type StartListSearchesResEvent struct {
	SearchID model.SearchID `json:"search_id"`
	UserID   model.UserID   `json:"user_id"`
}

func (h *handler) listSearchesHandler(ctx context.Context) (*StartListSearchesRes, error) {
	usc := h.registry.NewSearchUsecase()
	searches, err := usc.ListShouldUpdateSearches(ctx)

	res := &StartListSearchesRes{
		Events: []StartListSearchesResEvent{},
	}

	if errors.Is(err, repository.ErrNotFound) {
		return res, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list searches")
	}

	for _, search := range searches {
		res.Events = append(res.Events, StartListSearchesResEvent{
			SearchID: search.SearchID,
			UserID:   search.UserID,
		})
	}

	return res, nil
}

func (h *handler) StartCollectTweets() {
	lambda.Start(h.collectTweetsHandler)
}

func (h *handler) collectTweetsHandler(ctx context.Context, event StartListSearchesResEvent) error {
	return h.registry.NewBatchUsecase().CollectTweets(ctx, event.SearchID, event.UserID)
}
