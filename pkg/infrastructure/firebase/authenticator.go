package firebase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	firebase_auth "firebase.google.com/go/auth"
	"github.com/hareku/emosearch-api/internal/ctxval"
	"github.com/hareku/emosearch-api/pkg/domain/auth"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"google.golang.org/api/iterator"
)

type firebaseAuthenticator struct {
	firebaseAuth *firebase_auth.Client
}

// NewFirebaseAuthenticator creates Authenticator which implemented by Firebase-Authentication.
func NewFirebaseAuthenticator(firebaseAuth *firebase_auth.Client) auth.Authenticator {
	return &firebaseAuthenticator{firebaseAuth}
}

func (fa *firebaseAuthenticator) Authenticate(ctx context.Context) (context.Context, error) {
	authHeader, ok := ctxval.GetAuthHeader(ctx)
	if !ok {
		return nil, errors.New("`Authorization` header was not found")
	}

	idToken, err := fa.resolveIDToken(authHeader)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	userID, err := fa.checkIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return ctxval.SetUserID(ctx, userID), nil
}

func (fa *firebaseAuthenticator) IsAuthenticated(ctx context.Context) bool {
	_, ok := ctxval.GetUserID(ctx)
	return ok
}

func (fa *firebaseAuthenticator) UserID(ctx context.Context) (model.UserID, error) {
	userID, ok := ctxval.GetUserID(ctx)
	if !ok {
		return "", errors.New("context user is not authenticated")
	}

	return userID, nil
}

func (fa *firebaseAuthenticator) resolveIDToken(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("authorization type should be %q", "Bearer")
	}

	return strings.Replace(authHeader, "Bearer ", "", 1), nil
}

func (fa *firebaseAuthenticator) checkIDToken(ctx context.Context, idToken string) (model.UserID, error) {
	token, err := fa.firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", fmt.Errorf("firebase-authentication could not verify ID token: %w", err)
	}

	return model.UserID(token.UID), nil
}

func (fa *firebaseAuthenticator) ListUserID(ctx context.Context, pageToken string) ([]model.UserID, string, error) {
	var ids []model.UserID
	it := fa.firebaseAuth.Users(ctx, pageToken)
	nextPageToken := it.PageInfo().Token

	for {
		user, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", fmt.Errorf("firebase-auth failed to get users: %w", err)
		}
		ids = append(ids, model.UserID(user.UID))
	}

	return ids, nextPageToken, nil
}
