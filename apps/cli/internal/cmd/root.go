/*
Copyright © 2025 Leonardo Campitelli
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "card-watcher-cli",
	Short: "CLI to interact with the card watcher server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		v := viper.New()
		v.SetEnvPrefix("CARDWATCHER")
		v.AutomaticEnv()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	WatcherServerAddress string
	WatcherServerPort    int
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&WatcherServerAddress, "server", "s", "", "address of the watcher server to connect to")
	rootCmd.PersistentFlags().IntVarP(&WatcherServerPort, "port", "p", 3000, "port of the watcher server to connect to")
}
