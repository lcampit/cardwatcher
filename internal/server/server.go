// Package server implements gRPC endpoints using generated
// code from models files.
//
// It handles everything related to gRPC operations, delegating
// actual logic and database manipolation to the underlying service
package server

import (
	"log/slog"

	"card-watcher/internal/models"
	"card-watcher/internal/service"
)

type server struct {
	models.UnimplementedCardWatcherServer
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
