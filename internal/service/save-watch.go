package service

import (
	"card-watcher/internal/entities"
	"card-watcher/internal/models"
	"context"
)

func (s *service) SaveWatch(ctx context.Context, accessToken, blueprintId, expansionId string, condition models.Condition, foil bool) (string, error) {
	userId := HashAccessToken(accessToken)
	// TODO: consider calling cardtrader endpoint to retrieve card name starting from blueprint id and
	// saving it to db
	newWatchId, err := s.mongoAdapter.SaveWatch(ctx, &entities.Watch{
		UserId:      userId,
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
