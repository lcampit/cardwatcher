// Package cardtrader contains implementations
// to use cardtrader APIs
//
// First create an adapter specifying base url and access token
package cardtrader

import (
	"context"
	"log/slog"

	"card-watcher/internal/entities"
)

type CardtraderAdapter interface {
	GetGames(ctx context.Context) ([]*game, error)
	GetBlueprintNameByExpansionID(ctx context.Context, expansionID, blueprintID uint64) (string, error)
	GetExpansionNameByID(ctx context.Context, expansionID uint64) (string, error)
	GetExpansions(ctx context.Context) ([]*expansion, error)
	GetBlueprints(ctx context.Context, expansionID uint64) ([]*blueprint, error)
	GetCurrentPricingCents(ctx context.Context, watch *entities.Watch) (uint64, error)
}

type cardtraderAdapter struct {
	logger      *slog.Logger
	baseURL     string
	accessToken string
}

type CardtraderAdapterConfig struct {
	Logger      *slog.Logger
	AccessToken string
	BaseURL     string
}

func NewCardtraderAdapter(config CardtraderAdapterConfig) CardtraderAdapter {
	return &cardtraderAdapter{
		logger:      config.Logger,
		baseURL:     config.BaseURL,
		accessToken: config.AccessToken,
	}
}
