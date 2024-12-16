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
		fx.Provide(fx.Annotate(
			NewSlogger,
		)),
	)
}

// Logger Interface for Application Logger
type Logger interface {
	Info(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

// NewSlogger creates Slog Logger in JSON format
func NewSlogger() Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
	}))

	return logger

}
