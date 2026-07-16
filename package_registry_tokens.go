package buildkite

import (
	"context"
	"fmt"
)

// PackageRegistryTokensService handles communication with the package
// registry tokens related Buildkite APIs
type PackageRegistryTokensService struct {
	client *Client
}

// PackageRegistryToken represents a token scoped to a package registry.
type PackageRegistryToken struct {
	ID           string                       `json:"id,omitempty"`
	GraphQLID    string                       `json:"graphql_id,omitempty"`
	Description  string                       `json:"description,omitempty"`
	URL          string                       `json:"url,omitempty"`
	CreatedAt    *Timestamp                   `json:"created_at,omitempty"`
	CreatedBy    PackageRegistryTokenCreator  `json:"created_by,omitempty"`
	Organization PackageRegistryTokenOrg      `json:"organization,omitempty"`
	Registry     PackageRegistryTokenRegistry `json:"registry,omitempty"`

	// Token holds the raw secret token value. It is only ever populated in
	// the response from Create; Get, List, and Update never include it.
	Token string `json:"token,omitempty"`
}

// PackageRegistryTokenCreator represents the user who created a package
// registry token.
type PackageRegistryTokenCreator struct {
	ID        string     `json:"id,omitempty"`
	GraphQLID string     `json:"graphql_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// PackageRegistryTokenOrg represents the organization a package registry
// token belongs to.
type PackageRegistryTokenOrg struct {
	ID     string `json:"id,omitempty"`
	Slug   string `json:"slug,omitempty"`
	URL    string `json:"url,omitempty"`
	WebURL string `json:"web_url,omitempty"`
}

// PackageRegistryTokenRegistry represents the registry a package registry
// token belongs to.
type PackageRegistryTokenRegistry struct {
	ID        string `json:"id,omitempty"`
	GraphQLID string `json:"graphql_id,omitempty"`
	Slug      string `json:"slug,omitempty"`
	URL       string `json:"url,omitempty"`
	WebURL    string `json:"web_url,omitempty"`
}

// CreatePackageRegistryTokenInput represents the input to create a package
// registry token.
type CreatePackageRegistryTokenInput struct {
	Description string `json:"description"`
}

// UpdatePackageRegistryTokenInput represents the request body for updating
// a package registry token.
type UpdatePackageRegistryTokenInput struct {
	Description Optional[string] `json:"description,omitzero"`
}

// Get retrieves a package registry token for an organization's registry.
func (ts *PackageRegistryTokensService) Get(ctx context.Context, organizationSlug, registrySlug, tokenID string) (PackageRegistryToken, *Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/tokens/%s", organizationSlug, registrySlug, tokenID)
	req, err := ts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return PackageRegistryToken{}, nil, fmt.Errorf("creating GET package registry token request: %w", err)
	}

	var t PackageRegistryToken
	resp, err := ts.client.Do(req, &t)
	if err != nil {
		return PackageRegistryToken{}, resp, err
	}

	return t, resp, err
}

// List retrieves all package registry tokens for an organization's registry.
func (ts *PackageRegistryTokensService) List(ctx context.Context, organizationSlug, registrySlug string) ([]PackageRegistryToken, *Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/tokens", organizationSlug, registrySlug)
	req, err := ts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("creating GET package registry token request: %w", err)
	}

	var tokens []PackageRegistryToken
	resp, err := ts.client.Do(req, &tokens)
	if err != nil {
		return nil, resp, err
	}

	return tokens, resp, err
}

// Create creates a package registry token for an organization's registry.
func (ts *PackageRegistryTokensService) Create(ctx context.Context, organizationSlug, registrySlug string, input CreatePackageRegistryTokenInput) (PackageRegistryToken, *Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/tokens", organizationSlug, registrySlug)
	req, err := ts.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return PackageRegistryToken{}, nil, fmt.Errorf("creating POST package registry token request: %w", err)
	}

	var t PackageRegistryToken
	resp, err := ts.client.Do(req, &t)
	if err != nil {
		return PackageRegistryToken{}, resp, err
	}

	return t, resp, err
}

// Update updates a package registry token for an organization's registry.
// Note: unlike most Update methods in this client, this issues a POST, not
// a PATCH — that is what the Buildkite API route actually expects here.
func (ts *PackageRegistryTokensService) Update(ctx context.Context, organizationSlug, registrySlug, tokenID string, input UpdatePackageRegistryTokenInput) (PackageRegistryToken, *Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/tokens/%s", organizationSlug, registrySlug, tokenID)
	req, err := ts.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return PackageRegistryToken{}, nil, fmt.Errorf("creating POST package registry token request: %w", err)
	}

	var t PackageRegistryToken
	resp, err := ts.client.Do(req, &t)
	if err != nil {
		return PackageRegistryToken{}, resp, err
	}

	return t, resp, err
}
