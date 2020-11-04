package ctxval

import (
	"context"
)

type authHeaderContextKey string

const authHeaderKey = authHeaderContextKey("authorization-header")

// GetAuthHeader returns the `Authorization` header value in a context.
func GetAuthHeader(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(authHeaderKey).(string)
	return s, ok
}

// SetAuthHeader returns new context which has the `Authorization` header value in a context.
func SetAuthHeader(ctx context.Context, authHeader string) context.Context {
	return context.WithValue(ctx, authHeaderKey, authHeader)
}
