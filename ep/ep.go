package main

import (
	"flag"
	"fmt"

	_ "github.com/wallnutkraken/ep"
	"github.com/wallnutkraken/ep/ep/cmds"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		if args[0] == "help" {
			cmds.Help(args[1:])
		}

		foundCommand := false
		for _, cmd := range cmds.Commands {
			if cmd.UsageLine == args[0] {
				foundCommand = true
				if cmd.CanRun() {
					cmd.Run(args[1:])
				} else {
					fmt.Println("This command cannot be run")
					return
				}
			}
		}
		if !foundCommand {
			fmt.Println("No such command can be found")
		}
	} else {
		cmds.Help(args)
	}
}
