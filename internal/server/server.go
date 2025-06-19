package server

import (
	"card-watcher/internal/models"
	"card-watcher/internal/service"
)

type server struct {
	models.UnimplementedCardWatcherServer
	service service.Service
}

func NewServer(
	service service.Service,
) *server {
	server := &server{
		service: service,
	}

	return server
}
