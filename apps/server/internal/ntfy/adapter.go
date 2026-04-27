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

type NtfyAdapterConfig struct {
	Logger *slog.Logger
	Host   string
	Port   string
}

func NewNtfyAdapter(config NtfyAdapterConfig) NtfyAdapter {
	return &ntfyAdapter{
		logger:   config.Logger,
		ntfyHost: config.Host,
		ntfyPort: config.Port,
	}
}
