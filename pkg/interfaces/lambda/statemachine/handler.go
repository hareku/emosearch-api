package statemachine

import (
	"context"
	"errors"
	"time"

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

// NewHandler returns an instance of Handler.
func NewHandler(registry registry.Registry) Handler {
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
}

func (h *handler) listSearchesHandler(ctx context.Context) (*StartListSearchesRes, error) {
	r := h.registry.NewSearchRepository()
	ids, err := r.List(ctx, &repository.SearchRepositoryListInput{
		Limit:              100,
		UntilLastUpdatedAt: time.Now().Add(30 * time.Minute),
	})

	res := &StartListSearchesRes{
		Events: []StartListSearchesResEvent{},
	}

	if errors.Is(err, repository.ErrNotFound) {
		return res, nil
	}
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		res.Events = append(res.Events, StartListSearchesResEvent{id})
	}

	return res, nil
}

// CollectTweetsEvent is the event of CollectTweets lambda function.
type CollectTweetsEvent struct {
	SearchID model.SearchID `json:"search_id"`
}

func (h *handler) StartCollectTweets() {
	lambda.Start(h.collectTweetsHandler)
}

func (h *handler) collectTweetsHandler(ctx context.Context, eve CollectTweetsEvent) error {
	return nil
}
