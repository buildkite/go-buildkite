package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTestRunsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			[
				{
					"id": "3c90a8ad-8e86-4e78-87b4-acae5e808de4",
					"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
					"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
					"branch": "main",
					"commit_sha": "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
					"created_at": "2023-05-20T10:25:50.264Z",
					"state": "finished",
					"result": "passed"
				},
				{
					"id": "70fe7e45-c9e4-446b-95e3-c50d61519b87",
					"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/70fe7e45-c9e4-446b-95e3-c50d61519b87",
					"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/70fe7e45-c9e4-446b-95e3-c50d61519b87",
					"branch": "main",
					"commit_sha": "109f4b3c50d7b0df729d299bc6f8e9ef9066971f",
					"created_at": "2023-05-20T10:52:22.254Z",
					"state": "running",
					"result": "pending"
				}
			]`)
	})

	runs, _, err := client.TestRuns.List(context.Background(), "my-great-org", "suite-example", nil)
	if err != nil {
		t.Errorf("TestSuites.List returned error: %v", err)
	}

	// Create Time instances from strings in BuildKiteDateFormat friendly format
	parsedTime1 := must(time.Parse(BuildKiteDateFormat, "2023-05-20T10:25:50.264Z"))
	parsedTime2 := must(time.Parse(BuildKiteDateFormat, "2023-05-20T10:52:22.254Z"))

	if err != nil {
		t.Errorf("TestSuites.List time.Parse error: %v", err)
	}

	want := []TestRun{
		{
			ID:        "3c90a8ad-8e86-4e78-87b4-acae5e808de4",
			URL:       "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
			WebURL:    "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
			Branch:    "main",
			CommitSHA: "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
			CreatedAt: NewTimestamp(parsedTime1),
			State:     "finished",
			Result:    "passed",
		},
		{
			ID:        "70fe7e45-c9e4-446b-95e3-c50d61519b87",
			URL:       "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/70fe7e45-c9e4-446b-95e3-c50d61519b87",
			WebURL:    "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/70fe7e45-c9e4-446b-95e3-c50d61519b87",
			Branch:    "main",
			CommitSHA: "109f4b3c50d7b0df729d299bc6f8e9ef9066971f",
			CreatedAt: NewTimestamp(parsedTime2),
			State:     "running",
			Result:    "pending",
		},
	}

	if diff := cmp.Diff(runs, want); diff != "" {
		t.Errorf("TestRuns.List diff: (-got +want)\n%s", diff)
	}
}

func TestTestRunsService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			{
				"id": "3c90a8ad-8e86-4e78-87b4-acae5e808de4",
				"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
				"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
				"branch": "main",
				"commit_sha": "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
				"created_at": "2023-05-20T10:25:50.264Z",
				"state": "finished",
				"result": "failed"
			}`)
	})

	run, _, err := client.TestRuns.Get(context.Background(), "my-great-org", "suite-example", "3c90a8ad-8e86-4e78-87b4-acae5e808de4")
	if err != nil {
		t.Errorf("TestSuites.Get returned error: %v", err)
	}

	// Create Time instance from string in BuildKiteDateFormat friendly format
	parsedTime, err := time.Parse(BuildKiteDateFormat, "2023-05-20T10:25:50.264Z")
	if err != nil {
		t.Errorf("TestSuites.Get time.Parse error: %v", err)
	}

	want := TestRun{
		ID:        "3c90a8ad-8e86-4e78-87b4-acae5e808de4",
		URL:       "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
		WebURL:    "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
		Branch:    "main",
		CommitSHA: "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
		CreatedAt: NewTimestamp(parsedTime),
		State:     "finished",
		Result:    "failed",
	}

	if diff := cmp.Diff(run, want); diff != "" {
		t.Errorf("TestRuns.Get diff: (-got +want)\n%s", diff)
	}
}

