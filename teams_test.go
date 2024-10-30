package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTeamsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	teams, _, err := client.Teams.List(context.Background(), "my-great-org", nil)
	if err != nil {
		t.Errorf("Teams.List returned error: %v", err)
	}

	want := []Team{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(teams, want); diff != "" {
		t.Errorf("Teams.List diff: (-got +want)\n%s", diff)
	}
}

func TestTeamsService_ListForUser(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"user_id": "abc",
		})
		fmt.Fprint(w, `[{"id":"123"}]`)
	})

	opt := &TeamsListOptions{UserID: "abc"}
	teams, _, err := client.Teams.List(context.Background(), "my-great-org", opt)
	if err != nil {
		t.Errorf("Teams.List returned error: %v", err)
	}

	want := []Team{{ID: "123"}}
	if diff := cmp.Diff(teams, want); diff != "" {
		t.Errorf("Teams.List diff: (-got +want)\n%s", diff)
	}
}

func TestTeamsService_GetTeam(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"123"}`)
	})

	team, err := client.Teams.GetTeam(context.Background(), "my-great-org", "123")
	if err != nil {
		t.Errorf("Teams.GetTeam returned error: %v", err)
	}

	want := Team{ID: "123"}
	if diff := cmp.Diff(team, want); diff != "" {
		t.Errorf("Teams.GetTeam diff: (-got +want)\n%s", diff)
	}
}

func TestTeamsService_CreateTeam(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":"123"}`)
	})

	team, _, err := client.Teams.CreateTeam(context.Background(), "my-great-org", CreateTeam{})
	if err != nil {
		t.Errorf("Teams.CreateTeam returned error: %v", err)
	}

	want := Team{ID: "123"}
	if diff := cmp.Diff(team, want); diff != "" {
		t.Errorf("Teams.CreateTeam diff: (-got +want)\n%s", diff)
	}
}

func TestTeamsService_UpdateTeam(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":"123"}`)
	})

	team, _, err := client.Teams.UpdateTeam(context.Background(), "my-great-org", "123", CreateTeam{})
	if err != nil {
		t.Errorf("Teams.UpdateTeam returned error: %v", err)
	}

	want := Team{ID: "123"}
	if diff := cmp.Diff(team, want); diff != "" {
		t.Errorf("Teams.UpdateTeam diff: (-got +want)\n%s", diff)
	}
}

func TestTeamsService_DeleteTeam(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Teams.DeleteTeam(context.Background(), "my-great-org", "123")
	if err != nil {
		t.Errorf("Teams.DeleteTeam returned error: %v", err)
	}
}

func TestTeamsService_ListTeamMembers(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/teams/123/members", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"user_id":"123"}]`)
	})

	members, _, err := client.Teams.ListTeamMembers(context.Background(), "my-great-org", "123", nil)
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
		fmt.Fprint(w, `[{"user_id":"123"},{"user_id":"456"}]`)
	})

	opt := &TeamMembersListOptions{
		ListOptions: ListOptions{Page: 2},
	}
	members, _, err := client.Teams.ListTeamMembers(context.Background(), "my-great-org", "123", opt)
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
		fmt.Fprint(w, `{"user_id":"456"}`)
	})

	member, err := client.Teams.GetTeamMember(context.Background(), "my-great-org", "123", "456")
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
		fmt.Fprint(w, `{"user_id":"456"}`)
	})

	createTeamMember := CreateTeamMember{UserID: "456"}
	m, _, err := client.Teams.CreateTeamMember(context.Background(), "my-great-org", "123", createTeamMember)
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
		fmt.Fprint(w, `{"user_id":"456"}`)
	})

	m, _, err := client.Teams.UpdateTeamMember(context.Background(), "my-great-org", "123", "456", "maintainer")
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

	_, err := client.Teams.DeleteTeamMember(context.Background(), "my-great-org", "123", "456")
	if err != nil {
		t.Errorf("Teams.DeleteTeamMember returned error: %v", err)
	}
}
