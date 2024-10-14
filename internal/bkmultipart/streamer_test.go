package bkmultipart_test

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"strings"
	"testing"

	"github.com/buildkite/go-buildkite/v3/internal/bkmultipart"
)

func TestEncodingFormFields(t *testing.T) {
	t.Parallel()

	s := bkmultipart.NewStreamer()
	fields := map[string]string{
		"mountain": "cotopaxi",
		"city":     "guayaquil",
	}
	if err := s.WriteFields(fields); err != nil {
		t.Fatalf("WriteFields(%# v) = %v, want nil", fields, err)
	}

	form := parseForm(t, s)
	for key, values := range form.Value {
		if len(values) != 1 {
			t.Fatalf("got %d values for %s, want 1", len(values), key)
		}
		if got, want := values[0], fields[key]; got != want {
			t.Fatalf("form.Value[%q] = [%q], want [%q]", key, got, want)
		}
	}
}

func TestEncodingFile(t *testing.T) {
	t.Parallel()

	tempFile, err := os.CreateTemp("", "buildkite-go-sdk-test")
	if err != nil {
		t.Fatalf("os.CreateTemp() = %v, want nil", err)
	}

	t.Cleanup(func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	})

	if _, err := tempFile.WriteString("hello world"); err != nil {
		t.Fatalf(`tempFile.WriteString("hello world") = %v, want nil`, err)
	}

	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("tempFile.Seek(0, io.SeekStart) = %v, want nil", err)
	}

	s := bkmultipart.NewStreamer()

	fileKey := "file"
	if err := s.WriteFile(fileKey, tempFile, "test.txt"); err != nil {
		t.Fatalf("WriteFile() = %v, want nil", err)
	}

	form := parseForm(t, s)

	if got, want := len(form.File[fileKey]), 1; got != want {
		t.Fatalf("len(form.File[\"file\"]) = %d, want %d", got, want)
	}

	file := form.File[fileKey][0]
	if got, want := file.Filename, "test.txt"; got != want {
		t.Fatalf("file.Filename = %q, want %q", got, want)
	}

	f, err := file.Open()
	if err != nil {
		t.Fatalf("file.Open() = %v, want nil", err)
	}

	b := &strings.Builder{}
	if _, err := io.Copy(b, f); err != nil {
		t.Fatalf("io.Copy(strings.Builder, file) = %v, want nil", err)
	}

	if got, want := b.String(), "hello world"; got != want {
		t.Fatalf("file contents = %q, want %q", got, want)
	}
}

func parseForm(t *testing.T, streamer *bkmultipart.Streamer) *multipart.Form {
	_, params, err := mime.ParseMediaType(streamer.ContentType)
	if err != nil {
		t.Fatalf("mime.ParseMediaType(%q) = %v, want nil", streamer.ContentType, err)
	}

	boundary, ok := params["boundary"]
	if !ok {
		t.Fatalf("want boundary in content type")
	}

	b := &bytes.Buffer{}
	_, err = io.Copy(b, streamer.Reader())
	if err != nil {
		t.Fatalf("io.Copy(bytes.Buffer, s) = %v, want nil", err)
	}

	mr := multipart.NewReader(b, boundary)
	form, err := mr.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("mr.ReadForm() = %v, want nil", err)
	}

	return form
}
