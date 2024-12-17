package logger

import (
	"log/slog"
	"os"
)

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
