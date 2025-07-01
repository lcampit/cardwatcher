package entities

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WatchCondition string

const (
	WATCH_CONDITION_NEAR_MINT         WatchCondition = "Near Mint"
	WATCH_CONDITION_SLIGHTLY_PLAYED   WatchCondition = "Slightly Played"
	WATCH_CONDITION_MODERATELY_PLAYED WatchCondition = "Moderately Played"
	WATCH_CONDITION_PLAYED            WatchCondition = "Played"
	WATCH_CONDITION_POOR              WatchCondition = "Poor"
)

type Watch struct {
	WatchId       bson.ObjectID  `bson:"_id"`
	Name          string         `bson:"name"`
	ExpansionId   int            `bson:"expansionId"`
	ExpansionName string         `bson:"expansionName"`
	BlueprintId   int            `bson:"blueprintId"`
	Condition     WatchCondition `bson:"condition"`
	Foil          bool           `bson:"foil"`
}
