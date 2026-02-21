package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMemberService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w,
			`
			[
				{
					"id": "528000d8-4ee1-4479-8af1-032b143185f0",
					"name": "Cool Dev",
					"email": "cool.dev@my-great.org"
				},
				{
					"id": "3847t888-4e31-4e79-8ai1-03237587en80",
					"name": "Dool Cev",
					"email": "dool.cev@my-great.org"
				}
			]`)
	})

	members, _, err := client.Members.List(context.Background(), "my-great-org", nil)
	if err != nil {
		t.Errorf("TestMembers.List returned error: %v", err)
	}

	want := []Member{
		{
			UUID:  "528000d8-4ee1-4479-8af1-032b143185f0",
			Name:  "Cool Dev",
			Email: "cool.dev@my-great.org",
		},
		{
			UUID:  "3847t888-4e31-4e79-8ai1-03237587en80",
			Name:  "Dool Cev",
			Email: "dool.cev@my-great.org",
		},
	}

	if diff := cmp.Diff(members, want); diff != "" {
		t.Errorf("TestClusters.List diff: (-got +want)\n%s", diff)
	}
}

func TestMemberService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/members/528000d8-4ee1-4479-8af1-032b143185f0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w,
			`
				{
					"id": "528000d8-4ee1-4479-8af1-032b143185f0",
					"name": "Cool Dev",
					"email": "cool.dev@my-great.org"
				},
			`)
	})

	member, _, err := client.Members.Get(context.Background(), "my-great-org", "528000d8-4ee1-4479-8af1-032b143185f0")
	if err != nil {
		t.Errorf("TestMembers.List returned error: %v", err)
	}

	want := Member{
		UUID:  "528000d8-4ee1-4479-8af1-032b143185f0",
		Name:  "Cool Dev",
		Email: "cool.dev@my-great.org",
	}

	if diff := cmp.Diff(member, want); diff != "" {
		t.Errorf("TestClusters.Get diff: (-got +want)\n%s", diff)
	}
}
