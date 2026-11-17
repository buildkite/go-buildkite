package buildkite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClusterQueuesService_Insights(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		opt   *ClusterQueueInsightsOptions
		query string
	}{
		{name: "defaults"},
		{name: "empty options", opt: &ClusterQueueInsightsOptions{}},
		{
			name:  "history",
			opt:   &ClusterQueueInsightsOptions{Include: "history", TimePeriod: "24h"},
			query: "include=history&time_period=24h",
		},
		{
			name:  "history comparison without history include",
			opt:   &ClusterQueueInsightsOptions{Include: "history_comparison", TimePeriod: "7d"},
			query: "include=history_comparison&time_period=7d",
		},
		{
			name:  "hourly changes without a time period",
			opt:   &ClusterQueueInsightsOptions{Include: "hourly_changes"},
			query: "include=hourly_changes",
		},
		{
			name:  "all includes",
			opt:   &ClusterQueueInsightsOptions{Include: "history,history_comparison,hourly_changes", TimePeriod: "1h"},
			query: "include=history%2Chistory_comparison%2Chourly_changes&time_period=1h",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			server.HandleFunc("/v2/organizations/my-great-org/clusters/cluster-id/queues/queue-id/insights", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				if got := r.URL.RawQuery; got != tt.query {
					t.Errorf("query = %q, want %q", got, tt.query)
				}
				_, _ = fmt.Fprint(w, queueInsightsCurrentJSON)
			})

			got, resp, err := client.ClusterQueues.Insights(context.Background(), "my-great-org", "cluster-id", "queue-id", tt.opt)
			if err != nil {
				t.Fatalf("Insights returned error: %v", err)
			}
			if resp == nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("response = %v, want HTTP 200", resp)
			}
			if got.Resource.Type != "queue" || got.Resource.UUID != "queue-id" {
				t.Errorf("resource = %+v, want queue/queue-id", got.Resource)
			}
			if got.History != nil || got.HourlyChanges != nil {
				t.Error("absent optional sections must remain nil")
			}
			zero, three := int64(0), int64(3)
			if diff := cmp.Diff(map[string]*int64{"connected": &zero, "busy": nil, "idle": &zero, "paused": &zero, "draining": &zero}, got.Current.Agents.Values); diff != "" {
				t.Errorf("agent values mismatch (-want +got):\n%s", diff)
			}
			if got.Current.Agents.ObservedAt != nil || got.Current.Jobs.ObservedAt == nil || got.Current.Jobs.ObservedAt.Nanosecond() != 123456000 {
				t.Error("nullable timestamps or microsecond precision were lost")
			}
			if diff := cmp.Diff(map[string]*int64{"awaiting_execution": &three, "in_progress": &zero}, got.Current.Jobs.Computed); diff != "" {
				t.Errorf("computed counts mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClusterQueuesService_InsightsSections(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		body string
	}{
		{name: "current and recent", body: queueInsightsCurrentJSON},
		{name: "all sections", body: queueInsightsAllJSON},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/my-great-org/clusters/cluster-id/queues/queue-id/insights", func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprint(w, tt.body)
			})

			got, _, err := client.ClusterQueues.Insights(context.Background(), "my-great-org", "cluster-id", "queue-id", nil)
			if err != nil {
				t.Fatalf("Insights returned error: %v", err)
			}

			// Re-encoding must preserve nulls, fractional averages, coverage and limitation codes for MCP consumers.
			encoded, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("Marshal returned error: %v", err)
			}
			var wantJSON, gotJSON any
			if err := json.Unmarshal([]byte(tt.body), &wantJSON); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(encoded, &gotJSON); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
				t.Errorf("Insights response mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClusterQueuesService_InsightsAPIError(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/my-great-org/clusters/cluster-id/queues/queue-id/insights", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprint(w, `{"message":"Your plan does not include Cluster Insights."}`)
	})

	_, resp, err := client.ClusterQueues.Insights(context.Background(), "my-great-org", "cluster-id", "queue-id", nil)
	var apiErr *ErrorResponse
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %v, want *ErrorResponse", err)
	}
	if apiErr.Message != "Your plan does not include Cluster Insights." {
		t.Errorf("API error message = %q", apiErr.Message)
	}
	if resp == nil || resp.StatusCode != http.StatusForbidden {
		t.Errorf("response = %v, want HTTP 403", resp)
	}
}

