package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUserService_Get(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"123","name":"Jane Doe","email":"jane@doe.com"}`)
	})

	user, _, err := client.User.Get(context.Background())
	if err != nil {
		t.Errorf("User.Get returned error: %v", err)
	}

	want := &User{ID: String("123"), Name: String("Jane Doe"), Email: String("jane@doe.com")}
	if !reflect.DeepEqual(user, want) {
		t.Errorf("User.Get returned %+v, want %+v", user, want)
	}
}
