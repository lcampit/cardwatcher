package watcher

import (
	"card-watcher/internal/cardtrader"
	"card-watcher/internal/mongo"

	"github.com/gofiber/fiber/v2"
)

type Watcher struct {
	*fiber.App

	cardtraderAdapter cardtrader.CardtraderAdapter
	mongoAdapter      mongo.MongoAdapter
}

func New(
	cardtraderAdapter cardtrader.CardtraderAdapter,
	mongoAdapter mongo.MongoAdapter,
) *Watcher {
	server := &Watcher{
		App: fiber.New(fiber.Config{
			ServerHeader: "card-watcher",
			AppName:      "card-watcher",
		}),
		cardtraderAdapter: cardtraderAdapter,

		mongoAdapter: mongoAdapter,
	}

	return server
}
