package buildkite

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestExecutionsService_GetTrace(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/executions/01867216-8478-7fde-a55a-0300f88bb49b/trace", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"view": "full"})

		_, _ = fmt.Fprint(w, `{
			"execution_id": "01867216-8478-7fde-a55a-0300f88bb49b",
			"run_id": "01867216-8478-7fde-a55a-0300f88bb49c",
			"test_id": "01867216-8478-7fde-a55a-0300f88bb49d",
			"test_name": "is correctly formatted",
			"location": "./spec/models/user_spec.rb:42",
			"result": "passed",
			"duration": 0.213,
			"created_at": "2026-07-23T01:02:03.456Z",
			"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/executions/01867216-8478-7fde-a55a-0300f88bb49b",
			"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/executions/01867216-8478-7fde-a55a-0300f88bb49b",
			"view": "full",
			"trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
			"span_count": 3,
			"truncated": true,
			"trace_duration_ns": 213000000,
			"categories": [
				{
					"category": "sql",
					"span_count": 2,
					"duration_ns": 120000000,
					"self_time_ns": 110000000
				},
				{
					"category": "test",
					"span_count": 1,
					"duration_ns": 213000000,
					"self_time_ns": 93000000
				}
			],
			"slowest_spans": [
				{
					"span_id": "00f067aa0ba902b7",
					"parent_span_id": "0020000000000001",
					"name": "SELECT users",
					"kind": "client",
					"service_name": "web",
					"scope_name": "opentelemetry-instrumentation-pg",
					"category": "sql",
					"badge": "SQL",
					"label": "SELECT users",
					"detail": "SELECT * FROM users WHERE id = $1",
					"started_at": "2026-07-23T01:02:03.500Z",
					"duration_ns": 100000000,
					"self_time_ns": 100000000,
					"error": true,
					"status_message": "connection reset"
				}
			],
			"spans": [
				{
					"span_id": "0020000000000001",
					"name": "User#email is correctly formatted",
					"category": "test",
					"badge": "TEST",
					"label": "is correctly formatted",
					"started_at": "2026-07-23T01:02:03.456Z",
					"duration_ns": 213000000,
					"self_time_ns": 93000000
				},
				{
					"span_id": "00f067aa0ba902b7",
					"parent_span_id": "0020000000000001",
					"name": "SELECT users",
					"category": "sql",
					"badge": "SQL",
					"label": "SELECT users",
					"started_at": "2026-07-23T01:02:03.500Z",
					"duration_ns": 100000000,
					"self_time_ns": 100000000
				}
			]
		}`)
	})

	got, _, err := client.Executions.GetTrace(context.Background(), "my-great-org", "suite-example", "01867216-8478-7fde-a55a-0300f88bb49b", &ExecutionTraceOptions{View: TraceViewFull})
	if err != nil {
		t.Fatalf("ExecutionsService.GetTrace returned error: %v", err)
	}

	duration := 0.213
	createdAt := NewTimestamp(time.Date(2026, time.July, 23, 1, 2, 3, 456000000, time.UTC))
	rootStartedAt := NewTimestamp(time.Date(2026, time.July, 23, 1, 2, 3, 456000000, time.UTC))
	sqlStartedAt := NewTimestamp(time.Date(2026, time.July, 23, 1, 2, 3, 500000000, time.UTC))

	want := ExecutionTrace{
		ExecutionID:   "01867216-8478-7fde-a55a-0300f88bb49b",
		RunID:         "01867216-8478-7fde-a55a-0300f88bb49c",
		TestID:        "01867216-8478-7fde-a55a-0300f88bb49d",
		TestName:      "is correctly formatted",
		Location:      "./spec/models/user_spec.rb:42",
		Result:        "passed",
		Duration:      &duration,
		CreatedAt:     createdAt,
		URL:           "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/executions/01867216-8478-7fde-a55a-0300f88bb49b",
		WebURL:        "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/executions/01867216-8478-7fde-a55a-0300f88bb49b",
		View:          TraceViewFull,
		TraceID:       "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanCount:     3,
		Truncated:     true,
		TraceDuration: 213 * time.Millisecond,
		Categories: []TraceCategoryRollup{
			{
				Category:  TraceCategorySQL,
				SpanCount: 2,
				Duration:  120 * time.Millisecond,
				SelfTime:  110 * time.Millisecond,
			},
			{
				Category:  TraceCategoryTest,
				SpanCount: 1,
				Duration:  213 * time.Millisecond,
				SelfTime:  93 * time.Millisecond,
			},
		},
		SlowestSpans: []TraceSpan{
			{
				SpanID:        "00f067aa0ba902b7",
				ParentSpanID:  "0020000000000001",
				Name:          "SELECT users",
				Kind:          "client",
				ServiceName:   "web",
				ScopeName:     "opentelemetry-instrumentation-pg",
				Category:      TraceCategorySQL,
				Badge:         "SQL",
				Label:         "SELECT users",
				Detail:        "SELECT * FROM users WHERE id = $1",
				StartedAt:     sqlStartedAt,
				Duration:      100 * time.Millisecond,
				SelfTime:      100 * time.Millisecond,
				Error:         true,
				StatusMessage: "connection reset",
			},
		},
		Spans: []TraceSpan{
			{
				SpanID:    "0020000000000001",
				Name:      "User#email is correctly formatted",
				Category:  TraceCategoryTest,
				Badge:     "TEST",
				Label:     "is correctly formatted",
				StartedAt: rootStartedAt,
				Duration:  213 * time.Millisecond,
				SelfTime:  93 * time.Millisecond,
			},
			{
				SpanID:       "00f067aa0ba902b7",
				ParentSpanID: "0020000000000001",
				Name:         "SELECT users",
				Category:     TraceCategorySQL,
				Badge:        "SQL",
				Label:        "SELECT users",
				StartedAt:    sqlStartedAt,
				Duration:     100 * time.Millisecond,
				SelfTime:     100 * time.Millisecond,
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("ExecutionsService.GetTrace diff: (-got +want)\n%s", diff)
	}
}

func TestExecutionsService_GetTrace_SummaryView(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/executions/01867216-8478-7fde-a55a-0300f88bb49b/trace", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("request query = %q, want empty", got)
		}

		_, _ = fmt.Fprint(w, `{
			"execution_id": "01867216-8478-7fde-a55a-0300f88bb49b",
			"view": "summary",
			"trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
			"span_count": 1,
			"truncated": false,
			"trace_duration_ns": 213000000,
			"duration": null,
			"categories": [
				{
					"category": "other",
					"span_count": 1,
					"duration_ns": 213000000,
					"self_time_ns": 213000000
				}
			],
			"slowest_spans": [],
			"spans": null
		}`)
	})

	got, _, err := client.Executions.GetTrace(context.Background(), "my-great-org", "suite-example", "01867216-8478-7fde-a55a-0300f88bb49b", nil)
	if err != nil {
		t.Fatalf("ExecutionsService.GetTrace returned error: %v", err)
	}

	want := ExecutionTrace{
		ExecutionID:   "01867216-8478-7fde-a55a-0300f88bb49b",
		View:          TraceViewSummary,
		TraceID:       "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanCount:     1,
		TraceDuration: 213 * time.Millisecond,
		Categories: []TraceCategoryRollup{
			{
				Category:  TraceCategoryOther,
				SpanCount: 1,
				Duration:  213 * time.Millisecond,
				SelfTime:  213 * time.Millisecond,
			},
		},
		SlowestSpans: []TraceSpan{},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("ExecutionsService.GetTrace diff: (-got +want)\n%s", diff)
	}
	if got.Spans != nil {
		t.Errorf("ExecutionsService.GetTrace Spans = %v, want nil in the summary view", got.Spans)
	}
}

