package repository

import "github.com/hareku/emosearch-api/pkg/domain"

// SearchRepository provides CRUD methods for Search domain.
type SearchRepository interface {
	Create(domain.Search) (bool, error)
	FindByID(domain.SearchID) (domain.Search, error)
}
