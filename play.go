package ep

import (
	"errors"
	"github.com/hajimehoshi/oto"
	"github.com/tcolgate/mp3"
	"net/http"
	"strings"
	"sync"
)

const (
	bufferSize = 0x2000
)

func Stream(episode Episode) error {
	cleanURL := strings.Trim(episode.URL, " ")
	if !strings.HasSuffix(cleanURL, ".mp3") {
		return errors.New("Currently only mp3 files are supported")
	}

	resp, err := http.Get(cleanURL)
	if err != nil {
		return err
	}

	var frame mp3.Frame
	var skipped int
	dec := mp3.NewDecoder(resp.Body)
	/* Decode one frame */
	err = dec.Decode(&frame, &skipped)
	header := frame.Header()

	player, err := oto.NewPlayer(int(header.SampleRate()), 1, 1, bufferSize)
	if err != nil {
		return err
	}

	beginStream(dec, player).Wait()
	return nil
}

func beginStream(decoder *mp3.Decoder, player *oto.Player) *sync.WaitGroup {
	defer player.Close()

	bufChan := make(chan []uint8, 64)
	defer close(bufChan)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go parseMp3(decoder, bufChan, wg)
	go writeToStream(player, bufChan, wg)

	return wg
}

func parseMp3(decoder *mp3.Decoder, buffChan chan []uint8, wg *sync.WaitGroup) {
	frame := mp3.Frame{}
	var skipped int
	for skipped != 0 {
		decoder.Decode(&frame, &skipped)

		reader := frame.Reader()
		buf := make([]uint8, bufferSize)
		var count int = -1
		var err error

		/* Write whole frame in bufferSize chunks */
		for count != 0 {
			count, err = reader.Read(buf)
			if err != nil {
				panic("stream read error: " + err.Error())
			}
			buffChan <- buf
		}
	}

	/* Finish waitgroup */
	wg.Add(-1)
}

func writeToStream(player *oto.Player, buffChan chan []uint8, wg *sync.WaitGroup) {
	var count int = -1
	var err error
	for count != 0 {
		count, err = player.Write(<-buffChan)
		if err != nil {
			panic("player error: " + err.Error())
		}
	}

	/* Finish waitgroup */
	wg.Add(-1)
}
