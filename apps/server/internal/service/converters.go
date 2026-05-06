package service

import (
	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

func convertModelConditionToEntityCondition(modelCondition apiv1.Condition) mongo.WatchCondition {
	switch modelCondition {
	case apiv1.Condition_CONDITION_UNSPECIFIED:
		return mongo.WatchConditionAny
	case apiv1.Condition_CONDITION_NEAR_MINT:
		return mongo.WatchConditionNM
	case apiv1.Condition_CONDITION_SLIGHTLY_PLAYED:
		return mongo.WatchConditionSP
	case apiv1.Condition_CONDITION_MODERATELY_PLAYED:
		return mongo.WatchConditionMP
	case apiv1.Condition_CONDITION_PLAYED:
		return mongo.WatchConditionPL
	case apiv1.Condition_CONDITION_POOR:
		return mongo.WatchConditionPO
	default:
		return mongo.WatchConditionAny
	}
}

func convertEntityConditionToModelCondition(entityCondition mongo.WatchCondition) apiv1.Condition {
	switch entityCondition {
	case mongo.WatchConditionAny:
		return apiv1.Condition_CONDITION_UNSPECIFIED
	case mongo.WatchConditionNM:
		return apiv1.Condition_CONDITION_NEAR_MINT
	case mongo.WatchConditionSP:
		return apiv1.Condition_CONDITION_SLIGHTLY_PLAYED
	case mongo.WatchConditionMP:
		return apiv1.Condition_CONDITION_SLIGHTLY_PLAYED
	case mongo.WatchConditionPL:
		return apiv1.Condition_CONDITION_PLAYED
	case mongo.WatchConditionPO:
		return apiv1.Condition_CONDITION_POOR
	}
	return apiv1.Condition_CONDITION_UNSPECIFIED
}

func convertEntityWatchToModelWatch(entityWatch *mongo.Watch) *apiv1.Watch {
	return &apiv1.Watch{
		WatchId:       entityWatch.WatchID.Hex(),
		Name:          entityWatch.Name,
		ExpansionId:   entityWatch.ExpansionID,
		ExpansionName: entityWatch.ExpansionName,
		BlueprintId:   entityWatch.BlueprintID,
		Condition:     convertEntityConditionToModelCondition(entityWatch.Condition),
		Foil:          entityWatch.Foil,
	}
}
