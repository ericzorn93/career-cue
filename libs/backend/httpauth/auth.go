package httpauth

import "context"

// Validator ensure handling of access token by service
type Validator interface {
	EnsureValidToken(ctx context.Context, accessToken string) (interface{}, error)
}