// An execution with no trace returns an empty trace rather than a 404.
func TestExecutionsService_GetTrace_NoTrace(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/executions/01867216-8478-7fde-a55a-0300f88bb49b/trace", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
			"execution_id": "01867216-8478-7fde-a55a-0300f88bb49b",
			"view": "full",
			"trace_id": null,
			"span_count": 0,
			"truncated": false,
			"trace_duration_ns": 0,
			"categories": [],
			"slowest_spans": [],
			"spans": []
		}`)
	})

	got, _, err := client.Executions.GetTrace(context.Background(), "my-great-org", "suite-example", "01867216-8478-7fde-a55a-0300f88bb49b", &ExecutionTraceOptions{View: TraceViewFull})
	if err != nil {
		t.Fatalf("ExecutionsService.GetTrace returned error: %v", err)
	}

	if got.TraceID != "" {
		t.Errorf("ExecutionsService.GetTrace TraceID = %q, want empty", got.TraceID)
	}
	if got.SpanCount != 0 {
		t.Errorf("ExecutionsService.GetTrace SpanCount = %d, want 0", got.SpanCount)
	}
	if got.Spans == nil {
		t.Error("ExecutionsService.GetTrace Spans = nil, want an empty slice in the full view")
	}
	if len(got.Spans) != 0 {
		t.Errorf("ExecutionsService.GetTrace returned %d spans, want 0", len(got.Spans))
	}
}

func TestExecutionsService_GetTrace_NotFound(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/executions/01867216-8478-7fde-a55a-0300f88bb49b/trace", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"message": "Not Found"}`)
	})

	_, resp, err := client.Executions.GetTrace(context.Background(), "my-great-org", "suite-example", "01867216-8478-7fde-a55a-0300f88bb49b", nil)
	if err == nil {
		t.Fatal("ExecutionsService.GetTrace returned no error, want one")
	}

	var errResp *ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("ExecutionsService.GetTrace error = %T, want *ErrorResponse", err)
	}
	if got, want := resp.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("response.StatusCode = %d, want %d", got, want)
	}
}

