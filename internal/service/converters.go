package service

import (
	"card-watcher/internal/entities"
	"card-watcher/internal/models"
)

func convertModelConditionToEntityCondition(modelCondition models.Condition) entities.WatchCondition {
	switch modelCondition {
	case models.Condition_CONDITION_UNSPECIFIED:
	case models.Condition_CONDITION_NEAR_MINT:
		return entities.WATCH_CONDITION_NEAR_MINT
	case models.Condition_CONDITION_SLIGHTLY_PLAYED:
		return entities.WATCH_CONDITION_SLIGHTLY_PLAYED
	case models.Condition_CONDITION_MODERATELY_PLAYED:
		return entities.WATCH_CONDITION_MODERATELY_PLAYED
	case models.Condition_CONDITION_PLAYED:
		return entities.WATCH_CONDITION_PLAYED
	case models.Condition_CONDITION_POOR:
		return entities.WATCH_CONDITION_POOR
	default:
		return entities.WATCH_CONDITION_NEAR_MINT
	}
	return entities.WATCH_CONDITION_NEAR_MINT
}

func convertEntityConditionToModelCondition(entityCondition entities.WatchCondition) models.Condition {
	switch entityCondition {
	case entities.WATCH_CONDITION_NEAR_MINT:
		return models.Condition_CONDITION_NEAR_MINT
	case entities.WATCH_CONDITION_SLIGHTLY_PLAYED:
		return models.Condition_CONDITION_SLIGHTLY_PLAYED
	case entities.WATCH_CONDITION_MODERATELY_PLAYED:
		return models.Condition_CONDITION_SLIGHTLY_PLAYED
	case entities.WATCH_CONDITION_PLAYED:
		return models.Condition_CONDITION_PLAYED
	case entities.WATCH_CONDITION_POOR:
		return models.Condition_CONDITION_POOR
	}
	return models.Condition_CONDITION_NEAR_MINT
}

func convertEntityWatchToModelWatch(entityWatch *entities.Watch) *models.Watch {
	return &models.Watch{
		WatchId:       entityWatch.WatchId.Hex(),
		Name:          entityWatch.Name,
		ExpansionId:   int32(entityWatch.ExpansionId),
		ExpansionName: entityWatch.ExpansionName,
		BlueprintId:   int32(entityWatch.BlueprintId),
		Condition:     convertEntityConditionToModelCondition(entityWatch.Condition),
		Foil:          entityWatch.Foil,
	}
}
