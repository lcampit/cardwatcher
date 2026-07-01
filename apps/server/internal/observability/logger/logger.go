package logger

import (
	"log/slog"
	"os"
	"strings"
)

type LoggerConfig struct {
	Level     string
	Service   string
	AddSource bool
}

func New(config LoggerConfig) *slog.Logger {
	level := getLogLevelFromString(config.Level)
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: config.AddSource,
	})

	handler := &ContextHandler{
		inner: baseHandler,
	}

	return slog.New(handler).
		With(slog.String(KeyService, config.Service))
}

func getLogLevelFromString(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
