package repository

import (
	"context"

	"github.com/hareku/emosearch-api/pkg/domain/model"
)

// SearchRepositoryListInput is the input of List method.
type SearchRepositoryListInput struct {
	Limit int64
}

// SearchRepository provides CRUD methods for Search domain.
type SearchRepository interface {
	List(ctx context.Context, input SearchRepositoryListInput) ([]*model.Search, error)
	ListByUserID(ctx context.Context, userID model.UserID) ([]*model.Search, error)
	Find(ctx context.Context, userID model.UserID, searchID model.SearchID) (*model.Search, error)
	Create(ctx context.Context, search *model.Search) error
	Update(ctx context.Context, search *model.Search) error
}
