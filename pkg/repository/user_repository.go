package repository

import "github.com/hareku/emosearch-api/pkg/domain"

// UserRepository provides CRUD methods for User domain.
type UserRepository interface {
	Store(domain.User) (bool, error)
	FindByID(domain.UserID) (domain.User, error)
}
