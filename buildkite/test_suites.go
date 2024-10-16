package buildkite

import (
	"context"
	"errors"
	"fmt"
)

// TestSuitesService handles communication with the test suite related
// methods of the Buildkite Test Analytics API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/analytics/suites
type TestSuitesService struct {
	client *Client
}

type TestSuiteCreate struct {
	Name          string   `json:"name"`
	DefaultBranch string   `json:"default_branch,omitempty"`
	ShowAPIToken  bool     `json:"show_api_token,omitempty"`
	TeamUUIDs     []string `json:"team_ids,omitempty"`
}

type TestSuite struct {
	ID            *string `json:"id,omitempty"`
	GraphQLID     *string `json:"graphql_id,omitempty"`
	Slug          *string `json:"slug,omitempty"`
	Name          *string `json:"name,omitempty"`
	URL           *string `json:"url,omitempty"`
	WebURL        *string `json:"web_url,omitempty"`
	DefaultBranch *string `json:"default_branch,omitempty"`
}

type TestSuiteListOptions struct {
	ListOptions
}

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

	testSuites := new([]TestSuite)

	resp, err := tss.client.Do(req, testSuites)

	if err != nil {
		return nil, resp, err
	}

	return *testSuites, resp, err
}

func (tss *TestSuitesService) Get(ctx context.Context, org, slug string) (*TestSuite, *Response, error) {

	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s", org, slug)

	req, err := tss.client.NewRequest(ctx, "GET", u, nil)

	if err != nil {
		return nil, nil, err
	}

	testSuite := new(TestSuite)

	resp, err := tss.client.Do(req, testSuite)

	if err != nil {
		return nil, resp, err
	}

	return testSuite, resp, err
}

func (tss *TestSuitesService) Create(ctx context.Context, org string, ts *TestSuiteCreate) (*TestSuite, *Response, error) {

	u := fmt.Sprintf("v2/analytics/organizations/%s/suites", org)

	req, err := tss.client.NewRequest(ctx, "POST", u, ts)

	if err != nil {
		return nil, nil, err
	}

	testSuite := new(TestSuite)
	resp, err := tss.client.Do(req, testSuite)

	if err != nil {
		return nil, resp, err
	}

	return testSuite, resp, err
}

func (tss *TestSuitesService) Update(ctx context.Context, org string, ts *TestSuite) (*Response, error) {

	if ts == nil {
		return nil, errors.New("Test suite must not be nil")
	}

	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s", org, *ts.Slug)

	req, err := tss.client.NewRequest(ctx, "PATCH", u, ts)

	if err != nil {
		return nil, nil
	}

	resp, err := tss.client.Do(req, ts)

	if err != nil {
		return resp, err
	}

	return resp, err
}

func (tss *TestSuitesService) Delete(ctx context.Context, org, slug string) (*Response, error) {

	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s", org, slug)

	req, err := tss.client.NewRequest(ctx, "DELETE", u, nil)

	if err != nil {
		return nil, err
	}

	return tss.client.Do(req, nil)
}
