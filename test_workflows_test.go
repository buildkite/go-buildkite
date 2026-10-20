package buildkite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTestWorkflowsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-org/suites/my-suite/workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2", "per_page": "5"})
		w.Header().Set("Link", `<https://api.buildkite.com/v2/analytics/organizations/my-org/suites/my-suite/workflows?page=3&per_page=5>; rel="next", <https://api.buildkite.com/v2/analytics/organizations/my-org/suites/my-suite/workflows?page=4&per_page=5>; rel="last"`)
		_, _ = fmt.Fprint(w, `[
			{"id":"transition", "monitor":{"type":"transition_count","window_size":100,"alarm_threshold":0.1,"recover_threshold":0}},
			{"id":"probabilistic", "monitor":{"type":"probabilistic_flakiness","alarm_threshold":0.15,"recover_threshold":0.04}},
			{"id":"duration", "monitor":{"type":"duration_threshold","aggregation":"p50","evaluation_window":10,"alarm_threshold_ms":5000,"recover_threshold_ms":3000}},
			{"id":"retry", "monitor":{"type":"passed_on_retry"}},
			{"id":"new", "monitor":{"type":"new_test"}}
		]`)
	})

	got, resp, err := client.TestWorkflows.List(context.Background(), "my-org", "my-suite", &TestWorkflowListOptions{
		ListOptions: ListOptions{Page: 2, PerPage: 5},
	})
	if err != nil {
		t.Fatal(err)
	}
	window, evaluationWindow, alarmMS, recoverMS := 100, 10, 5000, 3000
	transitionAlarm, probabilisticAlarm, probabilisticRecover, zero := 0.1, 0.15, 0.04, 0.0
	want := []TestWorkflow{
		{ID: "transition", Monitor: TestWorkflowMonitor{Type: "transition_count", WindowSize: &window, AlarmThreshold: &transitionAlarm, RecoverThreshold: &zero}},
		{ID: "probabilistic", Monitor: TestWorkflowMonitor{Type: "probabilistic_flakiness", AlarmThreshold: &probabilisticAlarm, RecoverThreshold: &probabilisticRecover}},
		{ID: "duration", Monitor: TestWorkflowMonitor{Type: "duration_threshold", Aggregation: "p50", EvaluationWindow: &evaluationWindow, AlarmThresholdMS: &alarmMS, RecoverThresholdMS: &recoverMS}},
		{ID: "retry", Monitor: TestWorkflowMonitor{Type: "passed_on_retry"}},
		{ID: "new", Monitor: TestWorkflowMonitor{Type: "new_test"}},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("workflows (-want +got):\n%s", diff)
	}
	if resp.NextPage != 3 || resp.LastPage != 4 {
		t.Errorf("pagination: next=%d last=%d, want 3 and 4", resp.NextPage, resp.LastPage)
	}
	// Consumers such as MCP must not lose zero thresholds when encoding results.
	if len(got) != len(want) {
		t.Fatalf("got %d workflows, want %d", len(got), len(want))
	}
	encoded, err := json.Marshal(got[0].Monitor)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, string(encoded), `{"type":"transition_count","window_size":100,"alarm_threshold":0.1,"recover_threshold":0}`)
}

func TestTestWorkflowsService_ListEmpty(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/analytics/organizations/my-org/suites/my-suite/workflows", func(w http.ResponseWriter, r *http.Request) {
		testFormValues(t, r, values{})
		_, _ = fmt.Fprint(w, `[]`)
	})
	got, _, err := client.TestWorkflows.List(context.Background(), "my-org", "my-suite", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("got %#v, want empty list", got)
	}
}

