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
	GetBlueprintNameByExpansionID(ctx context.Context, expansionID, blueprintID int) (string, error)
	GetExpansionNameByID(ctx context.Context, expansionID int) (string, error)
	GetExpansions(ctx context.Context) ([]*expansion, error)
	GetBlueprints(ctx context.Context, expansionID int) ([]*blueprint, error)
	GetCurrentPricing(ctx context.Context, watch *entities.Watch) (int, error)
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
