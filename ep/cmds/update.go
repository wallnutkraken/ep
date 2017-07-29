package cmds

import (
	"github.com/wallnutkraken/ep"
	"sync"
	"fmt"
)

func Update(args []string) {
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
