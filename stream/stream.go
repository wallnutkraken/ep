// Package stream is responsible for creating streams from URLs and retrieving any relevant data
// such as file type
package stream

import (
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"github.com/wallnutkraken/ep/player"

	"github.com/pkg/errors"
)

// File contains the information of the file to be streamed, and methods to begin streaming
type File struct {
	url       string
	filetype  string
	mimetype  string
	extension string
}

var (
	// ErrMediaTypeUnsupported is the error for an unsupported media type being detected
	ErrMediaTypeUnsupported = errors.New("This media type is not supported by ep, yet")
)

// FromURL creates a stream File info object from the given file URL
func FromURL(url string) File {
	// First make sure we got a clean URL, sometimes RSS feeds can be dirty
	return File{
		url: strings.TrimSpace(url),
	}
}

// GetStream starts streaming the file and returns the ReadCloser for the stream data
func (f *File) GetStream() (io.ReadCloser, error) {
	// Call GET on the URL
	resp, err := http.Get(f.url)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed calling GET on [%s]", f.url)
	}

	// Read the first 261 bytes to get the file header
	head := make([]byte, 261)
	_, err = resp.Body.Read(head)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed reading file header from file at [%s]", f.url)
	}

	// Get the mime type and extension
	f.mimetype = resp.Header.Get("Content-Type")
	filename := path.Base(resp.Request.URL.Path)
	filenameParts := strings.Split(filename, ".")
	f.extension = filenameParts[len(filenameParts)-1]

	// Return the body
	return resp.Body, nil
}

// GetDecoder returns the audio decoder function for the data type detected. This function first
// tries to detect the data type, then returns a decode function or an error.
func (f File) GetDecoder() (player.DecodeFunc, error) {
	if f.mimetype != "" {
		return f.typeFromMIME()
	}
	return f.typeFromExtension()
}

func wavDecoder(rc io.ReadCloser) (s beep.StreamSeekCloser, format beep.Format, err error) {
	return wav.Decode(rc)
}

func (f File) typeFromMIME() (player.DecodeFunc, error) {
	switch strings.ToLower(f.mimetype) {
	case "audio/mpeg":
		fallthrough
	case "audio/mp3":
		return mp3.Decode, nil
	case "audio/wav":
		return wavDecoder, nil
	default:
		return nil, ErrMediaTypeUnsupported
	}
}

func (f File) typeFromExtension() (player.DecodeFunc, error) {
	switch strings.ToLower(f.extension) {
	case "mp3":
		return mp3.Decode, nil
	case "wav":
		return wavDecoder, nil
	default:
		return nil, ErrMediaTypeUnsupported
	}
}
