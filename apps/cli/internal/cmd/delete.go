/*
Copyright © 2025 Leonardo Campitelli
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Command to delete objects from watcher server",
	Long:  `Use the delete command to delete objects like watches from the card watcher server`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
