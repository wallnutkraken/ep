package cmds

import (
	"fmt"
	"github.com/wallnutkraken/ep"
)

func Remove(args []string) {
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