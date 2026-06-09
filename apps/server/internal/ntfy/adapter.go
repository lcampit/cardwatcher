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

func (config NtfyAdapterConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", config.Host),
		slog.String("port", config.Port),
		slog.String("topic", config.Topic),
	)
}

func NewNtfyAdapter(config NtfyAdapterConfig) NtfyAdapter {
	config.Logger.Debug("creating ntfy adapter", slog.Any("config", config))
	return &ntfyAdapter{
		logger: config.Logger,
		client: resty.New().SetBaseURL(fmt.Sprintf("http://%s:%s", config.Host, config.Port)),
		topic:  config.Topic,
	}
}
