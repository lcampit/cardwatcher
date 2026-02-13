package cardtrader

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/carlmjohnson/requests"
)

type game struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

func (g *game) GetNormalizedName() string {
	return strings.TrimSpace(strings.ToLower(g.Name))
}

func (a *cardtraderAdapter) GetGames(ctx context.Context) ([]*game, error) {
	var response []*game
	endpoint := fmt.Sprintf("%s/%s", a.baseURL, "games")
	err := requests.URL(endpoint).Bearer(a.accessToken).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("cardtrader get expansions endpoint %w", err)
	}
	a.logger.Debug("received games", slog.Int("gamesCount", len(response)))
	return response, nil
}
