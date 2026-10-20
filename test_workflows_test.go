package buildkite

import (
	"context"
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

	server.HandleFunc("/v2/analytics/organizations/acme/suites/web/workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2", "per_page": "10"})
		w.Header().Set("Link", `<https://api.buildkite.com/v2/analytics/organizations/acme/suites/web/workflows?page=3&per_page=10>; rel="next", <https://api.buildkite.com/v2/analytics/organizations/acme/suites/web/workflows?page=1&per_page=10>; rel="prev"`)
		_, _ = fmt.Fprint(w, `[{
			"id":"workflow-1","url":"https://api.buildkite.com/workflow-1","web_url":"https://buildkite.com/workflow-1",
			"monitor":{"type":"transition_count","window_size":100,"alarm_threshold":0.075,"recover_threshold":0},
			"filters":[{"key":"scm.branch","operator":"is","value":"main"}],
			"alarm_actions":[{"type":"add_label","label":"flaky"}],
			"recover_actions":[{"type":"remove_label","label":null}],
			"created_at":"2026-10-12T01:02:03.123Z","updated_at":"2026-10-13T04:05:06Z",
			"created_by":{"id":"user-1","graphql_id":"graphql-user-1","name":"Sam","email":"sam@example.com"}
		},{
			"id":"workflow-2","monitor":{"type":"future_monitor","future_setting":42},
			"filters":[],"alarm_actions":[{"type":"future_action","future_setting":true}],
			"recover_actions":[],"created_by":null
		}]`)
	})

	got, resp, err := client.TestWorkflows.List(context.Background(), "acme", "web", &TestWorkflowListOptions{
		ListOptions: ListOptions{Page: 2, PerPage: 10},
	})
	if err != nil {
		t.Fatal(err)
	}
	window, alarm, recover, label := 100, 0.075, 0.0, "flaky"
	want := []TestWorkflow{{
		ID: "workflow-1", URL: "https://api.buildkite.com/workflow-1", WebURL: "https://buildkite.com/workflow-1",
		Monitor:        TestWorkflowMonitor{Type: "transition_count", WindowSize: &window, AlarmThreshold: &alarm, RecoverThreshold: &recover},
		Filters:        []TestWorkflowFilter{{Key: "scm.branch", Operator: "is", Value: "main"}},
		AlarmActions:   []TestWorkflowAction{{Type: "add_label", Label: &label}},
		RecoverActions: []TestWorkflowAction{{Type: "remove_label"}},
		CreatedAt:      NewTimestamp(time.Date(2026, 10, 12, 1, 2, 3, 123000000, time.UTC)),
		UpdatedAt:      NewTimestamp(time.Date(2026, 10, 13, 4, 5, 6, 0, time.UTC)),
		CreatedBy:      &TestWorkflowCreator{ID: "user-1", GraphQLID: "graphql-user-1", Name: "Sam", Email: "sam@example.com"},
	}, {
		ID: "workflow-2", Monitor: TestWorkflowMonitor{Type: "future_monitor"},
		Filters: []TestWorkflowFilter{}, AlarmActions: []TestWorkflowAction{{Type: "future_action"}}, RecoverActions: []TestWorkflowAction{},
	}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("workflows (-want +got):\n%s", diff)
	}
	if resp.NextPage != 3 || resp.PrevPage != 1 {
		t.Errorf("pagination: next=%d prev=%d", resp.NextPage, resp.PrevPage)
	}
}

func TestTestWorkflowsService_ListEmpty(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/analytics/organizations/acme/suites/web/workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		_, _ = fmt.Fprint(w, `[]`)
	})
	got, resp, err := client.TestWorkflows.List(context.Background(), "acme", "web", nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]TestWorkflow{}, got); diff != "" {
		t.Errorf("workflows (-want +got):\n%s", diff)
	}
	if resp.NextPage != 0 {
		t.Errorf("next page = %d", resp.NextPage)
	}
}

