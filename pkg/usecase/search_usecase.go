package usecase

import (
	"context"
	"fmt"

	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

// SearchUsecase provides usecases of Search domain.
type SearchUsecase interface {
	ListByUserID(ctx context.Context, userID model.UserID) ([]*model.Search, error)
	Create(ctx context.Context, input *SearchUsecaseCreateInput) (*model.Search, error)
}

type searchUsecase struct {
	searchRepository repository.SearchRepository
}

// NewSearchUsecase creates SearchUsecase.
func NewSearchUsecase(searchRepository repository.SearchRepository) SearchUsecase {
	return &searchUsecase{searchRepository}
}

func (u *searchUsecase) ListByUserID(ctx context.Context, userID model.UserID) ([]*model.Search, error) {
	searches, err := u.searchRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cloudn't get searches: %w", err)
	}

	return searches, nil
}

// SearchUsecaseCreateInput is the input of SearchUsecase.Create().
type SearchUsecaseCreateInput struct {
	UserID model.UserID
	Title  string
	Query  string
}

func (u *searchUsecase) Create(ctx context.Context, input *SearchUsecaseCreateInput) (*model.Search, error) {
	search := &model.Search{
		UserID: input.UserID,
		Title:  input.Title,
		Query:  input.Query,
	}

	err := u.searchRepository.Create(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("search creation error: %w", err)
	}

	return search, nil
}
