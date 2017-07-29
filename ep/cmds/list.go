package cmds

import (
	"github.com/wallnutkraken/ep"
	"fmt"
	"github.com/wallnutkraken/ep/ep/cmds/temp"
)

func list(args []string) {
	if len(args) == 0 {
		/* List all tags */
		podcasts, err := ep.ListAll()
		if err != nil {
			fmt.Println("error reading feeds:", err.Error())
		} else {
			err := temp.WriteTemplate(temp.PodcastTemplate, podcasts)
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
		err = temp.WriteTemplate(temp.EpissodesTemplate, podcast.EpisodicItems)
		if err != nil {
			fmt.Println("error:", err.Error())
		}
	}
}