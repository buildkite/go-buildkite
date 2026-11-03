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

// ClusterToken represents a token scoped to a cluster.
type ClusterToken struct {
	ID                 string         `json:"id,omitempty"`
	GraphQLID          string         `json:"graphql_id,omitempty"`
	Description        string         `json:"description,omitempty"`
	AllowedIPAddresses string         `json:"allowed_ip_addresses,omitempty"`
	URL                string         `json:"url,omitempty"`
	ClusterURL         string         `json:"cluster_url,omitempty"`
	CreatedAt          *Timestamp     `json:"created_at,omitempty"`
	CreatedBy          ClusterCreator `json:"created_by,omitempty"`
	ExpiresAt          *Timestamp     `json:"expires_at,omitempty"`
	Token              string         `json:"token,omitempty"`
	// These fields are absent when their server features are disabled. Nil
	// slices are omitted, while initialized empty slices are retained as [].
	Permissions     []string `json:"permissions,omitzero"`
	AllowAllQueues  *bool    `json:"allow_all_queues,omitempty"`
	AllowedQueueIDs []string `json:"allowed_queue_ids,omitzero"`
}

// ClusterTokenCreate represents the request body for creating a cluster token.
type ClusterTokenCreate struct {
	Description        string `json:"description,omitempty"`
	AllowedIPAddresses string `json:"allowed_ip_addresses,omitempty"`
	// ExpiresAt is when the token stops letting new agents register. The API
	// requires it to be at least 10 minutes in the future, and it cannot be
	// changed after creation. Leave it nil for a token that never expires.
	ExpiresAt *Timestamp `json:"expires_at,omitempty"`
	// Permissions defaults to the server's permissions when omitted. Use
	// Some([]string{}) for no actions; nil slices encode as null, not [].
	Permissions Optional[[]string] `json:"permissions,omitzero"`
	// AllowAllQueues defaults to true when omitted. If false, AllowedQueueIDs
	// must be explicit, including Some([]string{}) to allow no queues.
	AllowAllQueues Optional[bool] `json:"allow_all_queues,omitzero"`
	// AllowedQueueIDs contains queue UUIDs, not GraphQL IDs. A nonempty list
	// cannot be combined with AllowAllQueues set to true.
	AllowedQueueIDs Optional[[]string] `json:"allowed_queue_ids,omitzero"`
}

// ClusterTokenUpdate represents the request body for updating a cluster token.
type ClusterTokenUpdate struct {
	Description        Optional[string] `json:"description,omitzero"`
	AllowedIPAddresses Optional[string] `json:"allowed_ip_addresses,omitzero"`
	// Permissions is unchanged when omitted. Use Some([]string{}) to remove
	// all actions; nil slices encode as null, not [].
	Permissions Optional[[]string] `json:"permissions,omitzero"`
	// AllowAllQueues is unchanged when omitted. To switch a restricted token
	// to all queues, send Some(true) and AllowedQueueIDs: Some([]string{}).
	AllowAllQueues Optional[bool] `json:"allow_all_queues,omitzero"`
	// AllowedQueueIDs replaces queue UUIDs when set and preserves them when
	// omitted. Some([]string{}) clears them; with AllowAllQueues false this
	// allows no queues. A nonempty list cannot accompany AllowAllQueues true.
	AllowedQueueIDs Optional[[]string] `json:"allowed_queue_ids,omitzero"`
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

func (cts *ClusterTokensService) Create(ctx context.Context, org, clusterID string, ctc ClusterTokenCreate) (ClusterToken, *Response, error) {
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

func (cts *ClusterTokensService) Update(ctx context.Context, org, clusterID, tokenID string, ctc ClusterTokenUpdate) (ClusterToken, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/tokens/%s", org, clusterID, tokenID)
	req, err := cts.client.NewRequest(ctx, "PATCH", u, ctc)
	if err != nil {
		return ClusterToken{}, nil, err
	}

	var ct ClusterToken
	resp, err := cts.client.Do(req, &ct)
	if err != nil {
		return ClusterToken{}, resp, err
	}

	return ct, resp, err
}

func (cts *ClusterTokensService) Delete(ctx context.Context, org, clusterID, tokenID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/tokens/%s", org, clusterID, tokenID)
	req, err := cts.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return cts.client.Do(req, nil)
}
