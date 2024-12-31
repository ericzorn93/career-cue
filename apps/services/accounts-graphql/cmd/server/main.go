package main

import (
	"apps/services/accounts-graphql/graph"
	"apps/services/accounts-graphql/internal/config"
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"libs/backend/boot"
)

// serviceName is the name of the microservice
const serviceName = "accounts-graphql"

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
	_ = []connect.HandlerOption{
		connect.WithInterceptors(validationInterceptor),
	}

	// Initialize the gRPC Options
	bootService := boot.
		NewBuildServiceBuilder().
		SetServiceName(serviceName).
		SetLogger(logger).
		SetAMQPOptions(boot.AMQPOptions{
			ConnectionURI: config.AMQPUrl,
			OnConnectionCallback: func(params boot.AMQPCallBackParams) error {
				params.Logger.Info("AMQP connected successfully")

				// Set up all AMQP queues and exchanges

				params.Logger.Info("Set up all AMQP queues and exchanges")

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

					// Set up all ConnectRPC Handlers
					srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
					srv.AddTransport(transport.Options{})
					srv.AddTransport(transport.GET{})
					srv.AddTransport(transport.POST{})
					// srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

					// Middleware
					srv.Use(extension.Introspection{})
					// srv.Use(extension.AutomaticPersistedQuery{
					// Cache: lru.New[string](100),
					// })

					params.Mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
					params.Mux.Handle("/graphql", srv)

					return nil
				},
			},
		}).
		SetBootCallbacks([]boot.BootCallback{
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