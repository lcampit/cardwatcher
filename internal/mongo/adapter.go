// Package mongo exposes database related operations
// using a single adapter
package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"card-watcher/internal/entities"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoAdapter interface {
	SaveWatch(ctx context.Context, watch *entities.Watch) (string, error)
	GetWatches(ctx context.Context) ([]*entities.Watch, error)
	GetWatchByWatchID(ctx context.Context, watchID string) (*entities.Watch, error)
	DeleteWatchByID(ctx context.Context, watchID string) error

	Health() map[string]string
	Close()
}

type mongoAdapter struct {
	logger          *slog.Logger
	client          *mongo.Client
	database        string
	watchCollection string
	cancelContext   context.CancelFunc
}

type MongoAdapterConfig struct {
	Logger              *slog.Logger
	Host                string
	Port                string
	Database            string
	WatchCollectionName string
}

func NewMongoAdapter(config MongoAdapterConfig) MongoAdapter {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)))
	if err != nil {
		config.Logger.Error("error connecting to mongo instance", slog.Any("error", err))
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		config.Logger.Error("error reaching mongo instance", slog.Any("error", err))
	}
	return &mongoAdapter{
		logger:          config.Logger,
		client:          client,
		database:        config.Database,
		watchCollection: config.WatchCollectionName,
		cancelContext:   cancelFunc,
	}
}

func (a *mongoAdapter) Close() {
	a.cancelContext()
}

func (a *mongoAdapter) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := a.client.Ping(ctx, nil)
	if err != nil {
		a.logger.Error("error reaching mongo instance in healthcheck", slog.Any("error", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
