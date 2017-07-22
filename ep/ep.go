package main

import (
	"flag"
	"fmt"

	_ "github.com/wallnutkraken/ep"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		if args[0] == "help" {
			help(args[1:])
		}

		foundCommand := false
		for _, cmd := range commands {
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
		help(args)
	}
}
