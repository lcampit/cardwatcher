/*
Copyright © 2025 Leonardo Campitelli
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/lcampit/cardwatcher/apps/cli/internal/client"

	"github.com/spf13/cobra"
)

var blueprintName string

// expansionsCmd represents the expansions command
var getBlueprintsCmd = &cobra.Command{
	Use:     "blueprints",
	Aliases: []string{"b"},
	Short:   "get information about blueprints optionally using their name",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(args)
		client, err := client.NewClient(WatcherServerAddress, WatcherServerPort)
		if err != nil {
			return fmt.Errorf("error when creating card watcher client: %w", err)
		}

		expansionID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("please provide a valid integer expansion ID as obtained using the `get expansions` command. Received %s", args[0])
		}

		return client.GetBlueprints(expansionID, blueprintName)
	},
}

func init() {
	getBlueprintsCmd.Flags().StringVarP(&blueprintName, "name", "n", "", "blueprint name to get info about")
	getCmd.AddCommand(getBlueprintsCmd)
}
