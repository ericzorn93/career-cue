package app

import (
	"apps/services/accounts-worker/internal/app/ports"
	"libs/backend/auth/m2m"
)

// App is the main application struct
type App struct {
	AccountService ports.AccountService
	M2MService     m2m.M2MGenerator
}

// AppOption is a function that modifies the App struct
type AppOption func(*App)

// NewApp creates a new App struct
func NewApp(opts ...AppOption) App {
	app := App{}

	for _, opt := range opts {
		opt(&app)
	}

	return app
}

// WithAccountService sets the account service
func WithAccountService(accountService ports.AccountService) AppOption {
	return func(app *App) {
		app.AccountService = accountService
	}
}

// WithM2MGenerator will allow the application to access m2m tokens from the auth provider
func WithM2MGenerator(m2mGenerator m2m.M2MGenerator) AppOption {
	return func(a *App) {
		a.M2MService = m2mGenerator
	}
}
