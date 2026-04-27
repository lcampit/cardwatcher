package client

import (
	"context"
	"fmt"
	"time"

	"github.com/lcampit/cardwatcher/apps/cli/internal/printer"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

func (c *client) GetBlueprints(expansionID uint64, name string) error {
	request := apiv1.ListBlueprintsRequest{
		ExpansionId: expansionID,
		Name:        name,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := c.watcherClient.ListBlueprints(ctx, &request)
	if err != nil {
		return fmt.Errorf("error when calling list blueprints: %w", err)
	}

	printer.PrintBlueprintsTable(response.Blueprints)

	return nil
}
