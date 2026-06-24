package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/lcampit/cardwatcher/apps/server/internal/cardtrader"
	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
	"github.com/lcampit/cardwatcher/apps/server/internal/ntfy"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"

	"github.com/robfig/cron/v3"
)

type Service interface {
	CreateWatch(ctx context.Context, request *apiv1.CreateWatchRequest) (string, error)
	ListExpansions(ctx context.Context, gameName, expansionName, expansionCode string) (*apiv1.ListExpansionsResponse, error)
	ListBlueprints(ctx context.Context, expansionID uint64, name string) (*apiv1.ListBlueprintsResponse, error)
	ListWatches(ctx context.Context) (*apiv1.ListWatchesResponse, error)
	DeleteWatchByID(ctx context.Context, watchID string) error

	Health(ctx context.Context) error
	Close()
}

type service struct {
	logger            *slog.Logger
	cardtraderAdapter cardtrader.CardtraderAdapter
	mongoAdapter      mongo.MongoAdapter
	ntfyAdapter       ntfy.NtfyAdapter

	cron             *cron.Cron
	gameNameToIDMap  sync.Map
	expansionNameMap sync.Map
	expansionCodeMap sync.Map
}

type ServiceConfig struct {
	Logger               *slog.Logger
	CardtraderAdapter    cardtrader.CardtraderAdapter
	MongoAdapter         mongo.MongoAdapter
	NtfyAdapter          ntfy.NtfyAdapter
	NotificationSchedule string
	UpdateMapsSchedule   string
}

func (config ServiceConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("notification_schedule", config.NotificationSchedule),
		slog.String("update_maps_schedule", config.UpdateMapsSchedule),
	)
}

func NewService(ctx context.Context, config ServiceConfig) *service {
	config.Logger.Debug("creating service", slog.Any("config", config))
	service := &service{
		logger:            config.Logger,
		cardtraderAdapter: config.CardtraderAdapter,
		mongoAdapter:      config.MongoAdapter,
		ntfyAdapter:       config.NtfyAdapter,
	}
	service.gameNameToIDMap = sync.Map{}
	service.updateGamesMap()
	service.updateExpansionsMaps()

	loc, _ := time.LoadLocation("Europe/Rome")
	c := cron.New(cron.WithLocation(loc))
	_, err := c.AddFunc(config.NotificationSchedule, service.watchAndNotify)
	if err != nil {
		service.logger.Error("setting up notification cron job", slog.Any("error", err))
	}
	_, err = c.AddFunc(config.UpdateMapsSchedule, service.updateGamesMap)
	if err != nil {
		service.logger.Error("setting up update games map cron job", slog.Any("error", err))
	}
	_, err = c.AddFunc(config.UpdateMapsSchedule, service.updateExpansionsMaps)
	if err != nil {
		service.logger.Error("setting up update expansions maps cron job", slog.Any("error", err))
	}
	c.Start()
	service.cron = c

	return service
}

func (s *service) Close() {
	s.cron.Stop()
}

func (s *service) updateGamesMap() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	games, err := s.cardtraderAdapter.GetGames(ctx)
	if err != nil {
		s.logger.Error("getting games from cardtrader adapter", slog.Any("error", err))
		return
	}

	for _, game := range games {
		s.gameNameToIDMap.Store(game.GetNormalizedName(), game.ID)
	}
}

func (s *service) updateExpansionsMaps() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	expansions, err := s.cardtraderAdapter.GetExpansions(ctx)
	if err != nil {
		s.logger.Error("getting expansions from cardtrader adapter", slog.Any("error", err))
		return
	}

	for _, expansion := range expansions {
		s.expansionCodeMap.Store(expansion.GetNormalizedCode(), expansion)
		s.expansionNameMap.Store(expansion.GetNormalizedName(), expansion)
	}
}

// getExpansionFromMaps returns the preloaded expansion
// from either one of the preloaded expansion maps.
//
// Returns a control boolean value if the given expansion name or code
// was found in the preloaded maps, false otherwise
func (s *service) getExpansionFromMaps(expansionNameOrCode string) (*cardtrader.Expansion, bool) {
	nameOrCodeRequested := normalizeString(expansionNameOrCode)
	savedExpansion, ok := s.expansionCodeMap.Load(nameOrCodeRequested)
	if ok {
		result, _ := savedExpansion.(*cardtrader.Expansion)
		return result, true
	}

	savedExpansion, ok = s.expansionNameMap.Load(nameOrCodeRequested)
	if ok {
		result, _ := savedExpansion.(*cardtrader.Expansion)
		return result, true
	}

	return nil, false
}

func HashAccessToken(accessToken string) string {
	h := sha256.New()
	h.Write([]byte(accessToken))
	return string(h.Sum(nil))
}

func normalizeString(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func (s *service) Health(ctx context.Context) error {
	return errors.Join(
		s.mongoAdapter.Health(ctx),
		s.cardtraderAdapter.Health(ctx),
	)
}
