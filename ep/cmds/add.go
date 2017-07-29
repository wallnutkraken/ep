package cmds

import (
	"fmt"
	"github.com/wallnutkraken/ep"
)

func Add(args []string) {
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