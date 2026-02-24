package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	api "github.com/lcampit/cardwatcher/internal/api/v1"
)

func (s *service) ListExpansions(ctx context.Context, gameName, expansionName, expansionCode string) (*api.ListExpansionsResponse, error) {
	expansions, err := s.cardtraderAdapter.GetExpansions(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting expansions from cardtrader adapter: %w", err)
	}
	var resultingExpanions []*api.Expansion
	var gameID uint64
	normalizedExpanionName := strings.ToLower(expansionName)
	normalizedExpansionCode := strings.ToLower(expansionCode)

	if gameName != "" {
		normalizedGameName := strings.ToLower(gameName)
		gameIDFromMap, ok := s.gameIDMap.Load(normalizedGameName)
		if !ok {
			s.logger.Debug("filtering expansions for game name: game not found in map", slog.String("gameName", gameName))
		}
		gameID, ok = (gameIDFromMap).(uint64)
		if !ok {
			s.logger.Error("filtering expansions for game name: ID found in map is not an int",
				slog.String("gameName", gameName),
				slog.Any("gameIDFromMap", gameIDFromMap))
		}
	}
	for _, expansion := range expansions {
		if gameID != 0 && expansion.GameID != gameID {
			// expansion is not of the right game, skip it
			continue
		}

		if expansionName != "" {
			// filter via name
			if strings.Contains(strings.ToLower(expansion.Name), normalizedExpanionName) {
				resultingExpanions = append(resultingExpanions, &api.Expansion{
					Id:   expansion.ID,
					Code: expansion.Code,
					Name: expansion.Name,
				})
			}
		} else if expansionCode != "" {
			// filter via code
			if strings.Contains(expansion.Code, normalizedExpansionCode) {
				resultingExpanions = append(resultingExpanions, &api.Expansion{
					Id:   expansion.ID,
					Code: expansion.Code,
					Name: expansion.Name,
				})
			}
		} else {
			// no filter provided, return all expansions of the given game, if any
			resultingExpanions = append(resultingExpanions, &api.Expansion{
				Id:   expansion.ID,
				Code: expansion.Code,
				Name: expansion.Name,
			})
		}
	}
	s.logger.Debug("returning filtered expansions", slog.Int("expansionCount", len(resultingExpanions)))
	return &api.ListExpansionsResponse{
		Expansions: resultingExpanions,
	}, nil
}
