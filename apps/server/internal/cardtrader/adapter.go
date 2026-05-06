// Package cardtrader contains implementations
// to use cardtrader APIs
//
// First create an adapter specifying base url and access token
package cardtrader

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"

	"github.com/carlmjohnson/requests"
)

type CardtraderAdapter interface {
	GetGames(ctx context.Context) ([]*game, error)
	GetBlueprintNameByExpansionID(ctx context.Context, expansionID, blueprintID uint64) (string, error)
	GetExpansionNameByID(ctx context.Context, expansionID uint64) (string, error)
	GetExpansions(ctx context.Context) ([]*expansion, error)
	GetBlueprints(ctx context.Context, expansionID uint64) ([]*blueprint, error)
	GetProducts(ctx context.Context, blueprintID uint64, foil bool) ([]Product, error)
}

type cardtraderAdapter struct {
	logger *slog.Logger
	client *requests.Builder
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
	transport := http.Transport{
		ForceAttemptHTTP2: true,
		TLSClientConfig:   &tlsConfig,
	}
	builder := requests.New(
		func(rb *requests.Builder) {
			rb.BaseURL(config.BaseURL)
			rb.Bearer(config.AccessToken)
			rb.Transport(&transport)
		},
	)
	return &cardtraderAdapter{
		logger: config.Logger,
		client: builder,
	}
}
