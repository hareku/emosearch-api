package repository

import "github.com/hareku/emosearch-api/pkg/domain/model"

// SearchRepository provides CRUD methods for Search domain.
type SearchRepository interface {
	ListByUserID(userID model.UserID) ([]model.Search, error)
}
