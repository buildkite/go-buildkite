// Copyright 2014 Mark Wolfe. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildkite

import "fmt"

// ProjectsService handles communication with the project related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api/projects
type ProjectsService struct {
	client *Client
}

// Project represents a buildkite project.
type Project struct {
	ID         *string    `json:"id,omitempty"`
	URL        *string    `json:"url,omitempty"`
	WebURL     *string    `json:"web_url,omitempty"`
	Name       *string    `json:"name,omitempty"`
	Slug       *string    `json:"slug,omitempty"`
	Repository *string    `json:"repository,omitempty"`
	BuildsURL  *string    `json:"builds_url,omitempty"`
	CreatedAt  *Timestamp `json:"created_at,omitempty"`

	ScheduledBuildsCount *int `json:"scheduled_builds_count,omitempty"`
	RunningBuildsCount   *int `json:"running_builds_count,omitempty"`
	ScheduledJobsCount   *int `json:"scheduled_jobs_count,omitempty"`
	RunningJobsCount     *int `json:"running_jobs_count,omitempty"`
	WaitingJobsCount     *int `json:"waiting_jobs_count,omitempty"`

	// the provider of sources
	Provider *Provider `json:"provider,omitempty"`

	// build featured when you view the project
	FeaturedBuild *Build `json:"featured_build,omitempty"`

	// build steps
	Steps []*Step `json:"steps,omitempty"`
}

// Provider represents a source code provider.
type Provider struct {
	ID         *string `json:"id,omitempty"`
	WebhookURL *string `json:"webhook_url,omitempty"`
}

// Step represents a build step in buildkites build pipeline
type Step struct {
	Type                *string           `json:"type,omitempty"`
	Name                *string           `json:"name,omitempty"`
	Command             *string           `json:"command,omitempty"`
	ArtifactPaths       *string           `json:"artifact_paths,omitempty"`
	BranchConfiguration *string           `json:"branch_configuration,omitempty"`
	Env                 map[string]string `json:"env,omitempty"`
	TimeoutInMinutes    interface{}       `json:"timeout_in_minutes,omitempty"` // *shrug*
	AgentQueryRules     interface{}       `json:"agent_query_rules,omitempty"`  // *shrug*
}

// ProjectListOptions specifies the optional parameters to the
// ProjectsService.List method.
type ProjectListOptions struct {
	ListOptions
}

// List the projects for a given orginisation.
//
// buildkite API docs: https://buildkite.com/docs/api/projects#list-projects
func (ps *ProjectsService) List(org string, opt *ProjectListOptions) ([]Project, *Response, error) {
	var u string

	u = fmt.Sprintf("v1/organizations/%s/projects", org)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := ps.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projects := new([]Project)
	resp, err := ps.client.Do(req, projects)
	if err != nil {
		return nil, resp, err
	}

	return *projects, resp, err
}
