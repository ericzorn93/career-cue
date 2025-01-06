package auth

import (
	"context"
	"errors"
	"libs/backend/boot"
	"strings"

	"connectrpc.com/connect"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

const (
	tokenValueLength = 2
)

type AuthInerceptor struct {
	logger               boot.Logger
	accessTokenValidator Validator
}

// NewAuthInterceptor will intercept connectRPC requests and handle authentication
func NewAuthInterceptor(logger boot.Logger) AuthInerceptor {
	return AuthInerceptor{logger: logger, accessTokenValidator: NewAccessTokenValidator(logger)}
}

// Incoming will handle auth for incoming server requests
func (a AuthInerceptor) Incoming() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			authTokenHeaderValue := req.Header().Get(AuthorizationHeaderKey)

			// Check token in handlers.
			if authTokenHeaderValue == "" {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("no token provided"),
				)
			}

			// Validate token length
			authTokenValues := strings.Split(authTokenHeaderValue, " ")

			if len(authTokenValues) != tokenValueLength {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("token invalid"),
				)
			}

			// Validate incoming token
			accessToken := authTokenValues[1]
			claims, err := a.accessTokenValidator.EnsureValidToken(ctx, accessToken)

			// Validate claims
			if err != nil {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("token claims invalid"),
				)
			}

			ctx = SetClaims(ctx, claims.(*validator.ValidatedClaims))

			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
