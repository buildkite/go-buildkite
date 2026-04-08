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

type BuildTest struct {
	ID                      string                   `json:"id,omitempty"`
	Scope                   string                   `json:"scope,omitempty"`
	Name                    string                   `json:"name,omitempty"`
	Location                string                   `json:"location,omitempty"`
	State                   string                   `json:"state,omitempty"`
	Labels                  []TestLabel              `json:"labels,omitempty"`
	ExecutionsCount         int                      `json:"executions_count,omitempty"`
	ExecutionsCountByResult BuildTestExecutionsCount `json:"executions_count_by_result"`
	Reliability             float64                  `json:"reliability"`
	// Executions is only populated when Include is set to "executions".
	Executions []BuildTestExecution `json:"executions,omitempty"`
	LatestFail *BuildTestLatestFail `json:"latest_fail,omitempty"`
}

type TestLabel struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

type BuildTestExecutionsCount struct {
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped,omitempty"`
	Pending int `json:"pending,omitempty"`
	Unknown int `json:"unknown,omitempty"`
}

type BuildTestLatestFail struct {
	ID              string            `json:"id,omitempty"`
	Timestamp       *Timestamp        `json:"timestamp,omitempty"`
	Duration        float64           `json:"duration"`
	FailureReason   string            `json:"failure_reason,omitempty"`
	FailureExpanded []FailureExpanded `json:"failure_expanded,omitempty"`
}

type BuildTestExecution struct {
	ID              string            `json:"id,omitempty"`
	Status          string            `json:"status,omitempty"`
	Timestamp       *Timestamp        `json:"timestamp,omitempty"`
	Duration        float64           `json:"duration"`
	Location        string            `json:"location,omitempty"`
	FailureReason   string            `json:"failure_reason,omitempty"`
	FailureExpanded []FailureExpanded `json:"failure_expanded,omitempty"`
}

type BuildTestsListOptions struct {
	ListOptions

	// Result filters by execution result.
	// "failed" = has any failed execution, "^failed" = all executions failed.
	// Same for "passed" / "^passed".
	Result string `url:"result,omitempty"`

	// State filters by test state: "enabled", "muted", etc.
	State string `url:"state,omitempty"`

	// Include set to "latest_fail" inlines the most recent failed execution per test.
	// Include set to "executions" inlines the executions matching the current filter.
	Include string `url:"include,omitempty"`
}

func (bts *BuildTestsService) List(ctx context.Context, org, buildUUID string, opt *BuildTestsListOptions) ([]BuildTest, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/builds/%s/tests", org, buildUUID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := bts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var buildTests []BuildTest
	resp, err := bts.client.Do(req, &buildTests)
	if err != nil {
		return nil, resp, err
	}

	return buildTests, resp, err
}
