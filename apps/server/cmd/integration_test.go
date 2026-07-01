//go:build integration

package main_test

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/lcampit/cardwatcher/apps/server/internal/app"
	"github.com/lcampit/cardwatcher/apps/server/internal/cardtrader"
	"github.com/lcampit/cardwatcher/apps/server/internal/handler"
	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
	"github.com/lcampit/cardwatcher/apps/server/internal/ntfy"
	"github.com/lcampit/cardwatcher/apps/server/internal/observability/logger"
	"github.com/lcampit/cardwatcher/apps/server/internal/service"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

type ServerIntegrationTestSuite struct {
	suite.Suite
	ctx context.Context

	mongoTestContainer *mongodb.MongoDBContainer
	mongoAdapter       mongo.MongoAdapter
	lis                *bufconn.Listener

	app            *app.App
	service        service.Service
	cardtraderMock *cardtrader.MockCardtraderAdapter
	ntfyMock       *ntfy.MockNtfyAdapter

	conn   *grpc.ClientConn
	client apiv1.CardWatcherServiceClient
}

const (
	bufSize = 1024 * 1024

	mongoUser               = "user-test"
	mongoPassword           = "password-test"
	mongoDatabase           = "cardwatcher-test"
	mongoWatchCollectioName = "cardwatcher-collection-test"

	expansionID uint64 = 1
	blueprintID uint64 = 1
	condition          = apiv1.Condition_CONDITION_NEAR_MINT
	foil               = true
)

var expansion = cardtrader.Expansion{
	ID:     expansionID,
	GameID: 1,
	Code:   "exptest",
	Name:   "Test Expansion",
}

var blueprint = cardtrader.Blueprint{
	ID:          blueprintID,
	Name:        "super strong card",
	GameID:      1,
	ExpansionID: expansionID,
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

	logger := logger.New(logger.LoggerConfig{
		Level:     "debug",
		Service:   "cardwatcher-integration-test",
		AddSource: false,
	})

	mongoHost, err := mongoContainer.Host(suite.ctx)
	if err != nil {
		suite.FailNowf("error getting mongo host: %s", err.Error())
	}
	port, err := mongoContainer.MappedPort(suite.ctx, "27017")
	if err != nil {
		suite.FailNowf("error getting mongo exposed port: %s", err.Error())
	}

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
	suite.mongoAdapter, err = mongo.NewMongoAdapter(mongoAdapterConfig)
	if err != nil {
		suite.FailNowf("error connecting to mongo test container", err.Error())
	}
}

