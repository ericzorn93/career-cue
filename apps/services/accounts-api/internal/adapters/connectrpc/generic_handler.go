package connectrpc

import (
	"libs/boot/pkg/logger"
)

// AuthHandler handles all gRPC endpoints for inbound webhooks
type GenericHandler struct {
	Logger logger.Logger
}

// NewAuthHandler will return a pointer to the inbound webhooks API server
func NewAuthHandler(logger logger.Logger) *GenericHandler {
	return &GenericHandler{
		Logger: logger,
	}
}
