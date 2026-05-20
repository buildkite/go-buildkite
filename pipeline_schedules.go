package buildkite

import (
	"context"
	"fmt"
)

// PipelineSchedulesService handles communication with the pipeline schedule
// related methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/pipeline-schedules
type PipelineSchedulesService struct {
	client *Client
}

// PipelineSchedule represents a buildkite pipeline schedule.
type PipelineSchedule struct {
	ID            string            `json:"id,omitempty"`
	GraphQLID     string            `json:"graphql_id,omitempty"`
	URL           string            `json:"url,omitempty"`
	Label         string            `json:"label,omitempty"`
	Cronline      string            `json:"cronline,omitempty"`
	Message       string            `json:"message,omitempty"`
	Commit        string            `json:"commit,omitempty"`
	Branch        string            `json:"branch,omitempty"`
	Env           map[string]string `json:"env,omitempty"`
	Enabled       bool              `json:"enabled"`
	NextBuildAt   *Timestamp        `json:"next_build_at,omitempty"`
	FailedMessage string            `json:"failed_message,omitempty"`
	FailedAt      *Timestamp        `json:"failed_at,omitempty"`
	CreatedAt     *Timestamp        `json:"created_at,omitempty"`
	CreatedBy     *Creator          `json:"created_by,omitempty"`
	Pipeline      *Pipeline         `json:"pipeline,omitempty"`
}

// CreatePipelineSchedule represents the request body for creating a pipeline schedule.
type CreatePipelineSchedule struct {
	Cronline string            `json:"cronline"`
	Label    string            `json:"label,omitempty"`
	Message  string            `json:"message,omitempty"`
	Commit   string            `json:"commit,omitempty"`
	Branch   string            `json:"branch,omitempty"`
	Env      map[string]string `json:"env,omitempty"`
	Enabled  *bool             `json:"enabled,omitempty"`
}

// UpdatePipelineSchedule represents the request body for updating a pipeline schedule.
type UpdatePipelineSchedule struct {
	Cronline Optional[string]            `json:"cronline,omitzero"`
	Label    Optional[string]            `json:"label,omitzero"`
	Message  Optional[string]            `json:"message,omitzero"`
	Commit   Optional[string]            `json:"commit,omitzero"`
	Branch   Optional[string]            `json:"branch,omitzero"`
	Env      Optional[map[string]string] `json:"env,omitzero"`
	Enabled  Optional[bool]              `json:"enabled,omitzero"`
}

// PipelineScheduleListOptions specifies the optional parameters to the
// PipelineSchedulesService.List method.
type PipelineScheduleListOptions struct {
	ListOptions
}

// List the pipeline schedules for a pipeline.
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/pipeline-schedules#list-pipeline-schedules
func (pss *PipelineSchedulesService) List(ctx context.Context, org, pipelineSlug string, opt *PipelineScheduleListOptions) ([]PipelineSchedule, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/schedules", org, pipelineSlug)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := pss.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var schedules []PipelineSchedule
	resp, err := pss.client.Do(req, &schedules)
	if err != nil {
		return nil, resp, err
	}

	return schedules, resp, err
}

// Get a pipeline schedule by ID.
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/pipeline-schedules#get-a-pipeline-schedule
func (pss *PipelineSchedulesService) Get(ctx context.Context, org, pipelineSlug, id string) (PipelineSchedule, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/schedules/%s", org, pipelineSlug, id)

	req, err := pss.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return PipelineSchedule{}, nil, err
	}

	var schedule PipelineSchedule
	resp, err := pss.client.Do(req, &schedule)
	if err != nil {
		return PipelineSchedule{}, resp, err
	}

	return schedule, resp, err
}

// Create a new pipeline schedule.
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/pipeline-schedules#create-a-pipeline-schedule
func (pss *PipelineSchedulesService) Create(ctx context.Context, org, pipelineSlug string, in CreatePipelineSchedule) (PipelineSchedule, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/schedules", org, pipelineSlug)

	req, err := pss.client.NewRequest(ctx, "POST", u, in)
	if err != nil {
		return PipelineSchedule{}, nil, err
	}

	var schedule PipelineSchedule
	resp, err := pss.client.Do(req, &schedule)
	if err != nil {
		return PipelineSchedule{}, resp, err
	}

	return schedule, resp, err
}

// Update an existing pipeline schedule.
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/pipeline-schedules#update-a-pipeline-schedule
func (pss *PipelineSchedulesService) Update(ctx context.Context, org, pipelineSlug, id string, in UpdatePipelineSchedule) (PipelineSchedule, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/schedules/%s", org, pipelineSlug, id)

	req, err := pss.client.NewRequest(ctx, "PATCH", u, in)
	if err != nil {
		return PipelineSchedule{}, nil, err
	}

	var schedule PipelineSchedule
	resp, err := pss.client.Do(req, &schedule)
	if err != nil {
		return PipelineSchedule{}, resp, err
	}

	return schedule, resp, err
}

// Delete a pipeline schedule.
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/pipeline-schedules#delete-a-pipeline-schedule
func (pss *PipelineSchedulesService) Delete(ctx context.Context, org, pipelineSlug, id string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/schedules/%s", org, pipelineSlug, id)

	req, err := pss.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return pss.client.Do(req, nil)
}
