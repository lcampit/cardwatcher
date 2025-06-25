package service

import (
	"card-watcher/internal/models"
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func (s *service) ListBlueprints(ctx context.Context, expansionId int, name string) (models.ListBlueprintsResponse, error) {
	blueprints, err := s.cardtraderAdapter.GetBlueprints(ctx, expansionId)
	if err != nil {
		return models.ListBlueprintsResponse{}, fmt.Errorf("error getting blueprints from adapter: %w", err)
	}

	var resultingBlueprints []*models.Blueprint
	normalizedName := strings.ToLower(name)
	for _, blueprint := range blueprints {
		if name != "" {
			// filter via card name
			if strings.Contains(strings.ToLower(blueprint.Name), normalizedName) {
				resultingBlueprints = append(resultingBlueprints, &models.Blueprint{
					Id:          int32(blueprint.Id),
					Name:        blueprint.Name,
					ExpansionId: int32(blueprint.ExpansionId),
				})
			}
		} else {
			// no filter provided, return all cards from given expansion
			resultingBlueprints = append(resultingBlueprints, &models.Blueprint{
				Id:          int32(blueprint.Id),
				Name:        blueprint.Name,
				ExpansionId: int32(blueprint.ExpansionId),
			})
		}
	}
	log.Debug().Msgf("returning %d filtered blueprints", len(resultingBlueprints))
	return models.ListBlueprintsResponse{
		Blueprints: resultingBlueprints,
	}, nil
}
