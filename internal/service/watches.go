package service

import (
	"card-watcher/internal/entities"
	"card-watcher/internal/models"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (s *service) ListWatches(ctx context.Context) (models.ListWatchesResponse, error) {
	watches, err := s.mongoAdapter.GetWatches(ctx)
	if err != nil {
		return models.ListWatchesResponse{}, fmt.Errorf("error getting watches from adapter: %w", err)
	}

	var result []*models.Watch

	for _, entity := range watches {
		result = append(result, convertEntityWatchToModelWatch(entity))
	}

	log.Debug().Msgf("returning %d watches", len(result))
	return models.ListWatchesResponse{
		Watches: result,
	}, nil
}

func (s *service) SaveWatch(ctx context.Context, expansionId, blueprintId int, condition models.Condition, foil bool) (string, error) {
	blueprintName, err := s.cardtraderAdapter.GetBlueprintNameByExpansionId(ctx, expansionId, blueprintId)
	if err != nil {
		return "", fmt.Errorf("error finding name for expansion %d and blueprint %d: %w", expansionId, blueprintId, err)
	}
	newWatchId, err := s.mongoAdapter.SaveWatch(ctx, &entities.Watch{
		Name:        blueprintName,
		ExpansionId: expansionId,
		BlueprintId: blueprintId,
		Condition:   convertModelConditionToEntityCondition(condition),
		Foil:        foil,
	})
	if err != nil {
		return "", err
	}
	return newWatchId, nil
}
