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

func (a *cardtraderAdapter) GetExpansionNameByID(ctx context.Context, expansionID int) (string, error) {
	var response []*expansion
	endpoint := fmt.Sprintf("%s/%s", a.baseUrl, "expansions")
	err := requests.URL(endpoint).Bearer(a.accessToken).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return "", fmt.Errorf("error in cardtrader expansions endpoint %w", err)
	}
	for _, expansion := range response {
		if expansion.ID == expansionID {
			log.Debug().Msgf("found name for expansion id %d: %s", expansionID, expansion.Name)
			return expansion.Name, nil
		}
	}
	return "", fmt.Errorf("no name found for expansion id %d", expansionID)
}
