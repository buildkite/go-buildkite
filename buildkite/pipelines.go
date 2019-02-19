// Copyright 2014 Mark Wolfe. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildkite

import (
	"errors"
	"fmt"
)

// PipelinesService handles communication with the pipeline related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api/pipelines
type PipelinesService struct {
	client *Client
}

// CreatePipeline - Create a Pipeline.
type CreatePipeline struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Steps      []Step `json:"steps"`

	// Optional fields
	Description                     string            `json:"description,omitempty"`
	Env                             map[string]string `json:"env,omitempty"`
	ProviderSettings                ProviderSettings  `json:"provider_settings,omitempty"`
	BranchConfiguration             string            `json:"branch_configuration,omitempty"`
	SkipQueuedBranchBuilds          bool              `json:"skip_queued_branch_builds,omitempty"`
	SkipQueuedBranchBuildsFilter    string            `json:"skip_queued_branch_builds_filter,omitempty"`
	CancelRunningBranchBuilds       bool              `json:"cancel_running_branch_builds,omitempty"`
	CancelRunningBranchBuildsFilter string            `json:"cancel_running_branch_builds_filter,omitempty"`
	TeamUuids                       []string          `json:"team_uuids,omitempty"`
}

// Pipeline represents a buildkite pipeline.
type Pipeline struct {
	ID         *string    `json:"id,omitempty"`
	URL        *string    `json:"url,omitempty"`
	WebURL     *string    `json:"web_url,omitempty"`
	Name       *string    `json:"name,omitempty"`
	Slug       *string    `json:"slug,omitempty"`
	Repository *string    `json:"repository,omitempty"`
	BuildsURL  *string    `json:"builds_url,omitempty"`
	BadgeURL   *string    `json:"badge_url,omitempty"`
	CreatedAt  *Timestamp `json:"created_at,omitempty"`

	ScheduledBuildsCount *int `json:"scheduled_builds_count,omitempty"`
	RunningBuildsCount   *int `json:"running_builds_count,omitempty"`
	ScheduledJobsCount   *int `json:"scheduled_jobs_count,omitempty"`
	RunningJobsCount     *int `json:"running_jobs_count,omitempty"`
	WaitingJobsCount     *int `json:"waiting_jobs_count,omitempty"`

	// the provider of sources
	Provider *Provider `json:"provider,omitempty"`

	// build steps
	Steps []*Step `json:"steps,omitempty"`
}

// Step represents a build step in buildkites build pipeline
type Step struct {
	Type                *string           `json:"type,omitempty"`
	Name                *string           `json:"name,omitempty"`
	Command             *string           `json:"command,omitempty"`
	ArtifactPaths       *string           `json:"artifact_paths,omitempty"`
	BranchConfiguration *string           `json:"branch_configuration,omitempty"`
	Env                 map[string]string `json:"env,omitempty"`
	TimeoutInMinutes    *int              `json:"timeout_in_minutes,omitempty"`
	AgentQueryRules     []string          `json:"agent_query_rules,omitempty"`
}

// PipelineListOptions specifies the optional parameters to the
// PipelinesService.List method.
type PipelineListOptions struct {
	ListOptions
}

// Create - Creates a pipeline for a given organisation.
//
// buildkite API docs: https://buildkite.com/docs/rest-api/pipelines#create-a-pipeline
func (ps *PipelinesService) Create(org string, p *CreatePipeline) (*Pipeline, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines", org)

	req, err := ps.client.NewRequest("POST", u, p)
	if err != nil {
		return nil, nil, err
	}

	pipeline := new(Pipeline)
	resp, err := ps.client.Do(req, pipeline)
	if err != nil {
		return nil, resp, err
	}

	return pipeline, resp, err
}

// Get fetches a pipeline.
//
// buildkite API docs: https://buildkite.com/docs/rest-api/pipelines#get-a-pipeline
func (ps *PipelinesService) Get(org string, slug string) (*Pipeline, *Response, error) {

	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s", org, slug)

	req, err := ps.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	pipeline := new(Pipeline)
	resp, err := ps.client.Do(req, pipeline)
	if err != nil {
		return nil, resp, err
	}

	return pipeline, resp, err
}

// List the pipelines for a given organisation.
//
// buildkite API docs: https://buildkite.com/docs/api/pipelines#list-pipelines
func (ps *PipelinesService) List(org string, opt *PipelineListOptions) ([]Pipeline, *Response, error) {
	var u string

	u = fmt.Sprintf("v2/organizations/%s/pipelines", org)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := ps.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	pipelines := new([]Pipeline)
	resp, err := ps.client.Do(req, pipelines)
	if err != nil {
		return nil, resp, err
	}

	return *pipelines, resp, err
}

// Delete a pipeline.
//
// buildkite API docs: https://buildkite.com/docs/rest-api/pipelines#delete-a-pipeline
func (ps *PipelinesService) Delete(org string, slug string) (*Response, error) {

	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s", org, slug)

	req, err := ps.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return ps.client.Do(req, nil)
}

// Update - Updates a pipeline.
//
// buildkite API docs: https://buildkite.com/docs/rest-api/pipelines#update-a-pipeline
func (ps *PipelinesService) Update(org string, p *Pipeline) (*Response, error) {
	if p == nil {
		return nil, errors.New("pipeline must not be nil")
	}

	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s", org, *p.Slug)

	// There is quite a lot of properties that are not represented by the Client-side
	// Pipeline abstraction hence only a subset can be updated.
	cp := &CreatePipeline{
		Name:       *p.Name,
		Repository: *p.Repository,
		Steps:      make([]Step, len(p.Steps)),
	}

	for i := range p.Steps {
		cp.Steps[i] = *p.Steps[i]
	}

	req, err := ps.client.NewRequest("PATCH", u, cp)
	if err != nil {
		return nil, err
	}

	resp, err := ps.client.Do(req, p)
	if err != nil {
		return resp, err
	}

	return resp, err
}
