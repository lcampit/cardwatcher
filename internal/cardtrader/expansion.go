package cardtrader

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/rs/zerolog/log"
)

type expansion struct {
	ID     int    `json:"id"`
	GameID int    `json:"game_id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
}

func (a *cardtraderAdapter) GetExpansions(ctx context.Context) ([]*expansion, error) {
	var response []*expansion
	endpoint := fmt.Sprintf("%s/%s", a.baseUrl, "expansions")
	err := requests.URL(endpoint).Bearer(a.accessToken).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in cardtrader expansions endpoint %w", err)
	}
	log.Debug().Msgf("received %d expansions", len(response))
	return response, nil
}
