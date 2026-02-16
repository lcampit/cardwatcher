package mongo

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

var testHost, testPort, testDatabase string

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
	testDatabase = "testDatabase"

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
	logger := slog.New(slog.DiscardHandler)
	config := MongoAdapterConfig{
		logger,
		testHost,
		testPort,
		testDatabase,
		"watch-test",
		5,
	}
	srv, _ := NewMongoAdapter(config)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	config := MongoAdapterConfig{
		logger,
		testHost,
		testPort,
		testDatabase,
		"watch-test",
		5,
	}
	srv, _ := NewMongoAdapter(config)

	err := srv.Health()
	if err != nil {
		t.Fatalf("error in healthcheck: %v", err)
	}
}
