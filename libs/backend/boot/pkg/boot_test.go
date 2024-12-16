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
				Name: mockServiceName,
				GRPCOptions: boot.GRPCOptions{
					Port: uint64(mockPort),
				},
			}
		}),
		boot.NewBootServiceModule(), // Assuming this sets up the BootService
		fx.Populate(&bootService),   // Populates the BootService dependency
	)

	// Start the app and handle any errors
	err := app.Start(mockCtx)
	assert.Error(t, err, "Application should not start with error")
}
