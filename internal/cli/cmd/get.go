/*
Copyright © 2025 Leonardo Campitelli
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Command to get information from watcher server",
	Long:  `Use the get command to obtain information about expansions, blueprints or watches`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
