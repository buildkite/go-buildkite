package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAccessTokensService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/access-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"uuid": "b63254c0-3271-4a98-8270-7cfbd6c2f14e","scopes": ["read_build"]}`)
	})

	ats, _, err := client.AccessTokens.Get(context.Background())
	if err != nil {
		t.Errorf("AccessTokens.Get returned error: %v", err)
	}

	want := AccessToken{
		UUID:   "b63254c0-3271-4a98-8270-7cfbd6c2f14e",
		Scopes: []string{"read_build"},
	}
	if diff := cmp.Diff(ats, want); diff != "" {
		t.Errorf("AccessTokens.Get diff: (-got +want)\n%s", diff)
	}
}

func TestAccessTokensService_Revoke(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/access-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.AccessTokens.Revoke(context.Background())
	if err != nil {
		t.Errorf("AccessTokens.Revoke returned error: %v", err)
	}

	want := http.StatusNoContent
	got := resp.Response.StatusCode

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("AccessTokens.Revoke diff: (-got +want)\n%s", diff)
	}
}
