package buildkite

import (
	"context"
	"fmt"
	"time"
)

// TestsService handles communication with test related
// methods of the Buildkite Test Analytics API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/analytics/tests
type TestsService struct {
	client *Client
}

// testsMetricsAPIVersion opts requests in to the versioned response that
// includes test execution metrics.
const testsMetricsAPIVersion = "2026-08-01"

// Test represents a test in a Buildkite Test Engine suite.
type Test struct {
	ID       string   `json:"id,omitempty"`
	URL      string   `json:"url,omitempty"`
	WebURL   string   `json:"web_url,omitempty"`
	Scope    string   `json:"scope,omitempty"`
	Name     string   `json:"name,omitempty"`
	Location string   `json:"location,omitempty"`
	FileName string   `json:"file_name,omitempty"`
	Labels   []string `json:"labels,omitempty"`
}

// TestWithMetrics represents a test and its execution metrics aggregated over
// the time window used by [TestsService.List], [TestsService.Get], or
// [BuildTestsService.List].
type TestWithMetrics struct {
	Test

	// Reliability is the ratio of passed executions to passed and failed
	// executions. It is nil when there are no passed or failed executions.
	Reliability *float64 `json:"reliability"`

	// DurationAverage is the average execution duration in seconds.
	DurationAverage float64 `json:"duration_avg"`

	// DurationTotal is the total execution duration in seconds.
	DurationTotal float64 `json:"duration_sum"`

	// DurationMinimum is the shortest execution duration in seconds.
	DurationMinimum float64 `json:"duration_min"`

	// DurationMaximum is the longest execution duration in seconds.
	DurationMaximum float64 `json:"duration_max"`

	// ExecutionsCount is the number of executions in the aggregation window.
	ExecutionsCount int `json:"executions_count"`

	// ExecutionsCountByResult contains execution counts grouped by result.
	// Passed and failed are always present; skipped, pending, and unknown are
	// present only when non-zero.
	ExecutionsCountByResult map[string]int `json:"executions_count_by_result"`
}

// TestsListOptions specifies optional parameters for [TestsService.List].
// Invalid values and incompatible combinations are validated by the API.
type TestsListOptions struct {
	ListOptions

	// Period sets a relative aggregation window, such as "7days" or "28days".
	// Available periods depend on the organization's maximum time window. Period
	// cannot be combined with MinTimestamp or MaxTimestamp.
	Period string `url:"period,omitempty"`

	// MinTimestamp sets the start of the aggregation window. When omitted, it
	// defaults to the organization's default period before the current time.
	MinTimestamp time.Time `url:"min_timestamp,omitempty"`

	// MaxTimestamp sets the end of the aggregation window. It defaults to the
	// current time when omitted.
	MaxTimestamp time.Time `url:"max_timestamp,omitempty"`

	// Labels filters by comma-separated test labels. Prefix a label with "!" to
	// exclude it, for example "flaky,!slow".
	Labels string `url:"labels,omitempty"`

	// Branch filters the executions included in the aggregation. Prefix the
	// value with "!" to exclude an exact branch, or suffix it with "*" to match
	// branches by prefix; for example "!main" or "feature*". Use at most one
	// operator.
	Branch string `url:"branch,omitempty"`

	// Owners filters by comma-separated test owner slugs. Prefix an owner with
	// "!" to exclude it, for example "payments,!platform".
	Owners string `url:"owners,omitempty"`

	// State filters by test state: "enabled", "muted", or "skipped".
	State string `url:"state,omitempty"`

	// Tags filters by comma-separated execution tags in key:value form. Prefix a
	// value with "!" to exclude an exact value, or suffix it with "*" to match
	// by prefix. For the result tag, prefix the value with "~" to return tests
	// with at least one execution in the aggregation window having that result,
	// or "^" to return tests for which every execution in the window has that
	// result. Use at most one operator per value and specify result at most once;
	// for example
	// "framework:!rspec,scm.branch:feature*,result:^passed".
	Tags string `url:"tags,omitempty"`

	// SortBy specifies the metric used to sort results: "duration_avg",
	// "duration_sum", "duration_min", "duration_max", or "reliability". It
	// defaults to "duration_avg".
	SortBy string `url:"sort_by,omitempty"`

	// Order specifies the sort direction: "asc" or "desc". It defaults to
	// "desc".
	Order string `url:"order,omitempty"`
}

type FindTestOptions struct {
	Scope string `json:"scope"`
	Name  string `json:"name"`
}

// List returns tests with execution metrics aggregated over the requested time
// window. It opts in to the versioned metrics response, which includes only
// tests with executions in that window. Pagination is available through the
// returned Response.
func (ts *TestsService) List(ctx context.Context, org, slug string, opt *TestsListOptions) ([]TestWithMetrics, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/tests", org, slug)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := ts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Buildkite-Version", testsMetricsAPIVersion)

	var tests []TestWithMetrics
	resp, err := ts.client.Do(req, &tests)
	if err != nil {
		return nil, resp, err
	}

	return tests, resp, err
}

// TestsGetOptions specifies optional parameters for [TestsService.Get].
// Invalid values and incompatible combinations are validated by the API.
type TestsGetOptions struct {
	// Period sets a relative aggregation window, such as "7days" or "28days".
	// Available periods depend on the organization's maximum time window. Period
	// cannot be combined with MinTimestamp or MaxTimestamp.
	Period string `url:"period,omitempty"`

	// MinTimestamp sets the start of the aggregation window. When omitted, it
	// defaults to the organization's default period before the current time.
	MinTimestamp time.Time `url:"min_timestamp,omitempty"`

	// MaxTimestamp sets the end of the aggregation window. It defaults to the
	// current time when omitted.
	MaxTimestamp time.Time `url:"max_timestamp,omitempty"`
}

// Get returns a single test with its execution metrics aggregated over the
// requested time window. It opts in to the versioned metrics response.
func (ts *TestsService) Get(ctx context.Context, org, slug, testID string, opt *TestsGetOptions) (TestWithMetrics, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/tests/%s", org, slug, testID)
	u, err := addOptions(u, opt)
	if err != nil {
		return TestWithMetrics{}, nil, err
	}

	req, err := ts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TestWithMetrics{}, nil, err
	}
	req.Header.Set("Buildkite-Version", testsMetricsAPIVersion)

	var t TestWithMetrics
	resp, err := ts.client.Do(req, &t)
	if err != nil {
		return TestWithMetrics{}, resp, err
	}

	return t, resp, err
}

func (ts *TestsService) Find(ctx context.Context, org, slug string, find FindTestOptions) (Test, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/tests/find", org, slug)
	req, err := ts.client.NewRequest(ctx, "POST", u, find)
	if err != nil {
		return Test{}, nil, err
	}
	var t Test
	resp, err := ts.client.Do(req, &t)
	if err != nil {
		return Test{}, resp, err
	}
	return t, resp, err
}
