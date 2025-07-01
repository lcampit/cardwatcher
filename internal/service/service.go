package service

import (
	"card-watcher/internal/cardtrader"
	"card-watcher/internal/models"
	"card-watcher/internal/mongo"
	"card-watcher/internal/ntfy"
	"context"
	"crypto/sha256"
)

type Service interface {
	SaveWatch(ctx context.Context, expansionId, blueprintId int, condition models.Condition, foil bool) (string, error)
	ListExpansions(ctx context.Context, name, code string) (models.ListExpansionsResponse, error)
	ListBlueprints(ctx context.Context, expansionId int, name string) (models.ListBlueprintsResponse, error)
	ListWatches(ctx context.Context) (models.ListWatchesResponse, error)
	DeleteWatchByID(ctx context.Context, watchID string) error

	WatchAndNotify()
}

type service struct {
	cardtraderAdapter cardtrader.CardtraderAdapter
	mongoAdapter      mongo.MongoAdapter
	ntfyAdapter       ntfy.NtfyAdapter
}

func NewService(
	cardtraderAdapter cardtrader.CardtraderAdapter,
	mongoAdapter mongo.MongoAdapter,
	ntfyAdapter ntfy.NtfyAdapter,
) *service {
	return &service{
		cardtraderAdapter: cardtraderAdapter,
		mongoAdapter:      mongoAdapter,
		ntfyAdapter:       ntfyAdapter,
	}
}

func HashAccessToken(accessToken string) string {
	h := sha256.New()
	h.Write([]byte(accessToken))
	return string(h.Sum(nil))
}
