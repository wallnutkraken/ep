package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/wallnutkraken/ep/podsync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wallnutkraken/ep/poddata/subscription"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"update"},
	Short:   "Sync stored podcast episodes with their RSS feeds",
	Long: `Sync retrieves the latest list of podcast episodes. If no argument is provided, ep will sync all saved podcasts.
Otherwise, it will sync the podcasts for which the tags are provided.`,
	Run: syncPodcasts,
	Example: `ep sync
ep sync cortex
ep sync cortex quomodo`,
	Args: cobra.ArbitraryArgs,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func syncPodcasts(cmd *cobra.Command, args []string) {
	toSync := []subscription.Subscription{}
	if len(args) == 0 {
		// Sync everything
		syncList, err := data.Subscriptions().GetSubscriptions()
		if err != nil {
			logrus.Fatalf("Failed getting a list of subscribed podcasts: %s", err.Error())
		}
		toSync = syncList
	} else {
		// Get podcasts by tags
		syncList, err := data.Subscriptions().GetSubscriptionsByTags(args...)
		if err != nil {
			logrus.Fatalf("Failed getting a list of subscribed podcasts for the given tags: %s", err.Error())
		}
		if len(syncList) != len(args) {
			// Didn'f find tags, find the ones it didn't find and let the user know through a warn
			badTags := findNonmatchingTags(syncList, args)
			logrus.Warnf("The following tags were not found in podcast subscriptions: %s; skipped", strings.Join(badTags, ", "))
		}
		toSync = syncList
	}

	// Go through every one of them to sync and sync asynchronously
	wg := &sync.WaitGroup{}
	wg.Add(len(toSync))
	for _, subToSync := range toSync {
		go func(sub *subscription.Subscription, waiter *sync.WaitGroup) {
			// Get an updated episode list
			var err error
			newEpisodes := []subscription.Episode{}
			if newEpisodes, err = podsync.GetNewEpisodes(*sub); err != nil {
				logrus.WithError(err).Errorf("Failed updating podcast [%s]", sub.Tag)
				waiter.Done()
				return
			}
			if len(newEpisodes) != 0 {
				// Add the newest episodes to the database. First, get the slice that's new.
				if err := data.Subscriptions().AddEpisodes(*sub, newEpisodes); err != nil {
					logrus.WithError(err).Errorf("Failed saving newly retrieved episodes for [%s]", sub.Name)
				} else {
					fmt.Printf("Updated podcast [%s/%s] with %d new episodes", sub.Tag, sub.Name, len(newEpisodes))
				}
			}
			// And add the episodes to the subscription in memory
			sub.Episodes = append(sub.Episodes, newEpisodes...)

			waiter.Done()
		}(&subToSync, wg)
	}

	// Wait for all syncs to complete
	wg.Wait()

	// Save all to the database
	for _, sub := range toSync {
		if err := data.Subscriptions().UpdateSubscription(sub); err != nil {
			logrus.WithError(err).Error("Failed updating podcast post sync, new episodes will not be saved")
		}
	}
}

// findNonmatchingTags looks through all the subs elements given and checks which elements of tags
// aren't represented, then it returns said list of unrepresented tags
func findNonmatchingTags(subs []subscription.Subscription, tags []string) []string {
	unmatched := []string{}
	// Go through every tag
	for _, tag := range tags {
		found := false
		// Search if that tag can be found in the list of subs
		for _, sub := range subs {
			if sub.Tag == tag {
				// Found it, break the inner loop
				found = true
				break
			}
		}
		if !found {
			// Not found in the subs, add it to the unmatched list
			unmatched = append(unmatched, tag)
		}
	}

	return unmatched
}
