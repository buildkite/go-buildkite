package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMetricsService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v3/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"agents": {
				"idle": 1,
				"busy": 0,
				"total": 1,
				"queues": {
					"default": {
						"idle": 1,
						"busy": 0,
						"total": 1
					}
				}
			},
			"jobs": {
				"scheduled": 5,
				"running": 0,
				"waiting": 0,
				"total": 5,
				"queues": {
					"default": {
						"scheduled": 5,
						"running": 0,
						"waiting": 0,
						"total": 5
					}
				}
			},
			"organization": {
				"slug": "buildkite"
			}
		}`)
	})

	metrics, _, err := client.Metrics.Get(context.Background())
	if err != nil {
		t.Errorf("Metrics.Get returned error: %v", err)
	}

	want := &Metrics{
		Agents: AgentMetrics{
			Idle:  1,
			Busy:  0,
			Total: 1,
			Queues: map[string]QueueMetrics{
				"default": {
					Idle:  1,
					Busy:  0,
					Total: 1,
				},
			},
		},
		Jobs: JobMetrics{
			Scheduled: 5,
			Running:   0,
			Waiting:   0,
			Total:     5,
			Queues: map[string]JobQueueMetrics{
				"default": {
					Scheduled: 5,
					Running:   0,
					Waiting:   0,
					Total:     5,
				},
			},
		},
		Organization: OrganizationMetrics{
			Slug: "buildkite",
		},
	}

	if diff := cmp.Diff(metrics, want); diff != "" {
		t.Errorf("Metrics.Get returned diff: (-got +want)\n%s", diff)
	}
}
