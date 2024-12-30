package connectrpc

import (
	"apps/services/accounts-api/internal/models"
	"context"
	"libs/backend/boot"
	accountsapiv1 "libs/backend/proto-gen/go/accounts/accountsapi/v1"
	"libs/backend/proto-gen/go/accounts/accountsapi/v1/accountsapiv1connect"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthHandler handles all gRPC endpoints for inbound webhooks
type RegistrationServiceHandler struct {
	accountsapiv1connect.UnimplementedRegistrationServiceHandler
	Logger boot.Logger
	DB     *gorm.DB
}

// NewRegistrationHandler will return a pointer to the inbound webhooks API server
func NewRegistrationHandler(logger boot.Logger, db *gorm.DB) *RegistrationServiceHandler {
	return &RegistrationServiceHandler{
		Logger: logger,
		DB:     db,
	}
}

// CreateAccount handles user creation and saves them in the database
func (r *RegistrationServiceHandler) CreateAccount(
	ctx context.Context,
	req *connect.Request[accountsapiv1.CreateAccountRequest],
) (*connect.Response[accountsapiv1.CreateAcountResponse], error) {
	r.Logger.Info("CreateAccount called")
	r.DB.Save(&models.Account{
		EmailAddress: req.Msg.EmailAddress,
		UserName:     req.Msg.Username,
		CommonID:     uuid.MustParse(req.Msg.CommonId),
	})
	return connect.NewResponse(&accountsapiv1.CreateAcountResponse{IsSuccess: true}), nil
}
