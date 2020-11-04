package repository

import "github.com/hareku/emosearch-api/pkg/domain/model"

// SearchRepository provides CRUD methods for Search domain.
type SearchRepository interface {
	Create(model.Search) (bool, error)
	FindByID(model.SearchID) (model.Search, error)
}
