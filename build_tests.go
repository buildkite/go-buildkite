package buildkite

import (
	"context"
	"fmt"
)

// BuildTestsService handles communication with the build test related
// methods of the Buildkite Test Analytics API.
//
// List expects a Buildkite build UUID, not a pipeline build number.
type BuildTestsService struct {
	client *Client
}

// BuildTestsListOptions specifies optional parameters for
// [BuildTestsService.List]. Invalid values and incompatible combinations are
// validated by the API.
type BuildTestsListOptions struct {
	ListOptions

	// Labels filters by comma-separated test labels. Prefix a label with "!" to
	// exclude it, for example "flaky,!slow".
	Labels string `url:"labels,omitempty"`

	// Branch filters the executions included in the metrics. Prefix the value
	// with "!" to exclude an exact branch, or suffix it with "*" to match
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
	// with at least one execution in the build having that result, or "^" to
	// return tests for which every execution has that result. A build.id filter
	// cannot override the buildUUID argument; for example
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

// List returns tests with execution metrics aggregated over the build's time
// window. Pagination is available through the returned Response.
func (bts *BuildTestsService) List(ctx context.Context, org, buildUUID string, opt *BuildTestsListOptions) ([]TestWithMetrics, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/builds/%s/tests", org, buildUUID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := bts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var buildTests []TestWithMetrics
	resp, err := bts.client.Do(req, &buildTests)
	if err != nil {
		return nil, resp, err
	}

	return buildTests, resp, err
}