func TestTestWorkflowsService_Get(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/analytics/organizations/my-org/suites/my-suite/workflows/workflow-id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{
			"id":"workflow-id",
			"url":"https://api.buildkite.com/v2/analytics/organizations/my-org/suites/my-suite/workflows/workflow-id",
			"web_url":"https://buildkite.com/organizations/my-org/analytics/suites/my-suite/workflows/workflow-id/events",
			"monitor":{"type":"passed_on_retry"},
			"filters":[{"key":"scm.branch","operator":"is","value":"main"}],
			"alarm_actions":[
				{"type":"add_label","label":"flaky"},
				{"type":"change_state","state":"muted"},
				{"type":"send_webhook_notification","endpoint":"https://example.com/[REDACTED]"},
				{"type":"send_slack_notification","slack_workspace_id":"slack-id","channel":"#tests","message_template":"{{subject.test_location}} failed"},
				{"type":"create_linear_issue","linear_team_id":"team-id","title_template":"Flaky test","description_template":"Investigate {{subject.test_location}}"}
			],
			"recover_actions":[{"type":"remove_label","label":null},{"type":"trigger_pipeline","pipeline":"legacy-pipeline"},{"type":"future_action","future_field":true}],
			"created_at":"2026-10-19T01:00:00Z",
			"updated_at":"2026-10-20T02:30:00Z",
			"created_by":{"id":"user-id","graphql_id":"graphql-id","name":"Wendy","email":"wendy@example.com"},
			"future_field":true
		}`)
	})
	got, _, err := client.TestWorkflows.Get(context.Background(), "my-org", "my-suite", "workflow-id")
	if err != nil {
		t.Fatal(err)
	}
	label := "flaky"
	want := TestWorkflow{
		ID:      "workflow-id",
		URL:     "https://api.buildkite.com/v2/analytics/organizations/my-org/suites/my-suite/workflows/workflow-id",
		WebURL:  "https://buildkite.com/organizations/my-org/analytics/suites/my-suite/workflows/workflow-id/events",
		Monitor: TestWorkflowMonitor{Type: "passed_on_retry"},
		Filters: []TestWorkflowFilter{{Key: "scm.branch", Operator: "is", Value: "main"}},
		AlarmActions: []TestWorkflowAction{
			{Type: "add_label", Label: &label},
			{Type: "change_state", State: "muted"},
			{Type: "send_webhook_notification", Endpoint: "https://example.com/[REDACTED]"},
			{Type: "send_slack_notification", SlackWorkspaceID: "slack-id", Channel: "#tests", MessageTemplate: "{{subject.test_location}} failed"},
			{Type: "create_linear_issue", LinearTeamID: "team-id", TitleTemplate: "Flaky test", DescriptionTemplate: "Investigate {{subject.test_location}}"},
		},
		RecoverActions: []TestWorkflowAction{{Type: "remove_label"}, {Type: "trigger_pipeline", Pipeline: "legacy-pipeline"}, {Type: "future_action"}},
		CreatedAt:      NewTimestamp(time.Date(2026, 10, 19, 1, 0, 0, 0, time.UTC)),
		UpdatedAt:      NewTimestamp(time.Date(2026, 10, 20, 2, 30, 0, 0, time.UTC)),
		CreatedBy:      &TestWorkflowCreator{ID: "user-id", GraphQLID: "graphql-id", Name: "Wendy", Email: "wendy@example.com"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("workflow (-want +got):\n%s", diff)
	}
}

func TestTestWorkflowsService_Create(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		monitor TestWorkflowMonitorInput
		json    string
	}{
		{"transition_count", TestWorkflowMonitorInput{Type: Some("transition_count"), WindowSize: Some(100), AlarmThreshold: Some(0.1), RecoverThreshold: Some(0.0)}, `{"type":"transition_count","window_size":100,"alarm_threshold":0.1,"recover_threshold":0}`},
		{"probabilistic_flakiness", TestWorkflowMonitorInput{Type: Some("probabilistic_flakiness"), AlarmThreshold: Some(0.15), RecoverThreshold: Some(0.04)}, `{"type":"probabilistic_flakiness","alarm_threshold":0.15,"recover_threshold":0.04}`},
		{"duration_threshold", TestWorkflowMonitorInput{Type: Some("duration_threshold"), Aggregation: Some("p50"), EvaluationWindow: Some(10), AlarmThresholdMS: Some(5000), RecoverThresholdMS: Some(0)}, `{"type":"duration_threshold","aggregation":"p50","evaluation_window":10,"alarm_threshold_ms":5000,"recover_threshold_ms":0}`},
		{"passed_on_retry", TestWorkflowMonitorInput{Type: Some("passed_on_retry")}, `{"type":"passed_on_retry"}`},
		{"new_test", TestWorkflowMonitorInput{Type: Some("new_test")}, `{"type":"new_test"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/analytics/organizations/my-org/suites/my-suite/workflows", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				if got := r.Header.Get("Content-Type"); got != "application/json" {
					t.Errorf("Content-Type=%q, want application/json", got)
				}
				assertRequestJSON(t, r, `{"monitor":`+tc.json+`,"filters":[{"key":"scm.branch","operator":"is","value":"main"}]}`)
				w.WriteHeader(http.StatusCreated)
				_, _ = fmt.Fprint(w, `{"id":"created-id","monitor":`+tc.json+`,"created_by":null}`)
			})
			got, resp, err := client.TestWorkflows.Create(context.Background(), "my-org", "my-suite", TestWorkflowCreate{
				Monitor: tc.monitor,
				Filters: []TestWorkflowFilter{{Key: "scm.branch", Operator: "is", Value: "main"}},
			})
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != "created-id" || got.Monitor.Type != tc.name || got.CreatedBy != nil || resp.StatusCode != http.StatusCreated {
				t.Errorf("unexpected workflow or status: %#v, %d", got, resp.StatusCode)
			}
		})
	}
}

