package buildkite

import (
	"context"
)

// MetricsService handles communication with the metrics related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/apis/agent-api#metrics
type MetricsService struct {
	client *Client
}

// QueueMetrics represents metrics for a specific queue
type QueueMetrics struct {
	Idle  int `json:"idle"`
	Busy  int `json:"busy"`
	Total int `json:"total"`
}

// AgentMetrics represents metrics about agents
type AgentMetrics struct {
	Idle   int                     `json:"idle"`
	Busy   int                     `json:"busy"`
	Total  int                     `json:"total"`
	Queues map[string]QueueMetrics `json:"queues"`
}

// JobMetrics represents metrics about jobs
type JobMetrics struct {
	Scheduled int                        `json:"scheduled"`
	Running   int                        `json:"running"`
	Waiting   int                        `json:"waiting"`
	Total     int                        `json:"total"`
	Queues    map[string]JobQueueMetrics `json:"queues"`
}

// JobQueueMetrics represents metrics for jobs in a specific queue
type JobQueueMetrics struct {
	Scheduled int `json:"scheduled"`
	Running   int `json:"running"`
	Waiting   int `json:"waiting"`
	Total     int `json:"total"`
}

// OrganizationMetrics represents organization information
type OrganizationMetrics struct {
	Slug string `json:"slug"`
}

// Metrics represents the full metrics response
type Metrics struct {
	Agents       AgentMetrics        `json:"agents"`
	Jobs         JobMetrics          `json:"jobs"`
	Organization OrganizationMetrics `json:"organization"`
}

// Get fetches the current metrics.
//
// buildkite API docs: https://buildkite.com/docs/apis/agent-api#metrics
func (ms *MetricsService) Get(ctx context.Context) (*Metrics, *Response, error) {
	u := "v3/metrics"
	req, err := ms.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	metrics := new(Metrics)
	resp, err := ms.client.Do(req, metrics)
	if err != nil {
		return nil, resp, err
	}

	return metrics, resp, nil
}
