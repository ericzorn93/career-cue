package boot

import (
	"log/slog"
	"os"

	"go.uber.org/fx"
)

const (
	loggerModuleName = "logger"
)

// NewLoggerModule returns logger
func NewLoggerModule() fx.Option {
	return fx.Module(
		loggerModuleName,
		fx.Provide(NewLogger),
	)
}

// NewLogger creates Slog Logger in JSON format
func NewLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
	}))

	return logger

}
