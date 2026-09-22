package buildkite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTestsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	const linkHeader = `<https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests?page=3>; rel="next", <https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests?page=4>; rel="last"`

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got, want := r.Header.Get("Buildkite-Version"), testsMetricsAPIVersion; got != want {
			t.Errorf("Buildkite-Version header = %q, want %q", got, want)
		}
		testFormValues(t, r, values{
			"page":           "2",
			"per_page":       "50",
			"min_timestamp":  "2026-07-01T00:00:00Z",
			"max_timestamp":  "2026-07-23T00:00:00Z",
			"min_executions": "10",
			"labels":         "flaky,!slow",
			"branch":         "main*",
			"owners":         "payments,!platform",
			"state":          "enabled",
			"tags":           "framework:rspec,result:^failed",
			"sort_by":        "reliability",
			"order":          "asc",
		})

		w.Header().Set("Link", linkHeader)
		_, _ = fmt.Fprint(w, `[
			{
				"id": "01867216-8478-7fde-a55a-0300f88bb49b",
				"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests/01867216-8478-7fde-a55a-0300f88bb49b",
				"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/tests/01867216-8478-7fde-a55a-0300f88bb49b",
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
	got, resp, err := client.Tests.List(context.Background(), "my-great-org", "suite-example", &TestsListOptions{
		ListOptions: ListOptions{Page: 2, PerPage: 50},
		MinTimestamp: time.Date(
			2026, time.July, 1, 0, 0, 0, 0, time.UTC,
		),
		MaxTimestamp:  time.Date(2026, time.July, 23, 0, 0, 0, 0, time.UTC),
		MinExecutions: 10,
		Labels:        "flaky,!slow",
		Branch:        "main*",
		Owners:        "payments,!platform",
		State:         "enabled",
		Tags:          "framework:rspec,result:^failed",
		SortBy:        "reliability",
		Order:         "asc",
	})
	if err != nil {
		t.Fatalf("TestsService.List returned error: %v", err)
	}

	want := []TestWithMetrics{
		{
			Test: Test{
				ID:       "01867216-8478-7fde-a55a-0300f88bb49b",
				URL:      "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests/01867216-8478-7fde-a55a-0300f88bb49b",
				WebURL:   "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/tests/01867216-8478-7fde-a55a-0300f88bb49b",
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
		t.Errorf("TestsService.List diff: (-got +want)\n%s", diff)
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

func TestTestsService_List_Period(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests", func(w http.ResponseWriter, r *http.Request) {
		testFormValues(t, r, values{"period": "28days"})
		_, _ = fmt.Fprint(w, `[]`)
	})

	got, _, err := client.Tests.List(context.Background(), "my-great-org", "suite-example", &TestsListOptions{Period: "28days"})
	if err != nil {
		t.Fatalf("TestsService.List returned error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("TestsService.List returned %d items, want 0", len(got))
	}
}

func TestTestsService_List_NilOptionsAndReliability(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("request query = %q, want empty", got)
		}
		if got, want := r.Header.Get("Buildkite-Version"), testsMetricsAPIVersion; got != want {
			t.Errorf("Buildkite-Version header = %q, want %q", got, want)
		}
		_, _ = fmt.Fprint(w, `[
			{
				"id": "01867216-8478-7fde-a55a-0300f88bb49b",
				"reliability": null,
				"executions_count": 1,
				"executions_count_by_result": {
					"passed": 0,
					"failed": 0,
					"skipped": 1
				}
			}
		]`)
	})

	got, _, err := client.Tests.List(context.Background(), "my-great-org", "suite-example", nil)
	if err != nil {
		t.Fatalf("TestsService.List returned error: %v", err)
	}

	want := []TestWithMetrics{
		{
			Test:            Test{ID: "01867216-8478-7fde-a55a-0300f88bb49b"},
			Reliability:     nil,
			ExecutionsCount: 1,
			ExecutionsCountByResult: map[string]int{
				"passed":  0,
				"failed":  0,
				"skipped": 1,
			},
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestsService.List diff: (-got +want)\n%s", diff)
	}
}

func TestTestsService_List_ServerError(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprint(w, "{\"message\":\"period cannot be combined with min_timestamp\"}")
	})

	got, resp, err := client.Tests.List(context.Background(), "my-great-org", "suite-example", nil)
	if err == nil {
		t.Fatal("TestsService.List returned nil error, want API error")
	}
	if got != nil {
		t.Errorf("TestsService.List returned items %v, want nil", got)
	}
	if resp == nil || resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("TestsService.List response = %#v, want status %d", resp, http.StatusUnprocessableEntity)
	}
	var apiErr *ErrorResponse
	if !errors.As(err, &apiErr) {
		t.Errorf("TestsService.List error type = %T, want *ErrorResponse", err)
	} else if got, want := apiErr.Message, "period cannot be combined with min_timestamp"; got != want {
		t.Errorf("TestsService.List error message = %q, want %q", got, want)
	}
}

func TestTestsService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got, want := r.Header.Get("Buildkite-Version"), testsMetricsAPIVersion; got != want {
			t.Errorf("Buildkite-Version header = %q, want %q", got, want)
		}
		testFormValues(t, r, values{
			"min_timestamp": "2026-07-01T00:00:00Z",
			"max_timestamp": "2026-07-23T00:00:00Z",
		})

		_, _ = fmt.Fprint(w,
			`
			{
				"id": "b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"name": "TestExample1_Create",
				"scope": "User#email",
				"location": "./resources/test_example_test.go:123",
				"file_name": "./resources/test_example_test.go",
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
			}`)
	})

	got, _, err := client.Tests.Get(context.Background(), "my-great-org", "suite-example", "b3abe2e9-35c5-4905-85e1-8c9f2da3240f", &TestsGetOptions{
		MinTimestamp: time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
		MaxTimestamp: time.Date(2026, time.July, 23, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Errorf("TestsService.Get returned error: %v", err)
	}

	reliability := 0.9821
	want := TestWithMetrics{
		Test: Test{
			ID:       "b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
			URL:      "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
			WebURL:   "https://buildkite.com/organizations/my-great-org/analytics/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
			Name:     "TestExample1_Create",
			Scope:    "User#email",
			Location: "./resources/test_example_test.go:123",
			FileName: "./resources/test_example_test.go",
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
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestsService.Get diff: (-got +want)\n%s", diff)
	}
}

func TestTestsService_Get_NilOptionsAndPeriod(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests/nil-options", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("request query = %q, want empty", got)
		}
		_, _ = fmt.Fprint(w, `{"id": "nil-options", "reliability": null, "executions_count": 0}`)
	})
	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests/period", func(w http.ResponseWriter, r *http.Request) {
		testFormValues(t, r, values{"period": "28days"})
		_, _ = fmt.Fprint(w, `{"id": "period", "reliability": null, "executions_count": 0}`)
	})

	got, _, err := client.Tests.Get(context.Background(), "my-great-org", "suite-example", "nil-options", nil)
	if err != nil {
		t.Fatalf("TestsService.Get returned error: %v", err)
	}
	if diff := cmp.Diff(got, TestWithMetrics{Test: Test{ID: "nil-options"}}); diff != "" {
		t.Errorf("TestsService.Get diff: (-got +want)\n%s", diff)
	}

	got, _, err = client.Tests.Get(context.Background(), "my-great-org", "suite-example", "period", &TestsGetOptions{Period: "28days"})
	if err != nil {
		t.Fatalf("TestsService.Get returned error: %v", err)
	}
	if diff := cmp.Diff(got, TestWithMetrics{Test: Test{ID: "period"}}); diff != "" {
		t.Errorf("TestsService.Get diff: (-got +want)\n%s", diff)
	}
}

func TestTestsService_Find(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests/find", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body FindTestOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("TestsService.Find failed to decode request body: %v", err)
		}
		want := FindTestOptions{Scope: "User#email", Name: "TestExample1_Create"}
		if diff := cmp.Diff(body, want); diff != "" {
			t.Errorf("TestsService.Find request body diff: (-got +want)\n%s", diff)
		}

		_, _ = fmt.Fprint(w,
			`
			{
				"id": "b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"name": "TestExample1_Create",
				"scope": "User#email",
				"location": "./resources/test_example_test.go:123",
				"file_name": "./resources/test_example_test.go"
			}`)
	})

	got, _, err := client.Tests.Find(context.Background(), "my-great-org", "suite-example",
		FindTestOptions{Scope: "User#email", Name: "TestExample1_Create"})
	if err != nil {
		t.Errorf("TestsService.Find returned error: %v", err)
	}

	want := Test{
		ID:       "b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
		URL:      "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
		WebURL:   "https://buildkite.com/organizations/my-great-org/analytics/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
		Name:     "TestExample1_Create",
		Scope:    "User#email",
		Location: "./resources/test_example_test.go:123",
		FileName: "./resources/test_example_test.go",
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestsService.Find diff: (-got +want)\n%s", diff)
	}
}

func TestTestsService_ChangeState(t *testing.T) {
	t.Parallel()

	const testID = "b3abe2e9-35c5-4905-85e1-8c9f2da3240f"

	tests := []struct {
		name string
		call func(*Client) (Test, *Response, error)
		path string
	}{
		{
			name: "skip",
			call: func(c *Client) (Test, *Response, error) {
				return c.Tests.Skip(context.Background(), "my-great-org", "suite-example", testID)
			},
			path: "/v2/analytics/organizations/my-great-org/suites/suite-example/tests/" + testID + "/skip",
		},
		{
			name: "mute",
			call: func(c *Client) (Test, *Response, error) {
				return c.Tests.Mute(context.Background(), "my-great-org", "suite-example", testID)
			},
			path: "/v2/analytics/organizations/my-great-org/suites/suite-example/tests/" + testID + "/mute",
		},
		{
			name: "enable",
			call: func(c *Client) (Test, *Response, error) {
				return c.Tests.Enable(context.Background(), "my-great-org", "suite-example", testID)
			},
			path: "/v2/analytics/organizations/my-great-org/suites/suite-example/tests/" + testID + "/enable",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			server.HandleFunc(tc.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "PUT")

				// The quarantine endpoints take no request body; sending one
				// would be a change the API never asked for.
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("reading request body: %v", err)
				}
				if len(body) != 0 {
					t.Errorf("request body = %q, want empty", body)
				}
				if got := r.Header.Get("Content-Type"); got != "" {
					t.Errorf("Content-Type header = %q, want empty", got)
				}

				_, _ = fmt.Fprint(w, `
			{
				"id": "b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"scope": "Flaky test",
				"name": "passes only on the second try on BK CI",
				"location": "./spec/flaky_spec.rb:6",
				"file_name": "./spec/flaky_spec.rb",
				"labels": ["flaky"]
			}`)
			})

			got, resp, err := tc.call(client)
			if err != nil {
				t.Fatalf("TestsService.%s returned error: %v", tc.name, err)
			}

			want := Test{
				ID:       testID,
				URL:      "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				WebURL:   "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				Scope:    "Flaky test",
				Name:     "passes only on the second try on BK CI",
				Location: "./spec/flaky_spec.rb:6",
				FileName: "./spec/flaky_spec.rb",
				Labels:   []string{"flaky"},
			}
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("TestsService.%s diff: (-got +want)\n%s", tc.name, diff)
			}
			if resp == nil || resp.StatusCode != http.StatusOK {
				t.Errorf("TestsService.%s response = %#v, want status %d", tc.name, resp, http.StatusOK)
			}
		})
	}
}

// A suite without test state management enabled 404s the quarantine endpoints,
// which must surface as an error rather than an empty test.
func TestTestsService_ChangeState_StateManagementDisabled(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f/mute", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"message":"Not Found"}`)
	})

	got, resp, err := client.Tests.Mute(context.Background(), "my-great-org", "suite-example", "b3abe2e9-35c5-4905-85e1-8c9f2da3240f")
	if err == nil {
		t.Fatal("TestsService.Mute returned nil error, want API error")
	}
	if diff := cmp.Diff(got, Test{}); diff != "" {
		t.Errorf("TestsService.Mute returned test diff: (-got +want)\n%s", diff)
	}
	if resp == nil || resp.StatusCode != http.StatusNotFound {
		t.Errorf("TestsService.Mute response = %#v, want status %d", resp, http.StatusNotFound)
	}
	var apiErr *ErrorResponse
	if !errors.As(err, &apiErr) {
		t.Errorf("TestsService.Mute error type = %T, want *ErrorResponse", err)
	}
}

