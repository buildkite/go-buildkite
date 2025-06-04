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
