package buildkite

import (
	"context"
	"fmt"
)

// ClustersService handles communication with cluster related
// methods of the Buildkite API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/clusters#clusters
type ClustersService struct {
	client *Client
}

type ClusterCreate struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Emoji       string `json:"emoji,omitempty"`
	Color       string `json:"color,omitempty"`
}

type ClusterUpdate struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	Emoji          string `json:"emoji,omitempty"`
	Color          string `json:"color,omitempty"`
	DefaultQueueID string `json:"default_queue_id,omitempty"`
}

type Cluster struct {
	ID              string         `json:"id,omitempty"`
	GraphQLID       string         `json:"graphql_id,omitempty"`
	DefaultQueueID  string         `json:"default_queue_id,omitempty"`
	Name            string         `json:"name,omitempty"`
	Description     string         `json:"description,omitempty"`
	Emoji           string         `json:"emoji,omitempty"`
	Color           string         `json:"color,omitempty"`
	URL             string         `json:"url,omitempty"`
	WebURL          string         `json:"web_url,omitempty"`
	QueuesURL       string         `json:"queues_url,omitempty"`
	DefaultQueueURL string         `json:"default_queue_url,omitempty"`
	CreatedAt       *Timestamp     `json:"created_at,omitempty"`
	CreatedBy       ClusterCreator `json:"created_by,omitempty"`
}

type ClusterCreator struct {
	ID        string     `json:"id,omitempty"`
	GraphQLID string     `json:"graphql_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

type ClustersListOptions struct{ ListOptions }

func (cs *ClustersService) List(ctx context.Context, org string, opt *ClustersListOptions) ([]Cluster, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := cs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var clusters []Cluster
	resp, err := cs.client.Do(req, &clusters)
	if err != nil {
		return nil, resp, err
	}

	return clusters, resp, err
}

func (cs *ClustersService) Get(ctx context.Context, org, id string) (Cluster, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s", org, id)
	req, err := cs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return Cluster{}, nil, err
	}

	var cluster Cluster
	resp, err := cs.client.Do(req, &cluster)

	if err != nil {
		return Cluster{}, resp, err
	}

	return cluster, resp, err
}

func (cs *ClustersService) Create(ctx context.Context, org string, cc ClusterCreate) (Cluster, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters", org)
	req, err := cs.client.NewRequest(ctx, "POST", u, cc)
	if err != nil {
		return Cluster{}, nil, err
	}

	var cluster Cluster
	resp, err := cs.client.Do(req, &cluster)
	if err != nil {
		return Cluster{}, resp, err
	}

	return cluster, resp, err
}

func (cs *ClustersService) Update(ctx context.Context, org, id string, cu ClusterUpdate) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s", org, id)
	req, err := cs.client.NewRequest(ctx, "PATCH", u, cu)
	if err != nil {
		return nil, nil
	}

	resp, err := cs.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (cs *ClustersService) Delete(ctx context.Context, org, id string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s", org, id)
	req, err := cs.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return cs.client.Do(req, nil)
}
