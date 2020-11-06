package repository

import (
	"context"

	"github.com/hareku/emosearch-api/pkg/domain/model"
)

// UserRepository provides CRUD methods for User domain.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, userID model.UserID) (*model.User, error)
}
