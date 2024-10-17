package buildkite

import (
	"context"
	"fmt"
)

// JobsService handles communication with the job related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api/jobs
type JobsService struct {
	client *Client
}

// Job represents a job run during a build in buildkite
type Job struct {
	ID                 string          `json:"id,omitempty"`
	GraphQLID          string          `json:"graphql_id,omitempty"`
	Type               string          `json:"type,omitempty"`
	Name               string          `json:"name,omitempty"`
	Label              string          `json:"label,omitempty"`
	StepKey            string          `json:"step_key,omitempty"`
	GroupKey           string          `json:"group_key,omitempty"`
	State              string          `json:"state,omitempty"`
	LogsURL            string          `json:"logs_url,omitempty"`
	RawLogsURL         string          `json:"raw_log_url,omitempty"`
	Command            string          `json:"command,omitempty"`
	ExitStatus         *int            `json:"exit_status,omitempty"`
	ArtifactPaths      string          `json:"artifact_paths,omitempty"`
	ArtifactsURL       string          `json:"artifacts_url,omitempty"`
	CreatedAt          *Timestamp      `json:"created_at,omitempty"`
	ScheduledAt        *Timestamp      `json:"scheduled_at,omitempty"`
	RunnableAt         *Timestamp      `json:"runnable_at,omitempty"`
	StartedAt          *Timestamp      `json:"started_at,omitempty"`
	FinishedAt         *Timestamp      `json:"finished_at,omitempty"`
	UnblockedAt        *Timestamp      `json:"unblocked_at,omitempty"`
	Agent              Agent           `json:"agent,omitempty"`
	AgentQueryRules    []string        `json:"agent_query_rules,omitempty"`
	WebURL             string          `json:"web_url"`
	Retried            bool            `json:"retried,omitempty"`
	RetriedInJobID     string          `json:"retried_in_job_id,omitempty"`
	RetriesCount       int             `json:"retries_count,omitempty"`
	RetrySource        *JobRetrySource `json:"retry_source,omitempty"`
	RetryType          string          `json:"retry_type,omitempty"`
	SoftFailed         bool            `json:"soft_failed,omitempty"`
	UnblockedBy        *UnblockedBy    `json:"unblocked_by,omitempty"`
	Unblockable        bool            `json:"unblockable,omitempty"`
	UnblockURL         string          `json:"unblock_url,omitempty"`
	ParallelGroupIndex *int            `json:"parallel_group_index,omitempty"`
	ParallelGroupTotal *int            `json:"parallel_group_total,omitempty"`
	ClusterID          string          `json:"cluster_id,omitempty"`
	ClusterQueueID     string          `json:"cluster_queue_id,omitempty"`
	TriggeredBuild     *TriggeredBuild `json:"triggered_build,omitempty"`
	Priority           *JobPriority    `json:"priority"`
}

// JobRetrySource represents what triggered this retry.
type JobRetrySource struct {
	JobID     string `json:"job_id,omitempty"`
	RetryType string `json:"retry_type,omitempty"`
}

// UnblockedBy represents the unblocked status of a job, when present
type UnblockedBy struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// JobUnblockOptions specifies the optional parameters to UnblockJob
type JobUnblockOptions struct {
	Fields map[string]string `json:"fields,omitempty"`
}

// JobLog represents a job log output
type JobLog struct {
	URL         string  `json:"url"`
	Content     string  `json:"content"`
	Size        int     `json:"size"`
	HeaderTimes []int64 `json:"header_times"`
}

// JobEnvs represent job environments output
type JobEnvs struct {
	EnvironmentVariables map[string]string `json:"env"`
}

type TriggeredBuild struct {
	ID     string `json:"id,omitempty"`
	Number int    `json:"number,omitempty"`
	URL    string `json:"url,omitempty"`
	WebURL string `json:"web_url,omitempty"`
}

// JobPriority represents the priority of the job
type JobPriority struct {
	Number int `json:"number,omitempty"`
}

// UnblockJob - unblock a job
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/jobs#unblock-a-job
func (js *JobsService) UnblockJob(ctx context.Context, org string, pipeline string, buildNumber string, jobID string, opt *JobUnblockOptions) (Job, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/jobs/%s/unblock", org, pipeline, buildNumber, jobID)
	u, err := addOptions(u, opt)
	if err != nil {
		return Job{}, nil, err
	}

	if opt == nil {
		opt = &JobUnblockOptions{}
	}

	req, err := js.client.NewRequest(ctx, "PUT", u, opt)
	if err != nil {
		return Job{}, nil, err
	}

	var job Job
	resp, err := js.client.Do(req, &job)
	if err != nil {
		return Job{}, resp, err
	}

	return job, resp, err
}

// RetryJob - retry a job
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/jobs#retry-a-job
func (js *JobsService) RetryJob(ctx context.Context, org string, pipeline string, buildNumber string, jobID string) (Job, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/jobs/%s/retry", org, pipeline, buildNumber, jobID)
	req, err := js.client.NewRequest(ctx, "PUT", u, nil)
	if err != nil {
		return Job{}, nil, err
	}

	var job Job
	resp, err := js.client.Do(req, &job)
	if err != nil {
		return Job{}, resp, err
	}

	return job, resp, err
}

// GetJobLog - get a job’s log output
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/jobs#get-a-jobs-log-output
func (js *JobsService) GetJobLog(ctx context.Context, org string, pipeline string, buildNumber string, jobID string) (JobLog, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/jobs/%s/log", org, pipeline, buildNumber, jobID)
	req, err := js.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return JobLog{}, nil, err
	}

	req.Header.Set("Accept", "application/json")

	var jobLog JobLog
	resp, err := js.client.Do(req, &jobLog)
	if err != nil {
		return JobLog{}, resp, err
	}

	return jobLog, resp, err
}

// GetJobEnvironmentVariables - get a job’s environment variables
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/jobs#get-a-jobs-environment-variables
func (js *JobsService) GetJobEnvironmentVariables(ctx context.Context, org string, pipeline string, buildNumber string, jobID string) (JobEnvs, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/jobs/%s/env", org, pipeline, buildNumber, jobID)
	req, err := js.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return JobEnvs{}, nil, err
	}

	req.Header.Set("Accept", "application/json")

	var jobEnvs JobEnvs
	resp, err := js.client.Do(req, &jobEnvs)
	if err != nil {
		return JobEnvs{}, resp, err
	}

	return jobEnvs, resp, err
}
