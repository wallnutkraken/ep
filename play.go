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
	"time"
)

func StartStreaming(episode Episode) (chan bool, Controller, error) {
	cleanURL := strings.Trim(episode.URL, " ")
	ext, err := getExtension(cleanURL)
	if err != nil {
		return nil, nil, err
	}

	decodeFunc, ok := audioDecoders[ext]
	if !ok {
		return nil, nil, errors.New("There is no handler for ." + ext + " files")
	}

	resp, err := http.Get(cleanURL)
	if err != nil {
		return nil, nil, err
	}

	stream, format, err := decodeFunc(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	control := &beep.Ctrl{stream, false}

	/* Start audio */
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	/* Create a channel for this 'done' */
	done := make(chan bool)

	callback := func() {
		done <- true
		close(done)
	}

	speaker.Play(beep.Seq(control, beep.Callback(callback)))

	return done, ctrlWrapper{control}, nil
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

type ctrlWrapper struct {
	controls *beep.Ctrl
}

func (c ctrlWrapper) TogglePaused() {
	speaker.Lock()
	c.controls.Paused = !c.controls.Paused
	speaker.Unlock()
}

type Controller interface {
	TogglePaused()
}