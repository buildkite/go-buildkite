package buildkite

import (
	"fmt"
)

// TestSuitesService handles communication with the test suite related
// methods of the buildkite API.
//
// Buildkite API docs: https://buildkite.com/docs/api/pipelines
type TestSuitesService struct {
	client *Client
}

type TestSuite struct {
	ID                              *string    `json:"id,omitempty" yaml:"id,omitempty"`
	GraphQLID                       *string    `json:"graphql_id,omitempty" yaml:"graphql_id,omitempty"`
	Slug                            *string    `json:"slug,omitempty" yaml:"slug,omitempty"`
	Name                            *string    `json:"name,omitempty" yaml:"name,omitempty"`
	URL                             *string    `json:"url,omitempty" yaml:"url,omitempty"`
	WebURL                          *string    `json:"web_url,omitempty" yaml:"web_url,omitempty"`
	DefaultBranch                   *string    `json:"default_branch,omitempty" yaml:"default_branch,omitempty"`
}

type TestSuiteListOptions struct {
	ListOptions
}

func (tss *TestSuitesService) List(org string, opt *TestSuiteListOptions) ([]TestSuite, *Response, error) {
	
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites", org)

	u, err := addOptions(u, opt)

	if err != nil {
		return nil, nil, err
	}

	req, err := tss.client.NewRequest("GET", u, nil)

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
