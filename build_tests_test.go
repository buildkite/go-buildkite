package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestBuildTestsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/builds/abc123/tests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.Query().Get("include"); got != "" {
			t.Errorf("query param include = %q, want empty", got)
		}
		_, _ = fmt.Fprint(w,
			`
			[
				{
					"id": "a915535c-a8f1-4e1a-bd6a-a5589e09f349",
					"scope": "User#email",
					"name": "TestExample1_Create",
					"location": "./spec/models/text_example.rb:55",
					"state": "enabled",
					"labels": [{"name": "slow", "color": "#ff0000"}],
					"executions_count": 5,
					"executions_count_by_result": {
						"passed": 3,
						"failed": 2
					},
					"reliability": 60.0,
					"latest_fail": {
						"id": "exec-001",
						"timestamp": "2023-07-10T13:14:03.214Z",
						"duration": 1.234,
						"failure_reason": "Expected true got false",
						"failure_expanded": [
							{
								"expanded": ["line 1", "line 2"]
							}
						]
					}
				},
				{
					"id": "01867216-8478-7fde-a55a-0300f88bb49b",
					"scope": "User#email",
					"name": "TestExample1_Delete",
					"location": "./spec/models/text_example.rb:102",
					"state": "muted",
					"labels": [],
					"executions_count": 3,
					"executions_count_by_result": {
						"passed": 3,
						"failed": 0
					},
					"reliability": 100.0
				}
			]`)
	})

	buildTests, _, err := client.BuildTests.List(context.Background(), "my-great-org", "abc123", nil)
	if err != nil {
		t.Errorf("BuildTests.List returned error: %v", err)
	}

	parsedTime := must(time.Parse(BuildKiteDateFormat, "2023-07-10T13:14:03.214Z"))

	want := []BuildTest{
		{
			ID:              "a915535c-a8f1-4e1a-bd6a-a5589e09f349",
			Scope:           "User#email",
			Name:            "TestExample1_Create",
			Location:        "./spec/models/text_example.rb:55",
			State:           "enabled",
			Labels:          []TestLabel{{Name: "slow", Color: "#ff0000"}},
			ExecutionsCount: 5,
			ExecutionsCountByResult: BuildTestExecutionsCount{
				Passed: 3,
				Failed: 2,
			},
			Reliability: 60.0,
			LatestFail: &BuildTestLatestFail{
				ID:            "exec-001",
				Timestamp:     NewTimestamp(parsedTime),
				Duration:      1.234,
				FailureReason: "Expected true got false",
				FailureExpanded: []FailureExpanded{
					{Expanded: []string{"line 1", "line 2"}},
				},
			},
		},
		{
			ID:              "01867216-8478-7fde-a55a-0300f88bb49b",
			Scope:           "User#email",
			Name:            "TestExample1_Delete",
			Location:        "./spec/models/text_example.rb:102",
			State:           "muted",
			Labels:          []TestLabel{},
			ExecutionsCount: 3,
			ExecutionsCountByResult: BuildTestExecutionsCount{
				Passed: 3,
				Failed: 0,
			},
			Reliability: 100.0,
		},
	}

	if diff := cmp.Diff(buildTests, want); diff != "" {
		t.Errorf("BuildTests.List diff: (-got +want)\n%s", diff)
	}
}

func TestBuildTestsService_List_WithExecutionsIncluded(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/builds/abc123/tests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		if got, want := r.URL.Query().Get("result"), "^failed"; got != want {
			t.Errorf("query param result = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("state"), "enabled"; got != want {
			t.Errorf("query param state = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("include"), "executions"; got != want {
			t.Errorf("query param include = %q, want %q", got, want)
		}

		_, _ = fmt.Fprint(w, `[
			{
				"id": "a915535c-a8f1-4e1a-bd6a-a5589e09f349",
				"scope": "User#email",
				"name": "TestExample1_Create",
				"location": "./spec/models/text_example.rb:55",
				"state": "enabled",
				"executions_count": 1,
				"executions_count_by_result": {
					"passed": 0,
					"failed": 1
				},
				"reliability": 0.0,
				"executions": [
					{
						"id": "exec-001",
						"status": "failed",
						"timestamp": "2023-07-10T13:14:03.214Z",
						"duration": 1.234,
						"location": "./spec/models/text_example.rb:55",
						"failure_reason": "Expected true got false",
						"failure_expanded": [
							{
								"expanded": ["line 1", "line 2"]
							}
						]
					}
				]
			}
		]`)
	})

	opts := &BuildTestsListOptions{
		Result:  "^failed",
		State:   "enabled",
		Include: "executions",
	}

	parsedTime := must(time.Parse(BuildKiteDateFormat, "2023-07-10T13:14:03.214Z"))

	buildTests, _, err := client.BuildTests.List(context.Background(), "my-great-org", "abc123", opts)
	if err != nil {
		t.Errorf("BuildTests.List returned error: %v", err)
	}

	want := []BuildTest{
		{
			ID:              "a915535c-a8f1-4e1a-bd6a-a5589e09f349",
			Scope:           "User#email",
			Name:            "TestExample1_Create",
			Location:        "./spec/models/text_example.rb:55",
			State:           "enabled",
			ExecutionsCount: 1,
			ExecutionsCountByResult: BuildTestExecutionsCount{
				Passed: 0,
				Failed: 1,
			},
			Reliability: 0.0,
			Executions: []BuildTestExecution{
				{
					ID:            "exec-001",
					Status:        "failed",
					Timestamp:     NewTimestamp(parsedTime),
					Duration:      1.234,
					Location:      "./spec/models/text_example.rb:55",
					FailureReason: "Expected true got false",
					FailureExpanded: []FailureExpanded{
						{Expanded: []string{"line 1", "line 2"}},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(buildTests, want); diff != "" {
		t.Errorf("BuildTests.List diff: (-got +want)\n%s", diff)
	}
}
