/*
Copyright © 2025 Leonardo Campitelli
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Command to create objects",
	Long:    `Use the create command to save objects like watches in the card watcher server`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
