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
	mockGatewayPort := 3000

	var bootService boot.BootService
	app := fx.New(
		fx.Provide(func() boot.BootServiceParams {
			return boot.BootServiceParams{
				Name: mockServiceName,
				GRPCOptions: boot.GRPCOptions{
					Port:              uint64(mockPort),
					GatewayPort:       uint64(mockGatewayPort),
					ReflectionEnabled: true,
				},
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

	// Stop the app to clean up resources
	err = app.Stop(mockCtx)
	assert.NoError(t, err, "Application should stop without errors")
}
