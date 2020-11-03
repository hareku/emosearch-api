package lambda

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hareku/emosearch-api/internal/ctxval"
	"github.com/hareku/emosearch-api/pkg/infrastructure/firebase"
	"github.com/hareku/emosearch-api/pkg/repository"

	firebase_app "firebase.google.com/go"
	firebase_auth "firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var router *lmdrouter.Router
var authRepository repository.AuthRepository

func init() {
	router = lmdrouter.NewRouter("/")
	router.Route("GET", "/@me", fetchMe)

	firebaseAuth, err := makeFirebaseAuth()
	if err != nil {
		panic(fmt.Errorf("firebase error: %w", err))
	}
	authRepository = firebase.NewFirebaseAuthRepository(firebaseAuth)
}

func makeFirebaseAuth() (*firebase_auth.Client, error) {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_JSON_PATH"))
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

// Start Lambda function.
func Start() {
	lambda.Start(router.Handler)
}

type authInput struct {
	AuthorizationHeader string `lambda:"header.Authorization"`
}

// authMiddleware checks whether the request user is authenticated,
// and create User domain data if it doesn't exist yet.
func authMiddleware(next lmdrouter.Handler) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		var input authInput
		err = lmdrouter.UnmarshalRequest(req, false, &input)
		if err != nil {
			return lmdrouter.HandleError(lmdrouter.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "credentials is missing.",
			})
		}

		userID, err := authRepository.Authenticate(ctx, input.AuthorizationHeader)
		if err != nil {
			return lmdrouter.HandleError(lmdrouter.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized",
			})
		}

		child := ctxval.SetUserID(ctx, userID)
		return next(child, req)
	}
}

type fetchMeInput struct {
	Authorization string `lambda:"header.Authorization"`
}

func fetchMe(ctx context.Context, req events.APIGatewayProxyRequest) (
	res events.APIGatewayProxyResponse,
	err error,
) {
	var input fetchMeInput
	err = lmdrouter.UnmarshalRequest(req, false, &input)
	if err != nil {
		return lmdrouter.HandleError(err)
	}

	var data interface{}
	return lmdrouter.MarshalResponse(http.StatusOK, nil, data)
}
