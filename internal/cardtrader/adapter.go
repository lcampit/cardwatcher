package cardtrader

import "context"

type CardtraderAdapter interface {
	GetBlueprintNameByExpansionId(ctx context.Context, accessToken string, expansionId, blueprintId int) (string, error)
}

type cardtraderAdapter struct {
	baseUrl string
}

func NewCardtraderAdapter(accessToken, baseUrl string) CardtraderAdapter {
	return &cardtraderAdapter{
		baseUrl: baseUrl,
	}
}
