package repository

import "github.com/hareku/emosearch-api/pkg/domain/model"

// UserRepository provides CRUD methods for User domain.
type UserRepository interface {
	Store(model.User) (bool, error)
	FindByID(model.UserID) (model.User, error)
}
