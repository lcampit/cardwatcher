package service

import (
	api "github.com/lcampit/card-watcher-server/internal/api/v1"
	"github.com/lcampit/card-watcher-server/internal/server/entities"
)

func convertModelConditionToEntityCondition(modelCondition api.Condition) entities.WatchCondition {
	switch modelCondition {
	case api.Condition_CONDITION_UNSPECIFIED:
	case api.Condition_CONDITION_NEAR_MINT:
		return entities.WatchConditionNM
	case api.Condition_CONDITION_SLIGHTLY_PLAYED:
		return entities.WatchConditionSP
	case api.Condition_CONDITION_MODERATELY_PLAYED:
		return entities.WatchConditionMP
	case api.Condition_CONDITION_PLAYED:
		return entities.WatchConditionPL
	case api.Condition_CONDITION_POOR:
		return entities.WatchConditionPO
	default:
		return entities.WatchConditionNM
	}
	return entities.WatchConditionNM
}

func convertEntityConditionToModelCondition(entityCondition entities.WatchCondition) api.Condition {
	switch entityCondition {
	case entities.WatchConditionNM:
		return api.Condition_CONDITION_NEAR_MINT
	case entities.WatchConditionSP:
		return api.Condition_CONDITION_SLIGHTLY_PLAYED
	case entities.WatchConditionMP:
		return api.Condition_CONDITION_SLIGHTLY_PLAYED
	case entities.WatchConditionPL:
		return api.Condition_CONDITION_PLAYED
	case entities.WatchConditionPO:
		return api.Condition_CONDITION_POOR
	}
	return api.Condition_CONDITION_NEAR_MINT
}

func convertEntityWatchToModelWatch(entityWatch *entities.Watch) *api.Watch {
	return &api.Watch{
		WatchId:       entityWatch.WatchID.Hex(),
		Name:          entityWatch.Name,
		ExpansionId:   entityWatch.ExpansionID,
		ExpansionName: entityWatch.ExpansionName,
		BlueprintId:   entityWatch.BlueprintID,
		Condition:     convertEntityConditionToModelCondition(entityWatch.Condition),
		Foil:          entityWatch.Foil,
	}
}
