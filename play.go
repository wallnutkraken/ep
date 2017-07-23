package ep

import (
	"os/exec"
	"net/http"
	"strings"
	"errors"
)

func Stream(episode Episode) error {
	if !strings.HasSuffix(episode.URL, ".mp3") {
		return errors.New("currently only mp3 is supported")
	}
	resp, err := http.Get(strings.Trim(episode.URL, " "))
	if err != nil {
		return err
	}
	dataChan := make(chan []byte, 64)

	cmd := exec.Command("mpg123", "-q", "-")

	cmd.Stdin = resp.Body
	err = cmd.Start()
	if err != nil {
		close(dataChan)
		return err
	}

	return cmd.Wait()
}
