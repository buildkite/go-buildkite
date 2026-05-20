package buildkite

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

func assertJSONEqual(t *testing.T, got, want string) {
	t.Helper()

	var gotJSON any
	if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
		t.Fatalf("got invalid JSON %q: %v", got, err)
	}

	var wantJSON any
	if err := json.Unmarshal([]byte(want), &wantJSON); err != nil {
		t.Fatalf("want invalid JSON %q: %v", want, err)
	}

	if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
		t.Errorf("JSON mismatch (-want +got):\n%s", diff)
	}
}

func assertRequestJSON(t *testing.T, r *http.Request, want string) {
	t.Helper()

	got, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("reading request body: %v", err)
	}

	assertJSONEqual(t, string(got), want)
}
