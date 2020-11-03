package firebase

import (
	"context"
	"fmt"
	"strings"

	firebase_auth "firebase.google.com/go/auth"
	"github.com/hareku/emosearch-api/pkg/domain"
	"github.com/hareku/emosearch-api/pkg/repository"
)

type firebaseAuthRepository struct {
	firebaseAuth *firebase_auth.Client
}

// NewFirebaseAuthRepository creates AuthRepository which implemented by Firebase-Authentication.
func NewFirebaseAuthRepository(firebaseAuth *firebase_auth.Client) repository.AuthRepository {
	return &firebaseAuthRepository{firebaseAuth}
}

func (r *firebaseAuthRepository) Authenticate(ctx context.Context, authHeader string) (domain.UserID, error) {
	idToken, err := resolveIDToken(authHeader)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	userID, err := r.checkIDToken(ctx, idToken)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	return userID, nil
}

func resolveIDToken(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("authorization type should be %q", "Bearer")
	}

	return strings.Replace(authHeader, "Bearer ", "", 1), nil
}

func (r *firebaseAuthRepository) checkIDToken(ctx context.Context, idToken string) (domain.UserID, error) {
	token, err := r.firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", fmt.Errorf("firebase-authentication could not verify ID token: %w", err)
	}

	return domain.UserID(token.UID), nil
}
