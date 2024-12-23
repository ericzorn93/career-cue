package application

// Application is the main application struct
type Application struct {
	AuthService AuthService
}

// NewApplication creates a new application
func NewApplication(authService AuthService) Application {
	return Application{
		AuthService: authService,
	}
}
