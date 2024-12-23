package boot_test

import (
	boot "libs/backend/boot"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoot(t *testing.T) {
	mockServiceName := "testService"

	bootService := boot.NewBuildServiceBuilder().
		SetServiceName(mockServiceName).
		SetLogger(boot.NewSlogger()).Build()

	assert.Equal(t, mockServiceName, bootService.GetServiceName())
}
