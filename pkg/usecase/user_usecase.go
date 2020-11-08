package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/hareku/emosearch-api/pkg/domain/auth"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

var (
	// ErrUserAlreadyExist is returned when user already exist in registration.
	ErrUserAlreadyExist = errors.New("user already exist")
)

// UserUsecase provides usecases of User domain.
type UserUsecase interface {
	FetchAuthUser(ctx context.Context) (*model.User, error)
	Register(ctx context.Context, input UserUsecaseRegisterInput) (*model.User, error)
}

type userUsecase struct {
	authenticator  auth.Authenticator
	userRepository repository.UserRepository
}

// NewUserUsecase creates UserUsecase.
func NewUserUsecase(authenticator auth.Authenticator, userRepository repository.UserRepository) UserUsecase {
	return &userUsecase{
		authenticator,
		userRepository,
	}
}

func (u *userUsecase) FetchAuthUser(ctx context.Context) (*model.User, error) {
	userID, err := u.authenticator.UserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get user id: %w", err)
	}

	user, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("could not get user from repository: %w", err)
	}

	return user, nil
}

// UserUsecaseRegisterInput represents the input of Register method.
type UserUsecaseRegisterInput struct {
	TwitterAccessToken       string
	TwitterAccessTokenSecret string
}

func (u *userUsecase) Register(ctx context.Context, input UserUsecaseRegisterInput) (*model.User, error) {
	userID, err := u.authenticator.UserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching user id error: %w", err)
	}

	user, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching user error: %w", err)
	}
	if user != nil {
		return user, ErrUserAlreadyExist
	}

	user = &model.User{
		UserID:                   userID,
		TwitterAccessToken:       input.TwitterAccessToken,
		TwitterAccessTokenSecret: input.TwitterAccessTokenSecret,
	}

	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("user registration error: %w", err)
	}

	return user, nil
}
