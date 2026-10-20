package buildkite

import (
	"context"
	"fmt"
)

// TestWorkflowsService handles communication with the suite-scoped workflow
// methods of the Buildkite Test Engine API. Reads require read_suites; writes
// require write_suites.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/test-engine/workflows
type TestWorkflowsService struct {
	client *Client
}

// TestWorkflow monitors a suite's tests and runs actions on alarm or recovery.
// Filters and actions are sets; their response order need not match input order.
type TestWorkflow struct {
	ID             string               `json:"id,omitempty"`
	URL            string               `json:"url,omitempty"`
	WebURL         string               `json:"web_url,omitempty"`
	Monitor        TestWorkflowMonitor  `json:"monitor"`
	Filters        []TestWorkflowFilter `json:"filters"`
	AlarmActions   []TestWorkflowAction `json:"alarm_actions"`
	RecoverActions []TestWorkflowAction `json:"recover_actions"`
	CreatedAt      *Timestamp           `json:"created_at,omitempty"`
	UpdatedAt      *Timestamp           `json:"updated_at,omitempty"`
	CreatedBy      *TestWorkflowCreator `json:"created_by"`
}

// TestWorkflowCreator is the user who created a workflow, or nil if unavailable.
type TestWorkflowCreator struct {
	ID        string `json:"id,omitempty"`
	GraphQLID string `json:"graphql_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
}

// TestWorkflowMonitor contains the configuration for the monitor's Type.
// Types are transition_count, probabilistic_flakiness, duration_threshold,
// passed_on_retry, and new_test. Only fields belonging to Type are returned.
// Numeric pointers distinguish absent settings from explicit zero thresholds.
type TestWorkflowMonitor struct {
	Type               string   `json:"type"`
	WindowSize         *int     `json:"window_size,omitempty"`
	AlarmThreshold     *float64 `json:"alarm_threshold,omitempty"`
	RecoverThreshold   *float64 `json:"recover_threshold,omitempty"`
	Aggregation        string   `json:"aggregation,omitempty"`
	EvaluationWindow   *int     `json:"evaluation_window,omitempty"`
	AlarmThresholdMS   *int     `json:"alarm_threshold_ms,omitempty"`
	RecoverThresholdMS *int     `json:"recover_threshold_ms,omitempty"`
}

// TestWorkflowMonitorInput configures a monitor on create or partially updates
// it. Use Some for every field to send, including zero thresholds. On creation,
// Type and every configuration field for that type are required; no defaults
// are added. On update, omitted fields remain unchanged and Type cannot change.
// Supplied fields are sent as-is, leaving validation to the API.
type TestWorkflowMonitorInput struct {
	Type               Optional[string]  `json:"type,omitzero"`
	WindowSize         Optional[int]     `json:"window_size,omitzero"`
	AlarmThreshold     Optional[float64] `json:"alarm_threshold,omitzero"`
	RecoverThreshold   Optional[float64] `json:"recover_threshold,omitzero"`
	Aggregation        Optional[string]  `json:"aggregation,omitzero"`
	EvaluationWindow   Optional[int]     `json:"evaluation_window,omitzero"`
	AlarmThresholdMS   Optional[int]     `json:"alarm_threshold_ms,omitzero"`
	RecoverThresholdMS Optional[int]     `json:"recover_threshold_ms,omitzero"`
}

// TestWorkflowFilter matches an execution tag. Operator is is, is_not, or
// starts_with. A transition_count monitor requires an scm.branch filter.
type TestWorkflowFilter struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// TestWorkflowAction is an alarm or recovery action. Label can be nil if its
// label was deleted. Endpoint may be masked for callers who cannot edit the
// suite. Pipeline is returned only for legacy, read-only trigger_pipeline
// actions. Unknown response types are tolerated.
type TestWorkflowAction struct {
	Type                string  `json:"type"`
	Label               *string `json:"label,omitempty"`
	State               string  `json:"state,omitempty"`
	Endpoint            string  `json:"endpoint,omitempty"`
	SlackWorkspaceID    string  `json:"slack_workspace_id,omitempty"`
	Channel             string  `json:"channel,omitempty"`
	MessageTemplate     string  `json:"message_template,omitempty"`
	LinearTeamID        string  `json:"linear_team_id,omitempty"`
	TitleTemplate       string  `json:"title_template,omitempty"`
	DescriptionTemplate string  `json:"description_template,omitempty"`
	Pipeline            string  `json:"pipeline,omitempty"`
}

// TestWorkflowActionInput supplies an action's Type and all fields for that type:
//   - add_label, remove_label: Label (an already-normalized label name)
//   - change_state: State (enabled, muted, or skipped)
//   - send_webhook_notification: Endpoint (HTTP or HTTPS URL)
//   - send_slack_notification: SlackWorkspaceID, Channel, MessageTemplate
//   - create_linear_issue: LinearTeamID, TitleTemplate, DescriptionTemplate
//
// SlackWorkspaceID identifies a Buildkite notification service. Templates use
// Liquid. Use Some for supplied fields; validation is left to the API, without
// normalizing labels or discarding fields that do not belong to Type.
// Pipeline and agent trigger actions are not supported for writes.
type TestWorkflowActionInput struct {
	Type                string           `json:"type"`
	Label               Optional[string] `json:"label,omitzero"`
	State               Optional[string] `json:"state,omitzero"`
	Endpoint            Optional[string] `json:"endpoint,omitzero"`
	SlackWorkspaceID    Optional[string] `json:"slack_workspace_id,omitzero"`
	Channel             Optional[string] `json:"channel,omitzero"`
	MessageTemplate     Optional[string] `json:"message_template,omitzero"`
	LinearTeamID        Optional[string] `json:"linear_team_id,omitzero"`
	TitleTemplate       Optional[string] `json:"title_template,omitzero"`
	DescriptionTemplate Optional[string] `json:"description_template,omitzero"`
}

// TestWorkflowCreate creates a workflow with an explicitly configured monitor.
// The API allows up to four filters and seven actions per list. new_test
// monitors do not accept recovery actions.
type TestWorkflowCreate struct {
	Monitor        TestWorkflowMonitorInput  `json:"monitor"`
	Filters        []TestWorkflowFilter      `json:"filters,omitempty"`
	AlarmActions   []TestWorkflowActionInput `json:"alarm_actions,omitempty"`
	RecoverActions []TestWorkflowActionInput `json:"recover_actions,omitempty"`
}

// TestWorkflowUpdate partially updates a workflow. Omitted fields, including
// individual monitor fields, remain unchanged. Supplied lists replace the whole
// list. Use Some([]TestWorkflowActionInput{}) or Some([]TestWorkflowFilter{}) to
// clear a list; a nil slice encodes as null, which the API rejects.
type TestWorkflowUpdate struct {
	Monitor        Optional[TestWorkflowMonitorInput]  `json:"monitor,omitzero"`
	Filters        Optional[[]TestWorkflowFilter]      `json:"filters,omitzero"`
	AlarmActions   Optional[[]TestWorkflowActionInput] `json:"alarm_actions,omitzero"`
	RecoverActions Optional[[]TestWorkflowActionInput] `json:"recover_actions,omitzero"`
}

type TestWorkflowListOptions struct{ ListOptions }

// List returns workflows newest first. Suite may be a suite slug or UUID.
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

	return workflows, resp, err
}

// Get returns a workflow by UUID. Suite may be a suite slug or UUID.
func (s *TestWorkflowsService) Get(ctx context.Context, org, suite, id string) (TestWorkflow, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/workflows/%s", org, suite, id)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TestWorkflow{}, nil, err
	}

	var workflow TestWorkflow
	resp, err := s.client.Do(req, &workflow)
	if err != nil {
		return TestWorkflow{}, resp, err
	}

	return workflow, resp, err
}

// Create creates a workflow. Suite may be a suite slug or UUID.
func (s *TestWorkflowsService) Create(ctx context.Context, org, suite string, input TestWorkflowCreate) (TestWorkflow, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/workflows", org, suite)
	req, err := s.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return TestWorkflow{}, nil, err
	}

	var workflow TestWorkflow
	resp, err := s.client.Do(req, &workflow)
	if err != nil {
		return TestWorkflow{}, resp, err
	}

	return workflow, resp, err
}

// Update updates a workflow by UUID. Suite may be a suite slug or UUID.
func (s *TestWorkflowsService) Update(ctx context.Context, org, suite, id string, input TestWorkflowUpdate) (TestWorkflow, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/workflows/%s", org, suite, id)
	req, err := s.client.NewRequest(ctx, "PATCH", u, input)
	if err != nil {
		return TestWorkflow{}, nil, err
	}

	var workflow TestWorkflow
	resp, err := s.client.Do(req, &workflow)
	if err != nil {
		return TestWorkflow{}, resp, err
	}

	return workflow, resp, err
}

// Delete deletes a workflow by UUID. Suite may be a suite slug or UUID.
func (s *TestWorkflowsService) Delete(ctx context.Context, org, suite, id string) (*Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/workflows/%s", org, suite, id)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
