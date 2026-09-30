package buildkite

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// CacheRegistriesService handles communication with the cluster cache registry
// methods of the Buildkite API.
type CacheRegistriesService struct {
	client *Client
}

// CacheRegistryPolicy is a normalized cache registry policy JSON object.
type CacheRegistryPolicy map[string]any

// CacheRegistry represents a cache registry within a cluster.
type CacheRegistry struct {
	UUID        string              `json:"uuid,omitempty"`
	Slug        string              `json:"slug,omitempty"`
	Name        string              `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Emoji       *string             `json:"emoji,omitempty"`
	Color       *string             `json:"color,omitempty"`
	Policy      CacheRegistryPolicy `json:"policy,omitempty"`
	// Default reports whether this is the cluster's default cache registry.
	Default    bool       `json:"default"`
	CreatedAt  *Timestamp `json:"created_at,omitempty"`
	UpdatedAt  *Timestamp `json:"updated_at,omitempty"`
	URL        string     `json:"url,omitempty"`
	ClusterURL string     `json:"cluster_url,omitempty"`
}

// CacheRegistryCreate represents the request body for creating a cache
// registry. The optional nullable fields can be omitted with their zero value,
// sent with Some(&value), or sent as JSON null with Some[*string](nil) for
// metadata and Some[CacheRegistryPolicy](nil) for policy.
type CacheRegistryCreate struct {
	Name        string                        `json:"name"`
	Description Optional[*string]             `json:"description,omitzero"`
	Emoji       Optional[*string]             `json:"emoji,omitzero"`
	Color       Optional[*string]             `json:"color,omitzero"`
	Policy      Optional[CacheRegistryPolicy] `json:"policy,omitzero"`
}

// CacheRegistryUpdate represents the request body for updating a cache
// registry. Unset fields are omitted. Nullable fields can be cleared by passing
// Some[*string](nil) for metadata or Some[CacheRegistryPolicy](nil) for policy.
type CacheRegistryUpdate struct {
	Name        Optional[string]              `json:"name,omitzero"`
	Description Optional[*string]             `json:"description,omitzero"`
	Emoji       Optional[*string]             `json:"emoji,omitzero"`
	Color       Optional[*string]             `json:"color,omitzero"`
	Policy      Optional[CacheRegistryPolicy] `json:"policy,omitzero"`
}

// CacheRegistriesListOptions specifies cursor pagination parameters.
type CacheRegistriesListOptions struct {
	After   string `url:"after,omitempty"`
	Before  string `url:"before,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
}

type CacheRegistriesListLink string

// ToOptions extracts cursor pagination options from a cache registries link.
func (l CacheRegistriesListLink) ToOptions() (*CacheRegistriesListOptions, error) {
	u, err := url.Parse(string(l))
	if err != nil {
		return nil, fmt.Errorf("parsing link: %w", err)
	}

	opt := &CacheRegistriesListOptions{
		After:  u.Query().Get("after"),
		Before: u.Query().Get("before"),
	}

	if perPage := u.Query().Get("per_page"); perPage != "" {
		opt.PerPage, err = strconv.Atoi(perPage)
		if err != nil {
			return nil, fmt.Errorf("parsing per_page: %w", err)
		}
	}

	return opt, nil
}

type CacheRegistriesListLinks struct {
	Self     CacheRegistriesListLink `json:"self,omitempty"`
	First    CacheRegistriesListLink `json:"first,omitempty"`
	Previous CacheRegistriesListLink `json:"prev,omitempty"`
	Next     CacheRegistriesListLink `json:"next,omitempty"`
}

type CacheRegistriesList struct {
	Items []CacheRegistry          `json:"items"`
	Links CacheRegistriesListLinks `json:"links"`
}

// List returns the cache registries in a cluster.
func (s *CacheRegistriesService) List(ctx context.Context, org, clusterID string, opt *CacheRegistriesListOptions) (CacheRegistriesList, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/cache-registries", org, clusterID)
	u, err := addOptions(u, opt)
	if err != nil {
		return CacheRegistriesList{}, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return CacheRegistriesList{}, nil, err
	}

	var registries CacheRegistriesList
	resp, err := s.client.Do(req, &registries)
	if err != nil {
		return CacheRegistriesList{}, resp, err
	}

	return registries, resp, nil
}

// Get returns a cache registry by UUID.
func (s *CacheRegistriesService) Get(ctx context.Context, org, clusterID, registryUUID string) (CacheRegistry, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/cache-registries/%s", org, clusterID, registryUUID)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return CacheRegistry{}, nil, err
	}

	var registry CacheRegistry
	resp, err := s.client.Do(req, &registry)
	if err != nil {
		return CacheRegistry{}, resp, err
	}

	return registry, resp, nil
}

// Create creates a cache registry in a cluster.
func (s *CacheRegistriesService) Create(ctx context.Context, org, clusterID string, input CacheRegistryCreate) (CacheRegistry, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/cache-registries", org, clusterID)
	req, err := s.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return CacheRegistry{}, nil, err
	}

	var registry CacheRegistry
	resp, err := s.client.Do(req, &registry)
	if err != nil {
		return CacheRegistry{}, resp, err
	}

	return registry, resp, nil
}

// Update updates a cache registry by UUID.
func (s *CacheRegistriesService) Update(ctx context.Context, org, clusterID, registryUUID string, input CacheRegistryUpdate) (CacheRegistry, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/cache-registries/%s", org, clusterID, registryUUID)
	req, err := s.client.NewRequest(ctx, "PATCH", u, input)
	if err != nil {
		return CacheRegistry{}, nil, err
	}

	var registry CacheRegistry
	resp, err := s.client.Do(req, &registry)
	if err != nil {
		return CacheRegistry{}, resp, err
	}

	return registry, resp, nil
}

// Delete deletes a cache registry by UUID.
func (s *CacheRegistriesService) Delete(ctx context.Context, org, clusterID, registryUUID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/cache-registries/%s", org, clusterID, registryUUID)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
