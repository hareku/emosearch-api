package repository

import (
	"context"

	"github.com/hareku/emosearch-api/pkg/domain"
)

// AuthRepository provides authentication methods.
type AuthRepository interface {
	Authenticate(ctx context.Context, authHeader string) (domain.UserID, error)
}
