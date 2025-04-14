package buildkite

import (
	"context"
	"fmt"
)

// PackageRegistriesService handles communication with the package registries related Buildkite APIs
type PackageRegistriesService struct {
	client *Client
}

// PackageRegistry represents a package registry within Buildkite
type PackageRegistry struct {
	ID          string `json:"id"`
	GraphQLID   string `json:"graphql_id"`
	Slug        string `json:"slug"`
	URL         string `json:"url"`
	WebURL      string `json:"web_url"`
	Name        string `json:"name"`
	Ecosystem   string `json:"ecosystem"`
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
	Color       string `json:"color"`
	Public      bool   `json:"public"`
	OIDCPolicy  string `json:"oidc_policy"`
}

// CreatePackageRegistryInput represents the input to create a package registry.
type CreatePackageRegistryInput struct {
	Name        string                    `json:"name,omitempty"`        // The name of the package registry
	Ecosystem   string                    `json:"ecosystem,omitempty"`   // The ecosystem of the package registry
	Description string                    `json:"description,omitempty"` // A description for the package registry
	Emoji       string                    `json:"emoji,omitempty"`       // An emoji for the package registry, in buildkite format (eg ":rocket:")
	Color       string                    `json:"color,omitempty"`       // A color for the package registry, in hex format (eg "#FF0000")
	OIDCPolicy  PackageRegistryOIDCPolicy `json:"oidc_policy,omitempty"` // The OIDC policy for the package registry, as a YAML or JSON string
}

type PackageRegistryOIDCPolicy []OIDCPolicyStatement

type OIDCPolicyStatement struct {
	Issuer string               `json:"iss"`
	Scopes []string             `json:"scopes,omitzero"`
	Claims map[string]ClaimRule `json:"claims"`
}

type ClaimRule struct {
	Equals    any      `json:"equals,omitempty"`
	NotEquals any      `json:"not_equals,omitempty"`
	In        []any    `json:"in,omitempty"`
	NotIn     []any    `json:"not_in,omitempty"`
	Matches   []string `json:"matches,omitempty"`
}

// Create creates a package registry for an organization
func (rs *PackageRegistriesService) Create(ctx context.Context, organizationSlug string, cpri CreatePackageRegistryInput) (PackageRegistry, *Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries", organizationSlug)
	req, err := rs.client.NewRequest(ctx, "POST", u, cpri)
	if err != nil {
		return PackageRegistry{}, nil, fmt.Errorf("creating POST package registry request: %v", err)
	}

	var pr PackageRegistry
	resp, err := rs.client.Do(req, &pr)
	if err != nil {
		return PackageRegistry{}, resp, err
	}

	return pr, resp, err
}

type UpdatePackageRegistryInput struct {
	Name        string                    `json:"name,omitempty"`        // The name of the package registry
	Description string                    `json:"description,omitempty"` // A description for the package registry
	Emoji       string                    `json:"emoji,omitempty"`       // An emoji for the package registry, in buildkite format (eg ":rocket:")
	Color       string                    `json:"color,omitempty"`       // A color for the package registry, in hex format (eg "#FF0000")
	OIDCPolicy  PackageRegistryOIDCPolicy `json:"oidc_policy,omitempty"` // The OIDC policy for the package registry, as a YAML or JSON string
}

// Update updates a package registry for an organization
func (rs *PackageRegistriesService) Update(ctx context.Context, organizationSlug, registrySlug string, upri UpdatePackageRegistryInput) (PackageRegistry, *Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries/%s", organizationSlug, registrySlug)
	req, err := rs.client.NewRequest(ctx, "PATCH", u, upri)
	if err != nil {
		return PackageRegistry{}, nil, fmt.Errorf("creating PATCH package registry request: %v", err)
	}

	var pr PackageRegistry
	resp, err := rs.client.Do(req, &pr)
	if err != nil {
		return PackageRegistry{}, resp, err
	}

	return pr, resp, err
}

// Get retrieves a package registry for an organization
func (rs *PackageRegistriesService) Get(ctx context.Context, organizationSlug, registrySlug string) (PackageRegistry, *Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries/%s", organizationSlug, registrySlug)
	req, err := rs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return PackageRegistry{}, nil, fmt.Errorf("creating GET package registry request: %v", err)
	}

	var pr PackageRegistry
	resp, err := rs.client.Do(req, &pr)
	if err != nil {
		return PackageRegistry{}, resp, err
	}

	return pr, resp, err
}

// List retrieves a list of package all package registries for an organization
func (rs *PackageRegistriesService) List(ctx context.Context, organizationSlug string) ([]PackageRegistry, *Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries", organizationSlug)
	req, err := rs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("creating GET package registry request: %v", err)
	}

	var prs []PackageRegistry
	resp, err := rs.client.Do(req, &prs)
	if err != nil {
		return nil, resp, err
	}

	return prs, resp, err
}

// Delete deletes a package registry for an organization
func (rs *PackageRegistriesService) Delete(ctx context.Context, organizationSlug, registrySlug string) (*Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries/%s", organizationSlug, registrySlug)
	req, err := rs.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating DELETE package registry request: %v", err)
	}

	resp, err := rs.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