// SetupTest creates the whole service -> handler -> app chain so that
// every test is performed in isolation between each other
func (suite *ServerIntegrationTestSuite) SetupTest() {
	logger := logger.New(logger.LoggerConfig{
		Level:     "debug",
		Service:   "cardwatcher-integration-test",
		AddSource: false,
	})
	suite.lis = bufconn.Listen(bufSize)

	suite.ntfyMock = ntfy.NewMockNtfyAdapter(suite.T())
	suite.cardtraderMock = cardtrader.NewMockCardtraderAdapter(suite.T())

	// Mock cardtrader adapter calls used for service startup here. Each test
	// will then set up its own mocks
	suite.cardtraderMock.On("GetGames", mock.Anything).Return(nil, nil).Once()
	suite.cardtraderMock.On("GetExpansions", mock.Anything).Return([]*cardtrader.Expansion{&expansion}, nil).Once()

	serviceConfig := service.ServiceConfig{
		Logger:               logger,
		CardtraderAdapter:    suite.cardtraderMock,
		MongoAdapter:         suite.mongoAdapter,
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

	var err error
	suite.app, err = app.NewApp(
		handler,
		logger,
		suite.lis,
		0,
		true,
		time.Minute,
	)
	if err != nil {
		suite.FailNowf("error creating app", err.Error())
	}

	go func() {
		err = suite.app.Run()
		if err != nil {
			logger.Error("error while starting app", slog.Any("error", err))
		}
	}()

	suite.conn, err = grpc.NewClient("passthrough:///bufconn",
		grpc.WithContextDialer(func(ctx context.Context, target string) (net.Conn, error) {
			return suite.lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		suite.FailNowf("failed to dial bufnet", "error %v", err)
	}
	suite.client = apiv1.NewCardWatcherServiceClient(suite.conn)
}

func (suite *ServerIntegrationTestSuite) TestCreateWatchFromCachedMap() {
	ctx := context.Background()
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
	resp, err := suite.client.CreateWatch(ctx, &request)
	if err != nil {
		suite.FailNowf("create watch request failed", "error %v", err)
	}

	suite.Assert().NotEmpty(resp.WatchId, "create watch returned an empty watch ID")

	watches, err := suite.client.ListWatches(ctx, &emptypb.Empty{})
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

func (suite *ServerIntegrationTestSuite) TestCreateWatchFromCachedMapDefaultValues() {
	ctx := context.Background()
	suite.cardtraderMock.On("GetBlueprints", mock.Anything, expansion.ID).Return(
		[]*cardtrader.Blueprint{
			&blueprint,
		}, nil,
	)

	request := apiv1.CreateWatchRequest{
		ExpansionNameOrCode: expansion.Name,
		CardName:            blueprint.Name,
	}
	resp, err := suite.client.CreateWatch(ctx, &request)
	if err != nil {
		suite.FailNowf("create watch request failed", "error %v", err)
	}

	suite.Assert().NotEmpty(resp.WatchId, "create watch returned an empty watch ID")

	watches, err := suite.client.ListWatches(ctx, &emptypb.Empty{})
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

func (suite *ServerIntegrationTestSuite) TestCreateWatchWithoutExpansionFails() {
	ctx := context.Background()

	request := apiv1.CreateWatchRequest{
		CardName:  blueprint.Name,
		Condition: condition,
		Foil:      foil,
		Language:  apiv1.Language_LANGUAGE_EN,
	}
	_, err := suite.client.CreateWatch(ctx, &request)
	suite.Assert().NotNil(err, "create watch did not fail on request without expansion")

	grpcErr, ok := status.FromError(err)
	suite.Assert().True(ok, "create watch returned a non-grpc error %v", err)
	suite.Assert().Equal(grpcErr.Code(), codes.InvalidArgument, "create watch returned another error code: %d", grpcErr.Code())
}

func (suite *ServerIntegrationTestSuite) TestCreateWatchWithoutCardnameFails() {
	ctx := context.Background()

	request := apiv1.CreateWatchRequest{
		ExpansionNameOrCode: expansion.Name,
		Condition:           condition,
		Foil:                foil,
		Language:            apiv1.Language_LANGUAGE_EN,
	}
	_, err := suite.client.CreateWatch(ctx, &request)
	suite.Assert().NotNil(err, "create watch did not fail on request without cardname")

	grpcErr, ok := status.FromError(err)
	suite.Assert().True(ok, "create watch returned a non-grpc error %v", err)
	suite.Assert().Equal(grpcErr.Code(), codes.InvalidArgument, "create watch returned another error code: %d", grpcErr.Code())
}

func (suite *ServerIntegrationTestSuite) TestCreateWatchWithNonExistingExpansionReturnsNotFound() {
	ctx := context.Background()

	suite.cardtraderMock.On("GetExpansions", mock.Anything).
		Return(nil, nil)

	request := apiv1.CreateWatchRequest{
		ExpansionNameOrCode: "non-existing-expansion-name",
		CardName:            "non-existing-card-name",
		Condition:           condition,
		Foil:                foil,
		Language:            apiv1.Language_LANGUAGE_EN,
	}
	_, err := suite.client.CreateWatch(ctx, &request)
	suite.Assert().NotNil(err)

	grpcErr, ok := status.FromError(err)
	suite.Assert().True(ok, "create watch returned a non-grpc error %v", err)
	suite.Assert().Equal(grpcErr.Code(), codes.NotFound, "create watch returned another error code: %d", grpcErr.Code())
}

func (suite *ServerIntegrationTestSuite) TestCreateWatchReturnsInternalOnCardtraderExpansionError() {
	ctx := context.Background()

	suite.cardtraderMock.On("GetExpansions", mock.Anything).
		Return(nil, errors.New("internal error"))

	// We need to use non existing names here to make sure we do not
	// hit the inner maps
	request := apiv1.CreateWatchRequest{
		ExpansionNameOrCode: "non-existing-expansion-name",
		CardName:            "non-existing-card-name",
		Condition:           condition,
		Foil:                foil,
		Language:            apiv1.Language_LANGUAGE_EN,
	}
	_, err := suite.client.CreateWatch(ctx, &request)
	suite.Assert().NotNil(err)

	grpcErr, ok := status.FromError(err)
	suite.Assert().True(ok, "create watch returned a non-grpc error %v", err)
	suite.Assert().Equal(grpcErr.Code(), codes.Internal, "create watch returned another error code: %d", grpcErr.Code())
}

func (suite *ServerIntegrationTestSuite) TestCreateWatchWithNonExistingBlueprintReturnsNotFound() {
	ctx := context.Background()

	suite.cardtraderMock.On("GetBlueprints", mock.Anything, expansion.ID).
		Return(nil, nil)
	request := apiv1.CreateWatchRequest{
		ExpansionNameOrCode: expansion.Code,
		CardName:            "non-existing-card-name",
		Condition:           condition,
		Foil:                foil,
		Language:            apiv1.Language_LANGUAGE_EN,
	}
	_, err := suite.client.CreateWatch(ctx, &request)
	suite.Assert().NotNil(err)

	grpcErr, ok := status.FromError(err)
	suite.Assert().True(ok, "create watch returned a non-grpc error %v", err)
	suite.Assert().Equal(grpcErr.Code(), codes.NotFound, "create watch returned another error code: %d", grpcErr.Code())
}

func (suite *ServerIntegrationTestSuite) TestCreateWatchReturnsInternalOnCardtraderBlueprintError() {
	ctx := context.Background()

	suite.cardtraderMock.On("GetBlueprints", mock.Anything, expansion.ID).
		Return(nil, errors.New("internal error"))
	request := apiv1.CreateWatchRequest{
		ExpansionNameOrCode: expansion.Code,
		CardName:            blueprint.Name,
		Condition:           condition,
		Foil:                foil,
		Language:            apiv1.Language_LANGUAGE_EN,
	}
	_, err := suite.client.CreateWatch(ctx, &request)
	suite.Assert().NotNil(err)

	grpcErr, ok := status.FromError(err)
	suite.Assert().True(ok, "create watch returned a non-grpc error %v", err)
	suite.Assert().Equal(grpcErr.Code(), codes.Internal, "create watch returned another error code: %d", grpcErr.Code())
}

// TeardownTest stops the app created previously so that the new
// setupTest can recreate it from scratch
func (suite *ServerIntegrationTestSuite) TeardownTest() {
	suite.mongoAdapter.ClearCollection(suite.ctx)
	suite.conn.Close()
	suite.app.Shutdown(1 * time.Second)
	suite.service.Close()
}

func (suite *ServerIntegrationTestSuite) TearDownSuite() {
	suite.app.Shutdown(1 * time.Second)
	suite.service.Close()

	_ = suite.mongoTestContainer.Terminate(suite.ctx)
}

func TestServerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ServerIntegrationTestSuite))
}
