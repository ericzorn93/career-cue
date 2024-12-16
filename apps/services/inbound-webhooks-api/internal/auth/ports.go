package auth

import (
	"apps/services/inbound-webhooks-api/internal/entities"
	"context"
	commonv1 "libs/backend/proto-gen/go/common/v1"
	pb "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"
)

// AuthHandler is the interface for the service methods on the
// gRPC handler
type AuthHandler interface {
	UserRegistered(
		ctx context.Context,
		req *pb.UserRegisteredRequest,
	) (*commonv1.Empty, error)
}

// AuthServicePort is the interface that defines
// how the authentication provider interacts with
// the internal system
type AuthServicePort interface {
	// SendUserRegistered will convert the domain type type and send to
	// the message broker
	SendUserRegistered(ctx context.Context, user entities.User) error
}