func TestTestsService_ListByState(t *testing.T) {
	t.Parallel()

	const linkHeader = `<https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests/muted?page=3>; rel="next", <https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests/muted?page=4>; rel="last"`

	tests := []struct {
		name string
		call func(*Client, *QuarantinedTestsListOptions) ([]Test, *Response, error)
		path string
	}{
		{
			name: "muted",
			call: func(c *Client, opt *QuarantinedTestsListOptions) ([]Test, *Response, error) {
				return c.Tests.ListMuted(context.Background(), "my-great-org", "suite-example", opt)
			},
			path: "/v2/analytics/organizations/my-great-org/suites/suite-example/tests/muted",
		},
		{
			name: "skipped",
			call: func(c *Client, opt *QuarantinedTestsListOptions) ([]Test, *Response, error) {
				return c.Tests.ListSkipped(context.Background(), "my-great-org", "suite-example", opt)
			},
			path: "/v2/analytics/organizations/my-great-org/suites/suite-example/tests/skipped",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			server.HandleFunc(tc.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testFormValues(t, r, values{"page": "2", "per_page": "50"})

				w.Header().Set("Link", linkHeader)
				_, _ = fmt.Fprint(w, `[
			{
				"id": "160988e4-836e-88ab-af45-22170a169e23",
				"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests/160988e4-836e-88ab-af45-22170a169e23",
				"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/tests/160988e4-836e-88ab-af45-22170a169e23",
				"scope": "Flaky test",
				"name": "passes only on the second try on BK CI",
				"location": "flaky.spec.js:1",
				"file_name": "flaky.spec.js",
				"labels": []
			}
		]`)
			})

			got, resp, err := tc.call(client, &QuarantinedTestsListOptions{
				ListOptions: ListOptions{Page: 2, PerPage: 50},
			})
			if err != nil {
				t.Fatalf("TestsService.%s returned error: %v", tc.name, err)
			}

			want := []Test{
				{
					ID:       "160988e4-836e-88ab-af45-22170a169e23",
					URL:      "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/tests/160988e4-836e-88ab-af45-22170a169e23",
					WebURL:   "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/tests/160988e4-836e-88ab-af45-22170a169e23",
					Scope:    "Flaky test",
					Name:     "passes only on the second try on BK CI",
					Location: "flaky.spec.js:1",
					FileName: "flaky.spec.js",
					Labels:   []string{},
				},
			}
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("TestsService.%s diff: (-got +want)\n%s", tc.name, diff)
			}
			if got, want := resp.NextPage, 3; got != want {
				t.Errorf("response.NextPage = %d, want %d", got, want)
			}
			if got, want := resp.LastPage, 4; got != want {
				t.Errorf("response.LastPage = %d, want %d", got, want)
			}
		})
	}
}

