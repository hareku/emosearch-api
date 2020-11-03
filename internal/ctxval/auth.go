package ctxval

import (
	"context"

	"github.com/hareku/emosearch-api/pkg/domain"
)

type authContextKey string

const userIDKey = authContextKey("user-id")

// GetUserID returns UserID of the given context.
func GetUserID(ctx context.Context) (domain.UserID, bool) {
	s, ok := ctx.Value(userIDKey).(domain.UserID)
	return s, ok
}

// SetUserID returns new context which has UserID for the context value.
func SetUserID(ctx context.Context, userID domain.UserID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
