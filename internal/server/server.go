package server

import (
	"github.com/gofiber/fiber/v2"

	"card-watcher/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "card-watcher",
			AppName:      "card-watcher",
		}),

		db: database.New(),
	}

	return server
}
