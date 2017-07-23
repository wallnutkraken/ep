//+build windows

package ep

import "os"

const (
	dir_seperator = "\\"
)

func getStorageDir() string {
	return os.Getenv("APPDATA") + `\ep`
}