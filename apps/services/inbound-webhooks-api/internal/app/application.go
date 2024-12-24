package app

import "apps/services/inbound-webhooks-api/internal/app/ports"

// Application is the main application struct
type Application struct {
	AuthService ports.AuthService
}

// NewApplication creates a new application
func NewApplication(authService ports.AuthService) Application {
	return Application{
		AuthService: authService,
	}
}
