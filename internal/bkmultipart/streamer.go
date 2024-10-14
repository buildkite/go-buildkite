package bkmultipart

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

// A wrapper around the complexities of streaming a multipart file and fields to an http endpoint that infuriatingly
// requires a Content-Length
// Stolen/adapted from the multipart streamer in https://github.com/buildkite/agentv/v3/internal/artifact, which itself
// was derived from https://github.com/technoweenie/multipartstreamer
type Streamer struct {
	ContentType string

	bodyBuffer  *bytes.Buffer
	bodyWriter  *multipart.Writer
	closeBuffer *bytes.Reader
	file        *os.File

	fileLength  int64
	writtenFile bool
}

// NewStreamer initializes a new Streamer.
func NewStreamer() *Streamer {
	s := &Streamer{bodyBuffer: &bytes.Buffer{}}

	s.bodyWriter = multipart.NewWriter(s.bodyBuffer)
	boundary := s.bodyWriter.Boundary()
	s.ContentType = "multipart/form-data; boundary=" + boundary

	closeBoundary := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	s.closeBuffer = bytes.NewReader([]byte(closeBoundary))

	return s
}

// WriteField writes a form field to the multipart.Writer.
func (s *Streamer) WriteField(key, value string) error {
	return s.bodyWriter.WriteField(key, value)
}

func (s *Streamer) WriteFields(fields map[string]string) error {
	for key, value := range fields {
		if err := s.WriteField(key, value); err != nil {
			return err
		}
	}
	return nil
}

// WriteFile writes the multi-part preamble which will be followed by file data
// This can only be called once and must be the last thing written to the streamer
func (s *Streamer) WriteFile(key string, file *os.File, displayedFilename string) error {
	if s.file != nil {
		return errors.New("WriteFile can't be called multiple times")
	}

	fi, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	s.fileLength = fi.Size()

	_, err = s.bodyWriter.CreateFormFile(key, displayedFilename)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	s.file = file

	return nil
}

// Len calculates the byte size of the multipart content.
func (s *Streamer) Len() int64 {
	return int64(s.bodyBuffer.Len()) + s.fileLength + int64(s.closeBuffer.Len())
}

func (s *Streamer) Reader() io.Reader {
	if s.file == nil {
		return io.MultiReader(s.bodyBuffer, s.closeBuffer)
	} else {
		return io.MultiReader(s.bodyBuffer, s.file, s.closeBuffer)
	}
}
