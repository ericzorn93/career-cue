package ports

import (
	userEntities "libs/backend/domain/user/entities"
)

// AuthService will handle auth webhook
// interactions
type AuthService interface {
	RegisterUser(user userEntities.User) error
}
