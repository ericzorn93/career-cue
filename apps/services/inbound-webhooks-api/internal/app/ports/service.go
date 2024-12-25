package ports

import (
	"libs/backend/domain/user"
)

// AuthService will handle auth webhook
// interactions
type AuthService interface {
	RegisterUser(user user.User) error
}
