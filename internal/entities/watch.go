// Package entities contains definition for
// entities used in the application
package entities

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WatchCondition string

const (
	WatchConditionNM WatchCondition = "Near Mint"
	WatchConditionSP WatchCondition = "Slightly Played"
	WatchConditionMP WatchCondition = "Moderately Played"
	WatchConditionPL WatchCondition = "Played"
	WatchConditionPO WatchCondition = "Poor"
)

type Watch struct {
	WatchID       bson.ObjectID  `bson:"_id"`
	Name          string         `bson:"name"`
	ExpansionID   int            `bson:"expansionId"`
	ExpansionName string         `bson:"expansionName"`
	BlueprintID   int            `bson:"blueprintId"`
	Condition     WatchCondition `bson:"condition"`
	Foil          bool           `bson:"foil"`
}
