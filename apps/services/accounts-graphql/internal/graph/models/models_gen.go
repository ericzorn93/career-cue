// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	// The unique identifier for the account
	ID uuid.UUID `json:"id"`
	// The email address of the account
	EmailAddress string `json:"emailAddress"`
	// The createdAt time of the account
	CreatedAt time.Time `json:"createdAt"`
	// The updatedAt time of the account
	UpdatedAt time.Time `json:"updatedAt"`
}

type Mutation struct {
}

type Query struct {
}

type Viewer struct {
	Empty   bool     `json:"empty"`
	Account *Account `json:"account,omitempty"`
}
