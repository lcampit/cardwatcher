package service

import (
	"context"
	"crypto/sha256"
	"log/slog"

	"card-watcher/internal/cardtrader"
	"card-watcher/internal/models"
	"card-watcher/internal/mongo"
	"card-watcher/internal/ntfy"
)

type Service interface {
	SaveWatch(ctx context.Context, expansionID, blueprintID uint64, condition models.Condition, foil bool) (string, error)
	ListExpansions(ctx context.Context, name, code string) (*models.ListExpansionsResponse, error)
	ListBlueprints(ctx context.Context, expansionID uint64, name string) (*models.ListBlueprintsResponse, error)
	ListWatches(ctx context.Context) (*models.ListWatchesResponse, error)
	DeleteWatchByID(ctx context.Context, watchID string) error

	WatchAndNotify()
}

type service struct {
	logger            *slog.Logger
	cardtraderAdapter cardtrader.CardtraderAdapter
	mongoAdapter      mongo.MongoAdapter
	ntfyAdapter       ntfy.NtfyAdapter
}

type ServiceConfig struct {
	Logger            *slog.Logger
	CardtraderAdapter cardtrader.CardtraderAdapter
	MongoAdapter      mongo.MongoAdapter
	NtfyAdapter       ntfy.NtfyAdapter
}

func NewService(config ServiceConfig) *service {
	return &service{
		logger:            config.Logger,
		cardtraderAdapter: config.CardtraderAdapter,
		mongoAdapter:      config.MongoAdapter,
		ntfyAdapter:       config.NtfyAdapter,
	}
}

func HashAccessToken(accessToken string) string {
	h := sha256.New()
	h.Write([]byte(accessToken))
	return string(h.Sum(nil))
}
