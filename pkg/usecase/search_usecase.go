package usecase

import (
	"fmt"

	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

// SearchUsecase provides usecases of Search domain.
type SearchUsecase interface {
	ListByUserID(userID model.UserID) ([]model.Search, error)
}

type searchUsecase struct {
	searchRepository repository.SearchRepository
}

// NewSearchUsecase creates SearchUsecase.
func NewSearchUsecase(searchRepository repository.SearchRepository) SearchUsecase {
	return &searchUsecase{searchRepository}
}

func (u *searchUsecase) ListByUserID(userID model.UserID) ([]model.Search, error) {
	searches, err := u.searchRepository.ListByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("cloudn't get searches: %w", err)
	}

	return searches, nil
}
