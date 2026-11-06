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

type RetryAgentAffinity string

const (
	RetryAgentAffinityPreferWarmest   RetryAgentAffinity = "prefer-warmest"
	RetryAgentAffinityPreferDifferent RetryAgentAffinity = "prefer-different"
)

type ClusterQueue struct {
	ID                 string                    `json:"id,omitempty"`
	GraphQLID          string                    `json:"graphql_id,omitempty"`
	Key                string                    `json:"key,omitempty"`
	Description        string                    `json:"description,omitempty"`
	URL                string                    `json:"url,omitempty"`
	WebURL             string                    `json:"web_url,omitempty"`
	ClusterURL         string                    `json:"cluster_url,omitempty"`
	RetryAgentAffinity RetryAgentAffinity        `json:"retry_agent_affinity,omitempty"`
	DispatchPaused     bool                      `json:"dispatch_paused"`
	DispatchPausedBy   *ClusterCreator           `json:"dispatch_paused_by,omitempty"`
	DispatchPausedAt   *Timestamp                `json:"dispatch_paused_at,omitempty"`
	DispatchPausedNote string                    `json:"dispatch_paused_note,omitempty"`
	CreatedAt          *Timestamp                `json:"created_at,omitempty"`
	CreatedBy          ClusterCreator            `json:"created_by,omitempty"`
	Hosted             bool                      `json:"hosted,omitempty"`
	HostedAgents       *ClusterQueueHostedAgents `json:"hosted_agents,omitempty"`
}

// ClusterQueueHostedAgents describes a hosted queue's current configuration.
type ClusterQueueHostedAgents struct {
	InstanceShape struct {
		Name         string `json:"name"`
		MachineType  string `json:"machine_type"`
		Architecture string `json:"architecture"`
		CPU          int    `json:"cpu"`
		Memory       int    `json:"memory"`
	} `json:"instance_shape"`
	PlatformSettings struct {
		Linux *struct {
			AgentImageRef string `json:"agent_image_ref"`
		} `json:"linux,omitempty"`
		MacOS *struct {
			BaseImageXcodeVersion string `json:"base_image_xcode_version"`
			MacOSVersion          string `json:"macos_version"`
		} `json:"macos,omitempty"`
	} `json:"platform_settings"`
}

// ClusterQueueHostedAgentsCreate selects a hosted instance shape and optional macOS image.
type ClusterQueueHostedAgentsCreate struct {
	InstanceShape string `json:"instanceShape"`
	XcodeVersion  string `json:"xcodeVersion,omitempty"`
	MacOSVersion  string `json:"macosVersion,omitempty"`
}

// ClusterQueueCreate represents the request body for creating a cluster queue.
type ClusterQueueCreate struct {
	Key                string                          `json:"key,omitempty"`
	Description        string                          `json:"description,omitempty"`
	RetryAgentAffinity RetryAgentAffinity              `json:"retry_agent_affinity,omitempty"`
	HostedAgents       *ClusterQueueHostedAgentsCreate `json:"hostedAgents,omitempty"`
}

// ClusterQueueUpdate represents the request body for updating a cluster queue.
type ClusterQueueUpdate struct {
	Description        Optional[string]             `json:"description,omitzero"`
	RetryAgentAffinity Optional[RetryAgentAffinity] `json:"retry_agent_affinity,omitzero"`
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

// Recommendations returns eligible missing starter queues as creation payloads.
// Requires read_clusters scope and permission to manage the cluster.
func (cqs *ClusterQueuesService) Recommendations(ctx context.Context, org, clusterID string) ([]ClusterQueueCreate, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/queues/recommendations", org, clusterID)
	req, err := cqs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var recommendations []ClusterQueueCreate
	resp, err := cqs.client.Do(req, &recommendations)
	if err != nil {
		return nil, resp, err
	}

	return recommendations, resp, nil
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
