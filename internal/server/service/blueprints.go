// Package service implements all business logic
package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	api "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

func (s *service) ListBlueprints(ctx context.Context, expansionID uint64, name string) (*api.ListBlueprintsResponse, error) {
	blueprints, err := s.cardtraderAdapter.GetBlueprints(ctx, expansionID)
	if err != nil {
		return nil, fmt.Errorf("getting blueprints from cardtrader adapter: %w", err)
	}

	var resultingBlueprints []*api.Blueprint
	normalizedName := strings.ToLower(name)
	for _, blueprint := range blueprints {
		if name != "" {
			// filter via card name
			if strings.Contains(strings.ToLower(blueprint.Name), normalizedName) {
				resultingBlueprints = append(resultingBlueprints, &api.Blueprint{
					Id:          blueprint.ID,
					Name:        blueprint.Name,
					ExpansionId: blueprint.ExpansionID,
				})
			}
		} else {
			// no filter provided, return all cards from given expansion
			resultingBlueprints = append(resultingBlueprints, &api.Blueprint{
				Id:          blueprint.ID,
				Name:        blueprint.Name,
				ExpansionId: blueprint.ExpansionID,
			})
		}
	}
	s.logger.Debug("returning filtered blueprints", slog.Int("blueprintsCount", len(resultingBlueprints)))
	return &api.ListBlueprintsResponse{
		Blueprints: resultingBlueprints,
	}, nil
}
