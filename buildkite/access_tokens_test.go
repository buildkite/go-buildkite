package buildkite

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestGetToken(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/access-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"uuid": "b63254c0-3271-4a98-8270-7cfbd6c2f14e","scopes": ["read_build"]}`)
	})

	accessToken, _, err := client.GetToken()
	if err != nil {
		t.Errorf("GetToken returned error: %v", err)
	}

	want := AccessToken{
		UUID:   String("b63254c0-3271-4a98-8270-7cfbd6c2f14e"),
		Scopes: &[]Scope{"read_build"},
	}
	if !reflect.DeepEqual(want, accessToken) {
		t.Errorf("GetToken returned %+v, want %+v", accessToken, want)
	}
}
