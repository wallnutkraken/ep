package main

import (
	"flag"
	"strings"
	"fmt"
)

type command struct {
	Run func(args []string)
	UsageLine string
	Short string
	Long string
	Flag flag.FlagSet
	CustomFlags bool
}

// Name gets first word in UsageLine
func (c *command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *command) Usage() {
	fmt.Printf("usage: %s\n\n", c.UsageLine)
	fmt.Printf("%s\n", strings.TrimSpace(c.Long))
}

func (c *command) CanRun() bool {
	return c.Run != nil
}

