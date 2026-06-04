// Package handler implements gRPC endpoints using generated
// code from models files.
//
// It handles everything related to gRPC operations, delegating
// actual logic and database manipolation to the underlying service
package handler

import (
	"context"
	"log/slog"

	"github.com/lcampit/cardwatcher/apps/server/internal/service"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

type Handler interface {
	apiv1.CardWatcherServiceServer
	Health(ctx context.Context) error
}

type handler struct {
	apiv1.UnsafeCardWatcherServiceServer
	logger  *slog.Logger
	service service.Service
}

type HandlerConfig struct {
	Logger  *slog.Logger
	Service service.Service
}

func NewHandler(config HandlerConfig) Handler {
	handler := &handler{
		logger:  config.Logger,
		service: config.Service,
	}

	return handler
}

func (s *handler) Health(ctx context.Context) error {
	return s.service.Health(ctx)
}
