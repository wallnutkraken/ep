package ep

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/itchyny/volume-go"
	"github.com/wallnutkraken/ep/progress"
)

const (
	volumeIncrement = 5
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

	ctrl := ctrlWrapper{
		controls: control,
		progBar: progress.Bar{
			Increment: 1,
		},
	}
	ctrl.progBar.Value = ctrl.getVolume()
	return done, ctrl, nil
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
	progBar  progress.Bar
}

func (c ctrlWrapper) TogglePaused() {
	speaker.Lock()
	c.controls.Paused = !c.controls.Paused
	speaker.Unlock()
}

func (c ctrlWrapper) getVolume() int {
	current, err := volume.GetVolume()
	if err != nil {
		fmt.Printf("Error getting volume: [%s]", err.Error())
	}
	return current
}

func (c ctrlWrapper) VolumeUp() {
	next := c.getVolume()
	if next+volumeIncrement > 100 {
		next = 100
	} else {
		next += volumeIncrement
	}
	volume.SetVolume(next)
	c.progBar.Value = c.getVolume()
	c.progBar.Draw()
}

func (c ctrlWrapper) VolumeDown() {
	next := c.getVolume()
	if next-volumeIncrement < 0 {
		next = 0
	} else {
		next -= volumeIncrement
	}
	volume.SetVolume(next)
	c.progBar.Value = c.getVolume()
	c.progBar.Draw()
}

func (c ctrlWrapper) DrawVolume() {
	c.progBar.Draw()
}

type Controller interface {
	TogglePaused()
	VolumeUp()
	VolumeDown()
	DrawVolume()
}
