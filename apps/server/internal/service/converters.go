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

func convertEntityLanguageToModelLanguage(entityLanguage mongo.WatchLanguage) apiv1.Language {
	switch entityLanguage {
	case mongo.WatchLanguageAny:
		return apiv1.Language_LANGUAGE_UNSPECIFIED
	case mongo.WatchLanguageEn:
		return apiv1.Language_LANGUAGE_EN
	case mongo.WatchLanguageDe:
		return apiv1.Language_LANGUAGE_DE
	case mongo.WatchLanguageFr:
		return apiv1.Language_LANGUAGE_FR
	case mongo.WatchLanguageIt:
		return apiv1.Language_LANGUAGE_IT
	case mongo.WatchLanguageJp:
		return apiv1.Language_LANGUAGE_JP
	case mongo.WatchLanguageEs:
		return apiv1.Language_LANGUAGE_ES
	}
	return apiv1.Language_LANGUAGE_UNSPECIFIED
}

func convertModelLanguageToEntityLanguage(modelLanguage apiv1.Language) mongo.WatchLanguage {
	switch modelLanguage {
	case apiv1.Language_LANGUAGE_UNSPECIFIED:
		return mongo.WatchLanguageAny
	case apiv1.Language_LANGUAGE_EN:
		return mongo.WatchLanguageEn
	case apiv1.Language_LANGUAGE_DE:
		return mongo.WatchLanguageDe
	case apiv1.Language_LANGUAGE_FR:
		return mongo.WatchLanguageFr
	case apiv1.Language_LANGUAGE_IT:
		return mongo.WatchLanguageIt
	case apiv1.Language_LANGUAGE_JP:
		return mongo.WatchLanguageJp
	case apiv1.Language_LANGUAGE_ES:
		return mongo.WatchLanguageEs
	}
	return mongo.WatchLanguageAny
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
		Language:      convertEntityLanguageToModelLanguage(entityWatch.Language),
	}
}
