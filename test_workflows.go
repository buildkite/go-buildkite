package buildkite

import (
	"context"
	"fmt"
)

// TestWorkflowsService manages suite-scoped Test Engine workflows.
// Reads require read_suites; writes require write_suites and permission to edit
// the suite. Suite identifiers can be slugs or UUIDs.
type TestWorkflowsService struct {
	client *Client
}

// TestWorkflowMonitor holds the fields for a monitor's Type:
// transition_count, probabilistic_flakiness, duration_threshold,
// passed_on_retry, or new_test. Type values may be extended by the API.
// Create requires every field for the selected type. Update sends only supplied
// fields and cannot change the type. Pointers distinguish omission from zero.
type TestWorkflowMonitor struct {
	Type               string   `json:"type,omitempty"`
	WindowSize         *int     `json:"window_size,omitempty"`
	AlarmThreshold     *float64 `json:"alarm_threshold,omitempty"`
	RecoverThreshold   *float64 `json:"recover_threshold,omitempty"`
	Aggregation        *string  `json:"aggregation,omitempty"`
	EvaluationWindow   *int     `json:"evaluation_window,omitempty"`
	AlarmThresholdMS   *int     `json:"alarm_threshold_ms,omitempty"`
	RecoverThresholdMS *int     `json:"recover_threshold_ms,omitempty"`
}

// TestWorkflowFilter matches an execution tag. Operator is is, is_not, or
// starts_with. Mutable mut.* tags are not accepted by the API.
type TestWorkflowFilter struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// TestWorkflowAction holds fields for add_label, remove_label, change_state,
// send_webhook_notification, send_slack_notification, or create_linear_issue.
// Supply only the fields for the selected Type; the API validates them.
// Type values may be extended. Unrecognized response fields are ignored.
type TestWorkflowAction struct {
	Type                string  `json:"type"`
	Label               *string `json:"label,omitempty"`
	State               *string `json:"state,omitempty"`
	Endpoint            *string `json:"endpoint,omitempty"`
	SlackWorkspaceID    *string `json:"slack_workspace_id,omitempty"`
	Channel             *string `json:"channel,omitempty"`
	MessageTemplate     *string `json:"message_template,omitempty"`
	LinearTeamID        *string `json:"linear_team_id,omitempty"`
	TitleTemplate       *string `json:"title_template,omitempty"`
	DescriptionTemplate *string `json:"description_template,omitempty"`
}

type TestWorkflowCreator struct {
	ID        string `json:"id"`
	GraphQLID string `json:"graphql_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

type TestWorkflow struct {
	ID             string               `json:"id"`
	URL            string               `json:"url"`
	WebURL         string               `json:"web_url"`
	Monitor        TestWorkflowMonitor  `json:"monitor"`
	Filters        []TestWorkflowFilter `json:"filters"`
	AlarmActions   []TestWorkflowAction `json:"alarm_actions"`
	RecoverActions []TestWorkflowAction `json:"recover_actions"`
	CreatedAt      *Timestamp           `json:"created_at"`
	UpdatedAt      *Timestamp           `json:"updated_at"`
	CreatedBy      *TestWorkflowCreator `json:"created_by"`
}

// TestWorkflowCreate requires a monitor with all fields for its type.
// Omitted lists default to empty on the server. The client adds no defaults.
type TestWorkflowCreate struct {
	Monitor        TestWorkflowMonitor  `json:"monitor"`
	Filters        []TestWorkflowFilter `json:"filters,omitempty"`
	AlarmActions   []TestWorkflowAction `json:"alarm_actions,omitempty"`
	RecoverActions []TestWorkflowAction `json:"recover_actions,omitempty"`
}

// TestWorkflowUpdate leaves unset fields unchanged. Supplied lists replace
// existing lists. Use Some([]TestWorkflowAction{}) or
// Some([]TestWorkflowFilter{}) to clear a list; a nil slice sends null, not [].
// Monitor updates merge supplied fields with the existing monitor.
type TestWorkflowUpdate struct {
	Monitor        Optional[TestWorkflowMonitor]  `json:"monitor,omitzero"`
	Filters        Optional[[]TestWorkflowFilter] `json:"filters,omitzero"`
	AlarmActions   Optional[[]TestWorkflowAction] `json:"alarm_actions,omitzero"`
	RecoverActions Optional[[]TestWorkflowAction] `json:"recover_actions,omitzero"`
}

type TestWorkflowListOptions struct{ ListOptions }

// List returns workflows newest first, with pagination in the response's Link header.
func (s *TestWorkflowsService) List(ctx context.Context, org, suite string, opt *TestWorkflowListOptions) ([]TestWorkflow, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/workflows", org, suite)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var workflows []TestWorkflow
	resp, err := s.client.Do(req, &workflows)
	if err != nil {
		return nil, resp, err
	}

	return workflows, resp, nil
}

func (s *TestWorkflowsService) Get(ctx context.Context, org, suite, workflowID string) (TestWorkflow, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/workflows/%s", org, suite, workflowID)
	return s.do(ctx, "GET", u, nil)
}

func (s *TestWorkflowsService) Create(ctx context.Context, org, suite string, input TestWorkflowCreate) (TestWorkflow, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/workflows", org, suite)
	return s.do(ctx, "POST", u, input)
}

func (s *TestWorkflowsService) Update(ctx context.Context, org, suite, workflowID string, input TestWorkflowUpdate) (TestWorkflow, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/workflows/%s", org, suite, workflowID)
	return s.do(ctx, "PATCH", u, input)
}

func (s *TestWorkflowsService) Delete(ctx context.Context, org, suite, workflowID string) (*Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/workflows/%s", org, suite, workflowID)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

func (s *TestWorkflowsService) do(ctx context.Context, method, u string, body any) (TestWorkflow, *Response, error) {
	req, err := s.client.NewRequest(ctx, method, u, body)
	if err != nil {
		return TestWorkflow{}, nil, err
	}

	var workflow TestWorkflow
	resp, err := s.client.Do(req, &workflow)
	if err != nil {
		return TestWorkflow{}, resp, err
	}

	return workflow, resp, nil
}
