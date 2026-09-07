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
		ExecutionID:              "01867216-8478-7fde-a55a-0300f88bb49b",
		RunID:                    "01867216-8478-7fde-a55a-0300f88bb49c",
		TestID:                   "01867216-8478-7fde-a55a-0300f88bb49d",
		TestName:                 "is correctly formatted",
		Location:                 "./spec/models/user_spec.rb:42",
		Result:                   "passed",
		Duration:                 &duration,
		CreatedAt:                createdAt,
		URL:                      "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-example/executions/01867216-8478-7fde-a55a-0300f88bb49b",
		WebURL:                   "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-example/executions/01867216-8478-7fde-a55a-0300f88bb49b",
		View:                     TraceViewFull,
		TraceID:                  "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanCount:                3,
		Truncated:                true,
		TraceDurationNanoseconds: 213000000,
		Categories: []TraceCategoryRollup{
			{
				Category:            TraceCategorySQL,
				SpanCount:           2,
				DurationNanoseconds: 120000000,
				SelfTimeNanoseconds: 110000000,
			},
			{
				Category:            TraceCategoryTest,
				SpanCount:           1,
				DurationNanoseconds: 213000000,
				SelfTimeNanoseconds: 93000000,
			},
		},
		SlowestSpans: []TraceSpan{
			{
				SpanID:              "00f067aa0ba902b7",
				ParentSpanID:        "0020000000000001",
				Name:                "SELECT users",
				Kind:                "client",
				ServiceName:         "web",
				ScopeName:           "opentelemetry-instrumentation-pg",
				Category:            TraceCategorySQL,
				Badge:               "SQL",
				Label:               "SELECT users",
				Detail:              "SELECT * FROM users WHERE id = $1",
				StartedAt:           sqlStartedAt,
				DurationNanoseconds: 100000000,
				SelfTimeNanoseconds: 100000000,
				Error:               true,
				StatusMessage:       "connection reset",
			},
		},
		Spans: []TraceSpan{
			{
				SpanID:              "0020000000000001",
				Name:                "User#email is correctly formatted",
				Category:            TraceCategoryTest,
				Badge:               "TEST",
				Label:               "is correctly formatted",
				StartedAt:           rootStartedAt,
				DurationNanoseconds: 213000000,
				SelfTimeNanoseconds: 93000000,
			},
			{
				SpanID:              "00f067aa0ba902b7",
				ParentSpanID:        "0020000000000001",
				Name:                "SELECT users",
				Category:            TraceCategorySQL,
				Badge:               "SQL",
				Label:               "SELECT users",
				StartedAt:           sqlStartedAt,
				DurationNanoseconds: 100000000,
				SelfTimeNanoseconds: 100000000,
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
		ExecutionID:              "01867216-8478-7fde-a55a-0300f88bb49b",
		View:                     TraceViewSummary,
		TraceID:                  "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanCount:                1,
		TraceDurationNanoseconds: 213000000,
		Categories: []TraceCategoryRollup{
			{
				Category:            TraceCategoryOther,
				SpanCount:           1,
				DurationNanoseconds: 213000000,
				SelfTimeNanoseconds: 213000000,
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
