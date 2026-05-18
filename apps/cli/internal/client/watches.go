package client

import (
	"context"
	"fmt"
	"time"

	"github.com/lcampit/cardwatcher/apps/cli/internal/printer"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

func (c *client) SaveWatch(expansionID, blueprintID uint64, condition apiv1.Condition, language apiv1.Language, foil bool) error {
	request := apiv1.SaveWatchRequest{
		ExpansionId: expansionID,
		BlueprintId: blueprintID,
		Condition:   condition,
		Foil:        foil,
		Language:    language,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := c.watcherClient.SaveWatch(ctx, &request)
	if err != nil {
		return fmt.Errorf("error when calling save watch: %w", err)
	}

	fmt.Printf("Watch saved with ID: %s\n", response.WatchId)
	return nil
}

func (c *client) CreateWatch(expansionNameOrCode, cardName string, condition apiv1.Condition, language apiv1.Language, foil bool) error {
	request := apiv1.CreateWatchRequest{
		ExpansionNameOrCode: expansionNameOrCode,
		CardName:            cardName,
		Condition:           condition,
		Foil:                foil,
		Language:            language,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := c.watcherClient.CreateWatch(ctx, &request)
	if err != nil {
		return fmt.Errorf("error when calling create watch: %w", err)
	}

	fmt.Printf("Watch saved with ID: %s\n", response.WatchId)
	return nil
}

func (c *client) GetWatches() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := c.watcherClient.ListWatches(ctx, nil)
	if err != nil {
		return fmt.Errorf("error when calling list watches: %w", err)
	}

	printer.PrintWatchesTable(response.Watches)
	return nil
}

func (c *client) DeleteWatchByID(watchID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.watcherClient.DeleteWatchByID(ctx, &apiv1.DeleteWatchByIDRequest{
		WatchId: watchID,
	})
	if err != nil {
		return fmt.Errorf("error when calling delete watch: %w", err)
	}

	fmt.Printf("deleted watch with ID %s\n", watchID)
	return nil
}
