package service

import (
	"card-watcher/internal/entities"
	"card-watcher/internal/models"
	"context"
	"fmt"
)

func (s *service) SaveWatch(ctx context.Context, accessToken string, expansionId, blueprintId int, condition models.Condition, foil bool) (string, error) {
	userId := HashAccessToken(accessToken)
	blueprintName, err := s.cardtraderAdapter.GetBlueprintNameByExpansionId(ctx, accessToken, expansionId, blueprintId)
	if err != nil {
		return "", fmt.Errorf("error finding name for expansion %d and blueprint %d: %w", expansionId, blueprintId, err)
	}
	newWatchId, err := s.mongoAdapter.SaveWatch(ctx, &entities.Watch{
		UserId:      userId,
		Name:        blueprintName,
		ExpansionId: expansionId,
		BlueprintId: blueprintId,
		Condition:   entities.ConvertModelConditionToWatchCondition(condition),
		Foil:        foil,
	})
	if err != nil {
		return "", err
	}
	return newWatchId, nil
}
