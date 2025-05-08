package buildkite

import (
	"context"
	"fmt"
	"net/url"
)

// PackageRegistriesService handles communication with the package registries related Buildkite APIs
type PackageRegistriesService struct {
	client *Client
}

// PackageRegistry represents a package registry within Buildkite
type PackageRegistry struct {
	ID          string `json:"id,omitempty"`
	GraphQLID   string `json:"graphql_id,omitempty"`
	Slug        string `json:"slug,omitempty"`
	URL         string `json:"url,omitempty"`
	WebURL      string `json:"web_url,omitempty"`
	Name        string `json:"name,omitempty"`
	Ecosystem   string `json:"ecosystem,omitempty"`
	Description string `json:"description,omitempty"`
	Emoji       string `json:"emoji,omitempty"`
	Color       string `json:"color,omitempty"`
	Public      bool   `json:"public,omitempty"`
	OIDCPolicy  string `json:"oidc_policy,omitempty"`
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
		return PackageRegistry{}, nil, fmt.Errorf("creating POST package registry request: %w", err)
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
		return PackageRegistry{}, nil, fmt.Errorf("creating PATCH package registry request: %w", err)
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
		return PackageRegistry{}, nil, fmt.Errorf("creating GET package registry request: %w", err)
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
		return nil, nil, fmt.Errorf("creating GET package registry request: %w", err)
	}

	var prs []PackageRegistry
	resp, err := rs.client.Do(req, &prs)
	if err != nil {
		return nil, resp, err
	}

	return prs, resp, err
}

type RegistryPackagesOptions struct {
	Before  string `url:"before,omitempty"`
	After   string `url:"after,omitempty"`
	PerPage string `url:"per_page,omitempty"` // Should be a string of an int eg "50"
}

type RegistryPackagesLink string

func (l RegistryPackagesLink) ToOptions() (*RegistryPackagesOptions, error) {
	u, err := url.Parse(string(l))
	if err != nil {
		return nil, fmt.Errorf("parsing link: %w", err)
	}

	q := u.Query()

	return &RegistryPackagesOptions{
		Before:  q.Get("before"),
		After:   q.Get("after"),
		PerPage: q.Get("per_page"),
	}, nil
}

type RegistryPackagesLinks struct {
	First    RegistryPackagesLink `json:"first,omitempty"`
	Previous RegistryPackagesLink `json:"prev,omitempty"`
	Self     RegistryPackagesLink `json:"self,omitempty"`
	Next     RegistryPackagesLink `json:"next,omitempty"`
}

type RegistryPackages struct {
	Items []Package             `json:"items"`
	Links RegistryPackagesLinks `json:"links"`
}

func (rs *PackageRegistriesService) ListPackages(ctx context.Context, organizationSlug, registrySlug string, opts *RegistryPackagesOptions) (RegistryPackages, *Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages", organizationSlug, registrySlug)
	u, err := addOptions(u, opts)
	if err != nil {
		return RegistryPackages{}, nil, fmt.Errorf("adding query params to path: %w", err)
	}

	req, err := rs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return RegistryPackages{}, nil, fmt.Errorf("creating GET registry packages request: %w", err)
	}

	var packages RegistryPackages
	resp, err := rs.client.Do(req, &packages)
	if err != nil {
		return RegistryPackages{}, resp, err
	}

	return packages, resp, nil
}

// Delete deletes a package registry for an organization
func (rs *PackageRegistriesService) Delete(ctx context.Context, organizationSlug, registrySlug string) (*Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries/%s", organizationSlug, registrySlug)
	req, err := rs.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating DELETE package registry request: %w", err)
	}

	resp, err := rs.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
