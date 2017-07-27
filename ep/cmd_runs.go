package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/wallnutkraken/ep"
	"sync"
	"strconv"
)

var commands []*command = []*command{
	&command{
		UsageLine: "help",
		Short:     "Shows this help screen",
		Long:      "",
	}, &command{
		UsageLine: "add",
		Short:     "Adds a podcast via feed URL",
		Long:      "Adds a podcast via feed URL, usage: ep add PodcastTag https://example.com/podcast/rss",
		Run:       add,
	}, &command{
		UsageLine: "list",
		Short:     "Lists all tags/podcast epsiodes",
		Long: `
Running "ep list" will list all added tags using the following format:
	[Tag]	[Podcast Name]

Running "ep list [tag]" will retrieve and list all episodes for that podcast.
The numbers in the boxes [15] are used as the [episode] argument for the "play" action.`,
		Run: list,
	}, &command{
		UsageLine: "update",
		Short:     "Updates podcast episode entries",
		Long: `
Running "ep update" will update ALL podcast feeds that are added to the latest episodes.

If the command is run with an extra [tag] argument (like "ep update [tag]")
then ep will update that specific podcast.`,
		Run: update,
	}, &command{
		UsageLine:"play",
		Short:"Plays the latest or specific podcast episode",
		Long:`
"play" can be used in two ways:

"ep play [tag]" will play the latest episode of the podcast corresponding to [tag]

"ep play [tag] [episode]" will play the specific episode for the podcast corresponding
to [tag], to view a list of episodes for a podcast use "ep list [tag]"`,
		Run:play,
	}, &command{
		UsageLine:"remove",
		Short:"Removes the selected podcast from ep's memory",
		Long:`
"remove" can be used to remove a podcast from ep's memory. This operation cannot be undone.
Podcasts may be added again after removal, however, with the same tag.

Usage: ep remove [tag]`,
		Run: remove,
	},
}

func help(args []string) {
	defer os.Exit(0)
	if len(args) != 0 {
		cmdText := strings.ToLower(args[0])
		for _, c := range commands {
			if c.UsageLine == cmdText {
				fmt.Printf("\t%s: %s\n\n", c.UsageLine, c.Short)
				fmt.Printf("%s\n", c.Long)
				return
			}
		}
		/* If we don't find such command */
		fmt.Println("No such command is supported")
	} else {
		/* Run help itself */
		fmt.Println("ep is a tool for easily podcast categorization and playback")
		fmt.Println()
		fmt.Println("\tUsage: ep [action] [arguments]")
		fmt.Println()

		err := writeTemplate(cmdTemplate, commands)
		if err != nil {
			fmt.Println("Error:", err.Error())
		}
	}
}

func add(args []string) {
	if len(args) < 2 {
		fmt.Println(`error: some required arguments are missing; see "ep help add"`)
	} else {
		if err := ep.AddPodcast(args[0], args[1]); err != nil {
			fmt.Println("error:", err.Error())
		} else {
			fmt.Println("Podcast with tag", args[0], "added successfully")
		}
	}
}

func list(args []string) {
	if len(args) == 0 {
		/* List all tags */
		podcasts, err := ep.ListAll()
		if err != nil {
			fmt.Println("error reading feeds:", err.Error())
		} else {
			err := writeTemplate(podcastTemplate, podcasts)
			if err != nil {
				fmt.Println("error:", err.Error())
			}
		}
	} else {
		/* Specific tag search */
		podcast, err := ep.GetPodcast(args[0])
		if err != nil {
			fmt.Println("Could not find podcast with tag:", args[0])
			return
		}
		err = writeTemplate(epissodesTemplate, podcast.EpisodicItems)
		if err != nil {
			fmt.Println("error:", err.Error())
		}
	}
}

func update(args []string) {
	/* Func to actually update and write if failed */
	updateP := func(p ep.Podcast, wg *sync.WaitGroup) {
		oldEpCount := len(p.EpisodicItems)
		upErr := p.UpdateEpisodes()
		if upErr != nil {
			fmt.Println("failed updating tag:", p.Tag, "error:", upErr.Error())
			return
		}
		/* Updated; time to save */
		wErr := p.Write()
		if wErr != nil {
			fmt.Println("failed writing updated tag:", p.Tag, "error:", wErr.Error())
			return
		}

		newItemCount := len(p.EpisodicItems) - oldEpCount

		fmt.Printf("Successfully updated podcast [%s] %s (%d new items)\n", p.Tag, p.Name,
			newItemCount)
		if wg != nil {
			wg.Add(-1)
		}
	}

	if len(args) == 0 {
		/* Update ALL */
		fmt.Println("Updating all feeds...\n")
		podcasts, err := ep.ListAll()
		if err != nil {
			fmt.Println("error:", err.Error())
		}
		/* Create waitgroup so that we don't end up exiting before the update finishes */
		wg := sync.WaitGroup{}
		wg.Add(len(podcasts))
		for _, podcast := range podcasts {
			go updateP(podcast, &wg)
		}

		wg.Wait()
	} else {
		fmt.Println("Updating podcast with tag:", args[0] + "\n")
		podcast, err := ep.GetPodcast(args[0])
		if err != nil {
			fmt.Println("No podcast with tag", args[0], "was found")
		} else {
			updateP(podcast, nil)
		}
	}
	fmt.Println("\nDone.")
}

func play(args []string) {
	switch len(args) {
	case 0:
		fmt.Println(`play requires at least one argument; please see "ep help play"`)
	case 1:
		/* Play latest for podcast */
		p, err := ep.GetPodcast(args[0])
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}
		playBack(p.EpisodicItems[len(p.EpisodicItems) - 1])
	default:
		/* Play specific episode */
		p, err := ep.GetPodcast(args[0])
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}
		if args[1] == "random" {
			fmt.Println("random is not yet supported")
			return
		}
		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("error: second argument was not a number")
			return
		}
		/* Remember to do -1 due to us adding +1 when we show episodes in list */
		playBack(p.EpisodicItems[num - 1])
	}
}

func playBack(episode ep.Episode) {
	err := ep.Stream(episode)
	if err != nil {
		fmt.Println("playback error:", err.Error())
		os.Exit(1)
	}
}

func remove(args []string) {
	if len(args) == 0 {
		fmt.Println("You need to provide a tag for the podcast you wish to remove")
		return
	}

	p, err := ep.GetPodcast(args[0])
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}

	err = p.Remove()
	if err != nil {
		fmt.Println("remove error:", err.Error())
		return
	}

	fmt.Printf("Successfully removed podcast [%s]\n", p.Tag)
}