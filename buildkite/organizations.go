package buildkite

import (
	"context"
	"fmt"
)

// OrganizationsService handles communication with the organization related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api/organizations
type OrganizationsService struct {
	client *Client
}

// Organization represents a buildkite organization.
type Organization struct {
	ID           *string    `json:"id,omitempty"`
	GraphQLID    *string    `json:"graphql_id,omitempty"`
	URL          *string    `json:"url,omitempty"`
	WebURL       *string    `json:"web_url,omitempty"`
	Name         *string    `json:"name,omitempty"`
	Slug         *string    `json:"slug,omitempty"`
	Repository   *string    `json:"repository,omitempty"`
	PipelinesURL *string    `json:"pipelines_url,omitempty"`
	EmojisURL    *string    `json:"emojis_url,omitempty"`
	AgentsURL    *string    `json:"agents_url,omitempty"`
	CreatedAt    *Timestamp `json:"created_at,omitempty"`
}

// OrganizationListOptions specifies the optional parameters to the
// OrganizationsService.List method.
type OrganizationListOptions struct {
	ListOptions
}

// List the organizations for the current user.
//
// buildkite API docs: https://buildkite.com/docs/api/organizations#list-organizations
func (os *OrganizationsService) List(ctx context.Context, opt *OrganizationListOptions) ([]Organization, *Response, error) {
	var u string

	u = fmt.Sprintf("v2/organizations")

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := os.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	orgs := new([]Organization)
	resp, err := os.client.Do(req, orgs)
	if err != nil {
		return nil, resp, err
	}

	return *orgs, resp, err
}

// Get fetches an organization
//
// buildkite API docs: https://buildkite.com/docs/api/organizations#get-an-organization
func (os *OrganizationsService) Get(ctx context.Context, slug string) (*Organization, *Response, error) {

	u := fmt.Sprintf("v2/organizations/%s", slug)

	req, err := os.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	organization := new(Organization)
	resp, err := os.client.Do(req, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, err
}
