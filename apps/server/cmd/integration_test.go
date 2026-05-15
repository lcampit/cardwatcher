//go:build integration

package main_test

import (
	"context"
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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type ServerIntegrationTestSuite struct {
	suite.Suite
	ctx context.Context

	mongoTestContainer *mongodb.MongoDBContainer
	service            service.Service
	lis                *bufconn.Listener
	grpcServer         *grpc.Server

	cardtraderMock *cardtrader.MockCardtraderAdapter
	ntfyMock       *ntfy.MockNtfyAdapter
}

const (
	bufSize = 1024 * 1024

	mongoUser               = "user-test"
	mongoPassword           = "password-test"
	mongoDatabase           = "cardwatcher-test"
	mongoWatchCollectioName = "cardwatcher-collection-test"

	expansionId = 1
	blueprintId = 1
	condition   = apiv1.Condition_CONDITION_NEAR_MINT
	foil        = true
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

	suite.lis = bufconn.Listen(bufSize)

	suite.ntfyMock = ntfy.NewMockNtfyAdapter(suite.T())
	suite.cardtraderMock = setupCardtraderAdapterMock(suite.T())

	mongoHost, err := mongoContainer.Host(suite.ctx)
	if err != nil {
		suite.FailNowf("error getting mongo host: %s", err.Error())
	}
	port, err := mongoContainer.MappedPort(suite.ctx, "27017")

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
		CardtraderAdapter:    suite.cardtraderMock,
		MongoAdapter:         mongoAdapter,
		NtfyAdapter:          suite.ntfyMock,
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
		err = suite.grpcServer.Serve(suite.lis)
		if err != nil {
			logger.Error("error while listening", slog.Any("error", err))
		}
	}()
}

func (suite *ServerIntegrationTestSuite) TestCreateWatchAndGetNotify() {
	ctx := context.Background()
	conn, err := grpc.NewClient("passthrough:///bufconn",
		grpc.WithContextDialer(func(ctx context.Context, target string) (net.Conn, error) {
			return suite.lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		suite.FailNowf("failed to dial bufnet", "error %v", err)
	}
	defer conn.Close()
	client := apiv1.NewCardWatcherServiceClient(conn)

	watchRequest := apiv1.SaveWatchRequest{
		ExpansionId: expansionId,
		BlueprintId: blueprintId,
		Condition:   condition,
		Foil:        foil,
	}

	suite.cardtraderMock.On("GetBlueprintNameByExpansionID", mock.Anything, expansionId, blueprintId).Return()

	resp, err := client.SaveWatch(ctx, &watchRequest)
	if err != nil {
		suite.FailNowf("save watch request failed", "error %v", err)
	}

	suite.Assert().NotEmpty(resp.WatchId, "save watch returned an empty watch ID")
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
