package boot_test

import (
	boot "libs/boot/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoot(t *testing.T) {
	mockServiceName := "testService"

	bootService, err := boot.NewBootService(boot.BootServiceParams{
		Name: mockServiceName,
	})

	assert.NoError(t, err)
	assert.Equal(t, mockServiceName, bootService.GetServiceName())
}
