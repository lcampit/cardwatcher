// Package ntfy exposes methods to handle notifications
//
// All interactions are handled via an adapter that must be created first
package ntfy

import (
	"context"
	"fmt"
	"log/slog"

	"resty.dev/v3"
)

type NtfyAdapter interface {
	Notify(ctx context.Context, message string) error
}

type ntfyAdapter struct {
	logger *slog.Logger
	client *resty.Client
	topic  string
}

type NtfyAdapterConfig struct {
	Logger *slog.Logger
	Host   string
	Port   string
	Topic  string
}

func NewNtfyAdapter(config NtfyAdapterConfig) NtfyAdapter {
	return &ntfyAdapter{
		logger: config.Logger,
		client: resty.New().SetBaseURL(fmt.Sprintf("http://%s:%s", config.Host, config.Port)),
		topic:  config.Topic,
	}
}
