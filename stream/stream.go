// Package stream is responsible for creating streams from URLs and retrieving any relevant data
// such as file type
package stream

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/wallnutkraken/ep/stream/streamtype"
)

// File contains the information of the file to be streamed, and methods to begin streaming
type File struct {
	url       string
	filetype  string
	mimetype  string
	extension string
}

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
	fmt.Printf("[%s]\n", f.mimetype)
	filename := path.Base(resp.Request.URL.Path)
	filenameParts := strings.Split(filename, ".")
	f.extension = filenameParts[len(filenameParts)-1]

	// Return the body
	return resp.Body, nil
}

// GetType returns the file type for the stream. This function will only return anything meaningful after GetStream is called.
// Currently the supported types are:
// MP3
// WAV
func (f File) GetType() string {
	if f.mimetype != "" {
		return f.typeFromMIME()
	}
	return f.typeFromExtension()
}

func (f File) typeFromMIME() string {
	switch strings.ToLower(f.mimetype) {
	case "audio/mpeg":
		fallthrough
	case "audio/mp3":
		return streamtype.MP3
	case "audio/wav":
		return streamtype.WAV
	default:
		return streamtype.Unknown
	}
}

func (f File) typeFromExtension() string {
	switch strings.ToLower(f.extension) {
	case "mp3":
		return streamtype.MP3
	case "wav":
		return streamtype.WAV
	default:
		return streamtype.Unknown
	}
}
