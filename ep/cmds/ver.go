package cmds

import (
	"fmt"
	"github.com/wallnutkraken/ep/ep/version"
)

func Version(args []string) {
	fmt.Printf("ep version %s\n", version.String())
}