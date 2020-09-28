// Package podsync handles the sync of podcast episodes, as well as the retrieval of podcast data from
// a given RSS feed
package podsync

import (
	"sort"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/wallnutkraken/ep/poddata/subscription"
)

// GetPodcastFromRSS returns a Subscription object created from the given RSS URL.
// It will also add the provided tag to the Subscription object
func GetPodcastFromRSS(rssURL, tag string) (subscription.Subscription, error) {
	rss := gofeed.NewParser()
	feed, err := rss.ParseURL(rssURL)
	if err != nil {
		return subscription.Subscription{}, errors.WithMessagef(err, "Failed reading RSS feed at [%s]", rssURL)
	}
	sub := subscription.Subscription{
		Name:     feed.Title,
		RSSURL:   rssURL,
		Tag:      tag,
		Episodes: []subscription.Episode{},
	}
	// Get all the current episodes for this podcast
	episodes, err := getNewEpisodes(sub, feed)
	if err != nil {
		return sub, errors.WithMessage(err, "Failed getting podcast episodes")
	}
	sub.Episodes = episodes

	return sub, nil
}

// getNewEpisodes does the gruntwork for GetNewEpisodes, it works with the rss feed given instead of parsing it itself
func getNewEpisodes(sub subscription.Subscription, feed *gofeed.Feed) ([]subscription.Episode, error) {
	episodes := []subscription.Episode{}

	// Reverse loop, to start from earliest items
	for index := len(feed.Items) - 1; index > 0; index-- {
		item := feed.Items[index]
		episode := subscription.Episode{}
		episode.Title = item.Title
		episode.PublishedAt = *item.PublishedParsed

		// Find the actual URL to the audio
		for _, encl := range item.Enclosures {
			if strings.HasPrefix(encl.Type, "audio") {
				episode.URL = encl.URL
				// We got the one we want; break
				break
			}
		}
		episodes = append(episodes, episode)
	}
	// Sort the episodes by publication date
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].PublishedAt.Unix() < episodes[j].PublishedAt.Unix()
	})

	// Figure out which episodes are new
	if len(sub.Episodes) == 0 {
		sub.Episodes = episodes
	} else {
		// There are new episodes, add the new ones, don't change the existing ones
		// First, get the index of the last shared element
		lastSharedIndex := indexOfLastSharedEpisode(sub.Episodes, episodes)
		// Combine existing episodes with the episodes after the last shared one
		sub.Episodes = episodes[lastSharedIndex+1:]
	}
	return episodes, nil
}

// GetNewEpisodes retrieves new episodes for the given podcast subscription
func GetNewEpisodes(sub subscription.Subscription) ([]subscription.Episode, error) {
	// Get the RSS Data
	rss := gofeed.NewParser()
	feed, err := rss.ParseURL(sub.RSSURL)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed reading RSS feed at [%s]", sub.RSSURL)
	}
	return getNewEpisodes(sub, feed)
}

// indexOfLastShared episode returns the index in the newEpisodes array of the last element that is shared between it and
// the newEpisodes array. Returns -1 if there is no match
func indexOfLastSharedEpisode(originalEpisodes, newEpisodes []subscription.Episode) int {
	lastOriginalURL := originalEpisodes[len(originalEpisodes)-1].URL
	for index := len(newEpisodes) - 1; index > 0; index-- {
		if newEpisodes[index].URL == lastOriginalURL {
			// Found the match, return the index
			return index
		}
	}
	return -1
}
