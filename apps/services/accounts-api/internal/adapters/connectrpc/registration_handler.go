package connectrpc

import (
	"context"
	"libs/backend/boot"
	accountsapiv1 "libs/backend/proto-gen/go/accounts/accountsapi/v1"
	"libs/backend/proto-gen/go/accounts/accountsapi/v1/accountsapiv1connect"

	"connectrpc.com/connect"
)

// AuthHandler handles all gRPC endpoints for inbound webhooks
type RegistrationServiceHandler struct {
	accountsapiv1connect.UnimplementedRegistrationServiceHandler
	Logger boot.Logger
}

// NewRegistrationHandler will return a pointer to the inbound webhooks API server
func NewRegistrationHandler(logger boot.Logger) *RegistrationServiceHandler {
	return &RegistrationServiceHandler{
		Logger: logger,
	}
}

// CreateAccount handles user creation and saves them in the database
func (r *RegistrationServiceHandler) CreateAccount(
	ctx context.Context,
	req *connect.Request[accountsapiv1.CreateAccountRequest],
) (*connect.Response[accountsapiv1.CreateAcountResponse], error) {
	r.Logger.Info("CreateAccount called")
	return connect.NewResponse(&accountsapiv1.CreateAcountResponse{IsSuccess: true}), nil
}
