package client

import (
	"context"
	"fmt"
	"time"

	"github.com/lcampit/cardwatcher/apps/cli/internal/printer"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

func (c *client) GetExpansions(gameName, expansionName, expansionCode string) error {
	request := apiv1.ListExpansionsRequest{
		Name: expansionName,
		Code: expansionCode,
		Game: gameName,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := c.watcherClient.ListExpansions(ctx, &request)
	if err != nil {
		return fmt.Errorf("error when calling list expansions: %w", err)
	}

	printer.PrintExpansionTable(response.Expansions)

	return nil
}
