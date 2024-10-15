package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUserService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"123","name":"Jane Doe","email":"jane@doe.com"}`)
	})

	user, _, err := client.User.Get(context.Background())
	if err != nil {
		t.Errorf("User.Get returned error: %v", err)
	}

	want := User{ID: "123", Name: "Jane Doe", Email: "jane@doe.com"}
	if diff := cmp.Diff(user, want); diff != "" {
		t.Errorf("User.Get diff: (-got +want)\n%s", diff)
	}
}
