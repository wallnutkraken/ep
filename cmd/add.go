package cmd

import (
	"fmt"
	"strings"

	"github.com/wallnutkraken/ep/pprinter"

	"github.com/wallnutkraken/ep/poddata/subscription"

	"github.com/sirupsen/logrus"
	"github.com/wallnutkraken/ep/podsync"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new podcast feed.",
	Long: `Adds a new podcast to the ep subscriptions.

Usage: ep add [tag/name for podcast] [RSS URL]
The tag/name will be used for every other command involving this podcast, such as play.
The tag does not allow whitespace`,
	Run:  addPodcast,
	Args: cobra.ExactArgs(2),
	Example: `ep add cortex https://www.relay.fm/cortex/feed
ep add quomodo http://quomododicitur.com/feed/podcast/quomododicitur`,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func toLowerArray(v []string) []string {
	for index, elem := range v {
		v[index] = strings.ToLower(elem)
	}
	return v
}

func addPodcast(cmd *cobra.Command, args []string) {
	tag := strings.ToLower(args[0])
	url := args[1]
	// Check if this tag isn't already in use
	_, err := data.Subscriptions().GetSubscriptionByTag(tag)
	if err != subscription.ErrSubNotFound {
		if err != nil {
			// Some other error, not the not found one. Log error and assume it wasn't found.
			logrus.Errorf("Failed checking database for tag [%s]: %s", tag, err.Error())
		} else {
			// The subscription already exists
			fmt.Printf("A subscription with the tag [%s] already exists\n", tag)
			return
		}
	}

	// Get the podcast at that URL
	sub, err := podsync.GetPodcastFromRSS(url, tag)
	if err != nil {
		logrus.Fatalf("Error getting podcast feed: %s", err.Error())
	}
	// Print out the podcast's episodes here
	pprinter.PrintEpisodeList(sub.Episodes)

	if err := data.Subscriptions().NewSubscription(&sub); err != nil {
		logrus.Fatalf("Failed saving subscription for tag [%s] and url [%s]: %s", tag, url, err.Error())
	}
	fmt.Printf("%d\n", sub.ID)
	fmt.Printf("Added podcast [%s] using the tag [%s] to subscribed podcasts", sub.Name, tag)
}
