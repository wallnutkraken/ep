package cmds

import (
	"fmt"
	"os"
	"strconv"
	"unicode"

	"github.com/wallnutkraken/ep"
	"github.com/zetamatta/go-getch"
)

func Play(args []string) {
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
		playBack(p.EpisodicItems[len(p.EpisodicItems)-1])
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
		playBack(p.EpisodicItems[num-1])
	}
}

func playBack(episode ep.Episode) {
	onDone, control, err := ep.StartStreaming(episode)
	if err != nil {
		fmt.Println("playback error:", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Name: %s\n", episode.Title)
	fmt.Println("Beginning playback. Press 'P' to pause/resume. Press Q to exit.")
	fmt.Println("Use +/- to control the volume.")

	/* Start controls listener goroutine */
	go controls(control, onDone)

	/* Block until onDone has something sent through it or is closed. */
	/* Expecting close if user-closed. */
	<-onDone
}

func controls(c ep.Controller, onDone chan bool) {
	exit := false
	for !exit {
		key := unicode.ToLower(getch.Rune())
		switch key {
		case 'q':
			exit = true
			close(onDone)
		case 'p':
			c.TogglePaused()
		case '+':
			c.VolumeUp()
		case '-':
			c.VolumeDown()
		}
	}
}
