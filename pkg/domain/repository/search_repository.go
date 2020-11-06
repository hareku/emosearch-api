package repository

import (
	"context"

	"github.com/hareku/emosearch-api/pkg/domain/model"
)

// SearchRepository provides CRUD methods for Search domain.
type SearchRepository interface {
	ListByUserID(ctx context.Context, userID model.UserID) ([]*model.Search, error)
	Create(ctx context.Context, search *model.Search) error
}