func TestTestsService_ListByState_NilOptions(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests/skipped", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("request query = %q, want empty", got)
		}
		_, _ = fmt.Fprint(w, `[]`)
	})

	got, _, err := client.Tests.ListSkipped(context.Background(), "my-great-org", "suite-example", nil)
	if err != nil {
		t.Fatalf("TestsService.ListSkipped returned error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("TestsService.ListSkipped returned %d items, want 0", len(got))
	}
}

func TestTestsService_ListByState_ServerError(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests/muted", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"message":"Not Found"}`)
	})

	got, resp, err := client.Tests.ListMuted(context.Background(), "my-great-org", "suite-example", nil)
	if err == nil {
		t.Fatal("TestsService.ListMuted returned nil error, want API error")
	}
	if got != nil {
		t.Errorf("TestsService.ListMuted returned items %v, want nil", got)
	}
	if resp == nil || resp.StatusCode != http.StatusNotFound {
		t.Errorf("TestsService.ListMuted response = %#v, want status %d", resp, http.StatusNotFound)
	}
}

// A test ID or suite slug carrying path metacharacters must not be able to
// redirect the request to a different endpoint. Unescaped, a trailing "#"
// truncates the path at that point: a Skip reaches the Enable endpoint, a List
// reaches the suite endpoint, and a Get reaches a different test.
//
// Escaping makes the value inert, not verbatim: resolveURL decodes "%2F" back
// to a separator, so an escaped slash arrives as extra path segments rather
// than as one literal segment. Either way no route matches and the API 404s,
// which is why the expectations below are the paths this client really sends.
func TestTestsService_PathInjection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		call     func(*Client) error
		wantPath string
		body     string // the shape the method under test decodes
	}{
		{
			name: "test id escaping an action onto the path",
			call: func(c *Client) error {
				_, _, err := c.Tests.Skip(context.Background(), "my-great-org", "suite-example", "b3abe2e9/enable#")
				return err
			},
			// The escaped slash decodes back to a separator, but the "#" stays
			// escaped and the action stays "skip" — no route matches, so the API
			// 404s instead of enabling the test.
			wantPath: "/v2/analytics/organizations/my-great-org/suites/suite-example/tests/b3abe2e9/enable%23/skip",
			body:     "{}",
		},
		{
			name: "slug truncating a list onto the suite endpoint",
			call: func(c *Client) error {
				_, _, err := c.Tests.List(context.Background(), "my-great-org", "suite-example#", nil)
				return err
			},
			wantPath: "/v2/analytics/organizations/my-great-org/suites/suite-example%23/tests",
			body:     "[]",
		},
		{
			name: "test id truncating a get onto a different test",
			call: func(c *Client) error {
				_, _, err := c.Tests.Get(context.Background(), "my-great-org", "suite-example", "b3abe2e9#", nil)
				return err
			},
			wantPath: "/v2/analytics/organizations/my-great-org/suites/suite-example/tests/b3abe2e9%23",
			body:     "{}",
		},
		{
			name: "slug truncating a find onto the suite endpoint",
			call: func(c *Client) error {
				_, _, err := c.Tests.Find(context.Background(), "my-great-org", "suite-example#", FindTestOptions{Scope: "s", Name: "n"})
				return err
			},
			wantPath: "/v2/analytics/organizations/my-great-org/suites/suite-example%23/tests/find",
			body:     "{}",
		},
		{
			name: "percent-encoded separator in a test id",
			call: func(c *Client) error {
				_, _, err := c.Tests.Mute(context.Background(), "my-great-org", "suite-example", "a%2Fb")
				return err
			},
			wantPath: "/v2/analytics/organizations/my-great-org/suites/suite-example/tests/a%252Fb/mute",
			body:     "{}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var gotPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.EscapedPath()
				_, _ = fmt.Fprint(w, tc.body)
			}))
			t.Cleanup(server.Close)

			client, err := NewOpts(WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("unexpected NewOpts() error: %v", err)
			}
			if err := tc.call(client); err != nil {
				t.Fatalf("call returned error: %v", err)
			}
			if gotPath != tc.wantPath {
				t.Errorf("request path = %q, want %q", gotPath, tc.wantPath)
			}
		})
	}
}
