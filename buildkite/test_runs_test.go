package buildkite

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestTestRunsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`
			[
				{
					"id": "3c90a8ad-8e86-4e78-87b4-acae5e808de4",
					"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
					"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
					"branch": "main",
					"commit_sha": "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
					"created_at": "2023-05-20T10:25:50.264Z"
				},
				{
					"id": "70fe7e45-c9e4-446b-95e3-c50d61519b87",
					"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/70fe7e45-c9e4-446b-95e3-c50d61519b87",
					"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/70fe7e45-c9e4-446b-95e3-c50d61519b87",
					"branch": "main",
					"commit_sha": "109f4b3c50d7b0df729d299bc6f8e9ef9066971f",
					"created_at": "2023-05-20T10:52:22.254Z"
				}
			]`)
	})

	runs, _, err := client.TestRuns.List("my-great-org", "suite-example", nil)

	if err != nil {
		t.Errorf("TestSuites.List returned error: %v", err)
	}

	// Create Time instances from strings in BuildKiteDateFormat friendly format
	parsedTime1, err := time.Parse(BuildKiteDateFormat, "2023-05-20T10:25:50.264Z")
	parsedTime2, err := time.Parse(BuildKiteDateFormat, "2023-05-20T10:52:22.254Z")

	if err != nil {
		t.Errorf("TestSuites.List time.Parse error: %v", err)
	}

	want := []TestRun{
		{
			ID:        String("3c90a8ad-8e86-4e78-87b4-acae5e808de4"),
			URL:       String("https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4"),
			WebURL:    String("https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4"),
			Branch:    String("main"),
			CommitSHA: String("a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"),
			CreatedAt: NewTimestamp(parsedTime1),
		},
		{
			ID:        String("70fe7e45-c9e4-446b-95e3-c50d61519b87"),
			URL:       String("https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/70fe7e45-c9e4-446b-95e3-c50d61519b87"),
			WebURL:    String("https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/70fe7e45-c9e4-446b-95e3-c50d61519b87"),
			Branch:    String("main"),
			CommitSHA: String("109f4b3c50d7b0df729d299bc6f8e9ef9066971f"),
			CreatedAt: NewTimestamp(parsedTime2),
		},
	}

	if !reflect.DeepEqual(runs, want) {
		t.Errorf("TestRuns.List returned %+v, want %+v", runs, want)
	}
}

func TestTestRunsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`
			{
				"id": "3c90a8ad-8e86-4e78-87b4-acae5e808de4",
				"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
				"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4",
				"branch": "main",
				"commit_sha": "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
				"created_at": "2023-05-20T10:25:50.264Z"
			}`)
	})

	run, _, err := client.TestRuns.Get("my-great-org", "suite-example", "3c90a8ad-8e86-4e78-87b4-acae5e808de4")

	if err != nil {
		t.Errorf("TestSuites.Get returned error: %v", err)
	}

	// Create Time instance from string in BuildKiteDateFormat friendly format
	parsedTime, err := time.Parse(BuildKiteDateFormat, "2023-05-20T10:25:50.264Z")

	if err != nil {
		t.Errorf("TestSuites.Get time.Parse error: %v", err)
	}

	want := &TestRun{

		ID:        String("3c90a8ad-8e86-4e78-87b4-acae5e808de4"),
		URL:       String("https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4"),
		WebURL:    String("https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/runs/3c90a8ad-8e86-4e78-87b4-acae5e808de4"),
		Branch:    String("main"),
		CommitSHA: String("a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"),
		CreatedAt: NewTimestamp(parsedTime),
	}

	if !reflect.DeepEqual(run, want) {
		t.Errorf("TestRuns.Get returned %+v, want %+v", run, want)
	}
}
