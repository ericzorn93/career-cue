package httpauth

import (
	"context"
	"libs/backend/boot"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	validator.CustomClaims
	Scope string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

// AccessTokenValidator will implement the auth.Validator interface
// for access token validation
type AccessTokenValidator struct {
	logger boot.Logger
}

// NewAccessTokenValidator handles validator construction
func NewAccessTokenValidator(logger boot.Logger) AccessTokenValidator {
	return AccessTokenValidator{logger}
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func (v AccessTokenValidator) EnsureValidToken(ctx context.Context, accessToken string) (interface{}, error) {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		v.logger.Error("Failed to parse the issuer url", slog.Any("error", err))
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		v.logger.Error("Failed to set up the jwt validator")
	}

	return jwtValidator.ValidateToken(ctx, accessToken)
}
