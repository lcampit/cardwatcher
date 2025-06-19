package entities

import (
	"card-watcher/internal/models"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WatchCondition string

const (
	WATCH_CONDITION_NEAR_MINT         WatchCondition = "NEAR_MINT"
	WATCH_CONDITION_SLIGHTLY_PLAYED   WatchCondition = "SLIGHTLY_PLAYED"
	WATCH_CONDITION_MODERATELY_PLAYED WatchCondition = "MODERATELY_PLAYED"
	WATCH_CONDITION_PLAYED            WatchCondition = "PLAYED"
	WATCH_CONDITION_POOR              WatchCondition = "POOR"
)

type Watch struct {
	WatchId     bson.ObjectID  `bson:"_id"`
	UserId      string         `bson:"userId"`
	ExpansionId string         `bson:"expansionId"`
	BlueprintId string         `bson:"blueprintId"`
	Condition   WatchCondition `bson:"condition"`
	Foil        bool           `bson:"foil"`
}

func ConvertModelConditionToWatchCondition(modelCondition models.Condition) WatchCondition {
	switch modelCondition {
	case models.Condition_CONDITION_UNSPECIFIED:
	case models.Condition_CONDITION_NEAR_MINT:
		return WATCH_CONDITION_NEAR_MINT
	case models.Condition_CONDITION_SLIGHTLY_PLAYED:
		return WATCH_CONDITION_SLIGHTLY_PLAYED
	case models.Condition_CONDITION_MODERATELY_PLAYED:
		return WATCH_CONDITION_MODERATELY_PLAYED
	case models.Condition_CONDITION_PLAYED:
		return WATCH_CONDITION_PLAYED
	case models.Condition_CONDITION_POOR:
		return WATCH_CONDITION_POOR
	default:
		return WATCH_CONDITION_NEAR_MINT
	}
	log.Debug().Msgf("received unknown model condition %s, resorting to default NEAR MINT", modelCondition)
	return WATCH_CONDITION_NEAR_MINT
}
