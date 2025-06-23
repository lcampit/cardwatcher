package service

import (
	"card-watcher/internal/cardtrader"
	"card-watcher/internal/models"
	"card-watcher/internal/mongo"
	"context"
	"crypto/sha256"
)

type Service interface {
	SaveWatch(ctx context.Context, expansionId, blueprintId int, condition models.Condition, foil bool) (string, error)
	ListExpansions(ctx context.Context, name, code string) (models.ListExpansionsResponse, error)
	ListBlueprints(ctx context.Context, expansionId int, name string) (models.ListBlueprintsResponse, error)
}

type service struct {
	cardtraderAdapter cardtrader.CardtraderAdapter
	mongoAdapter      mongo.MongoAdapter
}

func NewService(
	cardtraderAdapter cardtrader.CardtraderAdapter,
	mongoAdapter mongo.MongoAdapter,
) *service {
	return &service{
		cardtraderAdapter: cardtraderAdapter,
		mongoAdapter:      mongoAdapter,
	}
}

func HashAccessToken(accessToken string) string {
	h := sha256.New()
	h.Write([]byte(accessToken))
	return string(h.Sum(nil))
}
