package buildkite

import (
	"context"
	"fmt"
	"time"
)

// ExecutionsService handles communication with test execution related
// methods of the Buildkite Test Analytics API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/analytics/executions
type ExecutionsService struct {
	client *Client
}

// TraceCategory is the kind of work a span did, derived from its OpenTelemetry
// semantic-convention attributes and, failing those, from the instrumentation
// library that emitted it.
type TraceCategory string

const (
	TraceCategorySQL       TraceCategory = "sql"
	TraceCategoryDB        TraceCategory = "db"
	TraceCategoryCache     TraceCategory = "cache"
	TraceCategoryHTTP      TraceCategory = "http"
	TraceCategoryRPC       TraceCategory = "rpc"
	TraceCategoryMessaging TraceCategory = "messaging"
	TraceCategoryTest      TraceCategory = "test"
	TraceCategoryOther     TraceCategory = "other"
)

// TraceView selects how much of a trace [ExecutionsService.GetTrace] returns.
type TraceView string

const (
	// TraceViewSummary returns the category rollup and the slowest spans alone.
	TraceViewSummary TraceView = "summary"

	// TraceViewFull additionally returns the span list.
	TraceViewFull TraceView = "full"
)

// TraceCategoryRollup is one category's share of a trace, across every span in
// it.
type TraceCategoryRollup struct {
	Category  TraceCategory `json:"category,omitempty"`
	SpanCount int           `json:"span_count,omitempty"`

	// Duration is the total duration of the category's spans.
	Duration time.Duration `json:"duration_ns,omitempty"`

	// SelfTime is the total self time of the category's spans.
	SelfTime time.Duration `json:"self_time_ns,omitempty"`
}

// TraceSpan is a single span in a trace. Name, Kind, ServiceName, ScopeName,
// Badge, Label, Detail and StatusMessage come from customer instrumentation and
// are truncated to 2000 characters by the API.
type TraceSpan struct {
	SpanID string `json:"span_id,omitempty"`

	// ParentSpanID is empty for a span that anchors the tree.
	ParentSpanID string `json:"parent_span_id,omitempty"`

	Name        string `json:"name,omitempty"`
	Kind        string `json:"kind,omitempty"`
	ServiceName string `json:"service_name,omitempty"`

	// ScopeName is the OpenTelemetry instrumentation scope name.
	ScopeName string `json:"scope_name,omitempty"`

	Category TraceCategory `json:"category,omitempty"`

	// Badge is a short uppercase tag naming the work, such as "SQL" or "HTTP".
	Badge string `json:"badge,omitempty"`

	// Label describes what the span did, with no repeat of Badge, such as
	// "SELECT users".
	Label string `json:"label,omitempty"`

	// Detail is supporting text, such as the query statement or full URL.
	Detail string `json:"detail,omitempty"`

	StartedAt *Timestamp `json:"started_at,omitempty"`

	Duration time.Duration `json:"duration_ns,omitempty"`

	// SelfTime is the span's own time: its duration less what its recorded
	// children spent, clamped at zero for concurrent children that outlast it.
	SelfTime time.Duration `json:"self_time_ns,omitempty"`

	Error         bool   `json:"error,omitempty"`
	StatusMessage string `json:"status_message,omitempty"`
}

// ExecutionTrace is the OpenTelemetry trace recorded for one test execution.
type ExecutionTrace struct {
	ExecutionID string `json:"execution_id,omitempty"`
	RunID       string `json:"run_id,omitempty"`
	TestID      string `json:"test_id,omitempty"`
	TestName    string `json:"test_name,omitempty"`
	Location    string `json:"location,omitempty"`
	Result      string `json:"result,omitempty"`

	// Duration is the execution duration in seconds. It is nil when the
	// execution recorded no duration.
	Duration *float64 `json:"duration"`

	CreatedAt *Timestamp `json:"created_at,omitempty"`
	URL       string     `json:"url,omitempty"`
	WebURL    string     `json:"web_url,omitempty"`

	// View is the view the trace was returned in.
	View TraceView `json:"view,omitempty"`

	// TraceID is empty when the execution recorded no trace.
	TraceID string `json:"trace_id,omitempty"`

	// SpanCount counts every span in the trace, not only the ones returned.
	SpanCount int `json:"span_count,omitempty"`

	// Truncated is true when the payload does not account for the whole trace:
	// the rollup hit its own ceiling, so the totals are lower bounds, or the
	// span list was clipped.
	Truncated bool `json:"truncated,omitempty"`

	// TraceDuration spans the earliest span start to the latest span end, rather
	// than the root span's own duration - a root can finish before a child it
	// never awaited.
	TraceDuration time.Duration `json:"trace_duration_ns,omitempty"`

	// Categories is the category rollup over every span, ordered by self time.
	Categories []TraceCategoryRollup `json:"categories,omitempty"`

	// SlowestSpans holds the trace's slowest spans by self time, ranked across
	// the whole trace even when the span list is truncated.
	SlowestSpans []TraceSpan `json:"slowest_spans,omitempty"`

	// Spans is the trace's leading window, in start time order. It is nil in the
	// summary view, so a caller can tell a trace whose spans it did not ask for
	// from one that has none.
	Spans []TraceSpan `json:"spans,omitempty"`
}

// ExecutionTraceOptions specifies optional parameters for
// [ExecutionsService.GetTrace].
type ExecutionTraceOptions struct {
	// View selects how much of the trace is returned. It defaults to
	// [TraceViewSummary], which omits the span list.
	View TraceView `url:"view,omitempty"`
}

// GetTrace returns the OpenTelemetry trace recorded for a single test
// execution: where the execution's time went, rolled up by category and ranked
// by the spans that spent it.
//
// Both the ranking and the category rollup use self time - a span's own duration
// less the time its children accounted for - so an outer span is never credited
// with work it merely waited on.
//
// Spans are retained for a shorter window than executions, and an execution
// recorded before the organization exported traces never had one, so an
// execution with no trace returns an empty trace rather than an error.
func (es *ExecutionsService) GetTrace(ctx context.Context, org, slug, executionID string, opt *ExecutionTraceOptions) (ExecutionTrace, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/executions/%s/trace", org, slug, executionID)
	u, err := addOptions(u, opt)
	if err != nil {
		return ExecutionTrace{}, nil, err
	}

	req, err := es.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return ExecutionTrace{}, nil, err
	}

	var trace ExecutionTrace
	resp, err := es.client.Do(req, &trace)
	if err != nil {
		return ExecutionTrace{}, resp, err
	}

	return trace, resp, err
}
