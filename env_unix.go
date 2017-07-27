//+build !windows

package ep

import "os/user"

const (
	dir_seperator = "/"
)

func getStorageDir() string {
	user, _ := user.Current()
	return user.HomeDir + "/.config/ep"
}