package httpauth

import "errors"

// Claims Fetch Errors
var (
	ErrCustomClaimsNotValid = errors.New("custom claims not valid")
)

// Claims Validation Errors
var (
	ErrCustomClaimsScopeEmpty = errors.New("empty custom claims scope")
)
