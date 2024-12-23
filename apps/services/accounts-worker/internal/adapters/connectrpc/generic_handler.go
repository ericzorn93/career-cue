package connectrpc

import (
	boot "libs/boot"
)

// AuthHandler handles all gRPC endpoints for inbound webhooks
type GenericHandler struct {
	Logger boot.Logger
}

// NewAuthHandler will return a pointer to the inbound webhooks API server
func NewAuthHandler(logger boot.Logger) *GenericHandler {
	return &GenericHandler{
		Logger: logger,
	}
}
