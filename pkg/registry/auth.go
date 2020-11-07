package registry

import (
	"context"
	"fmt"

	firebase_app "firebase.google.com/go"
	firebase_auth "firebase.google.com/go/auth"
	"github.com/hareku/emosearch-api/internal/google"
	"github.com/hareku/emosearch-api/pkg/domain/auth"
	firebase_infra "github.com/hareku/emosearch-api/pkg/infrastructure/firebase"
	"google.golang.org/api/option"
)

func makeFirebaseAuth() (*firebase_auth.Client, error) {
	bytes, err := google.GetGoogleServiceAccountKey()
	if err != nil {
		return nil, fmt.Errorf("google service account key is not found: %w", err)
	}

	opt := option.WithCredentialsJSON(bytes)
	var config *firebase_app.Config
	ctx := context.Background()

	app, err := firebase_app.NewApp(ctx, config, opt)
	if err != nil {
		return nil, fmt.Errorf("firebase error: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase-authentication error: %w", err)
	}

	return client, nil
}

func (r *registry) NewAuthenticator() auth.Authenticator {
	firebaseAuth, err := makeFirebaseAuth()
	if err != nil {
		panic(fmt.Errorf("firebase error: %w", err))
	}

	return firebase_infra.NewFirebaseAuthenticator(firebaseAuth)
}
