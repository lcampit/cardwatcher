// Package server implements gRPC endpoints using generated
// code from models files.
//
// It handles everything related to gRPC operations, delegating
// actual logic and database manipolation to the underlying service
package server

import (
	"log/slog"

	api "github.com/lcampit/card-watcher-server/internal/api/v1"
	"github.com/lcampit/card-watcher-server/internal/server/service"
)

type server struct {
	api.UnimplementedCardWatcherServer
	logger  *slog.Logger
	service service.Service
}

type ServerConfig struct {
	Logger  *slog.Logger
	Service service.Service
}

func NewServer(config ServerConfig) *server {
	server := &server{
		logger:  config.Logger,
		service: config.Service,
	}

	return server
}
