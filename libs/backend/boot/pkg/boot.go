package boot

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type Service struct {
	Name              string
	Ctx               context.Context
	GRPCPort          *string
	GRPCDialOptions   []grpc.DialOption
	ReflectionEnabled bool
}

type ServiceOption func(*Service) error

func WithGRPC[P ~int](port P) ServiceOption {
	return func(s *Service) error {
		portStr := fmt.Sprintf(":%d", port)
		s.GRPCPort = &portStr
		return nil
	}
}

func WithGRPCDialOpts(grpcDialOptions ...grpc.DialOption) ServiceOption {
	return func(s *Service) error {
		s.GRPCDialOptions = grpcDialOptions
		return nil
	}
}

func WithReflection() ServiceOption {
	return func(s *Service) error {
		s.ReflectionEnabled = true
		return nil
	}
}

func NewService(ctx context.Context, serviceName string, opts ...ServiceOption) (*Service, error) {
	s := &Service{
		Ctx:  ctx,
		Name: serviceName,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s Service) GetServiceName() string {
	return s.Name
}

func (s Service) GetGRPCPort() *string {
	return s.GRPCPort
}
