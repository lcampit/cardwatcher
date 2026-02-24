// Package handler implements gRPC endpoints using generated
// code from models files.
//
// It handles everything related to gRPC operations, delegating
// actual logic and database manipolation to the underlying service
package handler

import (
	"log/slog"

	api "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
	"github.com/lcampit/cardwatcher/internal/server/service"
)

type handler struct {
	api.UnimplementedCardWatcherServer
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