func TestTestWorkflowsService_CreateActions(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/analytics/organizations/my-org/suites/my-suite/workflows", func(w http.ResponseWriter, r *http.Request) {
		assertRequestJSON(t, r, `{
			"monitor":{"type":"passed_on_retry"},
			"filters":[{"key":"cloud.region","operator":"is_not","value":"us-east-1"},{"key":"build.label","operator":"starts_with","value":"nightly"}],
			"alarm_actions":[
				{"type":"add_label","label":"flaky"},
				{"type":"change_state","state":"muted"},
				{"type":"send_webhook_notification","endpoint":"https://example.com/hook"},
				{"type":"send_slack_notification","slack_workspace_id":"slack-id","channel":"#tests","message_template":"{{subject.test_location}} failed"},
				{"type":"create_linear_issue","linear_team_id":"team-id","title_template":"Flaky test","description_template":"Investigate {{subject.test_location}}"}
			],
			"recover_actions":[{"type":"remove_label","label":"flaky"}]
		}`)
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, `{"id":"created-id"}`)
	})
	_, _, err := client.TestWorkflows.Create(context.Background(), "my-org", "my-suite", TestWorkflowCreate{
		Monitor: TestWorkflowMonitorInput{Type: Some("passed_on_retry")},
		Filters: []TestWorkflowFilter{{Key: "cloud.region", Operator: "is_not", Value: "us-east-1"}, {Key: "build.label", Operator: "starts_with", Value: "nightly"}},
		AlarmActions: []TestWorkflowActionInput{
			{Type: "add_label", Label: Some("flaky")},
			{Type: "change_state", State: Some("muted")},
			{Type: "send_webhook_notification", Endpoint: Some("https://example.com/hook")},
			{Type: "send_slack_notification", SlackWorkspaceID: Some("slack-id"), Channel: Some("#tests"), MessageTemplate: Some("{{subject.test_location}} failed")},
			{Type: "create_linear_issue", LinearTeamID: Some("team-id"), TitleTemplate: Some("Flaky test"), DescriptionTemplate: Some("Investigate {{subject.test_location}}")},
		},
		RecoverActions: []TestWorkflowActionInput{{Type: "remove_label", Label: Some("flaky")}},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTestWorkflowsService_Update(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		input TestWorkflowUpdate
		json  string
	}{
		{"omitted", TestWorkflowUpdate{}, `{}`},
		{"threshold only", TestWorkflowUpdate{Monitor: Some(TestWorkflowMonitorInput{AlarmThreshold: Some(0.2)})}, `{"monitor":{"alarm_threshold":0.2}}`},
		{"zero score", TestWorkflowUpdate{Monitor: Some(TestWorkflowMonitorInput{RecoverThreshold: Some(0.0)})}, `{"monitor":{"recover_threshold":0}}`},
		{"zero duration", TestWorkflowUpdate{Monitor: Some(TestWorkflowMonitorInput{RecoverThresholdMS: Some(0)})}, `{"monitor":{"recover_threshold_ms":0}}`},
		{"clear lists", TestWorkflowUpdate{Filters: Some([]TestWorkflowFilter{}), AlarmActions: Some([]TestWorkflowActionInput{}), RecoverActions: Some([]TestWorkflowActionInput{})}, `{"filters":[],"alarm_actions":[],"recover_actions":[]}`},
		{"replace lists", TestWorkflowUpdate{
			Filters:        Some([]TestWorkflowFilter{{Key: "scm.branch", Operator: "is", Value: "release"}}),
			AlarmActions:   Some([]TestWorkflowActionInput{{Type: "change_state", State: Some("skipped")}}),
			RecoverActions: Some([]TestWorkflowActionInput{{Type: "change_state", State: Some("enabled")}}),
		}, `{"filters":[{"key":"scm.branch","operator":"is","value":"release"}],"alarm_actions":[{"type":"change_state","state":"skipped"}],"recover_actions":[{"type":"change_state","state":"enabled"}]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/analytics/organizations/my-org/suites/my-suite/workflows/workflow-id", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "PATCH")
				assertRequestJSON(t, r, tc.json)
				_, _ = fmt.Fprint(w, `{"id":"workflow-id","monitor":{"type":"probabilistic_flakiness","alarm_threshold":0.2,"recover_threshold":0},"filters":[],"alarm_actions":[],"recover_actions":[]}`)
			})
			got, _, err := client.TestWorkflows.Update(context.Background(), "my-org", "my-suite", "workflow-id", tc.input)
			if err != nil {
				t.Fatal(err)
			}
			alarm, recover := 0.2, 0.0
			want := TestWorkflow{
				ID: "workflow-id", Monitor: TestWorkflowMonitor{Type: "probabilistic_flakiness", AlarmThreshold: &alarm, RecoverThreshold: &recover},
				Filters: []TestWorkflowFilter{}, AlarmActions: []TestWorkflowAction{}, RecoverActions: []TestWorkflowAction{},
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("workflow (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTestWorkflowsService_Delete(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/analytics/organizations/my-org/suites/my-suite/workflows/workflow-id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	resp, err := client.TestWorkflows.Delete(context.Background(), "my-org", "my-suite", "workflow-id")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status=%d, want 204", resp.StatusCode)
	}
}

func TestTestWorkflowsService_ErrorsAndCancellation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		method string
		status int
		call   func(context.Context, *TestWorkflowsService) (*Response, error)
	}{
		{"GET list", 400, func(ctx context.Context, s *TestWorkflowsService) (*Response, error) {
			_, resp, err := s.List(ctx, "my-org", "my-suite", &TestWorkflowListOptions{ListOptions: ListOptions{Page: -1}})
			return resp, err
		}},
		{"GET", 404, func(ctx context.Context, s *TestWorkflowsService) (*Response, error) {
			_, resp, err := s.Get(ctx, "my-org", "my-suite", "workflow-id")
			return resp, err
		}},
		{"POST", 422, func(ctx context.Context, s *TestWorkflowsService) (*Response, error) {
			_, resp, err := s.Create(ctx, "my-org", "my-suite", TestWorkflowCreate{
				Monitor:      TestWorkflowMonitorInput{Type: Some("unsupported"), Aggregation: Some("")},
				AlarmActions: []TestWorkflowActionInput{{Type: "add_label", Label: Some("Not Normalized"), Channel: Some("")}},
			})
			return resp, err
		}},
		{"PATCH", 503, func(ctx context.Context, s *TestWorkflowsService) (*Response, error) {
			_, resp, err := s.Update(ctx, "my-org", "my-suite", "workflow-id", TestWorkflowUpdate{})
			return resp, err
		}},
		{"DELETE", 403, func(ctx context.Context, s *TestWorkflowsService) (*Response, error) {
			return s.Delete(ctx, "my-org", "my-suite", "workflow-id")
		}},
	} {
		t.Run(tc.method, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			body := `{"message":"Validation Failed","errors":[{"field":"/monitor/type","code":"invalid"}]}`
			server.HandleFunc("/v2/analytics/organizations/my-org/suites/my-suite/workflows", func(w http.ResponseWriter, r *http.Request) {
				if tc.method == "POST" {
					assertRequestJSON(t, r, `{"monitor":{"type":"unsupported","aggregation":""},"alarm_actions":[{"type":"add_label","label":"Not Normalized","channel":""}]}`)
				}
				w.WriteHeader(tc.status)
				_, _ = fmt.Fprint(w, body)
			})
			server.HandleFunc("/v2/analytics/organizations/my-org/suites/my-suite/workflows/workflow-id", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = fmt.Fprint(w, body)
			})
			resp, err := tc.call(context.Background(), client.TestWorkflows)
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) {
				t.Fatalf("got %v, want ErrorResponse", err)
			}
			if resp == nil || resp.StatusCode != tc.status || apiErr.Message != "Validation Failed" {
				t.Fatalf("unexpected response/error: %#v, %#v", resp, apiErr)
			}
			assertJSONEqual(t, string(apiErr.RawBody), body)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err = tc.call(ctx, client.TestWorkflows)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("got %v, want context.Canceled", err)
			}
		})
	}
}
