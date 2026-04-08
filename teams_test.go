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
		_, _ = fmt.Fprint(w, `[
			{
				"id": "123",
				"graphql_id": "VGVhbS0tLTEyMw==",
				"default_member_role": "member",
				"members_can_create_pipelines": true,
				"members_can_create_suites": false,
				"members_can_create_registries": true,
				"members_can_destroy_registries": false,
				"members_can_destroy_packages": false
			},
			{"id": "1234"}
		]`)
	})

	teams, _, err := client.Teams.List(context.Background(), "my-great-org", nil)
	if err != nil {
		t.Errorf("Teams.List returned error: %v", err)
	}

	want := []Team{
		{
			ID:                          "123",
			GraphQLID:                   "VGVhbS0tLTEyMw==",
			DefaultMemberRole:           "member",
			MembersCanCreatePipelines:   true,
			MembersCanCreateSuites:      false,
			MembersCanCreateRegistries:  true,
			MembersCanDestroyRegistries: false,
			MembersCanDestroyPackages:   false,
		},
		{ID: "1234"},
	}
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
		_, _ = fmt.Fprint(w, `[{
			"id": "123",
			"graphql_id": "VGVhbS0tLTEyMw==",
			"default_member_role": "member",
			"members_can_create_pipelines": true,
			"members_can_create_suites": false,
			"members_can_create_registries": false,
			"members_can_destroy_registries": false,
			"members_can_destroy_packages": false
		}]`)
	})

	opt := &TeamsListOptions{UserID: "abc"}
	teams, _, err := client.Teams.List(context.Background(), "my-great-org", opt)
	if err != nil {
		t.Errorf("Teams.List returned error: %v", err)
	}

	want := []Team{{
		ID:                          "123",
		GraphQLID:                   "VGVhbS0tLTEyMw==",
		DefaultMemberRole:           "member",
		MembersCanCreatePipelines:   true,
		MembersCanCreateSuites:      false,
		MembersCanCreateRegistries:  false,
		MembersCanDestroyRegistries: false,
		MembersCanDestroyPackages:   false,
	}}
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
		_, _ = fmt.Fprint(w, `{
			"id": "123",
			"graphql_id": "VGVhbS0tLTEyMw==",
			"name": "my-team",
			"slug": "my-team",
			"description": "A test team",
			"privacy": "visible",
			"default": true,
			"default_member_role": "maintainer",
			"members_can_create_pipelines": true,
			"members_can_create_suites": true,
			"members_can_create_registries": false,
			"members_can_destroy_registries": false,
			"members_can_destroy_packages": true
		}`)
	})

	team, err := client.Teams.GetTeam(context.Background(), "my-great-org", "123")
	if err != nil {
		t.Errorf("Teams.GetTeam returned error: %v", err)
	}

	want := Team{
		ID:                          "123",
		GraphQLID:                   "VGVhbS0tLTEyMw==",
		Name:                        "my-team",
		Slug:                        "my-team",
		Description:                 "A test team",
		Privacy:                     "visible",
		Default:                     true,
		DefaultMemberRole:           "maintainer",
		MembersCanCreatePipelines:   true,
		MembersCanCreateSuites:      true,
		MembersCanCreateRegistries:  false,
		MembersCanDestroyRegistries: false,
		MembersCanDestroyPackages:   true,
	}
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
		_, _ = fmt.Fprint(w, `{
			"id": "123",
			"graphql_id": "VGVhbS0tLTEyMw==",
			"name": "new-team",
			"slug": "new-team",
			"privacy": "visible",
			"default": false,
			"default_member_role": "member",
			"members_can_create_pipelines": true,
			"members_can_create_suites": false,
			"members_can_create_registries": false,
			"members_can_destroy_registries": false,
			"members_can_destroy_packages": false
		}`)
	})

	team, _, err := client.Teams.CreateTeam(context.Background(), "my-great-org", CreateTeam{
		Name:                      "new-team",
		Privacy:                   "visible",
		MembersCanCreatePipelines: true,
	})
	if err != nil {
		t.Errorf("Teams.CreateTeam returned error: %v", err)
	}

	want := Team{
		ID:                        "123",
		GraphQLID:                 "VGVhbS0tLTEyMw==",
		Name:                      "new-team",
		Slug:                      "new-team",
		Privacy:                   "visible",
		DefaultMemberRole:         "member",
		MembersCanCreatePipelines: true,
	}
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
		_, _ = fmt.Fprint(w, `{
			"id": "123",
			"graphql_id": "VGVhbS0tLTEyMw==",
			"name": "updated-team",
			"slug": "updated-team",
			"description": "Updated description",
			"privacy": "secret",
			"default": false,
			"default_member_role": "maintainer",
			"members_can_create_pipelines": true,
			"members_can_create_suites": true,
			"members_can_create_registries": false,
			"members_can_destroy_registries": false,
			"members_can_destroy_packages": false
		}`)
	})

	team, _, err := client.Teams.UpdateTeam(context.Background(), "my-great-org", "123", CreateTeam{
		Name:                      "updated-team",
		Description:               "Updated description",
		Privacy:                   "secret",
		DefaultMemberRole:         "maintainer",
		MembersCanCreatePipelines: true,
	})
	if err != nil {
		t.Errorf("Teams.UpdateTeam returned error: %v", err)
	}

	want := Team{
		ID:                        "123",
		GraphQLID:                 "VGVhbS0tLTEyMw==",
		Name:                      "updated-team",
		Slug:                      "updated-team",
		Description:               "Updated description",
		Privacy:                   "secret",
		DefaultMemberRole:         "maintainer",
		MembersCanCreatePipelines: true,
		MembersCanCreateSuites:    true,
	}
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
