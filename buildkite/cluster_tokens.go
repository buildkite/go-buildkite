package buildkite

import (
	"context"
	"fmt"
)

// ClusterTokensService handles communication with cluster token related
// methods of the Buildkite API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/clusters#cluster-tokens
type ClusterTokensService struct {
	client *Client
}

type ClusterToken struct {
	ID                 string         `json:"id,omitempty"`
	GraphQLID          string         `json:"graphql_id,omitempty"`
	Description        string         `json:"description,omitempty"`
	AllowedIPAddresses string         `json:"allowed_ip_addresses,omitempty"`
	URL                string         `json:"url,omitempty"`
	ClusterURL         string         `json:"cluster_url,omitempty"`
	CreatedAt          *Timestamp     `json:"created_at,omitempty"`
	CreatedBy          ClusterCreator `json:"created_by,omitempty"`
	Token              string         `json:"token,omitempty"`
}

type ClusterTokenCreateUpdate struct {
	Description        string `json:"description,omitempty"`
	AllowedIPAddresses string `json:"allowed_ip_addresses,omitempty"`
}

type ClusterTokensListOptions struct {
	ListOptions
}

func (cts *ClusterTokensService) List(ctx context.Context, org, clusterID string, opt *ClusterTokensListOptions) ([]ClusterToken, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/tokens", org, clusterID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := cts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tokens []ClusterToken
	resp, err := cts.client.Do(req, &tokens)
	if err != nil {
		return nil, resp, err
	}

	return tokens, resp, err
}

func (cts *ClusterTokensService) Get(ctx context.Context, org, clusterID, tokenID string) (ClusterToken, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/tokens/%s", org, clusterID, tokenID)
	req, err := cts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return ClusterToken{}, nil, err
	}

	var token ClusterToken
	resp, err := cts.client.Do(req, &token)
	if err != nil {
		return ClusterToken{}, resp, err
	}

	return token, resp, err
}

func (cts *ClusterTokensService) Create(ctx context.Context, org, clusterID string, ctc ClusterTokenCreateUpdate) (ClusterToken, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/tokens", org, clusterID)
	req, err := cts.client.NewRequest(ctx, "POST", u, ctc)
	if err != nil {
		return ClusterToken{}, nil, err
	}

	var token ClusterToken
	resp, err := cts.client.Do(req, &token)
	if err != nil {
		return ClusterToken{}, resp, err
	}

	return token, resp, err
}

func (cts *ClusterTokensService) Update(ctx context.Context, org, clusterID, tokenID string, ctc ClusterTokenCreateUpdate) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/tokens/%s", org, clusterID, tokenID)
	req, err := cts.client.NewRequest(ctx, "PATCH", u, ctc)
	if err != nil {
		return nil, err
	}

	resp, err := cts.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (cts *ClusterTokensService) Delete(ctx context.Context, org, clusterID, tokenID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/tokens/%s", org, clusterID, tokenID)
	req, err := cts.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return cts.client.Do(req, nil)
}
