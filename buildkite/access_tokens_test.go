package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAccessTokensService_Get(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/access-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"uuid": "b63254c0-3271-4a98-8270-7cfbd6c2f14e","scopes": ["read_build"]}`)
	})

	ats, _, err := client.AccessTokens.Get(context.Background())
	if err != nil {
		t.Errorf("AccessTokens.Get returned error: %v", err)
	}

	want := &AccessToken{
		UUID:   String("b63254c0-3271-4a98-8270-7cfbd6c2f14e"),
		Scopes: &[]string{"read_build"},
	}
	if !reflect.DeepEqual(ats, want) {
		t.Errorf("AccessTokens.Get returned %+v, want %+v", ats, want)
	}
}

func TestAccessTokensService_Revoke(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/access-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.AccessTokens.Revoke(context.Background())
	if err != nil {
		t.Errorf("AccessTokens.Revoke returned error: %v", err)
	}

	want := http.StatusNoContent
	got := resp.Response.StatusCode

	if !reflect.DeepEqual(want, got) {
		t.Errorf("AccessTokens.Revoke returned %+v, want %+v", got, want)
	}
}
