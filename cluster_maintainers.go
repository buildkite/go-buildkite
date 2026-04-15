package buildkite

import (
	"context"
	"fmt"
)

// ClusterMaintainersService handles API calls for cluster maintainer assignments.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/clusters/maintainers

type ClusterMaintainersService struct {
	client *Client
}

// ClusterMaintainersListOptions controls pagination for listing maintainers.
type ClusterMaintainersListOptions struct {
	ListOptions
}

// List returns the maintainers assigned to a cluster.
func (cms *ClusterMaintainersService) List(ctx context.Context, org, clusterID string, opt *ClusterMaintainersListOptions) ([]ClusterMaintainerEntry, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/maintainers", org, clusterID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := cms.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var maintainers []ClusterMaintainerEntry
	resp, err := cms.client.Do(req, &maintainers)
	if err != nil {
		return nil, resp, err
	}

	return maintainers, resp, err
}

// Get returns one maintainer assignment by ID.
func (cms *ClusterMaintainersService) Get(ctx context.Context, org, clusterID, id string) (ClusterMaintainerEntry, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/maintainers/%s", org, clusterID, id)

	req, err := cms.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return ClusterMaintainerEntry{}, nil, err
	}

	var maintainer ClusterMaintainerEntry
	resp, err := cms.client.Do(req, &maintainer)
	if err != nil {
		return ClusterMaintainerEntry{}, resp, err
	}

	return maintainer, resp, err
}

// Create assigns a user or team as a maintainer.
func (cms *ClusterMaintainersService) Create(ctx context.Context, org, clusterID string, input ClusterMaintainer) (ClusterMaintainerEntry, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/maintainers", org, clusterID)

	req, err := cms.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return ClusterMaintainerEntry{}, nil, err
	}

	var maintainer ClusterMaintainerEntry
	resp, err := cms.client.Do(req, &maintainer)
	if err != nil {
		return ClusterMaintainerEntry{}, resp, err
	}

	return maintainer, resp, err
}

// Delete removes a maintainer assignment from a cluster.
func (cms *ClusterMaintainersService) Delete(ctx context.Context, org, clusterID, id string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/maintainers/%s", org, clusterID, id)

	req, err := cms.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return cms.client.Do(req, nil)
}
