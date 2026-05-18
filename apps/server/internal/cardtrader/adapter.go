// Package cardtrader contains implementations
// to use cardtrader APIs
//
// First create an adapter specifying base url and access token
package cardtrader

import (
	"context"
	"crypto/tls"
	"log/slog"

	"resty.dev/v3"
)

type CardtraderAdapter interface {
	GetGames(ctx context.Context) ([]*game, error)
	GetBlueprintNameByExpansionID(ctx context.Context, expansionID, blueprintID uint64) (string, error)
	GetExpansionNameByID(ctx context.Context, expansionID uint64) (string, error)
	GetExpansions(ctx context.Context) ([]*Expansion, error)
	GetBlueprints(ctx context.Context, expansionID uint64) ([]*blueprint, error)
	GetProducts(ctx context.Context, blueprintID uint64, foil bool) ([]Product, error)
}

type cardtraderAdapter struct {
	logger *slog.Logger
	client *resty.Client
}

type CardtraderAdapterConfig struct {
	Logger      *slog.Logger
	AccessToken string
	BaseURL     string
	// This options should only be used for testing
	SkipVerify bool
}

func NewCardtraderAdapter(config CardtraderAdapterConfig) CardtraderAdapter {
	tlsConfig := tls.Config{
		InsecureSkipVerify: config.SkipVerify,
	}

	client := resty.New().
		SetTLSClientConfig(&tlsConfig).
		SetAuthToken(config.AccessToken).
		SetBaseURL(config.BaseURL)

	return &cardtraderAdapter{
		logger: config.Logger,
		client: client,
	}
}
