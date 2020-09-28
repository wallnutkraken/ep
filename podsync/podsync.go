// Package podsync handles the sync of podcast episodes, as well as the retrieval of podcast data from
// a given RSS feed
package podsync

import (
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
	if err := updateEpisodes(&sub, feed); err != nil {
		return sub, errors.WithMessage(err, "Failed getting podcast episodes")
	}

	return sub, nil
}

// updateEpisodes does the gruntwork for UpdateEpisodes, it works with the rss feed given instead of parsing it itself
func updateEpisodes(sub *subscription.Subscription, feed *gofeed.Feed) error {
	episodes := []subscription.Episode{}
	// Reverse loop, to start from earliest items
	for index := len(feed.Items) - 1; index > 0; index-- {
		item := feed.Items[index]
		episode := subscription.Episode{}

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

	// Assign the newly synced episodes to the subscription and return
	// TODO: check existing episodes and don't replace any (to avoid db mishaps)
	sub.Episodes = episodes
	return nil
}

// UpdateEpisodes takes a pointer to a Subscription object and populates its Episodes element
// with the latest episodes for this podcast
func UpdateEpisodes(sub *subscription.Subscription) error {
	// Get the RSS Data
	rss := gofeed.NewParser()
	feed, err := rss.ParseURL(sub.RSSURL)
	if err != nil {
		return errors.WithMessagef(err, "Failed reading RSS feed at [%s]", sub.RSSURL)
	}
	return updateEpisodes(sub, feed)
}
