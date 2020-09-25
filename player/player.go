// Package player is responsible for playback of the podcast audio
//
// NOTE: a package used in player requires libasound2-dev on Linux. If you are using a
// Debian-based distro, the following command will install it
//
// apt install libasound2-dev
package player

import (
	"io"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/pkg/errors"
)

// Player contains the controls for audio playback
type Player struct {
	dataInput          io.ReadCloser
	decoder            DecodeFunc
	lock               *sync.Mutex
	controller         *beep.Ctrl
	onFinishedPlayback chan bool
}

// DecodeFunc is the definition for the function signature to decode the audio stream.
// For example, an MP3 decoder.
type DecodeFunc func(io.ReadCloser) (s beep.StreamSeekCloser, format beep.Format, err error)

// New returns a new Player
func New(stream io.ReadCloser, decoder DecodeFunc) *Player {
	return &Player{
		dataInput: stream,
		decoder:   decoder,
		lock:      &sync.Mutex{},
	}
}

// Start begins playback
func (p *Player) Start() error {
	// First, decode the stream
	streamer, format, err := p.decoder(p.dataInput)
	if err != nil {
		return errors.WithMessage(err, "Failed to decode stream")
	}

	// Initialize the speaker with the stream data
	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		return errors.WithMessagef(err, "Failed to initialize speakers with Sample Rate [%d]", format.SampleRate)
	}

	// Create stream controller
	p.controller = &beep.Ctrl{
		Streamer: streamer,
		Paused:   false,
	}

	p.onFinishedPlayback = make(chan bool)
	speaker.Play(beep.Seq(p.controller, beep.Callback(func() {
		p.onFinishedPlayback <- true
		close(p.onFinishedPlayback)
	})))

	return nil
}

// Wait blocks until the player finishes playback
func (p *Player) Wait() {
	<-p.onFinishedPlayback
}
