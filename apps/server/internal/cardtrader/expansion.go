package cardtrader

import (
	"context"
	"fmt"
	"log/slog"
)

type expansion struct {
	ID     uint64 `json:"id"`
	GameID uint64 `json:"game_id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
}

func (a *cardtraderAdapter) GetExpansions(ctx context.Context) ([]*expansion, error) {
	var response []*expansion
	endpoint := "expansions"
	err := a.client.Path(endpoint).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("cardtrader get expansions endpoint %w", err)
	}
	a.logger.Debug("received expansions", slog.Int("expansionCount", len(response)))
	return response, nil
}

func (a *cardtraderAdapter) GetExpansionNameByID(ctx context.Context, expansionID uint64) (string, error) {
	var response []*expansion
	endpoint := "expansions"
	err := a.client.Path(endpoint).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return "", fmt.Errorf("cardtrader get expansions endpoint for id %d: %w ", expansionID, err)
	}
	for _, expansion := range response {
		if expansion.ID == expansionID {
			a.logger.Debug("found name for expansion id",
				slog.Uint64("expansionId", expansionID),
				slog.String("expansionName", expansion.Name))
			return expansion.Name, nil
		}
	}
	return "", fmt.Errorf("no name found for expansion id %d", expansionID)
}
