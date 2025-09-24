package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ServerIntegrationTestSuite struct {
	suite.Suite
	ctx context.Context

	mongoTestContainer *mongodb.MongoDBContainer
	ntfyTestContainer  testcontainers.Container
}

func (suite *ServerIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	mongoContainer, err := mongodb.Run(context.Background(), "mongo:latest")
	if err != nil {
		suite.FailNowf("error starting up mongo test container: %s", err.Error())
	}
	suite.mongoTestContainer = mongoContainer

	ntfyContainerRequest := testcontainers.ContainerRequest{
		Image:        "binwiederhier/ntfy",
		Cmd:          []string{"serve"},
		ExposedPorts: []string{"80"},
		WaitingFor:   wait.ForLog(".*Listening on.*").AsRegexp(),
	}

	ntfyContainer, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: ntfyContainerRequest,
		Started:          true,
	})
	if err != nil {
		suite.mongoTestContainer.Terminate(suite.ctx)
		suite.FailNowf("error starting up ntfy test container: %s", err.Error())
	}
	suite.ntfyTestContainer = ntfyContainer
}

func (suite *ServerIntegrationTestSuite) TestConnectToContainers() {
}

func (suite *ServerIntegrationTestSuite) TearDownSuite() {
	suite.mongoTestContainer.Terminate(suite.ctx)
	suite.ntfyTestContainer.Terminate(suite.ctx)
}

func TestServerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ServerIntegrationTestSuite))
}
