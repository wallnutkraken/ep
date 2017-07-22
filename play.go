package ep

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"os"
	"os/signal"
	"net/http"
	"github.com/tcolgate/mp3"
	"strings"
	"errors"
	"io"
)

type FrameInfo struct {
	Frame *mp3.Frame
	Skipped int
}

func Stream(episode Episode) error {
	/* TODO: support more than just mp3 */
	if !strings.HasSuffix(episode.URL, "mp3") {
		return errors.New("Currently only mp3 is supported")
	}

	initErr := portaudio.Initialize()
	if initErr != nil {
		return initErr
	}

	var stream *portaudio.Stream
	var err portaudio.Error
	var nErr error
	/* Assuming 44100 sample rate because I'm not very sure on how to obtain something more concrete */
	var sampleRate float64 = 44100

	out := make([]int32, 8192)
	stream, nErr = portaudio.OpenDefaultStream(
		0,
		1,
		sampleRate,
		len(out),
		&out)
	if nErr != nil {
		return nErr
	}
	defer stream.Close()

	/* Prepare the cleanup goroutine to wait patiently for an interrupt */
	interruptChan := make(chan os.Signal)
	go endStream(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, os.Kill)


	/* Begin download, TODO: stream this */
	resp, nErr := http.Get(episode.URL)
	if nErr != nil {
		return nErr
	}

	/* Start decoding the MP3 */
	mp3Decoder := mp3.NewDecoder(resp.Body)
	decodeChan := make(chan FrameInfo, 128)
	go readMp3Frames(mp3Decoder, decodeChan)
	go addFrames(&out, resp.Body, decodeChan)


	/* Before exiting, call cleanup */
	interruptChan <- os.Interrupt
	return nil
}

func endStream(sigChan chan os.Signal) {
	<-sigChan
	/* Recieved interrupt/kill, clean up! */
	portaudio.Terminate()
}

// readMp3Frames is meant to be run as a goroutine in Stream
func readMp3Frames(mp3Decoder *mp3.Decoder, frameChan chan FrameInfo) {
	var err error
	var skipped int
	for err == nil {
		frame := &mp3.Frame{}
		err = mp3Decoder.Decode(frame, &skipped)
		frameChan <- FrameInfo{frame, skipped}
	}
}

func addFrames(out *[]int32, r io.Reader, frameChan chan FrameInfo) {
	for {
		frame := <- frameChan
		frameData := out.
	}
}