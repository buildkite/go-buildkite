package buildkite

import (
	"context"
	"fmt"
)

// TeamMemberService handles communication with the teams related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api
type TeamMemberService struct {
	client *Client
}

// TeamMember represents a member of a team.
type TeamMember struct {
	ID        string    `json:"user_id,omitempty"`
	UserName  string    `json:"user_name,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt Timestamp `json:"created_at,omitempty"`
}

// CreateTeamMember represents a request to create a team member.
type CreateTeamMember struct {
	UserID string `json:"user_id,omitempty"`
	Role   string `json:"role,omitempty"`
}

// TeamMembersListOptions specifies the optional parameters to the ListTeamMembers method.
type TeamMembersListOptions struct {
	ListOptions
}

// ListTeamMembers lists the members of a team.
func (ts *TeamMemberService) ListTeamMembers(ctx context.Context, org string, id string, opt *TeamMembersListOptions) ([]TeamMember, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/members", org, id)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := ts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var teamMembers []TeamMember
	resp, err := ts.client.Do(req, &teamMembers)
	if err != nil {
		return nil, resp, err
	}

	return teamMembers, resp, err
}

// GetTeamMember gets a team member.
func (ts *TeamMemberService) GetTeamMember(ctx context.Context, org string, teamID string, userID string) (TeamMember, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/members/%s", org, teamID, userID)

	req, err := ts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TeamMember{}, err
	}

	var teamMember TeamMember
	_, err = ts.client.Do(req, &teamMember)
	if err != nil {
		return teamMember, err
	}

	return teamMember, err
}

// CreateTeamMember creates a team member.
func (ts *TeamMemberService) CreateTeamMember(ctx context.Context, org string, teamID string, t CreateTeamMember) (TeamMember, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/members", org, teamID)

	req, err := ts.client.NewRequest(ctx, "POST", u, t)
	if err != nil {
		return TeamMember{}, nil, err
	}

	var teamMember TeamMember
	resp, err := ts.client.Do(req, &teamMember)
	if err != nil {
		return teamMember, resp, err
	}

	return teamMember, resp, err
}

// UpdateTeamMember updates a team member.
func (ts *TeamMemberService) UpdateTeamMember(ctx context.Context, org string, teamID string, userID string, role string) (TeamMember, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/members/%s", org, teamID, userID)

	req, err := ts.client.NewRequest(ctx, "PATCH", u, map[string]string{"role": role})
	if err != nil {
		return TeamMember{}, nil, err
	}

	var teamMember TeamMember
	resp, err := ts.client.Do(req, &teamMember)
	if err != nil {
		return teamMember, resp, err
	}

	return teamMember, resp, err
}

// DeleteTeamMember deletes a team member.
func (ts *TeamMemberService) DeleteTeamMember(ctx context.Context, org string, teamID string, userID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/members/%s", org, teamID, userID)

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
