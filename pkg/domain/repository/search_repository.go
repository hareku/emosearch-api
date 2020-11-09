package repository

import (
	"context"
	"time"

	"github.com/hareku/emosearch-api/pkg/domain/model"
)

// SearchRepositoryListInput is the input of List method.
type SearchRepositoryListInput struct {
	Limit              int64
	UntilLastUpdatedAt time.Time
}

// SearchRepository provides CRUD methods for Search domain.
type SearchRepository interface {
	List(ctx context.Context, input *SearchRepositoryListInput) ([]model.SearchID, error)
	ListByUserID(ctx context.Context, userID model.UserID) ([]*model.Search, error)
	Find(ctx context.Context, userID model.UserID, searchID model.SearchID) (*model.Search, error)
	Create(ctx context.Context, search *model.Search) error
}
