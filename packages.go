package buildkite

import (
	"context"
	"fmt"
)

const fileFormKey = "file"

// PackagesService handles communication with packages Buildkite API endpoints
type PackagesService struct {
	client *Client
}

// Package represents a package which has been stored in a registry
type Package struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Version      string          `json:"version"`
	Variant      *string         `json:"variant"`
	Metadata     map[string]any  `json:"metadata"`
	URL          string          `json:"url"`
	WebURL       string          `json:"web_url"`
	Organization Organization    `json:"organization"`
	Registry     PackageRegistry `json:"registry"`
}

func (ps *PackagesService) Get(ctx context.Context, organizationSlug, registrySlug, packageID string) (Package, *Response, error) {
	url := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages/%s", organizationSlug, registrySlug, packageID)
	req, err := ps.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return Package{}, nil, fmt.Errorf("creating GET package request: %w", err)
	}

	var p Package
	resp, err := ps.client.Do(req, &p)
	if err != nil {
		return Package{}, resp, fmt.Errorf("executing GET package request: %w", err)
	}

	return p, resp, err
}

func (ps *PackagesService) Delete(ctx context.Context, organizationSlug, registrySlug, packageID string) (*Response, error) {
	u := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages/%s", organizationSlug, registrySlug, packageID)
	req, err := ps.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating DELETE package request: %w", err)
	}

	resp, err := ps.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
