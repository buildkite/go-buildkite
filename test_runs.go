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

// FailureExpanded is the expanded failure detail attached to a failed
// execution. The OpenAPI description leaves the shape open, so it is hand
// written and wired in through the codegen overlay.
type FailureExpanded struct {
	Backtrace []string `json:"backtrace,omitempty"`
	Expanded  []string `json:"expanded,omitempty"`
}

type TestRunsListOptions struct {
	ListOptions
}

type FailedExecutionsOptions struct {
	IncludeFailureExpanded bool `url:"include_failure_expanded,omitempty"`
	Page                   int  `url:"page,omitempty"`
	PerPage                int  `url:"per_page,omitempty"`
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