const queueInsightsCurrentJSON = `{
  "resource":{"type":"queue","uuid":"queue-id"},
  "current":{
    "agents":{"observed_at":null,"source_high_watermark_at":null,"values":{"connected":0,"busy":null,"idle":0,"paused":0,"draining":0},"limitations":{"snapshot_evidence_unavailable":true}},
    "jobs":{"observed_at":"2026-11-12T08:30:00.123456Z","source_high_watermark_at":"2026-11-12T08:29:59Z","values":{"pending":0,"blocked_on_dependency":1,"unblocked":0,"limiting":0,"limited":0,"platform_limiting":0,"platform_limited":0,"scheduled":3,"reserved":0,"assigned":0,"accepted":0,"running":0,"canceling":0,"timing_out":0},"computed":{"awaiting_execution":3,"in_progress":0},"limitations":{"population_coverage_unverified":true}}
  },
  "recent":{
    "job_completions":{"window":{"from":"2026-11-12T08:29:00Z","to":"2026-11-12T08:30:00Z"},"observed_at":"2026-11-12T08:30:00Z","values":{"total":4,"passed":3,"failed":1},"limitations":{}},
    "wait_times":{"window":{"from":null,"to":null},"observed_at":null,"unit":"seconds","values":{"min":null,"p50":null,"p95":null,"max":null},"samples":{"valid_count":null,"incomplete_count":null},"limitations":{"missing_publication":true}}
  }
}`

