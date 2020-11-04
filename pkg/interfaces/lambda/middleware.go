package lambda

import (
	"context"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/hareku/emosearch-api/pkg/domain/auth"
)

type authInput struct {
	AuthorizationHeader string `lambda:"header.Authorization"`
}

func authMiddleware(authenticator auth.Authenticator) lmdrouter.Middleware {
	return func(next lmdrouter.Handler) lmdrouter.Handler {
		return func(ctx context.Context, req events.APIGatewayProxyRequest) (
			res events.APIGatewayProxyResponse,
			err error,
		) {
			var input authInput
			err = lmdrouter.UnmarshalRequest(req, false, &input)
			if err != nil {
				return lmdrouter.HandleError(lmdrouter.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "credentials is missing",
				})
			}

			authCtx, err := authenticator.Authenticate(ctx)
			if err != nil {
				return lmdrouter.HandleError(lmdrouter.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "unauthorized",
				})
			}

			return next(authCtx, req)
		}
	}
}

// // authMiddleware checks whether the request user is authenticated.
// func authMiddleware(next lmdrouter.Handler) lmdrouter.Handler {
// 	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
// 		res events.APIGatewayProxyResponse,
// 		err error,
// 	) {
// 		var input authInput
// 		err = lmdrouter.UnmarshalRequest(req, false, &input)
// 		if err != nil {
// 			return lmdrouter.HandleError(lmdrouter.HTTPError{
// 				Code:    http.StatusUnauthorized,
// 				Message: "credentials is missing",
// 			})
// 		}

// 		userID, err := authRepository.Authenticate(ctx, input.AuthorizationHeader)
// 		if err != nil {
// 			return lmdrouter.HandleError(lmdrouter.HTTPError{
// 				Code:    http.StatusUnauthorized,
// 				Message: "unauthorized",
// 			})
// 		}

// 		child := auth.SetUserID(ctx, userID)
// 		return next(child, req)
// 	}
// }
