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
	ID         string `json:"id"`
	URL        string `json:"url"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Repository string `json:"repository"`
	BuildsURL  string `json:"builds_url"`
	CreatedAt  string `json:"created_at"`

	// the provider of sources
	Provider *Provider `json:"provider"`

	// build featured when you view the project
	FeaturedBuild *Build `json:"featured_build"`

	// build steps
	Steps []*Step `json:"steps"`
}

// Provider represents a source code provider.
type Provider struct {
	ID         string `json:"id"`
	WebhookURL string `json:"webhook_url"`
}

// Build represents a build which has run in buildkite
type Build struct {
	ID          string            `json:"id"`
	URL         string            `json:"url"`
	Number      int               `json:"number"`
	State       string            `json:"state"`
	Message     string            `json:"message"`
	Commit      string            `json:"commit"`
	Branch      string            `json:"branch"`
	Env         map[string]string `json:"env"`
	CreatedAt   string            `json:"created_at"`
	ScheduledAt string            `json:"scheduled_at"`
	StartedAt   string            `json:"started_at"`
	FinishedAt  string            `json:"finished_at"`
	MetaData    interface{}       `json:"meta_data"`
}

// Step represents a build step in buildkites build pipeline
type Step struct {
	Type                string            `json:"type"`
	Name                string            `json:"name"`
	Command             string            `json:"command"`
	ArtifactPaths       string            `json:"artifact_paths"`
	BranchConfiguration string            `json:"branch_configuration"`
	Env                 map[string]string `json:"env"`
	TimeoutInMinutes    interface{}       `json:"timeout_in_minutes"` // *shrug*
	AgentQueryRules     interface{}       `json:"agent_query_rules"`  // *shrug*
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
