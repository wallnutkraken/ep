package ep

import "os"

func init() {
	/* Create ep data folders if they don't exist */
	checkAndCreate := func (path string) {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			mkDirErr := os.Mkdir(path, os.ModePerm)
			if mkDirErr != nil {
				panic("failed to create directory " + path + " due to error:\n" + mkDirErr.Error() + "\n")
			}
		} else if err != nil {
			panic("init error: " + err.Error())
		}
	}
	checkAndCreate(getStorageDir())
	checkAndCreate(feedsDir)
}
func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
