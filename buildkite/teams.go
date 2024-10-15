// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildkite

import (
	"context"
	"fmt"
)

// TeamService handles communication with the teams related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api
type TeamsService struct {
	client *Client
}

// Team represents a buildkite team.
type Team struct {
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Slug        string     `json:"slug,omitempty"`
	Description string     `json:"description,omitempty"`
	Privacy     string     `json:"privacy,omitempty"`
	Default     bool       `json:"default,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	CreatedBy   *User      `json:"created_by,omitempty"`
}

// TeamsListOptions specifies the optional parameters to the
// TeamsService.List method.
type TeamsListOptions struct {
	ListOptions
	UserID string `url:"user_id,omitempty"`
}

// Get the teams for an org.
//
// buildkite API docs: https://buildkite.com/docs/api
func (ts *TeamsService) List(ctx context.Context, org string, opt *TeamsListOptions) ([]Team, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := ts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var teams []Team
	resp, err := ts.client.Do(req, &teams)
	if err != nil {
		return nil, resp, err
	}

	return teams, resp, err
}
