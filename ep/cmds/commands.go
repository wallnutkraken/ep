package cmds

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

var Commands []*command = []*command{
	&command{
		UsageLine: "help",
		Short:     "Shows this help screen",
		Long:      "",
	}, &command{
		UsageLine: "add",
		Short:     "Adds a podcast via feed URL",
		Long:      "Adds a podcast via feed URL, usage: ep add PodcastTag https://example.com/podcast/rss",
		Run:       Add,
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
		Run: Update,
	}, &command{
		UsageLine:"play",
		Short:"Plays the latest or specific podcast episode",
		Long:`
"play" can be used in two ways:

"ep play [tag]" will play the latest episode of the podcast corresponding to [tag]

"ep play [tag] [episode]" will play the specific episode for the podcast corresponding
to [tag], to view a list of episodes for a podcast use "ep list [tag]"`,
		Run:Play,
	}, &command{
		UsageLine:"remove",
		Short:"Removes the selected podcast from ep's memory",
		Long:`
"remove" can be used to remove a podcast from ep's memory. This operation cannot be undone.
Podcasts may be added again after removal, however, with the same tag.

Usage: ep remove [tag]`,
		Run: Remove,
	}, &command{
		UsageLine:"version",
		Short:"Prints the version of ep",
		Long:"Do you really need a long explaination of what 'version' does?",
		Run:Version,
	},
}


