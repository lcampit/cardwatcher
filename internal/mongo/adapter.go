package mongo

import (
	"card-watcher/internal/entities"
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAdapter interface {
	SaveWatch(ctx context.Context, watch *entities.Watch) (string, error)
	GetWatches(ctx context.Context) ([]*entities.Watch, error)
	GetWatchByWatchId(ctx context.Context, watchId string) (*entities.Watch, error)
	DeleteWatchById(ctx context.Context, watchId string) error
	Health() map[string]string
}

type mongoAdapter struct {
	client   *mongo.Client
	database string
}

func NewMongoAdapter(
	host, port, database string,
) MongoAdapter {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		log.Fatal(err)
	}
	return &mongoAdapter{
		client:   client,
		database: database,
	}
}

func (s *mongoAdapter) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
