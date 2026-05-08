package cardtrader

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

type game struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type gameResponse struct {
	GameList []*game `json:"array"`
}

func (g *game) GetNormalizedName() string {
	return strings.TrimSpace(strings.ToLower(g.Name))
}

func (a *cardtraderAdapter) GetGames(ctx context.Context) ([]*game, error) {
	var response gameResponse
	endpoint := "games"
	_, err := a.client.R().
		SetResult(&response).
		Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("cardtrader get games endpoint %w", err)
	}
	a.logger.Debug("received games", slog.Int("gamesCount", len(response.GameList)))
	return response.GameList, nil
}
