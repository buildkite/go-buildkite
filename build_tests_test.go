package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const testBuildUUID = "019d66fb-e8db-47eb-866c-94b85d42b9a1"

func TestBuildTestsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	const linkHeader = `<https://api.buildkite.com/v2/analytics/organizations/my-great-org/builds/019d66fb-e8db-47eb-866c-94b85d42b9a1/tests?page=3>; rel="next", <https://api.buildkite.com/v2/analytics/organizations/my-great-org/builds/019d66fb-e8db-47eb-866c-94b85d42b9a1/tests?page=4>; rel="last"`

	server.HandleFunc(fmt.Sprintf("/v2/analytics/organizations/my-great-org/builds/%s/tests", testBuildUUID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "2",
			"per_page": "50",
			"labels":   "flaky,!slow",
			"branch":   "main*",
			"owners":   "payments,!platform",
			"state":    "enabled",
			"tags":     "framework:rspec,result:^failed",
			"sort_by":  "reliability",
			"order":    "asc",
		})

		w.Header().Set("Link", linkHeader)
		_, _ = fmt.Fprint(w,
			`
			[
				{
					"id": "a915535c-a8f1-4e1a-bd6a-a5589e09f349",
					"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests/a915535c-a8f1-4e1a-bd6a-a5589e09f349",
					"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/tests/a915535c-a8f1-4e1a-bd6a-a5589e09f349",
					"scope": "User#email",
					"name": "is correctly formatted",
					"location": "./spec/models/user_spec.rb:42",
					"file_name": "./spec/models/user_spec.rb",
					"labels": ["flaky"],
					"reliability": 0.9821,
					"duration_avg": 0.213,
					"duration_sum": 23.856,
					"duration_min": 0.108,
					"duration_max": 1.942,
					"executions_count": 113,
					"executions_count_by_result": {
						"passed": 110,
						"failed": 2,
						"skipped": 1
					}
				}
			]`)
	})

	reliability := 0.9821
	got, resp, err := client.BuildTests.List(context.Background(), "my-great-org", testBuildUUID, &BuildTestsListOptions{
		ListOptions: ListOptions{Page: 2, PerPage: 50},
		Labels:      "flaky,!slow",
		Branch:      "main*",
		Owners:      "payments,!platform",
		State:       "enabled",
		Tags:        "framework:rspec,result:^failed",
		SortBy:      "reliability",
		Order:       "asc",
	})
	if err != nil {
		t.Fatalf("BuildTests.List returned error: %v", err)
	}

	want := []TestWithMetrics{
		{
			Test: Test{
				ID:       "a915535c-a8f1-4e1a-bd6a-a5589e09f349",
				URL:      "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests/a915535c-a8f1-4e1a-bd6a-a5589e09f349",
				WebURL:   "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/tests/a915535c-a8f1-4e1a-bd6a-a5589e09f349",
				Scope:    "User#email",
				Name:     "is correctly formatted",
				Location: "./spec/models/user_spec.rb:42",
				FileName: "./spec/models/user_spec.rb",
				Labels:   []string{"flaky"},
			},
			Reliability:     &reliability,
			DurationAverage: 0.213,
			DurationTotal:   23.856,
			DurationMinimum: 0.108,
			DurationMaximum: 1.942,
			ExecutionsCount: 113,
			ExecutionsCountByResult: map[string]int{
				"passed":  110,
				"failed":  2,
				"skipped": 1,
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("BuildTests.List diff: (-got +want)\n%s", diff)
	}
	if got, want := resp.NextPage, 3; got != want {
		t.Errorf("response.NextPage = %d, want %d", got, want)
	}
	if got, want := resp.LastPage, 4; got != want {
		t.Errorf("response.LastPage = %d, want %d", got, want)
	}
	if got := resp.Header.Get("Link"); got != linkHeader {
		t.Errorf("response Link header = %q, want %q", got, linkHeader)
	}
}
