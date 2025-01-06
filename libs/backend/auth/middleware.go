package auth

import (
	"context"
	"net/http"
)

// AuthMiddleware will pass forward the auth header to downstream services
// via context.Context
func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtain auth header from the request
		authHeader := r.Header.Get(AuthorizationHeaderKey)

		// Append to context
		ctx := r.Context()
		ctxWithAuth := context.WithValue(ctx, AuthMiddlewareContextKey, authHeader)

		// Reassign request
		r = r.WithContext(ctxWithAuth)

		next.ServeHTTP(w, r)
	}
}
