package auth

import (
	"context"

	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// accessTokenClaimsKey used for obtaining the access token from the context
type accessTokenClaimsKey struct{}

// ClaimsKey will ensure the claims get stored and retrieved under the same key
var ClaimsKey = accessTokenClaimsKey{}

// authMiddlewareContextKey ensures key is appened to and from context
type authMiddlewareContextKey struct{}

// AuthMiddlewareContextKey will be used to obtain and set context auth
var AuthMiddlewareContextKey authMiddlewareContextKey = authMiddlewareContextKey{}

// SetClaimsToContext will assign the custom claims to the conext
func SetClaimsToContext(ctx context.Context, claims *validator.ValidatedClaims) context.Context {
	return context.WithValue(ctx, ClaimsKey, claims.CustomClaims.(*CustomClaims))
}

// GetClaimsFromContext will pull the custom claims out of the context
func GetClaimsFromContext(ctx context.Context) (*CustomClaims, error) {
	claims := ctx.Value(ClaimsKey)
	customClaimsValidated, ok := claims.(*CustomClaims)

	if !ok {
		return nil, ErrCustomClaimsNotValid
	}

	return customClaimsValidated, nil
}

// GetAuthTokenFromContext will pull context value from middleware
func GetAuthTokenFromContext(ctx context.Context) string {
	val, ok := ctx.Value(AuthMiddlewareContextKey).(string)
	if !ok {
		return ""
	}

	return val
}
