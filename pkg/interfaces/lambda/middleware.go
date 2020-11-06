package lambda

import (
	"context"
	"log"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/hareku/emosearch-api/internal/ctxval"
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

			headerCtx := ctxval.SetAuthHeader(ctx, input.AuthorizationHeader)

			authCtx, err := authenticator.Authenticate(headerCtx)
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

func corsMiddleware() lmdrouter.Middleware {
	return func(next lmdrouter.Handler) lmdrouter.Handler {
		return func(ctx context.Context, req events.APIGatewayProxyRequest) (
			res events.APIGatewayProxyResponse,
			err error,
		) {
			res, err = next(ctx, req)

			if res.Headers == nil {
				res.Headers = map[string]string{}
			}

			res.Headers["Access-Control-Allow-Origin"] = "*"
			res.Headers["Access-Control-Allow-Methods"] = "*"
			res.Headers["Access-Control-Allow-Headers"] = "*"

			return
		}
	}
}

func loggerMiddleware(next lmdrouter.Handler) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		// [LEVEL] [METHOD PATH] [CODE] EXTRA
		format := "[%s] [%s %s] [%d] %s"
		level := "INF"
		var code int
		var extra string

		res, err = next(ctx, req)
		if err != nil {
			level = "ERR"
			code = http.StatusInternalServerError
			extra = " " + err.Error()
		} else {
			code = res.StatusCode
			if code >= 400 {
				level = "ERR"
			}
		}

		log.Printf(format, level, req.HTTPMethod, req.Path, code, extra)

		return res, err
	}
}