func TestTestRunsService_GetFailedExecutions(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4/failed_executions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			[
				{
					"execution_id": "60f0e64c-ae4b-870e-b41f-5431205caf06",
					"run_id": "075bbcd9-662c-86f5-9d40-adfa6549eff1",
					"test_id": "f6cb6c43-df94-8b60-81ed-14f9db7bbfd8",
					"run_name": "075bbcd9-662c-86f5-9d40-adfa6549eff1",
					"commit_sha": "1c3214fcceb2c14579a2c3c50cd78f1442fd8936",
					"created_at": "2025-02-03T05:32:53.228Z",
					"branch": "main",
					"failure_reason": "it didn't work",
					"duration": 3.79073,
					"location": "./spec/models/user.rb:23",
					"test_name": "Deploy should be available",
					"run_url": "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/runs/075bbcd9-662c-86f5-9d40-adfa6549eff1",
					"test_url": "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/tests/f6cb6c43-df94-8b60-81ed-14f9db7bbfd8",
					"test_execution_url": "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/tests/f6cb6c43-df94-8b60-81ed-14f9db7bbfd8?execution_id=60f0e64c-ae4b-870e-b41f-5431205caf06"
				}
			]`)
	})

	failedExecutions, _, err := client.TestRuns.GetFailedExecutions(context.Background(), "my-great-org", "suite-example", "3c90a8ad-8e86-4e78-87b4-acae5e808de4", nil)
	if err != nil {
		t.Errorf("TestRuns.GetFailedExecutions returned error: %v", err)
	}

	// Create Time instance from string in BuildKiteDateFormat friendly format
	parsedTime, err := time.Parse(BuildKiteDateFormat, "2025-02-03T05:32:53.228Z")
	if err != nil {
		t.Errorf("TestRuns.GetFailedExecutions time.Parse error: %v", err)
	}

	want := []FailedExecution{
		{
			ExecutionID:      "60f0e64c-ae4b-870e-b41f-5431205caf06",
			RunID:            "075bbcd9-662c-86f5-9d40-adfa6549eff1",
			TestID:           "f6cb6c43-df94-8b60-81ed-14f9db7bbfd8",
			RunName:          "075bbcd9-662c-86f5-9d40-adfa6549eff1",
			CommitSHA:        "1c3214fcceb2c14579a2c3c50cd78f1442fd8936",
			CreatedAt:        NewTimestamp(parsedTime),
			Branch:           "main",
			FailureReason:    "it didn't work",
			Duration:         3.79073,
			Location:         "./spec/models/user.rb:23",
			TestName:         "Deploy should be available",
			RunURL:           "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/runs/075bbcd9-662c-86f5-9d40-adfa6549eff1",
			TestURL:          "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/tests/f6cb6c43-df94-8b60-81ed-14f9db7bbfd8",
			TestExecutionURL: "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/tests/f6cb6c43-df94-8b60-81ed-14f9db7bbfd8?execution_id=60f0e64c-ae4b-870e-b41f-5431205caf06",
		},
	}

	if diff := cmp.Diff(failedExecutions, want); diff != "" {
		t.Errorf("TestRuns.GetFailedExecutions diff: (-got +want)\n%s", diff)
	}
}

func TestTestRunsService_GetFailedExecutions_WithFailureExpanded(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4/failed_executions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		// Check that include_failure_expanded parameter is included
		if r.URL.Query().Get("include_failure_expanded") != "true" {
			t.Errorf("Expected include_failure_expanded=true in query parameters, got: %s", r.URL.Query().Get("include_failure_expanded"))
		}

		_, _ = fmt.Fprint(w,
			`
			[
				{
					"execution_id": "60f0e64c-ae4b-870e-b41f-5431205caf06",
					"run_id": "075bbcd9-662c-86f5-9d40-adfa6549eff1",
					"test_id": "f6cb6c43-df94-8b60-81ed-14f9db7bbfd8",
					"run_name": "075bbcd9-662c-86f5-9d40-adfa6549eff1",
					"commit_sha": "1c3214fcceb2c14579a2c3c50cd78f1442fd8936",
					"created_at": "2025-02-03T05:32:53.228Z",
					"branch": "main",
					"failure_reason": "it didn't work",
					"failure_expanded": [
						{
							"backtrace": [
								"./spec/models/user.rb:23:in 'block (2 levels) in <top (required)>'"
							],
							"expanded": [
								"RuntimeError:",
								"it didn't work"
							]
						}
					],
					"duration": 3.79073,
					"location": "./spec/models/user.rb:23",
					"test_name": "Deploy should be available",
					"run_url": "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/runs/075bbcd9-662c-86f5-9d40-adfa6549eff1",
					"test_url": "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/tests/f6cb6c43-df94-8b60-81ed-14f9db7bbfd8",
					"test_execution_url": "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/tests/f6cb6c43-df94-8b60-81ed-14f9db7bbfd8?execution_id=60f0e64c-ae4b-870e-b41f-5431205caf06"
				}
			]`)
	})

	options := &FailedExecutionsOptions{
		IncludeFailureExpanded: true,
	}

	failedExecutions, _, err := client.TestRuns.GetFailedExecutions(context.Background(), "my-great-org", "suite-example", "3c90a8ad-8e86-4e78-87b4-acae5e808de4", options)
	if err != nil {
		t.Errorf("TestRuns.GetFailedExecutions returned error: %v", err)
	}

	// Create Time instance from string in BuildKiteDateFormat friendly format
	parsedTime, err := time.Parse(BuildKiteDateFormat, "2025-02-03T05:32:53.228Z")
	if err != nil {
		t.Errorf("TestRuns.GetFailedExecutions time.Parse error: %v", err)
	}

	want := []FailedExecution{
		{
			ExecutionID:   "60f0e64c-ae4b-870e-b41f-5431205caf06",
			RunID:         "075bbcd9-662c-86f5-9d40-adfa6549eff1",
			TestID:        "f6cb6c43-df94-8b60-81ed-14f9db7bbfd8",
			RunName:       "075bbcd9-662c-86f5-9d40-adfa6549eff1",
			CommitSHA:     "1c3214fcceb2c14579a2c3c50cd78f1442fd8936",
			CreatedAt:     NewTimestamp(parsedTime),
			Branch:        "main",
			FailureReason: "it didn't work",
			FailureExpanded: []FailureExpanded{
				{
					Backtrace: []string{"./spec/models/user.rb:23:in 'block (2 levels) in <top (required)>'"},
					Expanded:  []string{"RuntimeError:", "it didn't work"},
				},
			},
			Duration:         3.79073,
			Location:         "./spec/models/user.rb:23",
			TestName:         "Deploy should be available",
			RunURL:           "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/runs/075bbcd9-662c-86f5-9d40-adfa6549eff1",
			TestURL:          "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/tests/f6cb6c43-df94-8b60-81ed-14f9db7bbfd8",
			TestExecutionURL: "https://buildkite.com/organizations/buildkite/analytics/suites/my-test-suite/tests/f6cb6c43-df94-8b60-81ed-14f9db7bbfd8?execution_id=60f0e64c-ae4b-870e-b41f-5431205caf06",
		},
	}

	if diff := cmp.Diff(failedExecutions, want); diff != "" {
		t.Errorf("TestRuns.GetFailedExecutions diff: (-got +want)\n%s", diff)
	}
}
