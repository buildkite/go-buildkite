// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildkite

import (
	"context"
	"fmt"
)

// TeamsService handles communication with the teams related
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
	Default     bool       `json:"default"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	CreatedBy   *User      `json:"created_by,omitempty"`
}

// TeamsListOptions specifies the optional parameters to the
// TeamsService.List method.
type TeamsListOptions struct {
	ListOptions
	UserID string `url:"user_id,omitempty"`
}

// CreateTeam represents a request to create a team.
type CreateTeam struct {
	Name                      string `json:"name,omitempty"`
	Description               string `json:"description,omitempty"`
	Privacy                   string `json:"privacy,omitempty"`
	IsDefaultTeam             bool   `json:"is_default_team"`
	DefaultMemberRole         string `json:"default_member_role,omitempty"`
	MembersCanCreatePipelines bool   `json:"members_can_create_pipelines"`
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

// GetTeam gets a team.
func (ts *TeamsService) GetTeam(ctx context.Context, org string, id string) (Team, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s", org, id)

	req, err := ts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return Team{}, err
	}

	var team Team
	_, err = ts.client.Do(req, &team)
	if err != nil {
		return team, err
	}

	return team, err
}

// CreateTeam creates a team.
func (ts *TeamsService) CreateTeam(ctx context.Context, org string, t CreateTeam) (Team, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams", org)

	req, err := ts.client.NewRequest(ctx, "POST", u, t)
	if err != nil {
		return Team{}, nil, err
	}

	var team Team
	resp, err := ts.client.Do(req, &team)
	if err != nil {
		return team, nil, err
	}

	return team, resp, err
}

// UpdateTeam updates a team.
func (ts *TeamsService) UpdateTeam(ctx context.Context, org string, id string, t CreateTeam) (Team, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s", org, id)

	req, err := ts.client.NewRequest(ctx, "PATCH", u, t)
	if err != nil {
		return Team{}, nil, err
	}

	var team Team
	resp, err := ts.client.Do(req, &team)
	if err != nil {
		return team, nil, err
	}

	return team, resp, err
}

// DeleteTeam deletes a team.
func (ts *TeamsService) DeleteTeam(ctx context.Context, org string, id string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s", org, id)

	req, err := ts.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := ts.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
