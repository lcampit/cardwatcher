package watcher

import (
	"card-watcher/internal/cardtrader"
	"card-watcher/internal/database"

	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App

	cardtraderAdapter cardtrader.CardtraderAdapter
	db                database.Service
}

func New(
	cardtraderAdapter cardtrader.CardtraderAdapter,
) *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "card-watcher",
			AppName:      "card-watcher",
		}),
		cardtraderAdapter: cardtraderAdapter,

		db: database.New(),
	}

	return server
}
