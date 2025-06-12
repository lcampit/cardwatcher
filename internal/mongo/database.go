package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAdapter interface {
	Health() map[string]string
}

type mongoAdapter struct {
	db *mongo.Client
}

func New(
	host, port string,
) MongoAdapter {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		log.Fatal(err)
	}
	return &mongoAdapter{
		db: client,
	}
}

func (s *mongoAdapter) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
