package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hareku/emosearch-api/pkg/domain/auth"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
	"github.com/hareku/emosearch-api/pkg/domain/validator"
)

// SearchUsecase provides usecases of Search domain.
type SearchUsecase interface {
	ListShouldUpdateSearches(ctx context.Context) ([]*model.Search, error)
	ListByUserID(ctx context.Context, userID model.UserID) ([]*model.Search, error)
	ListUserSearches(ctx context.Context) ([]*model.Search, error)
	Find(ctx context.Context, searchID model.SearchID, userID model.UserID) (*model.Search, error)
	GetUserSearch(ctx context.Context, searchID model.SearchID) (*model.Search, error)
	DeleteUserSearch(ctx context.Context, searchID model.SearchID) error
	Create(ctx context.Context, input *SearchUsecaseCreateInput) (*model.Search, error)
	UpdateNextUpdateAt(ctx context.Context, search *model.Search) error
}

type searchUsecase struct {
	authenticator    auth.Authenticator
	validator        validator.Validator
	searchRepository repository.SearchRepository
}

// NewSearchUsecase creates SearchUsecase.
func NewSearchUsecase(authenticator auth.Authenticator, validator validator.Validator, searchRepository repository.SearchRepository) SearchUsecase {
	return &searchUsecase{authenticator, validator, searchRepository}
}

func (u *searchUsecase) ListShouldUpdateSearches(ctx context.Context) ([]*model.Search, error) {
	now := time.Now()
	searches, err := u.searchRepository.List(ctx, repository.SearchRepositoryListInput{
		Limit:                   100,
		UntilNextSearchUpdateAt: &now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list searches: %w", err)
	}

	return searches, nil
}

func (u *searchUsecase) ListByUserID(ctx context.Context, userID model.UserID) ([]*model.Search, error) {
	searches, err := u.searchRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user(id: %s) searches: %w", string(userID), err)
	}

	return searches, nil
}

func (u *searchUsecase) ListUserSearches(ctx context.Context) ([]*model.Search, error) {
	userID, err := u.authenticator.UserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user id: %w", err)
	}

	return u.ListByUserID(ctx, userID)
}

func (u *searchUsecase) GetUserSearch(ctx context.Context, searchID model.SearchID) (*model.Search, error) {
	userID, err := u.authenticator.UserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user id: %w", err)
	}

	search, err := u.searchRepository.Find(ctx, userID, searchID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search (id: %v): %w", searchID, err)
	}
	return search, nil
}

func (u *searchUsecase) DeleteUserSearch(ctx context.Context, searchID model.SearchID) error {
	userID, err := u.authenticator.UserID(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch user id: %w", err)
	}

	search, err := u.searchRepository.Find(ctx, userID, searchID)
	if errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("search was not found (id: %v): %w", searchID, err)
	}
	if err != nil {
		return fmt.Errorf("failed to fetch search (id: %v): %w", searchID, err)
	}

	err = u.searchRepository.Delete(ctx, search)
	if err != nil {
		return fmt.Errorf("failed to delete search (id: %v): %w", searchID, err)
	}
	return nil
}

func (u *searchUsecase) Find(ctx context.Context, searchID model.SearchID, userID model.UserID) (*model.Search, error) {
	search, err := u.searchRepository.Find(ctx, userID, searchID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search (id: %v): %w", searchID, err)
	}
	return search, nil
}

// SearchUsecaseCreateInput is the input of SearchUsecase.Create().
type SearchUsecaseCreateInput struct {
	Query string `validate:"required,gte=3,lte=100"`
}

func (u *searchUsecase) Create(ctx context.Context, input *SearchUsecaseCreateInput) (*model.Search, error) {
	err := u.validator.StructCtx(ctx, input)
	if err != nil {
		return nil, err
	}

	userID, err := u.authenticator.UserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching user id error: %w", err)
	}

	search := &model.Search{
		UserID:              userID,
		Title:               "",
		Query:               input.Query,
		LastSearchUpdatedAt: nil,
		NextSearchUpdateAt:  time.Now().AddDate(-1, 0, 0),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	err = u.searchRepository.Create(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("creating search error: %w", err)
	}

	return search, nil
}

func (u *searchUsecase) UpdateNextUpdateAt(ctx context.Context, search *model.Search) error {
	now := time.Now()
	search.LastSearchUpdatedAt = &now
	search.NextSearchUpdateAt = time.Now().Add(30 * time.Minute)
	err := u.searchRepository.Update(ctx, search)

	if err != nil {
		return fmt.Errorf("failed to update search next updated at: %w", err)
	}

	return nil
}
