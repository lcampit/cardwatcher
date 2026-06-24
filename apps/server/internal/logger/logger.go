// Package logger provides means to create a logger
// used throughout the whole application
package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type logger struct {
	inner *slog.Logger
}

// Debug implements [Logger].
func (l logger) Debug(msg string, args ...any) {
	l.inner.DebugContext(context.Background(), msg, args)
}

// DebugContext implements [Logger].
func (l logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.inner.DebugContext(ctx, msg, args)
}

// Error implements [Logger].
func (l logger) Error(msg string, args ...any) {
	l.inner.ErrorContext(context.Background(), msg, args)
}

// ErrorContext implements [Logger].
func (l logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.inner.ErrorContext(ctx, msg, args)
}

// Info implements [Logger].
func (l logger) Info(msg string, args ...any) {
	l.inner.InfoContext(context.Background(), msg, args)
}

// InfoContext implements [Logger].
func (l logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.inner.InfoContext(ctx, msg, args)
}

// Warn implements [Logger].
func (l logger) Warn(msg string, args ...any) {
	l.inner.WarnContext(context.Background(), msg, args)
}

// WarnContext implements [Logger].
func (l logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.inner.WarnContext(ctx, msg, args)
}

func (l *logger) With(args ...any) Logger {
	return &logger{inner: l.inner.With(args...)}
}

func (l *logger) Slog() *slog.Logger {
	return l.inner
}

type Logger interface {
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
	Slog() *slog.Logger
}

func NewLogger(level slog.Level, serviceName string) Logger {
	handler := &contextHandler{
		inner: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		}),
	}
	return &logger{
		inner: slog.New(handler).With(slog.String(KeyService,
			serviceName)),
	}
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
