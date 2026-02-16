package client

import (
	"context"
	"fmt"
	"time"

	api "github.com/lcampit/card-watcher-server/internal/api/v1"
	"github.com/lcampit/card-watcher-server/internal/cli/printer"
)

func (c *client) GetExpansions(expansionName, expansionCode string) error {
	request := api.ListExpansionsRequest{
		Name: expansionName,
		Code: expansionCode,
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