func TestExecutionsService_ListSlowestByBuild(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc(fmt.Sprintf("/v2/analytics/organizations/my-great-org/builds/%s/executions/slowest", testBuildUUID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"limit": "2",
		})

		_, _ = fmt.Fprint(w,
			`
			[
				{
					"id": "019d66fc-1a2b-7c3d-8e4f-5a6b7c8d9e0f",
					"suite_slug": "suite-example",
					"test_id": "a915535c-a8f1-4e1a-bd6a-a5589e09f349",
					"duration": 12.345,
					"has_trace": true
				},
				{
					"id": "019d66fc-0000-7c3d-8e4f-5a6b7c8d9e0f",
					"suite_slug": "other-suite",
					"test_id": "b0e3a5b8-2b7c-4a4e-9d9e-1f2a3b4c5d6e",
					"duration": 2.5,
					"has_trace": false
				}
			]`)
	})

	got, _, err := client.Executions.ListSlowestByBuild(context.Background(), "my-great-org", testBuildUUID, &SlowestExecutionsOptions{Limit: 2})
	if err != nil {
		t.Fatalf("ExecutionsService.ListSlowestByBuild returned error: %v", err)
	}

	want := []BuildExecution{
		{
			ID:        "019d66fc-1a2b-7c3d-8e4f-5a6b7c8d9e0f",
			SuiteSlug: "suite-example",
			TestID:    "a915535c-a8f1-4e1a-bd6a-a5589e09f349",
			Duration:  12.345,
			HasTrace:  true,
		},
		{
			ID:        "019d66fc-0000-7c3d-8e4f-5a6b7c8d9e0f",
			SuiteSlug: "other-suite",
			TestID:    "b0e3a5b8-2b7c-4a4e-9d9e-1f2a3b4c5d6e",
			Duration:  2.5,
			HasTrace:  false,
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("ExecutionsService.ListSlowestByBuild diff: (-got +want)\n%s", diff)
	}
}

func TestExecutionsService_ListSlowestByBuild_ZeroLimitOmitted(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc(fmt.Sprintf("/v2/analytics/organizations/my-great-org/builds/%s/executions/slowest", testBuildUUID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		_, _ = fmt.Fprint(w, `[]`)
	})

	got, _, err := client.Executions.ListSlowestByBuild(context.Background(), "my-great-org", testBuildUUID, &SlowestExecutionsOptions{})
	if err != nil {
		t.Fatalf("ExecutionsService.ListSlowestByBuild returned error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("ExecutionsService.ListSlowestByBuild returned %v, want no executions", got)
	}
}

func TestExecutionsService_ListSlowestByBuild_ServerError(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc(fmt.Sprintf("/v2/analytics/organizations/my-great-org/builds/%s/executions/slowest", testBuildUUID), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprint(w, `{"message": "limit must be an integer between 1 and 100"}`)
	})

	got, resp, err := client.Executions.ListSlowestByBuild(context.Background(), "my-great-org", testBuildUUID, &SlowestExecutionsOptions{Limit: 500})
	if err == nil {
		t.Fatal("ExecutionsService.ListSlowestByBuild returned nil error, want API error")
	}
	if got != nil {
		t.Errorf("ExecutionsService.ListSlowestByBuild returned items %v, want nil", got)
	}
	if resp == nil || resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("ExecutionsService.ListSlowestByBuild response = %#v, want status %d", resp, http.StatusUnprocessableEntity)
	}
	var apiErr *ErrorResponse
	if !errors.As(err, &apiErr) {
		t.Errorf("ExecutionsService.ListSlowestByBuild error type = %T, want *ErrorResponse", err)
	} else if got, want := apiErr.Message, "limit must be an integer between 1 and 100"; got != want {
		t.Errorf("ExecutionsService.ListSlowestByBuild error message = %q, want %q", got, want)
	}
}
