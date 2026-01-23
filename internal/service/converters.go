package service

import (
	"card-watcher/internal/entities"
	"card-watcher/internal/models"
)

func convertModelConditionToEntityCondition(modelCondition models.Condition) entities.WatchCondition {
	switch modelCondition {
	case models.Condition_CONDITION_UNSPECIFIED:
	case models.Condition_CONDITION_NEAR_MINT:
		return entities.WatchConditionNM
	case models.Condition_CONDITION_SLIGHTLY_PLAYED:
		return entities.WatchConditionSP
	case models.Condition_CONDITION_MODERATELY_PLAYED:
		return entities.WatchConditionMP
	case models.Condition_CONDITION_PLAYED:
		return entities.WatchConditionPL
	case models.Condition_CONDITION_POOR:
		return entities.WatchConditionPO
	default:
		return entities.WatchConditionNM
	}
	return entities.WatchConditionNM
}

func convertEntityConditionToModelCondition(entityCondition entities.WatchCondition) models.Condition {
	switch entityCondition {
	case entities.WatchConditionNM:
		return models.Condition_CONDITION_NEAR_MINT
	case entities.WatchConditionSP:
		return models.Condition_CONDITION_SLIGHTLY_PLAYED
	case entities.WatchConditionMP:
		return models.Condition_CONDITION_SLIGHTLY_PLAYED
	case entities.WatchConditionPL:
		return models.Condition_CONDITION_PLAYED
	case entities.WatchConditionPO:
		return models.Condition_CONDITION_POOR
	}
	return models.Condition_CONDITION_NEAR_MINT
}

func convertEntityWatchToModelWatch(entityWatch *entities.Watch) *models.Watch {
	return &models.Watch{
		WatchId:       entityWatch.WatchID.Hex(),
		Name:          entityWatch.Name,
		ExpansionId:   int32(entityWatch.ExpansionID),
		ExpansionName: entityWatch.ExpansionName,
		BlueprintId:   int32(entityWatch.BlueprintID),
		Condition:     convertEntityConditionToModelCondition(entityWatch.Condition),
		Foil:          entityWatch.Foil,
	}
}
