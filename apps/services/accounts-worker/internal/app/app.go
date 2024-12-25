package app

import "apps/services/accounts-worker/internal/app/ports"

// App is the main application struct
type App struct {
	AccountService ports.AccountService
}

// AppOption is a function that modifies the App struct
type AppOption func(*App)

// NewApp creates a new App struct
func NewApp(opts ...AppOption) *App {
	app := new(App)

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// WithAccountService sets the account service
func WithAccountService(accountService ports.AccountService) AppOption {
	return func(app *App) {
		app.AccountService = accountService
	}
}
