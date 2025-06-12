package mongo

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

var testHost, testPort string

func mustStartMongoContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	dbContainer, err := mongodb.Run(context.Background(), "mongo:latest")
	if err != nil {
		return nil, fmt.Errorf("error when starting test container: %w", err)
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, fmt.Errorf("error when getting test container host: %w", err)
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "27017/tcp")
	if err != nil {
		return dbContainer.Terminate, fmt.Errorf("error when getting test container port: %w", err)
	}

	testHost = dbHost
	testPort = dbPort.Port()

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartMongoContainer()
	if err != nil {
		log.Fatalf("could not start mongodb container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown mongodb container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New(testHost, testPort)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New(testHost, testPort)

	stats := srv.Health()

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}
