package firebase

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
)

// AuthFirebase authenticates by Firebase-Authentication.
func AuthFirebase(ctx context.Context, authHeader string) error {
	idToken, err := resolveIDToken(authHeader)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	err = checkIDToken(ctx, idToken)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	return nil
}

func resolveIDToken(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("authorization type should be %q", "Bearer")
	}

	return strings.Replace(authHeader, "Bearer ", "", 1), nil
}

func checkIDToken(ctx context.Context, idToken string) error {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_JSON_PATH"))
	var config *firebase.Config

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return fmt.Errorf("firebase error: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("firebase-authentication error: %w", err)
	}

	_, err = client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return fmt.Errorf("firebase-authentication could not verify ID token: %w", err)
	}

	return nil
}
