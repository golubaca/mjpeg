// Package mjpeg handles mjpeg video stream from IP cameras
package mjpeg

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"
)

// Reader holds data needed for parsing mjpeg stream
type Reader struct {
	Reader   bufio.Reader
	Stream   http.Response
	Boundary string
	Image    []byte
	Headers  http.Header
	Length   string
	Height   int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// New instantiate reader that handles the stream.
// It handles opening new connection to the camera
// and setting initial values
func (mjpeg *Reader) New(streamURL string, auth ...map[string]string) (err error) {
	var stream *http.Response
	if len(auth) > 0 {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", streamURL, nil)
		req.SetBasicAuth(auth[0]["username"], auth[0]["passwd"])
		stream, err = client.Do(req)
	} else {
		stream, err = http.Get(streamURL)
	}

	if err != nil {
		return err
	}
	mjpeg.Stream = *stream
	reader := bufio.NewReader(stream.Body)
	mjpeg.Reader = *reader
	// Bondary can be found in headers so we split value to extract boundary
	mjpeg.Headers = stream.Header
	mjpeg.Boundary = mjpeg.GetBoundary()

	return nil

}

// Close just closes stream body
func (mjpeg *Reader) Close() {
	mjpeg.Stream.Body.Close()
}

// GetBoundary is a helper function which provides simple way of
// accessing boundary from header of a stream
func (mjpeg *Reader) GetBoundary() (boundary string) {
	return strings.Split(mjpeg.Headers.Get("Content-Type"), "=")[1]
}

// GetHeader is a helper function which returns value of requested header.
// If specific header is not found, empty string is returned.
func (mjpeg *Reader) GetHeader(headerType string) string {
	return mjpeg.Headers.Get(headerType)
}

// GetFrame extracts image from stream and returns image bytes
func (mjpeg *Reader) GetFrame() (image []byte, length int, imageType string) {
	mjpeg.Image = []byte{}
	for {
		line, _ := mjpeg.Reader.ReadBytes('\n')

		// if boundary is found, it's the end of our image, so we break here
		if bytes.Contains([]byte(line), []byte(mjpeg.Boundary)) {
			break
		} else if bytes.Equal(line, []byte{13, 10}) {
			// if we find line feed and carriage return, we will discard it because it's not part of images
			continue
		} else if bytes.Contains([]byte(line), []byte("Content")) {
			// @TODO set length and type of an image
			continue
		}

		mjpeg.Image = append(mjpeg.Image, []byte(line)...)
	}
	return mjpeg.Image, 0, ""
}
