package buildkite

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAccessTokensService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/access-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"uuid": "b63254c0-3271-4a98-8270-7cfbd6c2f14e","scopes": ["read_build"]}`)
	})

	ats, _, err := client.AccessTokens.Get()
	if err != nil {
		t.Errorf("AccessTokens.Get returned error: %v", err)
	}

	fmt.Print(ats)

	want := &AccessToken{
		UUID:   String("b63254c0-3271-4a98-8270-7cfbd6c2f14e"),
		Scopes: &[]string{"read_build"},
	}
	if !reflect.DeepEqual(ats, want) {
		t.Errorf("AccessTokens.Get returned %+v, want %+v", ats, want)
	}
}

func TestAccessTokensService_Revoke(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/access-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{}`)
	})

	ats, _, err := client.AccessTokens.Revoke()
	if err != nil {
		t.Errorf("AccessTokens.Revoke returned error: %v", err)
	}

	fmt.Print(ats)

	want := &AccessToken{}

	if !reflect.DeepEqual(ats, want) {
		t.Errorf("AccessTokens.Revoke returned %+v, want %+v", ats, want)
	}
}
