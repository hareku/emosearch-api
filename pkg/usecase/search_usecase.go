package usecase

import (
	"context"
	"fmt"

	"github.com/hareku/emosearch-api/pkg/domain/auth"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

// SearchUsecase provides usecases of Search domain.
type SearchUsecase interface {
	ListUserSearches(ctx context.Context) ([]*model.Search, error)
	Create(ctx context.Context, input *SearchUsecaseCreateInput) (*model.Search, error)
}

type searchUsecase struct {
	authenticator    auth.Authenticator
	searchRepository repository.SearchRepository
}

// NewSearchUsecase creates SearchUsecase.
func NewSearchUsecase(authenticator auth.Authenticator, searchRepository repository.SearchRepository) SearchUsecase {
	return &searchUsecase{authenticator, searchRepository}
}

func (u *searchUsecase) ListUserSearches(ctx context.Context) ([]*model.Search, error) {
	userID, err := u.authenticator.UserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching user id error: %w", err)
	}

	searches, err := u.searchRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching user's searches error: %w", err)
	}

	return searches, nil
}

// SearchUsecaseCreateInput is the input of SearchUsecase.Create().
type SearchUsecaseCreateInput struct {
	Title string
	Query string
}

func (u *searchUsecase) Create(ctx context.Context, input *SearchUsecaseCreateInput) (*model.Search, error) {
	userID, err := u.authenticator.UserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching user id error: %w", err)
	}

	search := &model.Search{
		UserID: userID,
		Title:  input.Title,
		Query:  input.Query,
	}

	err = u.searchRepository.Create(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("creating search error: %w", err)
	}

	return search, nil
}
