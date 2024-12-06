package boot_test

import (
	boot "packages/boot/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoot(t *testing.T) {
	mockServiceName := "testService"

	result := boot.NewService(mockServiceName)
	assert.Equal(t, mockServiceName, result.GetServiceName())

}
