package auth

import (
	"context"
	commonv1 "libs/backend/proto-gen/go/common/v1"
	pb "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"

	"connectrpc.com/connect"
)

// AuthHandler is the interface for the service methods on the
// gRPC handler
type AuthHandler interface {
	UserRegistered(
		ctx context.Context,
		req *connect.Request[pb.UserRegisteredRequest],
	) (*connect.Response[commonv1.Empty], error)
}
