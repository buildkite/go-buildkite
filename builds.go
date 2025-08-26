package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
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
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (a *Author) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as simple string first
	var email string
	if err := json.Unmarshal(data, &email); err == nil {
		a.Email = email
		return nil
	}

	// Otherwise unmarshal as object using an anonymous struct to avoid recursion
	type authorAlias Author
	aux := (*authorAlias)(a)
	return json.Unmarshal(data, aux)
}

// CreateBuild - Create a build.
type CreateBuild struct {
	Commit  string `json:"commit"`
	Branch  string `json:"branch"`
	Message string `json:"message"`

	// Optional fields
	Author                      Author            `json:"author,omitempty"`
	CleanCheckout               bool              `json:"clean_checkout"`
	Env                         map[string]string `json:"env,omitempty"`
	MetaData                    map[string]string `json:"meta_data,omitempty"`
	IgnorePipelineBranchFilters bool              `json:"ignore_pipeline_branch_filters"`
	PullRequestBaseBranch       string            `json:"pull_request_base_branch,omitempty"`
	PullRequestID               int64             `json:"pull_request_id,omitempty"`
	PullRequestRepository       string            `json:"pull_request_repository,omitempty"`
}

// Creator represents who created a build
type Creator struct {
	AvatarURL string     `json:"avatar_url"`
	CreatedAt *Timestamp `json:"created_at"`
	Email     string     `json:"email"`
	ID        string     `json:"id"`
	Name      string     `json:"name"`
}

// RebuiltFrom references a previous build
type RebuiltFrom struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	URL    string `json:"url"`
}

// PullRequest represents a Github PR
type PullRequest struct {
	ID         string `json:"id,omitempty"`
	Base       string `json:"base,omitempty"`
	Repository string `json:"repository,omitempty"`
}

type TriggeredFrom struct {
	BuildID           string `json:"build_id,omitempty"`
	BuildNumber       int    `json:"build_number,omitempty"`
	BuildPipelineSlug string `json:"build_pipeline_slug,omitempty"`
}

type TestEngineSuite struct {
	ID   string `json:"id,omitempty"`
	Slug string `json:"slug,omitempty"`
}

type TestEngineRun struct {
	ID    string          `json:"id,omitempty"`
	Suite TestEngineSuite `json:"suite"`
}

type TestEngineProperty struct {
	Runs []TestEngineRun `json:"runs,omitempty"`
}

// Build represents a build which has run in buildkite
type Build struct {
	ID          string                 `json:"id,omitempty"`
	GraphQLID   string                 `json:"graphql_id,omitempty"`
	URL         string                 `json:"url,omitempty"`
	WebURL      string                 `json:"web_url,omitempty"`
	Number      int                    `json:"number,omitempty"`
	State       string                 `json:"state,omitempty"`
	Blocked     bool                   `json:"blocked"`
	Message     string                 `json:"message,omitempty"`
	Commit      string                 `json:"commit,omitempty"`
	Branch      string                 `json:"branch,omitempty"`
	Author      Author                 `json:"author,omitempty"`
	Env         map[string]interface{} `json:"env,omitempty"`
	CreatedAt   *Timestamp             `json:"created_at,omitempty"`
	ScheduledAt *Timestamp             `json:"scheduled_at,omitempty"`
	StartedAt   *Timestamp             `json:"started_at,omitempty"`
	FinishedAt  *Timestamp             `json:"finished_at,omitempty"`
	MetaData    map[string]string      `json:"meta_data,omitempty"`
	Creator     Creator                `json:"creator,omitempty"`
	Source      string                 `json:"source,omitempty"`

	// jobs run during the build
	Jobs []Job `json:"jobs,omitempty"`

	// the pipeline this build is associated with
	Pipeline *Pipeline `json:"pipeline,omitempty"`

	// the build this build is a rebuild of
	RebuiltFrom *RebuiltFrom `json:"rebuilt_from,omitempty"`

	// the pull request this build is associated with
	PullRequest *PullRequest `json:"pull_request,omitempty"`

	// the build that this build is triggered from
	// https://buildkite.com/docs/pipelines/trigger-step
	TriggeredFrom *TriggeredFrom `json:"triggered_from,omitempty"`

	TestEngine *TestEngineProperty `json:"test_engine,omitempty"`
}

type MetaDataFilters struct {
	MetaData map[string]string
}

