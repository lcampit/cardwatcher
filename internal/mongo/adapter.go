// Package mongo exposes database related operations
// using a single adapter
package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"card-watcher/internal/entities"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoAdapter interface {
	SaveWatch(ctx context.Context, watch *entities.Watch) (string, error)
	GetWatches(ctx context.Context) ([]*entities.Watch, error)
	GetWatchByWatchID(ctx context.Context, watchID string) (*entities.Watch, error)
	DeleteWatchByID(ctx context.Context, watchID string) error

	Health() error
}

type mongoAdapter struct {
	logger          *slog.Logger
	client          *mongo.Client
	database        string
	watchCollection string
}

type MongoAdapterConfig struct {
	Logger              *slog.Logger
	Host                string
	Port                string
	Database            string
	WatchCollectionName string
	ConnectionRetries   int
}

func NewMongoAdapter(config MongoAdapterConfig) (MongoAdapter, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	client, err := mongo.Connect(options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)))
	if err != nil {
		return nil, fmt.Errorf("creating mongo client: %w", err)
	}
	for retry := range config.ConnectionRetries {
		err = client.Ping(ctx, nil)
		if err != nil {
			config.Logger.Error("connecting to mongo client",
				slog.Any("error", err),
				slog.Int("retryCount", retry))
			time.Sleep(1 * time.Second)
		}
	}
	// all tries have failed, return error
	if err != nil {
		return nil, fmt.Errorf("connecting to mongo client: %w", err)
	}
	return &mongoAdapter{
		logger:          config.Logger,
		client:          client,
		database:        config.Database,
		watchCollection: config.WatchCollectionName,
	}, nil
}

func (a *mongoAdapter) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return a.client.Ping(ctx, nil)
}
