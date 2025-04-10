package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTeamSuitesService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/testorg/teams/c6fa9b07-efeb-4aea-b5ad-c4aa01e91038/suites", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			[{
				"suite_id": "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
				"suite_url": "https://api.buildkite.com/v2/analytics/organizations/testorg/suites/suite-dreams",
				"access_level": ["read"],
				"created_at": "2023-08-10T05:24:08.651Z"
			},
			{
				"suite_id": "4569",
				"suite_url": "https://api.buildkite.com/v2/analytics/organizations/testorg/suites/suite-and-sour",
				"access_level": ["read", "edit"],
				"created_at": "2023-08-10T05:24:08.663Z"
			}]
			`)
	})

	got, _, err := client.TeamSuites.List(context.Background(), "testorg", "c6fa9b07-efeb-4aea-b5ad-c4aa01e91038", nil)
	if err != nil {
		t.Errorf("TeamSuitesService.List returned error: %v", err)
	}

	suite1CreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-10T05:24:08.651Z"))
	suite2CreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-10T05:24:08.663Z"))

	want := []TeamSuites{
		{
			ID:          "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
			URL:         "https://api.buildkite.com/v2/analytics/organizations/testorg/suites/suite-dreams",
			AccessLevel: []string{"read"},
			CreatedAt:   NewTimestamp(suite1CreatedAt),
		},
		{
			ID:          "4569",
			URL:         "https://api.buildkite.com/v2/analytics/organizations/testorg/suites/suite-and-sour",
			AccessLevel: []string{"read", "edit"},
			CreatedAt:   NewTimestamp(suite2CreatedAt),
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TeamSuitesService.List diff: (-got +want)\n%s", diff)
	}
}

func TestTeamSuitesService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/testorg/teams/c6fa9b07-efeb-4aea-b5ad-c4aa01e91038/suites/1239d7f9-394a-4d99-badf-7c3d8577a8ff", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			{
				"suite_id": "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
				"suite_url": "https://api.buildkite.com/v2/analytics/organizations/testorg/suites/suite-dreams",
				"access_level": ["read"],
				"created_at": "2023-08-10T05:24:08.651Z"
			}`)
	})

	got, _, err := client.TeamSuites.Get(context.Background(), "testorg", "c6fa9b07-efeb-4aea-b5ad-c4aa01e91038", "1239d7f9-394a-4d99-badf-7c3d8577a8ff")
	if err != nil {
		t.Errorf("TeamSuitesService.List returned error: %v", err)
	}

	suiteCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-10T05:24:08.651Z"))
	want := TeamSuites{
		ID:          "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
		URL:         "https://api.buildkite.com/v2/analytics/organizations/testorg/suites/suite-dreams",
		AccessLevel: []string{"read"},
		CreatedAt:   NewTimestamp(suiteCreatedAt),
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TeamSuitesService.List diff: (-got +want)\n%s", diff)
	}
}

func TestTeamSuitesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := CreateTeamSuites{
		SuiteID:     "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
		AccessLevel: []string{"read", "edit"},
	}

	server.HandleFunc("/v2/organizations/testorg/teams/c6fa9b07-efeb-4aea-b5ad-c4aa01e91038/suite", func(w http.ResponseWriter, r *http.Request) {
		var v CreateTeamSuites
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Errorf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("create Team Suites input diff: (-got +want)\n%s", diff)
		}

		_, _ = fmt.Fprint(w, `{
			"suite_id": "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
			"suite_url": "https://api.buildkite.com/v2/analytics/organizations/testorg/suites/suite-dreams",
			"access_level": ["read", "edit"],
			"created_at": "2023-08-10T05:24:08.651Z"}`)
	})

	got, _, err := client.TeamSuites.Create(context.Background(), "testorg", "c6fa9b07-efeb-4aea-b5ad-c4aa01e91038", input)
	if err != nil {
		t.Errorf("TeamSuitesService.Create returned error: %v", err)
	}

	suiteCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-10T05:24:08.651Z"))
	want := TeamSuites{
		ID:          "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
		URL:         "https://api.buildkite.com/v2/analytics/organizations/testorg/suites/suite-dreams",
		AccessLevel: []string{"read", "edit"},
		CreatedAt:   NewTimestamp(suiteCreatedAt),
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TeamSuitesService.Create diff: (-got +want)\n%s", diff)
	}
}

func TestTeamSuitesService_Update(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/testorg/teams/c6fa9b07-efeb-4aea-b5ad-c4aa01e91038/suites/1239d7f9-394a-4d99-badf-7c3d8577a8ff", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		_, _ = fmt.Fprint(w,
			`
			{
				"suite_id": "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
				"suite_url": "https://api.buildkite.com/v2/analytics/organizations/testorg/suites/suite-dreams",
				"access_level": ["read", "edit"],
				"created_at": "2023-08-10T05:24:08.651Z"
			}`)
	})

	wantUpdate := UpdateTeamSuites{
		AccessLevel: []string{"read", "edit"},
	}

	got, _, err := client.TeamSuites.Update(context.Background(), "testorg", "c6fa9b07-efeb-4aea-b5ad-c4aa01e91038", "1239d7f9-394a-4d99-badf-7c3d8577a8ff", wantUpdate)
	if err != nil {
		t.Errorf("TeamSuitesService.Get returned error: %v", err)
	}

	pipelineCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-10T05:24:08.651Z"))
	want := TeamSuites{
		ID:          "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
		URL:         "https://api.buildkite.com/v2/analytics/organizations/testorg/suites/suite-dreams",
		AccessLevel: []string{"read", "edit"},
		CreatedAt:   NewTimestamp(pipelineCreatedAt),
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TeamSuitesService.Get diff: (-got +want)\n%s", diff)
	}
}

func TestTeamSuitesService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/testorg/teams/c6fa9b07-efeb-4aea-b5ad-c4aa01e91038/suites/1239d7f9-394a-4d99-badf-7c3d8577a8ff", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.TeamSuites.Delete(context.Background(), "testorg", "c6fa9b07-efeb-4aea-b5ad-c4aa01e91038", "1239d7f9-394a-4d99-badf-7c3d8577a8ff")

	if err != nil {
		t.Errorf("TeamSuitesService.Delete returned error: %v", err)
	}

}
