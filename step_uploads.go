package buildkite

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// StepUploadsService handles communication with the step upload related
// methods of the buildkite API. Step uploads are the dynamic pipeline
// configurations uploaded by `buildkite-agent pipeline upload` while a
// build runs.
//
// Step uploads retrieval are only available while the build is within its maximum
// lifetime; older builds return 410 Gone.
type StepUploadsService struct {
	client *Client
}

// StepUpload represents one `buildkite-agent pipeline upload` performed
// during a build.
//
// CreatedJobsCount is nil until the upload has been applied; 0 means the
// upload was applied but created no jobs.
//
// The 'DefinitionXXXX' fields are only populated by Get: DefinitionYAML carries
// the uploaded configuration re-rendered as YAML, unless the definition
// exceeds the API's render limit, in which case DefinitionYAML is nil and
// DefinitionYAMLOmitted is true. DefinitionStoredBytes is the compressed
// stored size; DefinitionBytes (Get only) is the exact serialized size.
type StepUpload struct {
	UUID                  string     `json:"uuid,omitempty"`
	GraphQLID             string     `json:"graphql_id,omitempty"`
	State                 string     `json:"state,omitempty"`
	Source                string     `json:"source,omitempty"`
	SourceJobID           string     `json:"source_job_id,omitempty"`
	ReplaceExistingSteps  bool       `json:"replace_existing_steps"`
	CreatedJobsCount      *int       `json:"created_jobs_count,omitempty"`
	RejectionType         *string    `json:"rejection_type,omitempty"`
	Message               *string    `json:"message,omitempty"`
	URL                   string     `json:"url,omitempty"`
	CreatedAt             *Timestamp `json:"created_at,omitempty"`
	ProcessedAt           *Timestamp `json:"processed_at,omitempty"`
	DefinitionStoredBytes int        `json:"definition_stored_bytes,omitempty"`

	// Present on Get only.
	CreatedJobUUIDs       []string `json:"created_job_uuids,omitempty"`
	DefinitionBytes       *int     `json:"definition_bytes,omitempty"`
	DefinitionYAML        *string  `json:"definition_yaml,omitempty"`
	DefinitionYAMLOmitted *bool    `json:"definition_yaml_omitted,omitempty"`
}

// StepUploadsListOptions specifies the optional parameters to the
// StepUploadsService.ListByBuild method.
type StepUploadsListOptions struct {
	// SourceJobID filters to uploads performed by one job (a job UUID).
	SourceJobID string `url:"filter[source_job_id],omitempty"`

	PerPage int    `url:"per_page,omitempty"`
	After   string `url:"after,omitempty"`
	Before  string `url:"before,omitempty"`
}

type StepUploadsListLink string

func (l StepUploadsListLink) ToOptions() (*StepUploadsListOptions, error) {
	u, err := url.Parse(string(l))
	if err != nil {
		return nil, fmt.Errorf("parsing link: %w", err)
	}

	q := u.Query()

	opts := &StepUploadsListOptions{
		SourceJobID: q.Get("filter[source_job_id]"),
		After:       q.Get("after"),
		Before:      q.Get("before"),
	}

	if perPage := q.Get("per_page"); perPage != "" {
		opts.PerPage, err = strconv.Atoi(perPage)
		if err != nil {
			return nil, fmt.Errorf("parsing per_page: %w", err)
		}
	}

	return opts, nil
}

type StepUploadsListLinks struct {
	First    StepUploadsListLink `json:"first,omitempty"`
	Previous StepUploadsListLink `json:"prev,omitempty"`
	Self     StepUploadsListLink `json:"self,omitempty"`
	Next     StepUploadsListLink `json:"next,omitempty"`
}

type StepUploadsList struct {
	Items []StepUpload         `json:"items"`
	Links StepUploadsListLinks `json:"links"`
}

// ListByBuild lists a build's step uploads, newest first, without their
// definitions.
func (s *StepUploadsService) ListByBuild(ctx context.Context, org, pipeline, buildNumber string, opt *StepUploadsListOptions) (StepUploadsList, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/step-uploads", org, pipeline, buildNumber)
	u, err := addOptions(u, opt)
	if err != nil {
		return StepUploadsList{}, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return StepUploadsList{}, nil, err
	}

	var uploads StepUploadsList
	resp, err := s.client.Do(req, &uploads)
	if err != nil {
		return StepUploadsList{}, resp, err
	}

	return uploads, resp, err
}

// Get returns a single step upload, including its uploaded definition
// rendered as YAML (subject to the API's render limit — see StepUpload).
func (s *StepUploadsService) Get(ctx context.Context, org, pipeline, buildNumber, uploadUUID string) (StepUpload, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/step-uploads/%s", org, pipeline, buildNumber, uploadUUID)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return StepUpload{}, nil, err
	}

	var upload StepUpload
	resp, err := s.client.Do(req, &upload)
	if err != nil {
		return StepUpload{}, resp, err
	}

	return upload, resp, err
}
