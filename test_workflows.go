package buildkite

import (
	"context"
	"fmt"
)

// TestWorkflowsService reads suite-scoped Test Engine workflows.
// Requests require read_suites. Suite identifiers can be slugs or UUIDs.
type TestWorkflowsService struct {
	client *Client
}

// TestWorkflowMonitor holds the fields for a monitor's Type:
// transition_count, probabilistic_flakiness, duration_threshold,
// passed_on_retry, or new_test. Type values may be extended by the API.
// Pointers distinguish absent fields from zero thresholds.
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
// starts_with.
type TestWorkflowFilter struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// TestWorkflowAction holds fields for add_label, remove_label, change_state,
// send_webhook_notification, send_slack_notification, or create_linear_issue.
// Older workflows can also return trigger_pipeline with a Pipeline value.
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
	Pipeline            *string `json:"pipeline,omitempty"`
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
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
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
