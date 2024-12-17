package boot_test

import (
	boot "libs/boot/pkg"
	"libs/boot/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoot(t *testing.T) {
	mockServiceName := "testService"

	bootService, err := boot.NewBootService(boot.BootServiceParams{
		Logger: logger.NewSlogger(),
		Name:   mockServiceName,
	})

	assert.NoError(t, err)
	assert.Equal(t, mockServiceName, bootService.GetServiceName())
}
