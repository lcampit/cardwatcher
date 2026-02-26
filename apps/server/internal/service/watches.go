package service

import (
	"context"
	"fmt"
	"log/slog"

	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
)

func (s *service) ListWatches(ctx context.Context) (*apiv1.ListWatchesResponse, error) {
	watches, err := s.mongoAdapter.GetWatches(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting watches from mongo adapter: %w", err)
	}

	var result []*apiv1.Watch

	for _, entity := range watches {
		result = append(result, convertEntityWatchToModelWatch(entity))
	}

	s.logger.Debug("returning watches", slog.Int("watchCount", len(result)))
	return &apiv1.ListWatchesResponse{
		Watches: result,
	}, nil
}

func (s *service) SaveWatch(ctx context.Context, expansionID, blueprintID uint64, condition apiv1.Condition, foil bool) (string, error) {
	blueprintName, err := s.cardtraderAdapter.GetBlueprintNameByExpansionID(ctx, expansionID, blueprintID)
	if err != nil {
		return "", fmt.Errorf("finding name for expansion %d and blueprint %d: %w", expansionID, blueprintID, err)
	}
	expansionName, err := s.cardtraderAdapter.GetExpansionNameByID(ctx, expansionID)
	if err != nil {
		return "", fmt.Errorf("finding name for expansion %d: %w", expansionID, err)
	}
	newWatchID, err := s.mongoAdapter.SaveWatch(ctx, &mongo.Watch{
		Name:          blueprintName,
		ExpansionID:   expansionID,
		ExpansionName: expansionName,
		BlueprintID:   blueprintID,
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
