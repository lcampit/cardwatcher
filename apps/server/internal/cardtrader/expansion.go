package cardtrader

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

type Expansion struct {
	ID     uint64 `json:"id"`
	GameID uint64 `json:"game_id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
}

func (e *Expansion) GetNormalizedName() string {
	return strings.TrimSpace(strings.ToLower(e.Name))
}

func (e *Expansion) GetNormalizedCode() string {
	return strings.TrimSpace(strings.ToLower(e.Code))
}

func (a *cardtraderAdapter) GetExpansions(ctx context.Context) ([]*Expansion, error) {
	var response []*Expansion
	endpoint := "expansions"
	_, err := a.client.R().
		SetResult(&response).
		Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("cardtrader get expansions endpoint %w", err)
	}
	a.logger.Debug("received expansions", slog.Int("expansion_count", len(response)))
	return response, nil
}

func (a *cardtraderAdapter) GetExpansionNameByID(ctx context.Context, expansionID uint64) (string, error) {
	var response []*Expansion
	endpoint := "expansions"
	_, err := a.client.R().
		SetResult(&response).
		Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("cardtrader get expansions endpoint for id %d: %w ", expansionID, err)
	}
	for _, expansion := range response {
		if expansion.ID == expansionID {
			a.logger.Debug("found name for expansion id",
				slog.Uint64("expansion_id", expansionID),
				slog.String("expansion_name", expansion.Name))
			return expansion.Name, nil
		}
	}
	return "", fmt.Errorf("no name found for expansion id %d", expansionID)
}
