package cmd

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wallnutkraken/ep/player"
	"github.com/wallnutkraken/ep/poddata/subscription"
	"github.com/wallnutkraken/ep/stream"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Plays the specified podcast or episode",
	Long: `This command will start playback of the selected podcast or episode.

It can be used with a tag to play the latest episode of the relevant podcast,
or with a tag and podcast number to play a specific episode.

The podcast number can be found via the "ep list [tag]" command in the following format:` +
		"\n\n\t[Number] [Title]",
	Run: playPodcast,
	Example: `ep play cortex
ep play cortex 50
ep play quomodo 1`,
	Args: cobra.RangeArgs(1, 2),
}

func init() {
	rootCmd.AddCommand(playCmd)
}

func playPodcast(cmd *cobra.Command, args []string) {
	args = toLowerArray(args)
	// Get the podcast
	sub, err := data.Subscriptions().GetSubscriptionByTag(args[0])
	if err != nil {
		if err == subscription.ErrSubNotFound {
			logrus.Errorf("No podcast found with tag [%s]", args[0])
			return
		}
		logrus.WithError(err).Fatalf("Failed finding podcast with tag [%s]", args[0])
	}
	var episode subscription.Episode
	if len(args) == 1 {
		// Play the latest episode
		episode = sub.Episodes[len(sub.Episodes)-1]
	} else {
		// Get the specified episode
		podNumber, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			logrus.WithError(err).Fatalf("Podcast number argument %s is not a valid number", args[1])
		}
		if podNumber < 1 || podNumber >= int64(len(sub.Episodes)) {
			logrus.Fatal("Podcast episode number is out of range")
		}
		episode = sub.Episodes[podNumber-1]
	}

	// Start streaming the episode
	streamHandler := stream.FromURL(episode.URL)
	streamReader, err := streamHandler.GetStream()
	if err != nil {
		logrus.WithError(err).Fatalf("Failed reading stream at [%s]", episode.URL)
	}
	// Get the decoder function, now that the stream is open
	decoder, err := streamHandler.GetDecoder()
	if err != nil {
		// This can only be the ErrMediaTypeUnsupported error, so just print it out
		logrus.WithError(err).Fatal("Cannot decode file")
	}

	// Create the player
	play := player.New(streamReader, decoder)
	// Begin playback
	if err := play.Start(); err != nil {
		logrus.WithError(err).Fatal("Could not play audio file")
	}

	fmt.Printf("Now playing: %s (%s)", episode.Title, episode.PublishedAt.String())

	// And call the Wait blocking call, as the audio is running on a separate goroutine
	play.Wait()
}
