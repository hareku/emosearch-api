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

var firebaseAuthClient *firebase_auth.Client

func getFirebaseAuthClient() (*firebase_auth.Client, error) {
	if firebaseAuthClient == nil {
		ctx := context.Background()

		key, err := google.GetGoogleServiceAccountKey()
		if err != nil {
			return nil, fmt.Errorf("google service account key error: %w", err)
		}

		app, err := firebase_app.NewApp(ctx, nil, option.WithCredentialsJSON(key))
		if err != nil {
			return nil, fmt.Errorf("firebase app error: %w", err)
		}

		client, err := app.Auth(ctx)
		if err != nil {
			return nil, fmt.Errorf("firebase authentication error: %w", err)
		}

		firebaseAuthClient = client
	}

	return firebaseAuthClient, nil
}

func (r *registry) NewAuthenticator() auth.Authenticator {
	firebaseAuth, err := getFirebaseAuthClient()
	if err != nil {
		panic(fmt.Errorf("failed to get firebase auth client: %w", err))
	}

	return firebase_infra.NewFirebaseAuthenticator(firebaseAuth)
}
