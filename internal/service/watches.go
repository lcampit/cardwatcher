package service

import (
	"context"
	"fmt"
	"log/slog"

	"card-watcher/internal/entities"
	"card-watcher/internal/models"
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

	s.logger.Debug("returning watches", slog.Int("watchCount", len(result)))
	return models.ListWatchesResponse{
		Watches: result,
	}, nil
}

func (s *service) SaveWatch(ctx context.Context, expansionID, blueprintID int, condition models.Condition, foil bool) (string, error) {
	blueprintName, err := s.cardtraderAdapter.GetBlueprintNameByExpansionID(ctx, expansionID, blueprintID)
	if err != nil {
		return "", fmt.Errorf("error finding name for expansion %d and blueprint %d: %w", expansionID, blueprintID, err)
	}
	expansionName, err := s.cardtraderAdapter.GetExpansionNameByID(ctx, expansionID)
	if err != nil {
		return "", fmt.Errorf("error finding name for expansion %d: %w", expansionID, err)
	}
	newWatchID, err := s.mongoAdapter.SaveWatch(ctx, &entities.Watch{
		Name:          blueprintName,
		ExpansionId:   expansionID,
		ExpansionName: expansionName,
		BlueprintId:   blueprintID,
		Condition:     convertModelConditionToEntityCondition(condition),
		Foil:          foil,
	})
	if err != nil {
		return "", err
	}
	return newWatchID, nil
}

func (s *service) DeleteWatchByID(ctx context.Context, watchID string) error {
	return s.mongoAdapter.DeleteWatchByID(ctx, watchID)
}
