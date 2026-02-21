package buildkite

import (
	"context"
	"fmt"
	"net/http"
)

// MembersService handles communication with the members related
// methods of the Buildkite API
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/organizations/members
type MembersService struct {
	client *Client
}

// Member represents a Buildkite member
type Member struct {
	// UUID is a more accurate representation of the value
	UUID  string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// MemberListOptions specifies the optional parameters to the
// MembersService.List method
type MemberListOptions struct {
	ListOptions
}

// List members for a Buildkite organization
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/organizations/members#list-organization-members
func (ms *MembersService) List(ctx context.Context, org string, opt *MemberListOptions) ([]Member, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/members", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := ms.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var members []Member
	resp, err := ms.client.Do(req, &members)
	if err != nil {
		return nil, resp, err
	}
	return members, resp, err
}

// Get a specific member of a Buildkite organization
//
// https://buildkite.com/docs/apis/rest-api/organizations/members#get-an-organization-member
func (ms *MembersService) Get(ctx context.Context, org, memberUUID string) (Member, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/members/%s", org, memberUUID)
	req, err := ms.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return Member{}, nil, err
	}

	var member Member
	resp, err := ms.client.Do(req, &member)
	if err != nil {
		return Member{}, resp, err
	}

	return member, resp, err
}
