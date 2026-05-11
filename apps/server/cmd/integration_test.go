//go:build integration

package main_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"testing"

	"github.com/lcampit/cardwatcher/apps/server/internal/cardtrader"
	"github.com/lcampit/cardwatcher/apps/server/internal/handler"
	"github.com/lcampit/cardwatcher/apps/server/internal/logger"
	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
	"github.com/lcampit/cardwatcher/apps/server/internal/ntfy"
	"github.com/lcampit/cardwatcher/apps/server/internal/service"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"google.golang.org/grpc"
)

type ServerIntegrationTestSuite struct {
	suite.Suite
	ctx context.Context

	mongoTestContainer *mongodb.MongoDBContainer
	service            service.Service
	grpcServer         *grpc.Server
}

const (
	mongoUser               = "user-test"
	mongoPassword           = "password-test"
	mongoDatabase           = "cardwatcher-test"
	mongoWatchCollectioName = "cardwatcher-collection-test"
)

func (suite *ServerIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	mongoContainer, err := mongodb.Run(
		suite.ctx, "mongo:latest",
		mongodb.WithReplicaSet("rs0"),
		mongodb.WithUsername(mongoUser),
		mongodb.WithPassword(mongoPassword),
	)
	if err != nil {
		suite.FailNowf("error starting up mongo test container: %s", err.Error())
	}
	suite.mongoTestContainer = mongoContainer

	logger := logger.NewLogger("debug")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8000))
	if err != nil {
		suite.FailNowf("error listening: %s", err.Error())
	}

	ntfyMock := ntfy.NewMockNtfyAdapter(suite.T())
	cardtraderMock := setupCardtraderAdapterMock(suite.T())

	mongoHost, err := mongoContainer.Host(suite.ctx)
	if err != nil {
		suite.FailNowf("error getting mongo host: %s", err.Error())
	}
	port, err := mongoContainer.MappedPort(suite.ctx, "27017")
	logger.Info("mongo host", slog.String("host", mongoHost), slog.String("port", port.Port()))

	mongoAdapterConfig := mongo.MongoAdapterConfig{
		Logger:              logger,
		Host:                mongoHost,
		Port:                port.Port(),
		Username:            mongoUser,
		Password:            mongoPassword,
		Database:            mongoDatabase,
		AuthDatabase:        "admin",
		WatchCollectionName: mongoWatchCollectioName,
		UseReplicaSet:       true,
	}
	mongoAdapter, err := mongo.NewMongoAdapter(mongoAdapterConfig)
	if err != nil {
		suite.FailNowf("error connecting to mongo test container", err.Error())
	}

	serviceConfig := service.ServiceConfig{
		Logger:               logger,
		CardtraderAdapter:    cardtraderMock,
		MongoAdapter:         mongoAdapter,
		NtfyAdapter:          ntfyMock,
		NotificationSchedule: "10 * * * *",
		UpdateMapsSchedule:   "10 * * * *",
	}
	suite.service = service.NewService(suite.ctx, serviceConfig)

	handlerConfig := handler.HandlerConfig{
		Logger:  logger,
		Service: suite.service,
	}
	handler := handler.NewHandler(handlerConfig)
	suite.grpcServer = grpc.NewServer()
	apiv1.RegisterCardWatcherServiceServer(suite.grpcServer, handler)
	go func() {
		err = suite.grpcServer.Serve(lis)
		if err != nil {
			logger.Error("error while listening", slog.Any("error", err))
		}
	}()
}

func (suite *ServerIntegrationTestSuite) TestConnectToContainers() {
}

func (suite *ServerIntegrationTestSuite) TearDownSuite() {
	suite.grpcServer.GracefulStop()
	suite.service.Close()

	_ = suite.mongoTestContainer.Terminate(suite.ctx)
}

func TestServerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ServerIntegrationTestSuite))
}

func setupCardtraderAdapterMock(t *testing.T) *cardtrader.MockCardtraderAdapter {
	cardtraderMock := cardtrader.NewMockCardtraderAdapter(t)
	cardtraderMock.On("GetGames", mock.Anything).Return(nil, nil)
	return cardtraderMock
}
