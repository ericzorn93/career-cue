package boot_test

import (
	"context"
	boot "libs/boot/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestBoot(t *testing.T) {
	mockServiceName := "testService"
	mockCtx := context.Background()
	mockPort := 5000

	var bootService boot.BootService
	app := fx.New(
		fx.Provide(func() boot.BootServiceParams {
			return boot.BootServiceParams{
				Name:              mockServiceName,
				GRPCPort:          uint64(mockPort),
				ReflectionEnabled: true,
			}
		}),
		boot.NewBootServiceModule(), // Assuming this sets up the BootService
		fx.Populate(&bootService),   // Populates the BootService dependency
	)

	// Start the app and handle any errors
	err := app.Start(mockCtx)
	assert.NoError(t, err, "Application should start without errors")

	// Check the populated BootService
	assert.NotNil(t, bootService, "BootService should be populated")
	assert.Equal(t, mockServiceName, bootService.GetServiceName(), "Service name should match")
	assert.Equal(t, ":5000", bootService.GetGRPCPort(), "GRPC port should match")

	// Stop the app to clean up resources
	err = app.Stop(mockCtx)
	assert.NoError(t, err, "Application should stop without errors")
}
