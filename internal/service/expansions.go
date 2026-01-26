package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"card-watcher/internal/models"
)

func (s *service) ListExpansions(ctx context.Context, name, code string) (*models.ListExpansionsResponse, error) {
	expansions, err := s.cardtraderAdapter.GetExpansions(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting expansions from adapter: %w", err)
	}
	var resultingExpanions []*models.Expansion
	normalizedName := strings.ToLower(name)
	normalizedCode := strings.ToLower(code)
	for _, expansion := range expansions {
		// filter via name
		if name != "" {
			if strings.Contains(strings.ToLower(expansion.Name), normalizedName) {
				resultingExpanions = append(resultingExpanions, &models.Expansion{
					Id:   int32(expansion.ID),
					Code: expansion.Code,
					Name: expansion.Name,
				})
			}
		} else if code != "" {
			// filter via code
			if strings.Contains(expansion.Code, normalizedCode) {
				resultingExpanions = append(resultingExpanions, &models.Expansion{
					Id:   int32(expansion.ID),
					Code: expansion.Code,
					Name: expansion.Name,
				})
			}
		} else {
			// no filter provided, return all expansions
			resultingExpanions = append(resultingExpanions, &models.Expansion{
				Id:   int32(expansion.ID),
				Code: expansion.Code,
				Name: expansion.Name,
			})
		}
	}
	s.logger.Debug("returning filtered expansions", slog.Int("expansionCount", len(resultingExpanions)))
	return &models.ListExpansionsResponse{
		Expansions: resultingExpanions,
	}, nil
}
