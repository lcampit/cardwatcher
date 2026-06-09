package cmd

import (
	"fmt"

	"github.com/lcampit/cardwatcher/apps/cli/internal/client"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"

	"github.com/spf13/cobra"
)

var (
	isFoil                                                   bool
	nearMint, slightlyPlayed, moderatelyPlayed, played, poor bool
	en, de, fr, it, jp, pt, es                               bool
)

var createWatchCmd = &cobra.Command{
	Use:     "watch",
	Aliases: []string{"w"},
	Short:   "create a watch",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.NewClient(WatcherServerAddress, WatcherServerPort)
		if err != nil {
			return fmt.Errorf("error when creating card watcher client: %w", err)
		}
		defer client.Close()

		condition := extractCondition()
		language := extractLanguage()

		return client.CreateWatch(args[1], args[0], condition, language, isFoil)
	},
}

func extractLanguage() apiv1.Language {
	if en {
		return apiv1.Language_LANGUAGE_EN
	}
	if de {
		return apiv1.Language_LANGUAGE_DE
	}
	if fr {
		return apiv1.Language_LANGUAGE_FR
	}
	if it {
		return apiv1.Language_LANGUAGE_IT
	}
	if jp {
		return apiv1.Language_LANGUAGE_JP
	}
	if pt {
		return apiv1.Language_LANGUAGE_PT
	}
	if es {
		return apiv1.Language_LANGUAGE_ES
	}
	return apiv1.Language_LANGUAGE_UNSPECIFIED
}

func extractCondition() apiv1.Condition {
	if nearMint {
		return apiv1.Condition_CONDITION_NEAR_MINT
	}
	if slightlyPlayed {
		return apiv1.Condition_CONDITION_SLIGHTLY_PLAYED
	}
	if moderatelyPlayed {
		return apiv1.Condition_CONDITION_MODERATELY_PLAYED
	}
	if played {
		return apiv1.Condition_CONDITION_PLAYED
	}
	if poor {
		return apiv1.Condition_CONDITION_POOR
	}
	return apiv1.Condition_CONDITION_UNSPECIFIED
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
	createWatchCmd.Flags().BoolVarP(&isFoil, "foil", "f", false, "whether to look for foil cards")
	// add condition flags
	createWatchCmd.Flags().BoolVar(&nearMint, "nm", false, "whether to look for near mint condition cards")
	createWatchCmd.Flags().BoolVar(&slightlyPlayed, "sp", false, "whether to look for slightly played condition cards")
	createWatchCmd.Flags().BoolVar(&moderatelyPlayed, "mp", false, "whether to look for moderately played condition cards")
	createWatchCmd.Flags().BoolVar(&played, "pl", false, "whether to look for played condition cards")
	createWatchCmd.Flags().BoolVar(&poor, "po", false, "whether to look for poor condition cards")

	// add language flags
	createWatchCmd.Flags().BoolVar(&en, "en", false, "whether to look for english cards")
	createWatchCmd.Flags().BoolVar(&de, "de", false, "whether to look for german cards")
	createWatchCmd.Flags().BoolVar(&fr, "fr", false, "whether to look for french cards")
	createWatchCmd.Flags().BoolVar(&it, "it", false, "whether to look for italian cards")
	createWatchCmd.Flags().BoolVar(&jp, "jp", false, "whether to look for japanese cards")
	createWatchCmd.Flags().BoolVar(&pt, "pt", false, "whether to look for portuguese cards")
	createWatchCmd.Flags().BoolVar(&es, "es", false, "whether to look for spanish cards")

	createCmd.AddCommand(createWatchCmd)
	getCmd.AddCommand(getWatchesCmd)
	deleteCmd.AddCommand(deleteWatchCmd)
}
