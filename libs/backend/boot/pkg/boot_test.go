package boot_test

import (
	"context"
	boot "libs/boot/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoot(t *testing.T) {
	mockServiceName := "testService"
	mockCtx := context.Background()

	mockPort := 5000
	service, err := boot.NewService(
		mockCtx,
		mockServiceName,
		boot.WithGRPC(mockPort),
	)

	assert.NoError(t, err)
	assert.Equal(t, mockServiceName, service.GetServiceName())
	assert.Equal(t, ":5000", *service.GRPCPort)
}
