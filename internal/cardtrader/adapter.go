package cardtrader

import "context"

type CardtraderAdapter interface {
	Info(ctx context.Context) (*InfoResponse, error)
	GetExpansions(ctx context.Context) (*[]Expansion, error)
}

type cardtraderAdapter struct {
	// TODO: move token to watcher requests to allow for multiple users
	accessToken string
	baseUrl     string
}

func NewCardtraderAdapter(accessToken, baseUrl string) CardtraderAdapter {
	return &cardtraderAdapter{
		accessToken: accessToken,
		baseUrl:     baseUrl,
	}
}
