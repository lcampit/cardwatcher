// Package handler implements gRPC endpoints using generated
// code from models files.
//
// It handles everything related to gRPC operations, delegating
// actual logic and database manipolation to the underlying service
package handler

import (
	"log/slog"

	"github.com/lcampit/cardwatcher/apps/server/internal/service"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

type handler struct {
	apiv1.UnsafeCardWatcherServiceServer
	logger  *slog.Logger
	service service.Service
}

type HandlerConfig struct {
	Logger  *slog.Logger
	Service service.Service
}

func NewHandler(config HandlerConfig) *handler {
	handler := &handler{
		logger:  config.Logger,
		service: config.Service,
	}

	return handler
}
