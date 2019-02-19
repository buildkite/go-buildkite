package buildkite

import (
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
	ID              *string    `json:"id,omitempty"`
	Type            *string    `json:"type,omitempty"`
	Name            *string    `json:"name,omitempty"`
	State           *string    `json:"state,omitempty"`
	LogsURL         *string    `json:"logs_url,omitempty"`
	RawLogsURL      *string    `json:"raw_log_url,omitempty"`
	Command         *string    `json:"command,omitempty"`
	ExitStatus      *int       `json:"exit_status,omitempty"`
	ArtifactPaths   *string    `json:"artifact_paths,omitempty"`
	CreatedAt       *Timestamp `json:"created_at,omitempty"`
	ScheduledAt     *Timestamp `json:"scheduled_at,omitempty"`
	RunnableAt      *Timestamp `json:"runnable_at,omitempty"`
	StartedAt       *Timestamp `json:"started_at,omitempty"`
	FinishedAt      *Timestamp `json:"finished_at,omitempty"`
	Agent           Agent      `json:"agent,omitempty"`
	AgentQueryRules []string   `json:"agent_query_rules,omitempty"`
	WebURL          string     `json:"web_url"`
}

// JobUnblockOptions specifies the optional parameters to UnblockJob
type JobUnblockOptions struct {
	Fields map[string]string `json:"fields,omitempty"`
}

// UnblockJob - unblock a job
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/jobs#unblock-a-job
func (js *JobsService) UnblockJob(org string, pipeline string, buildNumber string, jobID string, opt *JobUnblockOptions) (*Job, *Response, error) {
	var u string

	u = fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/jobs/%s/unblock", org, pipeline, buildNumber, jobID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := js.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	job := new(Job)
	resp, err := js.client.Do(req, job)
	if err != nil {
		return nil, resp, err
	}

	return job, resp, err
}
