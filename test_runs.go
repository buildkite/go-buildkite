package buildkite

import (
	"context"
	"fmt"
)

// TestRunsService handles communication with test run related
// methods of the Buildkite Test Analytics API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/analytics/runs
type TestRunsService struct {
	client *Client
}

type TestRun struct {
	ID        string     `json:"id,omitempty"`
	URL       string     `json:"url,omitempty"`
	WebURL    string     `json:"web_url,omitempty"`
	Branch    string     `json:"branch,omitempty"`
	CommitSHA string     `json:"commit_sha,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

type FailureExpanded struct {
	Backtrace []string `json:"backtrace,omitempty"`
	Expanded  []string `json:"expanded,omitempty"`
}

type FailedExecution struct {
	ExecutionID      string            `json:"execution_id,omitempty"`
	RunID            string            `json:"run_id,omitempty"`
	TestID           string            `json:"test_id,omitempty"`
	RunName          string            `json:"run_name,omitempty"`
	CommitSHA        string            `json:"commit_sha,omitempty"`
	CreatedAt        *Timestamp        `json:"created_at,omitempty"`
	Branch           string            `json:"branch,omitempty"`
	FailureReason    string            `json:"failure_reason,omitempty"`
	FailureExpanded  []FailureExpanded `json:"failure_expanded,omitempty"`
	Duration         float64           `json:"duration,omitempty"`
	Location         string            `json:"location,omitempty"`
	TestName         string            `json:"test_name,omitempty"`
	RunURL           string            `json:"run_url,omitempty"`
	TestURL          string            `json:"test_url,omitempty"`
	TestExecutionURL string            `json:"test_execution_url,omitempty"`
}

type TestRunsListOptions struct {
	ListOptions
}

type FailedExecutionsOptions struct {
	IncludeFailureExpanded bool `url:"include_failure_expanded,omitempty"`
}

func (trs *TestRunsService) List(ctx context.Context, org, slug string, opt *TestRunsListOptions) ([]TestRun, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/runs", org, slug)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := trs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var testRuns []TestRun
	resp, err := trs.client.Do(req, &testRuns)
	if err != nil {
		return nil, resp, err
	}

	return testRuns, resp, err
}

func (trs *TestRunsService) Get(ctx context.Context, org, slug, runID string) (TestRun, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/runs/%s", org, slug, runID)
	req, err := trs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TestRun{}, nil, err
	}

	var testRun TestRun
	resp, err := trs.client.Do(req, &testRun)
	if err != nil {
		return TestRun{}, resp, err
	}

	return testRun, resp, err
}

func (trs *TestRunsService) GetFailedExecutions(ctx context.Context, org, slug, runID string, opt *FailedExecutionsOptions) ([]FailedExecution, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/runs/%s/failed_executions", org, slug, runID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := trs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var failedExecutions []FailedExecution
	resp, err := trs.client.Do(req, &failedExecutions)
	if err != nil {
		return nil, resp, err
	}

	return failedExecutions, resp, err
}
