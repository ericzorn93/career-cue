package main

import (
	"apps/services/accounts-api/internal/app"
	"apps/services/accounts-api/internal/config"
	"apps/services/accounts-api/internal/domain/services"
	"apps/services/accounts-api/internal/models"
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	connectrpcadapter "apps/services/accounts-api/internal/adapters/connectrpc"
	"apps/services/accounts-api/internal/adapters/database/repositories"
	"libs/backend/boot"
	auth "libs/backend/httpauth"
	"libs/backend/proto-gen/go/accounts/accountsapi/v1/accountsapiv1connect"
)

// serviceName is the name of the microservice
const serviceName = "accounts-api"

func run() error {
	// Application Context
	ctx := context.Background()

	// Create logger
	logger := boot.NewSlogger()

	// Construct config
	config, err := config.NewConfig()
	if err != nil {
		logger.Error("Trouble constructing config")
		os.Exit(1)
	}

	// Connect Interceptors
	validationInterceptor, err := validate.NewInterceptor()
	if err != nil {
		logger.Error("Cannot set up validation interceptor", slog.Any("error", err))
		return err
	}

	// Custom interceptors
	authInterceptor := auth.NewAuthInterceptor(logger)

	options := []connect.HandlerOption{
		connect.WithInterceptors(
			validationInterceptor,
			authInterceptor.Incoming(),
		),
	}

	// Initialize the gRPC Options
	bootService := boot.
		NewBuildServiceBuilder().
		SetServiceName(serviceName).
		SetLogger(logger).
		SetDBOptions(boot.DBOptions{
			Host:     config.DBHost,
			Name:     config.DBName,
			User:     config.DBUser,
			Password: config.DBPassword,
			Port:     config.DBPort,
			SSLMode:  config.DBSSLMode,
			TimeZone: config.DBTimeZone,
		}).
		SetAMQPOptions(boot.AMQPOptions{
			ConnectionURI: config.AMQPUrl,
			OnConnectionCallback: func(params boot.AMQPCallBackParams) error {
				params.Logger.Info("AMQP connected successfully")
				return nil
			},
		}).
		SetConnectRPCOptions(boot.ConnectRPCOptions{
			Port: 3000,
			TransportCredentials: []credentials.TransportCredentials{
				insecure.NewCredentials(),
			},
			Handlers: []boot.ConnectRPCHandler{
				func(params boot.ConnectRPCHandlerParams) error {
					if !params.AMQPController.IsConnected() {
						errMsg := "AMQP not conntected"
						logger.Error(errMsg)
						return errors.New(errMsg)
					}

					// Create repositories
					accountRepo := repositories.NewAccountRespository(params.Logger, params.DB)

					// Create services
					registrationService := services.NewAccountService(params.Logger, accountRepo)

					// Create Application
					app := app.NewApp(
						app.WithRegistrationService(registrationService),
					)

					// Register all ConnectRPC handlers
					registrationHandler := connectrpcadapter.NewRegistrationHandler(params.Logger, app)

					// Assign the handler to the HTTP path
					path, httpHandler := accountsapiv1connect.NewAccountServiceHandler(
						registrationHandler,
						options...,
					)

					// HTTP Handlers and reflection registered with Mux
					params.Mux.Handle(path, httpHandler)
					reflector := grpcreflect.NewStaticReflector(accountsapiv1connect.AccountServiceName)
					params.Mux.Handle(grpcreflect.NewHandlerV1(reflector))
					params.Mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

					return nil
				},
			},
		}).
		SetBootCallbacks([]boot.BootCallback{
			func(params boot.BootCallbackParams) error {
				// Run DB migrations
				if err := params.DB.AutoMigrate(&models.Account{}); err != nil {
					params.Logger.Error("Failed to run DB migrations", slog.Any("error", err))
					return err
				}

				params.Logger.Info("Ran DB migrations", slog.String("serviceName", serviceName))
				return nil
			},
			func(params boot.BootCallbackParams) error {
				params.Logger.Info("Service booted successfully", slog.String("serviceName", serviceName))
				return nil
			},
		}).
		Build()

	return bootService.Start(ctx)
}

func main() {
	if err := run(); err != nil {
		log.Printf("Cannot start service %s", serviceName)
		os.Exit(1)
	}
}
