package ep

import (
	"errors"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	bufferSize = 0x2000
)

func Stream(episode Episode) error {
	cleanURL := strings.Trim(episode.URL, " ")
	ext, err := getExtension(cleanURL)
	if err != nil {
		return err
	}

	decodeFunc, ok := audioDecoders[ext]
	if !ok {
		return errors.New("There is no handler for ." + ext + " files")
	}

	resp, err := http.Get(cleanURL)
	if err != nil {
		return err
	}

	stream, format, err := decodeFunc(resp.Body)
	if err != nil {
		return err
	}

	/* Start audio */
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	/* Create a waitgroup for this thread */
	wg := &sync.WaitGroup{}
	wg.Add(1)

	callback := func() {
		wg.Add(-1)
	}

	speaker.Play(beep.Seq(stream, beep.Callback(callback)))
	wg.Wait()

	return nil
}

func getExtension(url string) (string, error) {
	sections := strings.Split(url, "/")
	if len(sections) == 0 {
		return "", errors.New("URL has no sections")
	}
	dotted := strings.Split(sections[len(sections)-1], ".")
	if len(dotted) == 0 {
		return "", errors.New("No extension can be found")
	}
	return dotted[len(dotted)-1], nil
}

var audioDecoders = map[string]func(io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error){
	"wav": wav.Decode,
	"mp3": mp3.Decode,
}
