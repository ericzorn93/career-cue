package app

import "apps/services/accounts-api/internal/app/ports"

// App is the application struct
type App struct {
	// RegistrationService is the registration service
	RegistrationService ports.RegistrationService
}

// AppOption is the option for the application
type AppOption func(*App)

// NewApp creates a new application
func NewApp(opts ...AppOption) App {
	app := App{}

	for _, opt := range opts {
		opt(&app)
	}

	return app
}

// WithRegistrationService sets the registration service in the application
func WithRegistrationService(service ports.RegistrationService) AppOption {
	return func(a *App) {
		a.RegistrationService = service
	}
}
