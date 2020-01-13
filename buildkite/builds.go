package buildkite

import (
	"fmt"
	"time"
)

// BuildsService handles communication with the build related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api/builds
type BuildsService struct {
	client *Client
}

// Author of a commit (used in CreateBuild)
type Author struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// CreateBuild - Create a build.
type CreateBuild struct {
	Commit  string `json:"commit" yaml:"commit"`
	Branch  string `json:"branch" yaml:"branch"`
	Message string `json:"message" yaml:"message"`

	// Optional fields
	Author                      Author            `json:"author,omitempty" yaml:"author,omitempty"`
	Env                         map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	MetaData                    map[string]string `json:"meta_data,omitempty" yaml:"meta_data,omitempty"`
	IgnorePipelineBranchFilters bool              `json:"ignore_pipeline_branch_filters,omitempty" yaml:"ignore_pipeline_branch_filters,omitempty"`
	PullRequestBaseBranch       string            `json:"pull_request_base_branch,omitempty" yaml:"pull_request_base_branch,omitempty"`
	PullRequestID               int64             `json:"pull_request_id,omitempty" yaml:"pull_request_id,omitempty"`
	PullRequestRepository       string            `json:"pull_request_repository,omitempty" yaml:"pull_request_repository,omitempty"`
}

// Creator represents who created a build
type Creator struct {
	AvatarURL string     `json:"avatar_url" yaml:"avatar_url"`
	CreatedAt *Timestamp `json:"created_at" yaml:"created_at"`
	Email     string     `json:"email" yaml:"email"`
	ID        string     `json:"id" yaml:"id"`
	Name      string     `json:"name" yaml:"name"`
}

// PullRequest represents a Github PR
type PullRequest struct {
	ID         *string `json:"id,omitempty" yaml:"id,omitempty"`
	Base       *string `json:"base,omitempty" yaml:"base,omitempty"`
	Repository *string `json:"repository,omitempty" yaml:"repository,omitempty"`
}

