package cardtrader

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/carlmjohnson/requests"
)

type expansion struct {
	ID     int    `json:"id"`
	GameID int    `json:"game_id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
}

func (a *cardtraderAdapter) GetExpansions(ctx context.Context) ([]*expansion, error) {
	var response []*expansion
	endpoint := fmt.Sprintf("%s/%s", a.baseURL, "expansions")
	err := requests.URL(endpoint).Bearer(a.accessToken).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in cardtrader expansions endpoint %w", err)
	}
	a.logger.Debug("received expansions", slog.Int("expansionCount", len(response)))
	return response, nil
}

func (a *cardtraderAdapter) GetExpansionNameByID(ctx context.Context, expansionID int) (string, error) {
	var response []*expansion
	endpoint := fmt.Sprintf("%s/%s", a.baseURL, "expansions")
	err := requests.URL(endpoint).Bearer(a.accessToken).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return "", fmt.Errorf("error in cardtrader expansions endpoint %w", err)
	}
	for _, expansion := range response {
		if expansion.ID == expansionID {
			a.logger.Debug("found name for expansion id",
				slog.Int("expansionId", expansionID),
				slog.String("expansionName", expansion.Name))
			return expansion.Name, nil
		}
	}
	return "", fmt.Errorf("no name found for expansion id %d", expansionID)
}
