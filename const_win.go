//+build windows

package ep

const (
	dir_seperator="\"
)

func getStorageDir() string {
	return `%appdata%\ep`
}