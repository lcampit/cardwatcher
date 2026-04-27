/*
Copyright © 2025 Leonardo Campitelli
*/
package cmd

import (
	"fmt"

	"github.com/lcampit/cardwatcher/apps/cli/internal/client"

	"github.com/spf13/cobra"
)

var (
	gameName      string
	expansionName string
	expansionCode string
)

// getExpansionsCmd represents the expansions command
var getExpansionsCmd = &cobra.Command{
	Use:     "expansions",
	Aliases: []string{"e"},
	Short:   "get information about expansions optionally using expansion name or code",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.NewClient(WatcherServerAddress, WatcherServerPort)
		if err != nil {
			return fmt.Errorf("error when creating card watcher client: %w", err)
		}
		defer client.Close()

		return client.GetExpansions(gameName, expansionName, expansionCode)
	},
}

func init() {
	getExpansionsCmd.Flags().StringVarP(&gameName, "game", "g", "", "game name expansions belong to")
	getExpansionsCmd.Flags().StringVarP(&expansionName, "name", "n", "", "expansion name to get info about")
	getExpansionsCmd.Flags().StringVarP(&expansionCode, "code", "c", "", "expansion code to get info about")
	getCmd.AddCommand(getExpansionsCmd)
}
