package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/lcampit/card-watcher-server/internal/server/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (a *mongoAdapter) SaveWatch(ctx context.Context, watch *entities.Watch) (string, error) {
	watch.WatchID = bson.NewObjectID()
	_, err := a.client.Database(a.database).Collection(a.watchCollection).
		InsertOne(ctx, watch)
	if err != nil {
		return "", fmt.Errorf("inserting watch with id %s: %w", watch.WatchID, err)
	}
	return watch.WatchID.Hex(), nil
}

func (a *mongoAdapter) DeleteWatchByID(ctx context.Context, watchID string) error {
	convertedID, err := bson.ObjectIDFromHex(watchID)
	if err != nil && errors.Is(err, bson.ErrInvalidHex) {
		return errors.New("invalid watch ID provided")
	}
	if err != nil {
		return fmt.Errorf("converting object id %s in delete watch by id: %w", watchID, err)
	}
	_, err = a.client.Database(a.database).Collection(a.watchCollection).
		DeleteOne(ctx, bson.M{"_id": convertedID})
	if err != nil {
		return fmt.Errorf("deleting watch with id %s: %w", watchID, err)
	}
	return nil
}

func (a *mongoAdapter) GetWatches(ctx context.Context) ([]*entities.Watch, error) {
	var watches []*entities.Watch
	cursor, err := a.client.Database(a.database).Collection(a.watchCollection).
		Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("finding all watches: %w", err)
	}

	err = cursor.All(ctx, &watches)
	if err != nil {
		return nil, fmt.Errorf("decoding all watches: %w", err)
	}

	return watches, nil
}

func (a *mongoAdapter) GetWatchByWatchID(ctx context.Context, watchID string) (*entities.Watch, error) {
	convertedID, err := bson.ObjectIDFromHex(watchID)
	if err != nil {
		return nil, fmt.Errorf("converting object id %s in get watch by id: %w", watchID, err)
	}

	var watch entities.Watch
	filter := bson.M{"_id": convertedID}
	result := a.client.Database(a.database).Collection(a.watchCollection).
		FindOne(ctx, filter)

	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("getting watch by id %s: %w", watchID, result.Err())
	}

	err = result.Decode(&watch)
	if err != nil {
		return nil, fmt.Errorf("decoding watch with id %s from database: %w", watchID, err)
	}

	return &watch, nil
}
