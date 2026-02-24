package cmd

import (
	"fmt"
	"strconv"

	api "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
	"github.com/lcampit/cardwatcher/internal/cli/client"

	"github.com/spf13/cobra"
)

var (
	isFoil                                                   bool
	nearMint, slightlyPlayed, moderatelyPlayed, played, poor bool
)

var saveWatchCmd = &cobra.Command{
	Use:     "watch",
	Aliases: []string{"w"},
	Short:   "save a watch to the card watcher server",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.NewClient(WatcherServerAddress, WatcherServerPort)
		if err != nil {
			return fmt.Errorf("error when creating card watcher client: %w", err)
		}
		defer client.Close()

		expansionID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("please provide a valid integer expansion ID as obtained using the `get expansions` command. Received %s", args[0])
		}

		blueprintID, err := strconv.ParseUint(args[1], 10, 32)
		if err != nil {
			return fmt.Errorf("please provide a valid integer expansion ID as obtained using the `get expansions` command. Received %s", args[0])
		}

		condition := api.Condition_CONDITION_NEAR_MINT
		if slightlyPlayed {
			condition = api.Condition_CONDITION_SLIGHTLY_PLAYED
		}
		if moderatelyPlayed {
			condition = api.Condition_CONDITION_MODERATELY_PLAYED
		}
		if played {
			condition = api.Condition_CONDITION_PLAYED
		}
		if poor {
			condition = api.Condition_CONDITION_POOR
		}

		return client.SaveWatch(expansionID, blueprintID, condition, isFoil)
	},
}

var getWatchesCmd = &cobra.Command{
	Use:     "watches",
	Aliases: []string{"w"},
	Short:   "get all watches currently saved on the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.NewClient(WatcherServerAddress, WatcherServerPort)
		if err != nil {
			return fmt.Errorf("error when creating card watcher client: %w", err)
		}
		defer client.Close()

		return client.GetWatches()
	},
}

var deleteWatchCmd = &cobra.Command{
	Use:     "watch",
	Aliases: []string{"w"},
	Short:   "delete a watch currently saved on the server",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.NewClient(WatcherServerAddress, WatcherServerPort)
		if err != nil {
			return fmt.Errorf("error when creating card watcher client: %w", err)
		}
		defer client.Close()

		return client.DeleteWatchByID(args[0])
	},
}

func init() {
	saveWatchCmd.Flags().BoolVarP(&isFoil, "foil", "f", false, "whether to look for foil cards or not")
	saveWatchCmd.Flags().BoolVar(&nearMint, "nm", false, "whether to look for near mint condition cards")
	saveWatchCmd.Flags().BoolVar(&slightlyPlayed, "sp", false, "whether to look for slightly played condition cards")
	saveWatchCmd.Flags().BoolVar(&moderatelyPlayed, "mp", false, "whether to look for moderately played condition cards")
	saveWatchCmd.Flags().BoolVar(&played, "pl", false, "whether to look for played condition cards")
	saveWatchCmd.Flags().BoolVar(&poor, "po", false, "whether to look for poor condition cards")
	saveCmd.AddCommand(saveWatchCmd)
	getCmd.AddCommand(getWatchesCmd)
	deleteCmd.AddCommand(deleteWatchCmd)
}
