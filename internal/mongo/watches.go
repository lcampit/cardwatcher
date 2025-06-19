package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	gomongo "go.mongodb.org/mongo-driver/v2/mongo"
)

type WatcherCondition string

const (
	NEAR_MINT         WatcherCondition = "NEAR_MINT"
	SLIGHTLY_PLAYED   WatcherCondition = "SLIGHTLY_PLAYED"
	MODERATELY_PLAYED WatcherCondition = "MODERATELY_PLAYED"
	PLAYED            WatcherCondition = "PLAYED"
	POOR              WatcherCondition = "POOR"
)

const WATCH_COLLECTION string = "watches"

type Watch struct {
	Name        string           `bson:"name"`
	WatchId     string           `bson:"watchId"`
	UserId      string           `bson:"userId"`
	ExpansionId string           `bson:"expansionId"`
	BlueprintId string           `bson:"blueprintId"`
	Condition   WatcherCondition `bson:"condition"`
	Foil        bool             `bson:"foil"`
}

func (a *mongoAdapter) SaveWatch(ctx context.Context, watch *Watch) error {
	_, err := a.client.Database(a.database).Collection(WATCH_COLLECTION).
		InsertOne(ctx, watch)
	if err != nil {
		return fmt.Errorf("error inserting watch with id %s: %w", watch.WatchId, err)
	}
	return nil
}

func (a *mongoAdapter) DeleteWatch(ctx context.Context, watchId string) error {
	_, err := a.client.Database(a.database).Collection(WATCH_COLLECTION).
		DeleteOne(ctx, bson.M{"watchId": watchId})
	if err != nil {
		return fmt.Errorf("error deleting watch with id %s: %w", watchId, err)
	}
	return nil
}

func (a *mongoAdapter) DeleteWatchesByUserId(ctx context.Context, userId string) error {
	_, err := a.client.Database(a.database).Collection(WATCH_COLLECTION).
		DeleteMany(ctx, bson.M{"userId": userId})
	if err != nil {
		return fmt.Errorf("error deleting watch by user id %s: %w", userId, err)
	}
	return nil
}

func (a *mongoAdapter) GetWatchByWatchId(ctx context.Context, watchId string) (*Watch, error) {
	var watch Watch
	filter := bson.M{"watchId": watchId}
	result := a.client.Database(a.database).Collection(WATCH_COLLECTION).
		FindOne(ctx, filter)

	if result.Err() == gomongo.ErrNoDocuments {
		return nil, nil
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("error getting watch by id %s: %w", watchId, result.Err())
	}

	err := result.Decode(&watch)
	if err != nil {
		return nil, fmt.Errorf("error decoding watch with id %s from database: %w", watchId, err)
	}

	return &watch, nil
}
