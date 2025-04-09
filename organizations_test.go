package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	orgs, _, err := client.Organizations.List(context.Background(), nil)
	if err != nil {
		t.Errorf("Organizations.List returned error: %v", err)
	}

	want := []Organization{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(orgs, want); diff != "" {
		t.Errorf("Organizations.List diff: (-got +want)\n%s", diff)
	}
}

func TestOrganizationsService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/babelstoemp", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{"id":"123"}`)
	})

	org, _, err := client.Organizations.Get(context.Background(), "babelstoemp")
	if err != nil {
		t.Errorf("Organizations.Get returned error: %v", err)
	}

	want := Organization{ID: "123"}
	if diff := cmp.Diff(org, want); diff != "" {
		t.Errorf("Organizations.Get diff: (-got +want)\n%s", diff)
	}
}
