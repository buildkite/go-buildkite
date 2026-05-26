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
		t.Errorf("got invalid JSON %q: %v", got, err)
		return
	}

	var wantJSON any
	if err := json.Unmarshal([]byte(want), &wantJSON); err != nil {
		t.Errorf("want invalid JSON %q: %v", want, err)
		return
	}

	if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
		t.Errorf("JSON mismatch (-want +got):\n%s", diff)
	}
}

func assertRequestJSON(t *testing.T, r *http.Request, want string) {
	t.Helper()

	got, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("reading request body: %v", err)
		return
	}

	assertJSONEqual(t, string(got), want)
}
