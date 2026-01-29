// Package service implements all business logic
package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"card-watcher/internal/models"
)

func (s *service) ListBlueprints(ctx context.Context, expansionID uint64, name string) (*models.ListBlueprintsResponse, error) {
	blueprints, err := s.cardtraderAdapter.GetBlueprints(ctx, expansionID)
	if err != nil {
		return nil, fmt.Errorf("getting blueprints from cardtrader adapter: %w", err)
	}

	var resultingBlueprints []*models.Blueprint
	normalizedName := strings.ToLower(name)
	for _, blueprint := range blueprints {
		if name != "" {
			// filter via card name
			if strings.Contains(strings.ToLower(blueprint.Name), normalizedName) {
				resultingBlueprints = append(resultingBlueprints, &models.Blueprint{
					Id:          blueprint.ID,
					Name:        blueprint.Name,
					ExpansionId: blueprint.ExpansionID,
				})
			}
		} else {
			// no filter provided, return all cards from given expansion
			resultingBlueprints = append(resultingBlueprints, &models.Blueprint{
				Id:          blueprint.ID,
				Name:        blueprint.Name,
				ExpansionId: blueprint.ExpansionID,
			})
		}
	}
	s.logger.Debug("returning filtered blueprints", slog.Int("blueprintsCount", len(resultingBlueprints)))
	return &models.ListBlueprintsResponse{
		Blueprints: resultingBlueprints,
	}, nil
}
