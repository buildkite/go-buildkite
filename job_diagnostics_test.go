package buildkite

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJobsService_GetJobDiagnosticsByOrg(t *testing.T) {
	t.Parallel()

	payload, err := os.ReadFile("testdata/job_diagnostics.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name      string
		opt       *JobDiagnosticsOptions
		query     string
		overrides map[string]json.RawMessage
	}{
		{name: "defaults"},
		{name: "empty options", opt: &JobDiagnosticsOptions{}},
		{name: "includes", opt: &JobDiagnosticsOptions{Include: "events,dependencies,insights"}, query: "include=events%2Cdependencies%2Cinsights"},
		{name: "automatic evidence without guidance", opt: &JobDiagnosticsOptions{Include: "auto"}, query: "include=auto"},
		{name: "guidance with automatic and explicit evidence", opt: &JobDiagnosticsOptions{Include: "auto,events,guidance"}, query: "include=auto%2Cevents%2Cguidance", overrides: map[string]json.RawMessage{
			"guidance": json.RawMessage(`{"findings":[{
				"code":"manual_action_required",
				"kind":"finding",
				"message":"This manual job is blocked and requires user input or approval to continue.",
				"evidence":["/details/type","/details/lifecycle/state","/details/unblockable"],
				"limits":["The unblockable flag does not grant the caller permission or guarantee that an unblock request will succeed."],
				"next_check":"Ask an authorized operator to review the manual step's inputs and confirm that unblocking is appropriate."
			}]}`),
		}},
		{name: "guidance without findings", opt: &JobDiagnosticsOptions{Include: "guidance"}, query: "include=guidance", overrides: map[string]json.RawMessage{
			"guidance": json.RawMessage(`{"findings":[]}`),
		}},
		{name: "expanded dependencies and dynamic payloads", overrides: map[string]json.RawMessage{
			"dependencies": json.RawMessage(`{"selection":{"order":"breadth_first","max_steps":50,"max_edges":100},"items":[{"required_by":{"type":"step","uuid":"step-id","key":null},"step":{"type":"step","uuid":null,"key":"missing-step"},"source":"direct","allow_failure":false,"inherited_from":null},{"required_by":{"type":"step","uuid":"step-id","key":null},"step":{"type":"step","uuid":"upstream-id","key":"upstream"},"source":"upload","allow_failure":true,"inherited_from":{"type":"step","uuid":"group-id","key":null}}],"unexpanded_steps":[{"type":"step","uuid":"upstream-id","key":"upstream"}],"limitations":{"truncated":true},"links":[]}`),
			"errors":       json.RawMessage(`{"items":[{"future":{"count":9007199254740993,"values":[null,false,0,{},[]]}}],"limitations":{}}`),
			"history":      json.RawMessage(`{"recent_runs":{"items":[],"limitations":{}}}`),
			"capacity":     json.RawMessage(`{"requirements":{"agent_query_rules":[]},"evidence":{"future":{"value":null}},"limitations":{},"links":[]}`),
		}},
		{name: "unavailable sections", overrides: map[string]json.RawMessage{
			"retries":      json.RawMessage(`{"requirements":null,"evidence":null,"limitations":{"not_applicable":true},"links":[]}`),
			"capacity":     json.RawMessage(`{"requirements":null,"evidence":null,"limitations":{"not_applicable":true},"links":[]}`),
			"events":       json.RawMessage(`{"window":{"from":"2026-11-12T02:34:56.123456Z","to":"2026-11-19T02:34:56.123456Z"},"selection":{"order":"newest_first","max_items":100},"items":[],"limitations":{}}`),
			"dependencies": json.RawMessage(`{"selection":{"order":"breadth_first","max_steps":50,"max_edges":100},"items":[],"limitations":{},"links":[]}`),
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			responseBody := payload
			if tt.overrides != nil {
				var sections map[string]json.RawMessage
				if err := json.Unmarshal(payload, &sections); err != nil {
					t.Fatal(err)
				}
				for key, value := range tt.overrides {
					sections[key] = value
				}
				var err error
				responseBody, err = json.Marshal(sections)
				if err != nil {
					t.Fatal(err)
				}
			}
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/jobs/job-id/diagnostics", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				if r.URL.RawQuery != tt.query {
					t.Errorf("query = %q, want %q", r.URL.RawQuery, tt.query)
				}
				_, _ = w.Write(responseBody)
			})
			got, resp, err := client.Jobs.GetJobDiagnosticsByOrg(context.Background(), "acme", "job-id", tt.opt)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("response = %v, want HTTP 200", resp)
			}
			if got.Details.Priority == nil || *got.Details.Priority != 0 || got.Details.SoftFailed {
				t.Fatalf("zero priority or false soft_failed lost: %+v", got.Details)
			}
			if got.Details.ExitStatus != nil || got.Details.Timings.ExecutionSeconds != nil || got.Details.Unblockable == nil || *got.Details.Unblockable {
				t.Fatalf("nullable fields lost: %+v", got.Details)
			}
			if tt.overrides == nil {
				if got.Events.Items[0].Type != "future_event" || string(got.Events.Items[0].Details["large_integer"]) != "9007199254740993" {
					t.Fatalf("event details lost: %+v", got.Events.Items)
				}
				if got.Dependencies.Items != nil || got.Context.Builds == nil || len(got.Context.Builds) != 0 || got.Insights[0].Current != nil {
					t.Fatal("nulls and empty collections must remain distinct")
				}
			}
			encoded, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			var wantJSON, gotJSON any
			wantDecoder := json.NewDecoder(bytes.NewReader(responseBody))
			wantDecoder.UseNumber()
			if err := wantDecoder.Decode(&wantJSON); err != nil {
				t.Fatal(err)
			}
			gotDecoder := json.NewDecoder(bytes.NewReader(encoded))
			gotDecoder.UseNumber()
			if err := gotDecoder.Decode(&gotJSON); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
				t.Errorf("diagnostics mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
