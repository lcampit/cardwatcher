// Package logger provides means to create a logger
// used throughout the whole application
package logger

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger(levelStr string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     getLogLevelFromConfig(levelStr),
		AddSource: true,
	}))
}

func getLogLevelFromConfig(levelStr string) slog.Level {
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
