package cardtrader

import (
	"card-watcher/internal/entities"
	"context"
)

type CardtraderAdapter interface {
	GetBlueprintNameByExpansionId(ctx context.Context, expansionId, blueprintId int) (string, error)
	GetExpansionNameByID(ctx context.Context, expansionID int) (string, error)
	GetExpansions(ctx context.Context) ([]*expansion, error)
	GetBlueprints(ctx context.Context, expansionId int) ([]*blueprint, error)
	GetCurrentPricing(ctx context.Context, watch *entities.Watch) (int, error)
}

type cardtraderAdapter struct {
	baseUrl     string
	accessToken string
}

func NewCardtraderAdapter(accessToken, baseUrl string) CardtraderAdapter {
	return &cardtraderAdapter{
		baseUrl:     baseUrl,
		accessToken: accessToken,
	}
}