// Encodes MetaData in the format expected by the buildkite API
// Example: ?meta_data[some-key]=some-value
func (mo MetaDataFilters) EncodeValues(parent_key string, v *url.Values) error {
	for key, val := range mo.MetaData {
		keyString := fmt.Sprintf("%s[%s]", parent_key, key)
		v.Add(keyString, val)
	}
	return nil
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

	// Filters the results by branch name(s)
	Branch []string `url:"branch,brackets,omitempty"`

	// Filters the results by builds for the specific commit SHA (full, not shortened). Default is "".
	Commit string `url:"commit,omitempty"`

	// Include all retried jobs in each build’s jobs list
	IncludeRetriedJobs bool `url:"include_retried_jobs,omitempty"`

	// Filters results by metadata.
	MetaData MetaDataFilters `url:"meta_data,omitempty"`

	// Exclude jobs from the API response
	ExcludeJobs bool `url:"exclude_jobs,omitempty"`

	// Exclude pipeline data from the API response
	ExcludePipeline bool `url:"exclude_pipeline,omitempty"`

	ListOptions
}

type BuildGetOptions struct {
	BuildsListOptions

	// Include Test Engine data
	IncludeTestEngine bool `url:"include_test_engine,omitempty"`
}

// Cancel - Trigger a cancel for the target build
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/builds#cancel-a-build
func (bs *BuildsService) Cancel(ctx context.Context, org, pipeline, build string) (Build, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/cancel", org, pipeline, build)
	req, err := bs.client.NewRequest(ctx, "PUT", u, nil)
	if err != nil {
		return Build{}, err
	}

	var result Build
	_, err = bs.client.Do(req, &result)
	if err != nil {
		return Build{}, err
	}

	return result, nil
}

// Create - Create a pipeline
//
// buildkite API docs: https://buildkite.com/docs/api/builds#create-a-build
func (bs *BuildsService) Create(ctx context.Context, org string, pipeline string, b CreateBuild) (Build, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds", org, pipeline)
	req, err := bs.client.NewRequest(ctx, "POST", u, b)
	if err != nil {
		return Build{}, nil, err
	}

	var build Build
	resp, err := bs.client.Do(req, &build)
	if err != nil {
		return Build{}, resp, err
	}

	return build, resp, err
}

// Get fetches a build.
//
// buildkite API docs: https://buildkite.com/docs/api/builds#get-a-build
func (bs *BuildsService) Get(ctx context.Context, org string, pipeline string, id string, opt *BuildGetOptions) (Build, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s", org, pipeline, id)

	u, err := addOptions(u, opt)
	if err != nil {
		return Build{}, nil, err
	}

	req, err := bs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return Build{}, nil, err
	}

	var build Build
	resp, err := bs.client.Do(req, &build)
	if err != nil {
		return Build{}, resp, err
	}

	return build, resp, err
}

// List the builds for the current user.
//
// buildkite API docs: https://buildkite.com/docs/api/builds#list-all-builds
func (bs *BuildsService) List(ctx context.Context, opt *BuildsListOptions) ([]Build, *Response, error) {
	u := "v2/builds"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var builds []Build
	resp, err := bs.client.Do(req, &builds)
	if err != nil {
		return nil, resp, err
	}

	return builds, resp, err
}

// ListByOrg lists the builds within the specified orginisation.
//
// buildkite API docs: https://buildkite.com/docs/api/builds#list-builds-for-an-organization
func (bs *BuildsService) ListByOrg(ctx context.Context, org string, opt *BuildsListOptions) ([]Build, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/builds", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var builds []Build
	resp, err := bs.client.Do(req, &builds)
	if err != nil {
		return nil, resp, err
	}

	return builds, resp, err
}

// ListByPipeline lists the builds for a pipeline within the specified originisation.
//
// buildkite API docs: https://buildkite.com/docs/api/builds#list-builds-for-a-pipeline
func (bs *BuildsService) ListByPipeline(ctx context.Context, org string, pipeline string, opt *BuildsListOptions) ([]Build, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds", org, pipeline)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var builds []Build
	resp, err := bs.client.Do(req, &builds)
	if err != nil {
		return nil, resp, err
	}

	return builds, resp, err
}

// Rebuild triggers a rebuild for the target build
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/builds#rebuild-a-build
func (bs *BuildsService) Rebuild(ctx context.Context, org, pipeline, build string) (Build, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/rebuild", org, pipeline, build)
	req, err := bs.client.NewRequest(ctx, "PUT", u, nil)
	if err != nil {
		return Build{}, err
	}

	var result Build
	_, err = bs.client.Do(req, &result)
	if err != nil {
		return Build{}, err
	}

	return result, nil
}
