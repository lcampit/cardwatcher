package service

import (
	api "github.com/lcampit/card-watcher-server/internal/api/v1"
	"github.com/lcampit/card-watcher-server/internal/server/mongo"
)

func convertModelConditionToEntityCondition(modelCondition api.Condition) mongo.WatchCondition {
	switch modelCondition {
	case api.Condition_CONDITION_UNSPECIFIED:
	case api.Condition_CONDITION_NEAR_MINT:
		return mongo.WatchConditionNM
	case api.Condition_CONDITION_SLIGHTLY_PLAYED:
		return mongo.WatchConditionSP
	case api.Condition_CONDITION_MODERATELY_PLAYED:
		return mongo.WatchConditionMP
	case api.Condition_CONDITION_PLAYED:
		return mongo.WatchConditionPL
	case api.Condition_CONDITION_POOR:
		return mongo.WatchConditionPO
	default:
		return mongo.WatchConditionNM
	}
	return mongo.WatchConditionNM
}

func convertEntityConditionToModelCondition(entityCondition mongo.WatchCondition) api.Condition {
	switch entityCondition {
	case mongo.WatchConditionNM:
		return api.Condition_CONDITION_NEAR_MINT
	case mongo.WatchConditionSP:
		return api.Condition_CONDITION_SLIGHTLY_PLAYED
	case mongo.WatchConditionMP:
		return api.Condition_CONDITION_SLIGHTLY_PLAYED
	case mongo.WatchConditionPL:
		return api.Condition_CONDITION_PLAYED
	case mongo.WatchConditionPO:
		return api.Condition_CONDITION_POOR
	}
	return api.Condition_CONDITION_NEAR_MINT
}

func convertEntityWatchToModelWatch(entityWatch *mongo.Watch) *api.Watch {
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
