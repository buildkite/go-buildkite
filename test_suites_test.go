package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTestSuitesService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			[
				{
					"id": "7c202aaa-3165-4811-9813-173c4c285463",
					"graphql_id": "N2MyMDJhYWEtMzE2NS00ODExLTk4MTMtMTczYzRjMjg1NDYz=",
					"slug": "suite-1",
					"name": "suite-1",
					"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-1",
					"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-1",
					"default_branch": "main"
				},
				{
					"id": "38ed1d73-cea9-4aba-b223-def25e66ef51",
					"graphql_id": "MzhlZDFkNzMtY2VhOS00YWJhLWIyMjMtZGVmMjVlNjZlZjUx=",
					"slug": "suite-2",
					"name": "suite-2",
					"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-2",
					"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-2",
					"default_branch": "main"
				}
			]`)
	})

	suites, _, err := client.TestSuites.List(context.Background(), "my-great-org", nil)

	if err != nil {
		t.Errorf("TestSuites.List returned error: %v", err)
	}

	want := []TestSuite{
		{
			ID:            "7c202aaa-3165-4811-9813-173c4c285463",
			GraphQLID:     "N2MyMDJhYWEtMzE2NS00ODExLTk4MTMtMTczYzRjMjg1NDYz=",
			Slug:          "suite-1",
			Name:          "suite-1",
			URL:           "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-1",
			WebURL:        "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-1",
			DefaultBranch: "main",
		},
		{
			ID:            "38ed1d73-cea9-4aba-b223-def25e66ef51",
			GraphQLID:     "MzhlZDFkNzMtY2VhOS00YWJhLWIyMjMtZGVmMjVlNjZlZjUx=",
			Slug:          "suite-2",
			Name:          "suite-2",
			URL:           "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-2",
			WebURL:        "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-2",
			DefaultBranch: "main",
		},
	}
	if diff := cmp.Diff(suites, want); diff != "" {
		t.Errorf("TestSuites.List diff: (-got +want)\n%s", diff)
	}
}

func TestTestSuitesService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			{
				"id": "7c202aaa-3165-4811-9813-173c4c285463",
				"graphql_id": "N2MyMDJhYWEtMzE2NS00ODExLTk4MTMtMTczYzRjMjg1NDYz=",
				"slug": "suite-1",
				"name": "suite-1",
				"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-1",
				"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-1",
				"default_branch": "main"
			}`)
	})

	suite, _, err := client.TestSuites.Get(context.Background(), "my-great-org", "suite-1")

	if err != nil {
		t.Errorf("TestSuites.Get returned error: %v", err)
	}

	want := TestSuite{
		ID:            "7c202aaa-3165-4811-9813-173c4c285463",
		GraphQLID:     "N2MyMDJhYWEtMzE2NS00ODExLTk4MTMtMTczYzRjMjg1NDYz=",
		Slug:          "suite-1",
		Name:          "suite-1",
		URL:           "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-1",
		WebURL:        "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-1",
		DefaultBranch: "main",
	}

	if diff := cmp.Diff(suite, want); diff != "" {
		t.Errorf("TestSuites.Get diff: (-got +want)\n%s", diff)
	}
}

func TestTestSuitesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := TestSuiteCreate{
		Name:          "Suite 3",
		DefaultBranch: "main",
		TeamUUIDs:     []string{"8369b300-fff0-4ef1-91de-010f72f4458d"},
	}

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites", func(w http.ResponseWriter, r *http.Request) {
		var v TestSuiteCreate
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		_, _ = fmt.Fprint(w,
			`
			{
				"name" : "Suite 3",
				"default_branch": "main",
				"team_ids": ["8369b300-fff0-4ef1-91de-010f72f4458d"]
			}`)
	})

	suite, _, err := client.TestSuites.Create(context.Background(), "my-great-org", input)

	if err != nil {
		t.Errorf("TestSuites.Create returned error: %v", err)
	}

	want := TestSuite{
		Name:          "Suite 3",
		DefaultBranch: "main",
	}

	if diff := cmp.Diff(suite, want); diff != "" {
		t.Errorf("TestSuites.Create diff: (-got +want)\n%s", diff)
	}
}

func TestTestSuitesService_Update(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := TestSuiteCreate{
		Name:          "Suite 4",
		DefaultBranch: "main",
		TeamUUIDs:     []string{"818b0849-9718-4898-8de3-42d591a7fe26"},
	}

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites", func(w http.ResponseWriter, r *http.Request) {
		var v TestSuiteCreate
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		_, _ = fmt.Fprint(w,
			`
			{
				"name" : "Suite 4",
				"default_branch": "main",
				"team_ids": ["818b0849-9718-4898-8de3-42d591a7fe26"],
				"slug": "suite-4"
			}`)
	})

	suite, _, err := client.TestSuites.Create(context.Background(), "my-great-org", input)
	if err != nil {
		t.Errorf("TestSuites.Create returned error: %v", err)
	}

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-4", func(w http.ResponseWriter, r *http.Request) {
		var v TestSuiteCreate
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "PATCH")

		_, _ = fmt.Fprint(w,
			`
			{
				"name" : "Suite 4",
				"default_branch": "develop",
				"team_ids": ["818b0849-9718-4898-8de3-42d591a7fe26"],
				"slug": "suite-4"
			}`)
	})

	got, _, err := client.TestSuites.Update(context.Background(), "my-great-org", suite.Slug, TestSuite{DefaultBranch: "default"})
	if err != nil {
		t.Errorf("Pipelines.Update returned error: %v", err)
	}

	want := TestSuite{
		Name:          "Suite 4",
		Slug:          "suite-4",
		DefaultBranch: "develop",
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestSuites.Update diff: (-got +want)\n%s", diff)
	}
}

func TestTestSuitesService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-5", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.TestSuites.Delete(context.Background(), "my-great-org", "suite-5")

	if err != nil {
		t.Errorf("TestSuites.Delete returned error: %v", err)
	}
}
