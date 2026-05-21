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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
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

	expansionId uint64 = 1
	blueprintId uint64 = 1
	condition          = apiv1.Condition_CONDITION_NEAR_MINT
	foil               = true
)

var expansion = cardtrader.Expansion{
	ID:     expansionId,
	GameID: 1,
	Code:   "exptest",
	Name:   "Test Expansion",
}

var blueprint = cardtrader.Blueprint{
	ID:          blueprintId,
	Name:        "super strong card",
	GameID:      1,
	ExpansionID: expansionId,
}

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

func (suite *ServerIntegrationTestSuite) TestCreateWatch() {
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
	suite.cardtraderMock.On("GetBlueprints", mock.Anything, expansion.ID).Return(
		[]*cardtrader.Blueprint{
			&blueprint,
		}, nil,
	)

	request := apiv1.CreateWatchRequest{
		ExpansionNameOrCode: expansion.Name,
		CardName:            blueprint.Name,
		Condition:           condition,
		Foil:                foil,
		Language:            apiv1.Language_LANGUAGE_EN,
	}
	resp, err := client.CreateWatch(ctx, &request)
	if err != nil {
		suite.FailNowf("create watch request failed", "error %v", err)
	}

	suite.Assert().NotEmpty(resp.WatchId, "create watch returned an empty watch ID")

	watches, err := client.ListWatches(ctx, &emptypb.Empty{})
	suite.Assert().Nil(err, "list watches request failed: %v", err)

	found := false
	for _, watch := range watches.GetWatches() {
		if watch.WatchId == resp.GetWatchId() {
			found = true
			suite.Assert().Equal(request.Condition, watch.Condition)
			suite.Assert().Equal(request.Language, watch.Language)
			suite.Assert().Equal(request.Foil, watch.Foil)
			suite.Assert().Equal(blueprint.ID, watch.BlueprintId)
			suite.Assert().Equal(blueprint.ExpansionID, watch.ExpansionId)
			suite.Assert().Equal(expansion.Name, watch.ExpansionName)
		}
	}
	suite.Assert().True(found, "created watch not found in list watches response")
}

func (suite *ServerIntegrationTestSuite) TestCreateWatchDefaultValues() {
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
	suite.cardtraderMock.On("GetBlueprints", mock.Anything, expansion.ID).Return(
		[]*cardtrader.Blueprint{
			&blueprint,
		}, nil,
	)

	request := apiv1.CreateWatchRequest{
		ExpansionNameOrCode: expansion.Name,
		CardName:            blueprint.Name,
	}
	resp, err := client.CreateWatch(ctx, &request)
	if err != nil {
		suite.FailNowf("create watch request failed", "error %v", err)
	}

	suite.Assert().NotEmpty(resp.WatchId, "create watch returned an empty watch ID")

	watches, err := client.ListWatches(ctx, &emptypb.Empty{})
	suite.Assert().Nil(err, "list watches request failed: %v", err)

	found := false
	for _, watch := range watches.GetWatches() {
		if watch.WatchId == resp.GetWatchId() {
			found = true
			suite.Assert().Equal(apiv1.Condition_CONDITION_UNSPECIFIED, watch.Condition)
			suite.Assert().Equal(apiv1.Language_LANGUAGE_UNSPECIFIED, watch.Language)
			suite.Assert().Equal(false, watch.Foil)
			suite.Assert().Equal(blueprint.ID, watch.BlueprintId)
			suite.Assert().Equal(blueprint.ExpansionID, watch.ExpansionId)
			suite.Assert().Equal(expansion.Name, watch.ExpansionName)
		}
	}
	suite.Assert().True(found, "created watch not found in list watches response")
}

func (suite *ServerIntegrationTestSuite) TestCreateWatchInvalidRequest() {
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
	suite.cardtraderMock.On("GetBlueprints", mock.Anything, expansion.ID).Return(
		[]*cardtrader.Blueprint{
			&blueprint,
		}, nil,
	)

	request := apiv1.CreateWatchRequest{}
	_, err = client.CreateWatch(ctx, &request)
	suite.Assert().NotNil(err, "create watch did not fail on empty request")

	grpcErr, ok := status.FromError(err)
	suite.Assert().True(ok, "cretae watch returned a non-grpc error %v", err)
	suite.Assert().Equal(grpcErr.Code(), codes.InvalidArgument, "create watch returned another error code: %d", grpcErr.Code())
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
	cardtraderMock.On("GetExpansions", mock.Anything).Return([]*cardtrader.Expansion{&expansion}, nil)
	return cardtraderMock
}
