package cardtrader

import "context"

type CardtraderAdapter interface {
	Info(ctx context.Context) (*InfoResponse, error)
}

type cardtraderAdapter struct {
	// TODO: move token to watcher requests to allow for multiple users
	accessToken string
	// NOTE: can this be a constant?
	baseUrl string
}

func NewCardtraderAdapter(accessToken, baseUrl string) CardtraderAdapter {
	return &cardtraderAdapter{
		accessToken: accessToken,
		baseUrl:     baseUrl,
	}
}
