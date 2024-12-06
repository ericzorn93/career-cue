package boot

import (
	"fmt"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// NewBootServiceModule helps define the dependency injection that
// can easily be connected within any microservice system
func NewBootServiceModule() fx.Option {
	module := fx.Module(
		"bootService",
		fx.Provide(NewBootService),
	)

	return module
}

// NewBootService sets up constructor for the boot service
// without any functionality or options
func NewBootService(params BootServiceParams) BootService {
	return BootService{
		name:              params.Name,
		gRPCPort:          params.GRPCPort,
		gRPCDialOptions:   params.GRPCDialOptions,
		reflectionEnabled: params.ReflectionEnabled,
	}
}

// BootService primary struct that defines
// holding the data for all service modules
type BootService struct {
	name              string
	gRPCPort          uint64
	gRPCDialOptions   []grpc.DialOption
	reflectionEnabled bool
}

// GetServiceName returns the service name
func (s BootService) GetServiceName() string {
	return s.name
}

// GetGRPCPort returns the gRPC port
func (s BootService) GetGRPCPort() string {
	return fmt.Sprintf(":%d", s.gRPCPort)
}

// BootServiceParams are the incoming options
// for the boot service construction
type BootServiceParams struct {
	Name              string
	GRPCPort          uint64
	GRPCDialOptions   []grpc.DialOption
	ReflectionEnabled bool
}
