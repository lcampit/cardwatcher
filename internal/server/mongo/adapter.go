// Package mongo exposes database related operations
// using a single adapter
package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lcampit/card-watcher-server/internal/server/entities"

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
	Username            string
	Password            string
	Database            string
	WatchCollectionName string
}

func NewMongoAdapter(config MongoAdapterConfig) (MongoAdapter, error) {
	config.Logger.Debug("creating mongo adapter with config", slog.Any("config", config))
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	credentials := options.Credential{
		AuthSource:    config.Database,
		Username:      config.Username,
		Password:      config.Password,
		AuthMechanism: "SCRAM-SHA-256",
	}

	host := config.Host
	if host == "" {
		host = "localhost"
	}

	port := config.Port
	if port == "" {
		port = "27017"
	}

	clientOpts := options.Client().
		SetHosts([]string{fmt.Sprintf("%s:%s", host, port)}).
		SetAuth(credentials).
		SetRetryReads(true).
		SetRetryWrites(true).
		SetServerSelectionTimeout(5 * time.Second).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(50).
		SetMinPoolSize(5)

	if host == "localhost" || host == "127.0.0.1" {
		clientOpts.SetDirect(true)
	} else {
		clientOpts.SetReplicaSet("rs0")
	}

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("creating mongo client: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("pinging mongo instance: %w", err)
	}

	return &mongoAdapter{
		logger:          config.Logger,
		client:          client,
		database:        config.Database,
		watchCollection: config.WatchCollectionName,
	}, nil
}

func (a *mongoAdapter) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.client.Ping(ctx, nil)
}
