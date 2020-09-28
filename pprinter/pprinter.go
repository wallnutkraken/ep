// Package pprinter is the pretty printer package that handles printing content like episode lists
// to the console
package pprinter

import (
	"fmt"
	"time"

	"github.com/wallnutkraken/ep/poddata/subscription"
)

// PrintEpisodeList prints the list of episodes provided to stdout
// with a prettified index (index + 1) prior to the title for easier
// selection.
func PrintEpisodeList(ep []subscription.Episode) {
	fmt.Println()
	for index, episode := range ep {
		fmt.Printf("\t[%d] [%s]\n", index+1, episode.Title)
	}
	fmt.Println()
}

// PrintPodcastList prints a list of podcasts given, along with their tags
func PrintPodcastList(subs []subscription.Subscription) {
	fmt.Println()
	for _, podcast := range subs {
		fmt.Printf("\t[%s] %s (%d known episodes, updated %s ago)\n", podcast.Tag, podcast.Name, len(podcast.Episodes), time.Since(podcast.UpdatedAt).String())
	}
	fmt.Println()
}
