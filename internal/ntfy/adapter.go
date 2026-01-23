// Package ntfy exposes methods to handle notifications
//
// All interactions are handled via an adapter that must be created first
package ntfy

import (
	"context"
	"log/slog"
)

type NtfyAdapter interface {
	Notify(ctx context.Context, message string) error
}

type ntfyAdapter struct {
	logger   *slog.Logger
	ntfyHost string
	ntfyPort string
}

func NewNtfyAdapter(logger *slog.Logger, host, port string) NtfyAdapter {
	return &ntfyAdapter{
		logger:   logger,
		ntfyHost: host,
		ntfyPort: port,
	}
}