// Build represents a build which has run in buildkite
type Build struct {
	ID          *string                `json:"id,omitempty" yaml:"id,omitempty"`
	URL         *string                `json:"url,omitempty" yaml:"url,omitempty"`
	WebURL      *string                `json:"web_url,omitempty" yaml:"web_url,omitempty"`
	Number      *int                   `json:"number,omitempty" yaml:"number,omitempty"`
	State       *string                `json:"state,omitempty" yaml:"state,omitempty"`
	Blocked     *bool                  `json:"blocked,omitempty" yaml:"blocked,omitempty"`
	Message     *string                `json:"message,omitempty" yaml:"message,omitempty"`
	Commit      *string                `json:"commit,omitempty" yaml:"commit,omitempty"`
	Branch      *string                `json:"branch,omitempty" yaml:"branch,omitempty"`
	Env         map[string]interface{} `json:"env,omitempty" yaml:"env,omitempty"`
	CreatedAt   *Timestamp             `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ScheduledAt *Timestamp             `json:"scheduled_at,omitempty" yaml:"scheduled_at,omitempty"`
	StartedAt   *Timestamp             `json:"started_at,omitempty" yaml:"started_at,omitempty"`
	FinishedAt  *Timestamp             `json:"finished_at,omitempty" yaml:"finished_at,omitempty"`
	MetaData    interface{}            `json:"meta_data,omitempty" yaml:"meta_data,omitempty"`
	Creator     *Creator               `json:"creator,omitempty" yaml:"creator,omitempty"`
	Source      *string                `json:"source,omitempty" yaml:"source,omitempty"`

	// jobs run during the build
	Jobs []*Job `json:"jobs,omitempty" yaml:"jobs,omitempty"`

	// the pipeline this build is associated with
	Pipeline *Pipeline `json:"pipeline,omitempty" yaml:"pipeline,omitempty"`

	// the pull request this build is associated with
	PullRequest *PullRequest `json:"pull_request,omitempty" yaml:"pull_request,omitempty"`
}

// BuildsListOptions specifies the optional parameters to the
// BuildsService.List method.
type BuildsListOptions struct {

	// Filters the results by the user who created the build
	Creator string `url:"creator,omitempty"`

	// Filters the results by builds created on or after the given time
	CreatedFrom time.Time `url:"created_from,omitempty"`

	// Filters the results by builds created before the given time
	CreatedTo time.Time `url:"created_to,omitempty"`

	// Filters the results by builds finished on or after the given time
	FinishedFrom time.Time `url:"finished_from,omitempty"`

	// State of builds to list.  Possible values are: running, scheduled, passed,
	// failed, canceled, skipped and not_run. Default is "".
	State []string `url:"state,brackets,omitempty"`

	// Branch filter by the name of the branch. Default is "".
	Branch string `url:"branch,omitempty"`

	// Filters the results by builds for the specific commit SHA (full, not shortened). Default is "".
	Commit string `url:"commit,omitempty"`

	// Include all retried jobs in each build’s jobs list
	IncludeRetriedJobs bool `url:"include_retried_jobs,omitempty"`

	ListOptions
}

// Cancel - Trigger a cancel for the target build
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/builds#cancel-a-build
func (bs *BuildsService) Cancel(org, pipeline, build string) (*Build, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/cancel", org, pipeline, build)
	req, err := bs.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}
	result := Build{}
	_, err = bs.client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create - Create a pipeline
//
// buildkite API docs: https://buildkite.com/docs/api/builds#create-a-build
func (bs *BuildsService) Create(org string, pipeline string, b *CreateBuild) (*Build, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds", org, pipeline)

	req, err := bs.client.NewRequest("POST", u, b)
	if err != nil {
		return nil, nil, err
	}

	build := new(Build)
	resp, err := bs.client.Do(req, build)
	if err != nil {
		return nil, resp, err
	}

	return build, resp, err
}

// Get fetches a build.
//
// buildkite API docs: https://buildkite.com/docs/api/builds#get-a-build
func (bs *BuildsService) Get(org string, pipeline string, id string, opt *BuildsListOptions) (*Build, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s", org, pipeline, id)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	build := new(Build)
	resp, err := bs.client.Do(req, build)
	if err != nil {
		return nil, resp, err
	}

	return build, resp, err
}

// List the builds for the current user.
//
// buildkite API docs: https://buildkite.com/docs/api/builds#list-all-builds
func (bs *BuildsService) List(opt *BuildsListOptions) ([]Build, *Response, error) {
	var u string

	u = fmt.Sprintf("v2/builds")

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	orgs := new([]Build)
	resp, err := bs.client.Do(req, orgs)
	if err != nil {
		return nil, resp, err
	}

	return *orgs, resp, err
}

// ListByOrg lists the builds within the specified orginisation.
//
// buildkite API docs: https://buildkite.com/docs/api/builds#list-builds-for-an-organization
func (bs *BuildsService) ListByOrg(org string, opt *BuildsListOptions) ([]Build, *Response, error) {
	var u string

	u = fmt.Sprintf("v2/organizations/%s/builds", org)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	orgs := new([]Build)
	resp, err := bs.client.Do(req, orgs)
	if err != nil {
		return nil, resp, err
	}

	return *orgs, resp, err
}

// ListByPipeline lists the builds for a pipeline within the specified originisation.
//
// buildkite API docs: https://buildkite.com/docs/api/builds#list-builds-for-a-pipeline
func (bs *BuildsService) ListByPipeline(org string, pipeline string, opt *BuildsListOptions) ([]Build, *Response, error) {
	var u string

	u = fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds", org, pipeline)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	orgs := new([]Build)
	resp, err := bs.client.Do(req, orgs)
	if err != nil {
		return nil, resp, err
	}

	return *orgs, resp, err
}

// Rebuild triggers a rebuild for the target build
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/builds#rebuild-a-build
func (bs *BuildsService) Rebuild(org, pipeline, build string) (*Build, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/rebuild", org, pipeline, build)
	req, err := bs.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}
	result := Build{}
	_, err = bs.client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
