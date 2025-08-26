package buildkite

import (
	"context"
	"fmt"
)

// TestsService handles communication with test related
// methods of the Buildkite Test Analytics API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/analytics/tests
type TestsService struct {
	client *Client
}

type Test struct {
	ID       string `json:"id,omitempty"`
	URL      string `json:"url,omitempty"`
	WebURL   string `json:"web_url,omitempty"`
	Scope    string `json:"scope,omitempty"`
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
	FileName string `json:"file_name,omitempty"`
}

func (ts *TestsService) Get(ctx context.Context, org, slug, testID string) (Test, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/tests/%s", org, slug, testID)
	req, err := ts.client.NewRequest(ctx, "GET", u, nil)
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