// Series are shortened to keep the fixture readable; values follow the REST response contract.
const queueInsightsAllJSON = `{
  "resource":{"type":"queue","uuid":"queue-id"},
  "current":{
    "agents":{"observed_at":null,"source_high_watermark_at":null,"values":{"connected":0},"limitations":{}},
    "jobs":{"observed_at":null,"source_high_watermark_at":null,"values":{"running":null},"computed":{"awaiting_execution":null,"in_progress":null},"limitations":{"source_unavailable":true}}
  },
  "recent":{
    "job_completions":{"window":{"from":"2026-11-12T08:29:00Z","to":"2026-11-12T08:30:00Z"},"observed_at":"2026-11-12T08:30:00Z","values":{"total":0,"passed":0,"failed":0},"limitations":{}},
    "wait_times":{"window":{"from":"2026-11-12T08:29:00Z","to":"2026-11-12T08:30:00Z"},"observed_at":"2026-11-12T08:30:00Z","unit":"seconds","values":{"min":0,"p50":1.25,"p95":4.5,"max":6},"samples":{"valid_count":4,"incomplete_count":1},"limitations":{"incomplete_wait_samples":true}}
  },
  "history":{
    "selected_period":{
      "window":{"from":"2026-11-12T07:30:00Z","to":"2026-11-12T08:30:00Z","interval_seconds":60},
      "counts":{"aggregation":"last","agents":{"values":{"connected":[0,null]},"limitations":{"missing_publication":[1]}},"jobs":{"values":{"awaiting_execution":[2,null],"in_progress":[1,null]},"limitations":{"missing_publication":[1]}}},
      "job_completions":{"values":{"total":[4,0],"passed":[3,0],"failed":[1,0]},"limitations":{}},
      "wait_times":{"unit":"seconds","values":{"min":[0.25,null],"p50":[1.5,null],"p95":[3.75,null],"max":[5,null]},"samples":{"valid_count":[3,0],"incomplete_count":[1,0]},"limitations":{"no_samples":[1],"incomplete_wait_samples":[0]}}
    },
    "previous_period":{
      "window":{"from":"2026-11-12T06:30:00Z","to":"2026-11-12T07:30:00Z","interval_seconds":60},
      "counts":{"aggregation":"last","agents":{"values":{"connected":[1,0]},"limitations":{}},"jobs":{"values":{"awaiting_execution":[1,1],"in_progress":[0,0]},"limitations":{}}},
      "job_completions":{"values":{"total":[2,0],"passed":[1,0],"failed":[1,0]},"limitations":{}},
      "wait_times":{"unit":"seconds","values":{"min":[1,1],"p50":[2,2],"p95":[4,4],"max":[6,6]},"samples":{"valid_count":[2,1],"incomplete_count":[0,0]},"limitations":{}}
    },
    "changes":{
      "counts":{"basis":"last_bucket","agents":{"connected":{"selected":null,"previous":0,"difference":null,"relative_change_percent":null,"limitations":{"selected_value_unavailable":true}}},"jobs":{"awaiting_execution":{"selected":null,"previous":1,"difference":null,"relative_change_percent":null,"limitations":{"selected_value_unavailable":true}},"in_progress":{"selected":null,"previous":0,"difference":null,"relative_change_percent":null,"limitations":{"selected_value_unavailable":true}}}},
      "job_completions":{"basis":"period_total","total":{"selected":4,"previous":2,"difference":2,"relative_change_percent":100,"limitations":{}},"pass_rate":{"unit":"percent","selected":75,"previous":50,"percentage_point_change":25,"limitations":{}}},
      "wait_times":{"basis":"last_bucket","unit":"seconds","min":{"selected":null,"previous":1,"difference":null,"relative_change_percent":null,"limitations":{"selected_value_unavailable":true}},"p50":{"selected":null,"previous":2,"difference":null,"relative_change_percent":null,"limitations":{"selected_value_unavailable":true}},"p95":{"selected":null,"previous":4,"difference":null,"relative_change_percent":null,"limitations":{"selected_value_unavailable":true}},"max":{"selected":null,"previous":6,"difference":null,"relative_change_percent":null,"limitations":{"selected_value_unavailable":true}}}
    }
  },
  "hourly_changes":{
    "aggregation":"time_weighted_average","interval_seconds":3600,
    "current":{"reference_at":"2026-11-12T08:00:00Z","window":{"from":"2026-11-12T07:00:00Z","to":"2026-11-12T08:30:00Z"},"agents":{"values":{"connected":[1.25,0]},"covered_seconds":[3600,1799.25],"limitations":{"partial_bucket":[1]}},"jobs":{"values":{"awaiting_execution":[2.5,null],"in_progress":[0.75,null]},"covered_seconds":[3600,null],"limitations":{"missing_publication":[1]}}},
    "yesterday":{"reference_at":"2026-11-11T08:00:00Z","window":{"from":"2026-11-11T07:00:00Z","to":"2026-11-11T10:00:00Z"},"agents":{"values":{"connected":[1,0,2]},"covered_seconds":[3600,3600,3600],"limitations":{}},"jobs":{"values":{"awaiting_execution":[5,0,1],"in_progress":[0,1,2]},"covered_seconds":[3600,3600,3600],"limitations":{}}},
    "last_week":{"reference_at":"2026-11-05T08:00:00Z","window":{"from":"2026-11-05T07:00:00Z","to":"2026-11-05T10:00:00Z"},"agents":{"values":{"connected":[null,null,null]},"covered_seconds":[null,null,null],"limitations":{"missing_publication":[0,1,2]}},"jobs":{"values":{"awaiting_execution":[null,null,null],"in_progress":[null,null,null]},"covered_seconds":[null,null,null],"limitations":{"missing_publication":[0,1,2]}}},
    "changes":{"basis":"completed_hour","yesterday":{"agents":{"connected":{"selected":1.25,"previous":1,"difference":0.25,"relative_change_percent":25,"limitations":{}}},"jobs":{"awaiting_execution":{"selected":2.5,"previous":5,"difference":-2.5,"relative_change_percent":-50,"limitations":{}},"in_progress":{"selected":0.75,"previous":0,"difference":0.75,"relative_change_percent":null,"limitations":{}}}},"last_week":{"agents":{"connected":{"selected":1.25,"previous":null,"difference":null,"relative_change_percent":null,"limitations":{"previous_value_unavailable":true}}},"jobs":{"awaiting_execution":{"selected":2.5,"previous":null,"difference":null,"relative_change_percent":null,"limitations":{"previous_value_unavailable":true}},"in_progress":{"selected":0.75,"previous":null,"difference":null,"relative_change_percent":null,"limitations":{"previous_value_unavailable":true}}}}}
  }
}`