func TestTestWorkflowsService_Get(t *testing.T) {
	t.Parallel()
	window, alarm, recover := 100, 0.075, 0.0
	aggregation, evaluationWindow, alarmMS, recoverMS := "p50", 12, 1501, 0
	label, state, endpoint := "flaky", "muted", "https://example.com/[REDACTED]"
	workspace, channel, message := "workspace-1", "#tests", "{{ test.name }} alarmed"
	team, title, description, pipeline := "team-1", "Flaky test", "{{ test.name }}", "acme/tests"
	actions := []TestWorkflowAction{
		{Type: "add_label", Label: &label},
		{Type: "change_state", State: &state},
		{Type: "send_webhook_notification", Endpoint: &endpoint},
		{Type: "send_slack_notification", SlackWorkspaceID: &workspace, Channel: &channel, MessageTemplate: &message},
		{Type: "create_linear_issue", LinearTeamID: &team, TitleTemplate: &title, DescriptionTemplate: &description},
		{Type: "trigger_pipeline", Pipeline: &pipeline},
	}
	tests := []struct {
		name    string
		monitor TestWorkflowMonitor
		json    string
	}{
		{"transition_count", TestWorkflowMonitor{Type: "transition_count", WindowSize: &window, AlarmThreshold: &alarm, RecoverThreshold: &recover}, `{"type":"transition_count","window_size":100,"alarm_threshold":0.075,"recover_threshold":0}`},
		{"probabilistic_flakiness", TestWorkflowMonitor{Type: "probabilistic_flakiness", AlarmThreshold: &alarm, RecoverThreshold: &recover}, `{"type":"probabilistic_flakiness","alarm_threshold":0.075,"recover_threshold":0}`},
		{"duration_threshold", TestWorkflowMonitor{Type: "duration_threshold", Aggregation: &aggregation, EvaluationWindow: &evaluationWindow, AlarmThresholdMS: &alarmMS, RecoverThresholdMS: &recoverMS}, `{"type":"duration_threshold","aggregation":"p50","evaluation_window":12,"alarm_threshold_ms":1501,"recover_threshold_ms":0}`},
		{"passed_on_retry", TestWorkflowMonitor{Type: "passed_on_retry"}, `{"type":"passed_on_retry"}`},
		{"new_test", TestWorkflowMonitor{Type: "new_test"}, `{"type":"new_test"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			suiteID := "7c202aaa-3165-4811-9813-173c4c285463"
			server.HandleFunc("/v2/analytics/organizations/acme/suites/"+suiteID+"/workflows/workflow-1", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testFormValues(t, r, values{})
				_, _ = fmt.Fprint(w, `{"id":"workflow-1","monitor":`+tt.json+`,"filters":[{"key":"scm.branch","operator":"is_not","value":"main"},{"key":"test.scope","operator":"starts_with","value":"api"}],"alarm_actions":[
					{"type":"add_label","label":"flaky"},
					{"type":"change_state","state":"muted"},
					{"type":"send_webhook_notification","endpoint":"https://example.com/[REDACTED]"},
					{"type":"send_slack_notification","slack_workspace_id":"workspace-1","channel":"#tests","message_template":"{{ test.name }} alarmed"},
					{"type":"create_linear_issue","linear_team_id":"team-1","title_template":"Flaky test","description_template":"{{ test.name }}"},
					{"type":"trigger_pipeline","pipeline":"acme/tests"}
				],"recover_actions":[]}`)
			})
			got, resp, err := client.TestWorkflows.Get(context.Background(), "acme", suiteID, "workflow-1")
			if err != nil {
				t.Fatal(err)
			}
			want := TestWorkflow{
				ID: "workflow-1", Monitor: tt.monitor,
				Filters:      []TestWorkflowFilter{{Key: "scm.branch", Operator: "is_not", Value: "main"}, {Key: "test.scope", Operator: "starts_with", Value: "api"}},
				AlarmActions: actions, RecoverActions: []TestWorkflowAction{},
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("workflow (-want +got):\n%s", diff)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("get status = %d", resp.StatusCode)
			}
		})
	}
}

func TestTestWorkflowsService_Errors(t *testing.T) {
	t.Parallel()
	calls := []struct {
		name string
		call func(context.Context, *Client) (*Response, error)
	}{
		{"list", func(ctx context.Context, c *Client) (*Response, error) {
			got, r, e := c.TestWorkflows.List(ctx, "acme", "web", nil)
			if e != nil && got != nil {
				t.Errorf("list returned workflows on error: %+v", got)
			}
			return r, e
		}},
		{"get", func(ctx context.Context, c *Client) (*Response, error) {
			_, r, e := c.TestWorkflows.Get(ctx, "acme", "web", "workflow-1")
			return r, e
		}},
	}
	for _, call := range calls {
		t.Run(call.name, func(t *testing.T) {
			t.Parallel()
			for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound} {
				t.Run(fmt.Sprint(status), func(t *testing.T) {
					server, client, teardown := newMockServerAndClient(t)
					t.Cleanup(teardown)
					message := http.StatusText(status)
					body := fmt.Sprintf(`{"message":%q}`, message)
					server.HandleFunc("/v2/analytics/organizations/acme/", func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(status)
						_, _ = fmt.Fprint(w, body)
					})
					resp, err := call.call(context.Background(), client)
					var apiErr *ErrorResponse
					if !errors.As(err, &apiErr) {
						t.Fatalf("error = %v; want *ErrorResponse", err)
					}
					if resp == nil || resp.StatusCode != status || apiErr.Message != message || string(apiErr.RawBody) != body {
						t.Errorf("response=%+v error=%+v", resp, apiErr)
					}
				})
			}
			t.Run("canceled", func(t *testing.T) {
				_, client, teardown := newMockServerAndClient(t)
				t.Cleanup(teardown)
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_, err := call.call(ctx, client)
				if !errors.Is(err, context.Canceled) {
					t.Errorf("error = %v; want context.Canceled", err)
				}
			})
		})
	}
}
