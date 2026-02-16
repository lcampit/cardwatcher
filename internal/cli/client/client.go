package client

import (
	"fmt"

	api "github.com/lcampit/card-watcher-server/internal/api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	SaveWatch(expansionID, blueprintID uint64, condition api.Condition, foil bool) error
	GetWatches() error
	DeleteWatchByID(watchID string) error
	GetExpansions(expansionName, expansionCode string) error
	GetBlueprints(expansionID uint64, name string) error
	Close()
}

type client struct {
	connection    *grpc.ClientConn
	watcherClient api.CardWatcherClient
}

func NewClient(server string, port int) (Client, error) {
	serverAddress := fmt.Sprintf("%s:%d", server, port)
	fmt.Printf("creating client to server %s\n", serverAddress)
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	watcherClient := api.NewCardWatcherClient(conn)

	return &client{
		connection:    conn,
		watcherClient: watcherClient,
	}, nil
}

func (c *client) Close() {
	c.connection.Close()
}
