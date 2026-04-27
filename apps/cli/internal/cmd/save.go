/*
Copyright © 2025 Leonardo Campitelli
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Command to save objects to the server",
	Long:  `Use the save command to save objects like watches to the card watcher server`,
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
