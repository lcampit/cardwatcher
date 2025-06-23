package cardtrader

import "context"

type CardtraderAdapter interface {
	GetBlueprintNameByExpansionId(ctx context.Context, expansionId, blueprintId int) (string, error)
	GetExpansions(ctx context.Context) ([]*expansion, error)
	GetBlueprints(ctx context.Context, expansionId int) ([]*blueprint, error)
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
