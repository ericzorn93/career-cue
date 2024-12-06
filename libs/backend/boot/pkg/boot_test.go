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

	result, err := boot.NewService(mockCtx, mockServiceName)
	assert.NoError(t, err)
	assert.Equal(t, mockServiceName, result.GetServiceName())
}
