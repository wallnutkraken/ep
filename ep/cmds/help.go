package cmds

import (
	"os"
	"strings"
	"fmt"
	"github.com/wallnutkraken/ep/ep/cmds/temp"
)

func Help(args []string) {
	defer os.Exit(0)
	if len(args) != 0 {
		cmdText := strings.ToLower(args[0])
		for _, c := range Commands {
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

		err := temp.WriteTemplate(temp.CmdTemplate, Commands)
		if err != nil {
			fmt.Println("Error:", err.Error())
		}
	}
}