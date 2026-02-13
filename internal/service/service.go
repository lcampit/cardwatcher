package service

import (
	"context"
	"crypto/sha256"
	"log/slog"
	"sync"

	"card-watcher/internal/cardtrader"
	"card-watcher/internal/models"
	"card-watcher/internal/mongo"
	"card-watcher/internal/ntfy"
)

type Service interface {
	SaveWatch(ctx context.Context, expansionID, blueprintID uint64, condition models.Condition, foil bool) (string, error)
	ListExpansions(ctx context.Context, gameName, expansionName, expansionCode string) (*models.ListExpansionsResponse, error)
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

	gameIDMap sync.Map
}

type ServiceConfig struct {
	Logger            *slog.Logger
	CardtraderAdapter cardtrader.CardtraderAdapter
	MongoAdapter      mongo.MongoAdapter
	NtfyAdapter       ntfy.NtfyAdapter
}

func NewService(ctx context.Context, config ServiceConfig) *service {
	service := &service{
		logger:            config.Logger,
		cardtraderAdapter: config.CardtraderAdapter,
		mongoAdapter:      config.MongoAdapter,
		ntfyAdapter:       config.NtfyAdapter,
	}
	service.gameIDMap = sync.Map{}
	service.updateGamesMap(ctx)
	return service
}

func HashAccessToken(accessToken string) string {
	h := sha256.New()
	h.Write([]byte(accessToken))
	return string(h.Sum(nil))
}

func (s *service) updateGamesMap(ctx context.Context) {
	games, err := s.cardtraderAdapter.GetGames(ctx)
	if err != nil {
		s.logger.Error("getting games from cardtrader adapter", slog.Any("error", err))
	}

	for _, game := range games {
		s.gameIDMap.Store(game.GetNormalizedName(), game.ID)
	}
}
