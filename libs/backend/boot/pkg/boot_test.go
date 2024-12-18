package boot_test

import (
	boot "libs/boot/pkg"
	"libs/boot/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoot(t *testing.T) {
	mockServiceName := "testService"

	bootService := boot.NewBuildServiceBuilder().
		SetServiceName(mockServiceName).
		SetLogger(logger.NewSlogger()).Build()

	assert.Equal(t, mockServiceName, bootService.GetServiceName())
}
