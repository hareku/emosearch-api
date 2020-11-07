package auth

import (
	"context"

	"github.com/hareku/emosearch-api/pkg/domain/model"
)

// Authenticator provides authentication methods.
type Authenticator interface {
	Authenticate(ctx context.Context) (context.Context, error)
	IsAuthenticated(ctx context.Context) bool
	UserID(ctx context.Context) (model.UserID, error)
	UserIDs(ctx context.Context, pageToken string) (ids []model.UserID, nextPageToken string, err error)
}
