package auth

import (
	"context"

	"github.com/hareku/emosearch-api/pkg/domain/model"
)

type authContextKey string

const userIDKey = authContextKey("user-id")

// GetUserID returns UserID of the given context.
func getUserID(ctx context.Context) (model.UserID, bool) {
	s, ok := ctx.Value(userIDKey).(model.UserID)
	return s, ok
}

// SetUserID returns new context which has UserID for the context value.
func setUserID(ctx context.Context, userID model.UserID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
