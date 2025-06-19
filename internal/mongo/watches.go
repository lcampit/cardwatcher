package mongo

import (
	"card-watcher/internal/entities"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	gomongo "go.mongodb.org/mongo-driver/v2/mongo"
)

const WATCH_COLLECTION string = "watches"

func (a *mongoAdapter) SaveWatch(ctx context.Context, watch *entities.Watch) (string, error) {
	watch.WatchId = bson.NewObjectID()
	_, err := a.client.Database(a.database).Collection(WATCH_COLLECTION).
		InsertOne(ctx, watch)
	if err != nil {
		return "", fmt.Errorf("error inserting watch with id %s: %w", watch.WatchId, err)
	}
	return watch.WatchId.Hex(), nil
}

func (a *mongoAdapter) DeleteWatchById(ctx context.Context, watchId string) error {
	convertedId, err := bson.ObjectIDFromHex(watchId)
	if err != nil {
		return fmt.Errorf("error converting object id %s in delete watch by id: %w", watchId, err)
	}
	_, err = a.client.Database(a.database).Collection(WATCH_COLLECTION).
		DeleteOne(ctx, bson.M{"_id": convertedId})
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

func (a *mongoAdapter) GetWatchByWatchId(ctx context.Context, watchId string) (*entities.Watch, error) {
	convertedId, err := bson.ObjectIDFromHex(watchId)
	if err != nil {
		return nil, fmt.Errorf("error converting object id %s in get watch by id: %w", watchId, err)
	}

	var watch entities.Watch
	filter := bson.M{"_id": convertedId}
	result := a.client.Database(a.database).Collection(WATCH_COLLECTION).
		FindOne(ctx, filter)

	if result.Err() == gomongo.ErrNoDocuments {
		return nil, nil
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("error getting watch by id %s: %w", watchId, result.Err())
	}

	err = result.Decode(&watch)
	if err != nil {
		return nil, fmt.Errorf("error decoding watch with id %s from database: %w", watchId, err)
	}

	return &watch, nil
}
