package lambda

import (
	"context"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hareku/emosearch-api/pkg/infrastructure/firebase"
)

var router *lmdrouter.Router

func init() {
	router = lmdrouter.NewRouter("/")
	router.Route("GET", "/@me", fetchMe)
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

		err = firebase.AuthFirebase(ctx, input.AuthorizationHeader)
		if err != nil {
			return lmdrouter.HandleError(lmdrouter.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized",
			})
		}

		return next(ctx, req)
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
