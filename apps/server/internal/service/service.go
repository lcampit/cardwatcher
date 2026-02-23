package service

import (
    "context"
    "crypto/sha256"
    "log/slog"
    "sync"
    "time"

    api "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
    cardtraderpkg "github.com/lcampit/cardwatcher/apps/server/internal/cardtrader"
    mongopkg "github.com/lcampit/cardwatcher/apps/server/internal/mongo"
    ntfypkg "github.com/lcampit/cardwatcher/apps/server/internal/ntfy"

    "github.com/robfig/cron/v3"
)

type Service interface {
    SaveWatch(ctx context.Context, expansionID, blueprintID uint64, condition api.Condition, foil bool) (string, error)
    ListExpansions(ctx context.Context, gameName, expansionName, expansionCode string) (*api.ListExpansionsResponse, error)
    ListBlueprints(ctx context.Context, expansionID uint64, name string) (*api.ListBlueprintsResponse, error)
    ListWatches(ctx context.Context) (*api.ListWatchesResponse, error)
    DeleteWatchByID(ctx context.Context, watchID string) error

    Close()
}

type service struct {
    logger            *slog.Logger
    cardtraderAdapter cardtraderpkg.CardtraderAdapter
    mongoAdapter      mongopkg.MongoAdapter
    ntfyAdapter       ntfypkg.NtfyAdapter

    cron      *cron.Cron
    gameIDMap sync.Map
}

type ServiceConfig struct {
    Logger               *slog.Logger
    CardtraderAdapter    cardtraderpkg.CardtraderAdapter
    MongoAdapter         mongopkg.MongoAdapter
    NtfyAdapter          ntfypkg.NtfyAdapter
    NotificationSchedule string
    UpdateMapsSchedule   string
}

func NewService(ctx context.Context, config ServiceConfig) *service {
    service := &service{
        logger:            config.Logger,
        cardtraderAdapter: config.CardtraderAdapter,
        mongoAdapter:      config.MongoAdapter,
        ntfyAdapter:       config.NtfyAdapter,
    }
    service.gameIDMap = sync.Map{}
    service.updateGamesMap()

    loc, _ := time.LoadLocation("Europe/Rome")
    c := cron.New(cron.WithLocation(loc))
    _, err := c.AddFunc(config.NotificationSchedule, service.watchAndNotify)
    if err != nil {
        service.logger.Error("setting up notification cron job", slog.Any("error", err))
    }
    _, err = c.AddFunc(config.UpdateMapsSchedule, service.updateGamesMap)
    if err != nil {
        service.logger.Error("setting up notification cron job", slog.Any("error", err))
    }
    c.Start()
    service.cron = c

    return service
}

func (s *service) Close() {
    s.cron.Stop()
}

func HashAccessToken(accessToken string) string {
    h := sha256.New()
    h.Write([]byte(accessToken))
    return string(h.Sum(nil))
}

func (s *service) updateGamesMap() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    games, err := s.cardtraderAdapter.GetGames(ctx)
    if err != nil {
        s.logger.Error("getting games from cardtrader adapter", slog.Any("error", err))
    }

    for _, game := range games {
        s.gameIDMap.Store(game.GetNormalizedName(), game.ID)
    }
}
