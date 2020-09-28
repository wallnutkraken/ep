package cmd

import (
	"fmt"

	"github.com/wallnutkraken/ep/pprinter"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wallnutkraken/ep/poddata/subscription"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all podcasts or episodes",
	Long: `Lists all podcasts or episodes stored in ep's memory.
list can be called on its own to list all subscribed podcasts, or with one tag to list all the episodes ep knows of for that podcast.`,
	Run: listEntries,
	Example: `ep list
ep list cortex
ep list quomodo`,
	Args: cobra.MaximumNArgs(1),
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listEntries(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		listPodcasts()
	} else {
		listEpisodes(args[0])
	}
}

func listPodcasts() {
	subs, err := data.Subscriptions().GetSubscriptions()
	if err != nil {
		logrus.WithError(err).Fatal("Failed getting all podcast subscriptions")
	}
	pprinter.PrintPodcastList(subs)
}

func listEpisodes(tag string) {
	sub, err := data.Subscriptions().GetSubscriptionByTag(tag)
	if err == subscription.ErrSubNotFound {
		fmt.Printf("Tag [%s] was not found\n", tag)
		return
	} else if err != nil {
		logrus.WithError(err).Fatalf("Failed getting subscription for podcast [%s]", tag)
	}
	pprinter.PrintEpisodeList(sub.Episodes)
}
