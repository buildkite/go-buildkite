package buildkite

import (
	"context"
	"fmt"
)

// ClusterSecretsService handles communication with cluster secret related
// methods of the Buildkite API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/clusters#cluster-secrets
type ClusterSecretsService struct {
	client *Client
}

type ClusterSecret struct {
	ID           string         `json:"id,omitempty"`
	GraphQLID    string         `json:"graphql_id,omitempty"`
	Key          string         `json:"key,omitempty"`
	Description  string         `json:"description,omitempty"`
	Policy       string         `json:"policy,omitempty"`
	URL          string         `json:"url,omitempty"`
	ClusterURL   string         `json:"cluster_url,omitempty"`
	CreatedAt    *Timestamp     `json:"created_at,omitempty"`
	CreatedBy    ClusterCreator `json:"created_by"`
	UpdatedAt    *Timestamp     `json:"updated_at,omitempty"`
	UpdatedBy    ClusterCreator `json:"updated_by"`
	LastReadAt   *Timestamp     `json:"last_read_at,omitempty"`
	Organization string         `json:"organization,omitempty"`
}

type ClusterSecretCreate struct {
	Key         string `json:"key"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	Policy      string `json:"policy,omitempty"`
}

type ClusterSecretUpdate struct {
	Description string `json:"description,omitempty"`
	Policy      string `json:"policy,omitempty"`
}

type ClusterSecretValueUpdate struct {
	Value string `json:"value"`
}

type ClusterSecretsListOptions struct {
	ListOptions
}

func (css *ClusterSecretsService) List(ctx context.Context, org, clusterID string, opt *ClusterSecretsListOptions) ([]ClusterSecret, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/secrets", org, clusterID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := css.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var secrets []ClusterSecret
	resp, err := css.client.Do(req, &secrets)
	if err != nil {
		return nil, resp, err
	}

	return secrets, resp, err
}

func (css *ClusterSecretsService) Get(ctx context.Context, org, clusterID, secretID string) (ClusterSecret, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/secrets/%s", org, clusterID, secretID)
	req, err := css.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return ClusterSecret{}, nil, err
	}

	var secret ClusterSecret
	resp, err := css.client.Do(req, &secret)
	if err != nil {
		return ClusterSecret{}, resp, err
	}

	return secret, resp, err
}

func (css *ClusterSecretsService) Create(ctx context.Context, org, clusterID string, input ClusterSecretCreate) (ClusterSecret, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/secrets", org, clusterID)
	req, err := css.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return ClusterSecret{}, nil, err
	}

	var secret ClusterSecret
	resp, err := css.client.Do(req, &secret)
	if err != nil {
		return ClusterSecret{}, resp, err
	}

	return secret, resp, err
}

func (css *ClusterSecretsService) Update(ctx context.Context, org, clusterID, secretID string, input ClusterSecretUpdate) (ClusterSecret, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/secrets/%s", org, clusterID, secretID)
	req, err := css.client.NewRequest(ctx, "PUT", u, input)
	if err != nil {
		return ClusterSecret{}, nil, err
	}

	var secret ClusterSecret
	resp, err := css.client.Do(req, &secret)
	if err != nil {
		return ClusterSecret{}, resp, err
	}

	return secret, resp, err
}

func (css *ClusterSecretsService) UpdateValue(ctx context.Context, org, clusterID, secretID string, input ClusterSecretValueUpdate) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/secrets/%s/value", org, clusterID, secretID)
	req, err := css.client.NewRequest(ctx, "PUT", u, input)
	if err != nil {
		return nil, err
	}

	return css.client.Do(req, nil)
}

func (css *ClusterSecretsService) Delete(ctx context.Context, org, clusterID, secretID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/secrets/%s", org, clusterID, secretID)
	req, err := css.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return css.client.Do(req, nil)
}
