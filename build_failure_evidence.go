package buildkite

import (
	"context"
	"fmt"
)

// BuildFailureEvidenceOptions controls the bounded, log-free evidence response.
type BuildFailureEvidenceOptions struct {
	// JobLimit is 1..50; omitted defaults to 10. Invalid values are rejected by the API.
	JobLimit int `url:"job_limit,omitempty"`
	// Include accepts "test_results" and requires read_suites in addition to read_builds.
	Include string `url:"include,omitempty"`
}

// BuildFailureEvidence is a snapshot of terminal hard-failed script jobs,
// warning/error annotations and optionally caller-authorized Test Engine results.
// Soft failures and superseded retries are excluded. Logs are never included.
type BuildFailureEvidence struct {
	SchemaVersion        int                         `json:"schema_version"`
	Build                FailureEvidenceBuild        `json:"build"`
	Jobs                 []FailureEvidenceJob        `json:"jobs"`
	JobsTruncated        bool                        `json:"jobs_truncated"`
	JobLimit             int                         `json:"job_limit"`
	Annotations          []FailureEvidenceAnnotation `json:"annotations"`
	AnnotationsTruncated bool                        `json:"annotations_truncated"`
	CollectedAt          *Timestamp                  `json:"collected_at"`
}

type FailureEvidenceBuild struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	WebURL string `json:"web_url"`
	State  string `json:"state"`
}

type FailureEvidenceJob struct {
	ID                   string                      `json:"id"`
	Name                 string                      `json:"name"`
	StepKey              *string                     `json:"step_key"`
	State                string                      `json:"state"`
	ExitStatus           int                         `json:"exit_status"`
	StartedAt            *Timestamp                  `json:"started_at"`
	FinishedAt           *Timestamp                  `json:"finished_at"`
	WebURL               string                      `json:"web_url"`
	Annotations          []FailureEvidenceAnnotation `json:"annotations"`
	AnnotationsTruncated bool                        `json:"annotations_truncated"`
	TestResults          FailureEvidenceTestResults  `json:"test_results"`
	CollectedAt          *Timestamp                  `json:"collected_at"`
	Omissions            map[string]string           `json:"omissions"`
}

type FailureEvidenceAnnotation struct {
	ID            string  `json:"id"`
	JobID         *string `json:"job_id"`
	Scope         string  `json:"scope"`
	Context       *string `json:"context"`
	Style         string  `json:"style"`
	Body          string  `json:"body"`
	BodyTruncated bool    `json:"body_truncated"`
}

type FailureEvidenceTestResults struct {
	// Status is found, none_recorded, no_accessible_data, unavailable, or not_requested.
	Status            string                       `json:"status"`
	Counts            map[string]int               `json:"counts,omitempty"`
	Failures          []FailureEvidenceTestFailure `json:"failures,omitempty"`
	FailuresTruncated bool                         `json:"failures_truncated,omitempty"`
	Window            *FailureEvidenceTestWindow   `json:"window,omitempty"`
	CollectedAt       *Timestamp                   `json:"collected_at"`
}

type FailureEvidenceTestWindow struct {
	From      *Timestamp `json:"from"`
	To        *Timestamp `json:"to"`
	Truncated bool       `json:"truncated"`
}

type FailureEvidenceTestFailure struct {
	TestID           string     `json:"test_id"`
	SuiteID          string     `json:"suite_id"`
	RunID            string     `json:"run_id"`
	ExecutionID      *string    `json:"execution_id"`
	Timestamp        *Timestamp `json:"timestamp"`
	Name             string     `json:"name"`
	Scope            string     `json:"scope"`
	Location         string     `json:"location"`
	FailureReason    string     `json:"failure_reason"`
	FailedAttempts   int        `json:"failed_attempts"`
	DetailsTruncated bool       `json:"details_truncated"`
}

// GetFailureEvidence fetches one bounded snapshot. The endpoint requires read_builds
// and organization rollout access; including test_results also requires read_suites.
func (bs *BuildsService) GetFailureEvidence(ctx context.Context, org, pipeline, buildNumber string, opt *BuildFailureEvidenceOptions) (BuildFailureEvidence, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/failure-evidence", org, pipeline, buildNumber)
	u, err := addOptions(u, opt)
	if err != nil {
		return BuildFailureEvidence{}, nil, err
	}
	req, err := bs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return BuildFailureEvidence{}, nil, err
	}
	var evidence BuildFailureEvidence
	resp, err := bs.client.Do(req, &evidence)
	if err != nil {
		return BuildFailureEvidence{}, resp, err
	}
	return evidence, resp, nil
}
