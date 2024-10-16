package buildkite

import (
	"context"
	"fmt"
)

// ClusterQueuesService handles communication with cluster queue related
// methods of the Buildkite API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/clusters#cluster-queues
type ClusterQueuesService struct {
	client *Client
}

type ClusterQueue struct {
	ID                 string          `json:"id,omitempty"`
	GraphQLID          string          `json:"graphql_id,omitempty"`
	Key                string          `json:"key,omitempty"`
	Description        string          `json:"description,omitempty"`
	URL                string          `json:"url,omitempty"`
	WebURL             string          `json:"web_url,omitempty"`
	ClusterURL         string          `json:"cluster_url,omitempty"`
	DispatchPaused     bool            `json:"dispatch_paused,omitempty"`
	DispatchPausedBy   *ClusterCreator `json:"dispatch_paused_by,omitempty"`
	DispatchPausedAt   *Timestamp      `json:"dispatch_paused_at,omitempty"`
	DispatchPausedNote string          `json:"dispatch_paused_note,omitempty"`
	CreatedAt          *Timestamp      `json:"created_at,omitempty"`
	CreatedBy          ClusterCreator  `json:"created_by,omitempty"`
}

type ClusterQueueCreate struct {
	Key         string `json:"key,omitempty"`
	Description string `json:"description,omitempty"`
}

type ClusterQueueUpdate struct {
	Description string `json:"description,omitempty"`
}

type ClusterQueuePause struct {
	Note string `json:"dispatch_paused_note,omitempty"`
}

type ClusterQueuesListOptions struct {
	ListOptions
}

func (cqs *ClusterQueuesService) List(ctx context.Context, org, clusterID string, opt *ClusterQueuesListOptions) ([]ClusterQueue, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/queues", org, clusterID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := cqs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var queues []ClusterQueue
	resp, err := cqs.client.Do(req, &queues)
	if err != nil {
		return nil, resp, err
	}

	return queues, resp, err
}

func (cqs *ClusterQueuesService) Get(ctx context.Context, org, clusterID, queueID string) (ClusterQueue, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/queues/%s", org, clusterID, queueID)
	req, err := cqs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return ClusterQueue{}, nil, err
	}

	var queue ClusterQueue
	resp, err := cqs.client.Do(req, &queue)
	if err != nil {
		return ClusterQueue{}, resp, err
	}

	return queue, resp, err
}

func (cqs *ClusterQueuesService) Create(ctx context.Context, org, clusterID string, qc ClusterQueueCreate) (ClusterQueue, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/queues", org, clusterID)
	req, err := cqs.client.NewRequest(ctx, "POST", u, qc)
	if err != nil {
		return ClusterQueue{}, nil, err
	}

	var queue ClusterQueue
	resp, err := cqs.client.Do(req, &queue)

	if err != nil {
		return ClusterQueue{}, resp, err
	}

	return queue, resp, err
}

func (cqs *ClusterQueuesService) Update(ctx context.Context, org, clusterID, queueID string, qu ClusterQueueUpdate) (ClusterQueue, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/queues/%s", org, clusterID, queueID)
	req, err := cqs.client.NewRequest(ctx, "PATCH", u, qu)
	if err != nil {
		return ClusterQueue{}, nil, err
	}

	var cq ClusterQueue
	resp, err := cqs.client.Do(req, &cq)
	if err != nil {
		return ClusterQueue{}, resp, err
	}

	return cq, resp, err
}

func (cqs *ClusterQueuesService) Delete(ctx context.Context, org, clusterID, queueID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/queues/%s", org, clusterID, queueID)
	req, err := cqs.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return cqs.client.Do(req, nil)
}

func (cqs *ClusterQueuesService) Pause(ctx context.Context, org, clusterID, queueID string, qp ClusterQueuePause) (ClusterQueue, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/queues/%s/pause_dispatch", org, clusterID, queueID)
	req, err := cqs.client.NewRequest(ctx, "POST", u, qp)
	if err != nil {
		return ClusterQueue{}, nil, err
	}

	var cq ClusterQueue
	resp, err := cqs.client.Do(req, &cq)
	if err != nil {
		return ClusterQueue{}, resp, err
	}

	return cq, resp, err
}

func (cqs *ClusterQueuesService) Resume(ctx context.Context, org, clusterID, queueID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/queues/%s/resume_dispatch", org, clusterID, queueID)
	req, err := cqs.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, err
	}

	return cqs.client.Do(req, nil)
}
