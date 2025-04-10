package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTeamsService_ListTeamMembers(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams/123/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{"user_id":"123"}]`)
	})

	members, _, err := client.TeamMember.ListTeamMembers(context.Background(), "my-great-org", "123", nil)
	if err != nil {
		t.Errorf("Teams.ListTeamMembers returned error: %v", err)
	}

	want := []TeamMember{{ID: "123"}}
	if diff := cmp.Diff(members, want); diff != "" {
		t.Errorf("Teams.ListTeamMembers diff: (-got +want)\n%s", diff)
	}
}

func TestTeamsService_ListTeamMembersPaginated(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams/123/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		_, _ = fmt.Fprint(w, `[{"user_id":"123"},{"user_id":"456"}]`)
	})

	opt := &TeamMembersListOptions{
		ListOptions: ListOptions{Page: 2},
	}
	members, _, err := client.TeamMember.ListTeamMembers(context.Background(), "my-great-org", "123", opt)
	if err != nil {
		t.Errorf("Teams.ListTeamMembers returned error: %v", err)
	}

	want := []TeamMember{{ID: "123"}, {ID: "456"}}
	if diff := cmp.Diff(members, want); diff != "" {
		t.Errorf("Teams.ListTeamMembers diff: (-got +want)\n%s", diff)
	}
}

func TestTeamsService_GetTeamMember(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams/123/members/456", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{"user_id":"456"}`)
	})

	member, err := client.TeamMember.GetTeamMember(context.Background(), "my-great-org", "123", "456")
	if err != nil {
		t.Errorf("Teams.GetTeamMember returned error: %v", err)
	}

	want := TeamMember{ID: "456"}
	if diff := cmp.Diff(member, want); diff != "" {
		t.Errorf("Teams.GetTeamMember diff: (-got +want)\n%s", diff)
	}
}

func TestTeamsService_CreateTeamMember(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams/123/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = fmt.Fprint(w, `{"user_id":"456"}`)
	})

	createTeamMember := CreateTeamMember{UserID: "456"}
	m, _, err := client.TeamMember.CreateTeamMember(context.Background(), "my-great-org", "123", createTeamMember)
	if err != nil {
		t.Errorf("Teams.CreateTeamMember returned error: %v", err)
	}

	want := TeamMember{ID: "456"}
	if diff := cmp.Diff(m, want); diff != "" {
		t.Errorf("Teams.CreateTeamMember diff: (-got +want)\n%s", diff)
	}
}

func TestTeamsService_UpdateTeamMember(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams/123/members/456", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		_, _ = fmt.Fprint(w, `{"user_id":"456"}`)
	})

	m, _, err := client.TeamMember.UpdateTeamMember(context.Background(), "my-great-org", "123", "456", "maintainer")
	if err != nil {
		t.Errorf("Teams.UpdateTeamMember returned error: %v", err)
	}

	want := TeamMember{ID: "456"}
	if diff := cmp.Diff(m, want); diff != "" {
		t.Errorf("Teams.UpdateTeamMember diff: (-got +want)\n%s", diff)
	}
}

func TestTeamsService_DeleteTeamMember(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams/123/members/456", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.TeamMember.DeleteTeamMember(context.Background(), "my-great-org", "123", "456")
	if err != nil {
		t.Errorf("Teams.DeleteTeamMember returned error: %v", err)
	}
}
