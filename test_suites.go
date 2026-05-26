package buildkite

import (
	"context"
	"fmt"
)

// TestSuitesService handles communication with the test suite related
// methods of the Buildkite Test Analytics API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/analytics/suites
type TestSuitesService struct {
	client *Client
}

// TestSuiteCreate represents the request body for creating a test suite.
type TestSuiteCreate struct {
	Name          string   `json:"name"`
	DefaultBranch string   `json:"default_branch,omitempty"`
	ShowAPIToken  bool     `json:"show_api_token"`
	TeamUUIDs     []string `json:"team_ids,omitempty"`
}

// TestSuite represents a Buildkite Test Analytics suite.
type TestSuite struct {
	ID            string `json:"id,omitempty"`
	GraphQLID     string `json:"graphql_id,omitempty"`
	Slug          string `json:"slug,omitempty"`
	Name          string `json:"name,omitempty"`
	URL           string `json:"url,omitempty"`
	WebURL        string `json:"web_url,omitempty"`
	DefaultBranch string `json:"default_branch,omitempty"`
}

// TestSuiteUpdate represents the request body for updating a test suite.
type TestSuiteUpdate struct {
	Name          Optional[string] `json:"name,omitzero"`
	DefaultBranch Optional[string] `json:"default_branch,omitzero"`
}

type TestSuiteListOptions struct{ ListOptions }

func (tss *TestSuitesService) List(ctx context.Context, org string, opt *TestSuiteListOptions) ([]TestSuite, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := tss.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var testSuites []TestSuite
	resp, err := tss.client.Do(req, &testSuites)
	if err != nil {
		return nil, resp, err
	}

	return testSuites, resp, err
}

func (tss *TestSuitesService) Get(ctx context.Context, org, slug string) (TestSuite, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s", org, slug)
	req, err := tss.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TestSuite{}, nil, err
	}

	var testSuite TestSuite
	resp, err := tss.client.Do(req, &testSuite)
	if err != nil {
		return TestSuite{}, resp, err
	}

	return testSuite, resp, err
}

func (tss *TestSuitesService) Create(ctx context.Context, org string, ts TestSuiteCreate) (TestSuite, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites", org)
	req, err := tss.client.NewRequest(ctx, "POST", u, ts)
	if err != nil {
		return TestSuite{}, nil, err
	}

	var testSuite TestSuite
	resp, err := tss.client.Do(req, &testSuite)
	if err != nil {
		return TestSuite{}, resp, err
	}

	return testSuite, resp, err
}

func (tss *TestSuitesService) Update(ctx context.Context, org, slug string, ts TestSuiteUpdate) (TestSuite, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s", org, slug)
	req, err := tss.client.NewRequest(ctx, "PATCH", u, ts)
	if err != nil {
		return TestSuite{}, nil, err
	}

	var out TestSuite
	resp, err := tss.client.Do(req, &out)
	if err != nil {
		return TestSuite{}, resp, err
	}

	return out, resp, err
}

func (tss *TestSuitesService) Delete(ctx context.Context, org, slug string) (*Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s", org, slug)
	req, err := tss.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return tss.client.Do(req, nil)
}
