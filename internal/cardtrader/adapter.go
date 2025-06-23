package cardtrader

import "context"

type CardtraderAdapter interface {
	GetBlueprintNameByExpansionId(ctx context.Context, expansionId, blueprintId int) (string, error)
	GetExpansions(ctx context.Context) ([]*expansion, error)
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
