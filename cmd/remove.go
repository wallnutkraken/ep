package cmd

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wallnutkraken/ep/poddata/subscription"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "rm"},
	Short:   "Removes the specified podcasts from ep's memory",
	Long: `Removes the specified podcasts from ep's memory, as well as all episode data ep has downloaded.

Multiple tags can be provided.`,
	Run: removePodcast,
	Example: `ep remove cortex
ep remove cortex quomodo
ep delete cortex
ep rm cortex`,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func removePodcast(cmd *cobra.Command, args []string) {
	// Get podcasts by tags
	deleteList, err := data.Subscriptions().GetSubscriptionsByTags(args...)
	if err != nil {
		logrus.Fatalf("Failed getting a list of subscribed podcasts for the given tags: %s", err.Error())
	}
	if len(deleteList) != len(args) {
		// Didn'f find tags, find the ones it didn't find and let the user know through a warn
		badTags := findNonmatchingTags(deleteList, args)
		logrus.Warnf("The following tags were not found in podcast subscriptions: %s; skipped", strings.Join(badTags, ", "))
	}
	if len(deleteList) == 0 {
		fmt.Println("Nothing to remove.")
		return
	}

	// First remove all the episdoes for these podcasts. To do that,
	// put them all in one array
	episodesToDelete := []subscription.Episode{}
	for _, sub := range deleteList {
		episodesToDelete = append(episodesToDelete, sub.Episodes...)
	}
	// Call delete on all these episodes
	if err := data.Subscriptions().RemoveEpisodes(episodesToDelete); err != nil {
		logrus.WithError(err).Fatal("Failed removing episodes for the podcasts to be deleted")
	}

	// And then remove the subscriptions
	if err := data.Subscriptions().RemoveSubscriptions(deleteList); err != nil {
		logrus.WithError(err).Fatal("Failed removing the podcasts")
	}

	// Done, report the stats
	fmt.Printf("[%d] Podcasts and [%d] episodes removed", len(deleteList), len(episodesToDelete))
}
