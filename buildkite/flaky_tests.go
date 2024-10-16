package buildkite

import (
	"context"
	"fmt"
)

// FlakyTestsService handles communication with flaky test related
// methods of the Buildkite Test Analytics API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/analytics/flaky-tests
type FlakyTestsService struct {
	client *Client
}

type FlakyTest struct {
	ID                   string     `json:"id,omitempty"`
	WebURL               string     `json:"web_url,omitempty"`
	Scope                string     `json:"scope,omitempty"`
	Name                 string     `json:"name,omitempty"`
	Location             string     `json:"location,omitempty"`
	FileName             string     `json:"file_name,omitempty"`
	Instances            int        `json:"instances,omitempty"`
	MostRecentInstanceAt *Timestamp `json:"most_recent_instance_at,omitempty"`
}

type FlakyTestsListOptions struct{ ListOptions }

func (fts *FlakyTestsService) List(ctx context.Context, org, slug string, opt *FlakyTestsListOptions) ([]FlakyTest, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/flaky-tests", org, slug)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := fts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var flakyTests []FlakyTest
	resp, err := fts.client.Do(req, &flakyTests)
	if err != nil {
		return nil, resp, err
	}

	return flakyTests, resp, err
}
