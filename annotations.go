package buildkite

import (
	"context"
	"fmt"
	"net/http"
)

// AnnotationsService handles communication with the annotation related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api/annotations
type AnnotationsService struct {
	client *Client
}

// Annotation represents an annotation which has been stored from a build
type Annotation struct {
	ID        string     `json:"id,omitempty"`
	Context   string     `json:"context,omitempty"`
	Style     string     `json:"style,omitempty"`
	Scope     string     `json:"scope,omitempty"`
	Priority  int        `json:"priority,omitempty"`
	BodyHTML  string     `json:"body_html,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	UpdatedAt *Timestamp `json:"updated_at,omitempty"`
}

type AnnotationCreate struct {
	Body     string `json:"body,omitempty"`
	Context  string `json:"context,omitempty"`
	Style    string `json:"style,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Append   bool   `json:"append"`
}

// AnnotationListOptions specifies the optional parameters to the
// AnnoationsService.List method.
type AnnotationListOptions struct {
	ListOptions
}

// ListByBuild gets annotations for a specific build
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/annotations#list-annotations-for-a-build
func (as *AnnotationsService) ListByBuild(ctx context.Context, org string, pipeline string, build string, opt *AnnotationListOptions) ([]Annotation, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/annotations", org, pipeline, build)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := as.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var annotations []Annotation
	resp, err := as.client.Do(req, &annotations)
	if err != nil {
		return nil, resp, err
	}
	return annotations, resp, err
}

func (as *AnnotationsService) Create(ctx context.Context, org, pipeline, build string, ac AnnotationCreate) (Annotation, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/annotations", org, pipeline, build)
	req, err := as.client.NewRequest(ctx, "POST", u, ac)
	if err != nil {
		return Annotation{}, nil, err
	}

	var annotation Annotation
	resp, err := as.client.Do(req, &annotation)
	if err != nil {
		return Annotation{}, resp, err
	}

	return annotation, resp, err
}

func (as *AnnotationsService) Delete(ctx context.Context, org, pipeline, build, annotationUUID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/annotations/%s", org, pipeline, build, annotationUUID)
	req, err := as.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return as.client.Do(req, nil)
}

// ListByJob gets annotations for a specific job
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/annotations#list-annotations-for-a-job
func (as *AnnotationsService) ListByJob(ctx context.Context, org, pipeline, build, jobID string, opt *AnnotationListOptions) ([]Annotation, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/jobs/%s/annotations", org, pipeline, build, jobID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := as.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var annotations []Annotation
	resp, err := as.client.Do(req, &annotations)
	if err != nil {
		return nil, resp, err
	}
	return annotations, resp, err
}

// CreateForJob creates an annotation scoped to a specific job
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/annotations#create-an-annotation-on-a-job
func (as *AnnotationsService) CreateForJob(ctx context.Context, org, pipeline, build, jobID string, ac AnnotationCreate) (Annotation, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/jobs/%s/annotations", org, pipeline, build, jobID)
	req, err := as.client.NewRequest(ctx, "POST", u, ac)
	if err != nil {
		return Annotation{}, nil, err
	}

	var annotation Annotation
	resp, err := as.client.Do(req, &annotation)
	if err != nil {
		return Annotation{}, resp, err
	}

	return annotation, resp, err
}

// DeleteForJob deletes an annotation on a job
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/annotations#delete-an-annotation-on-a-job
func (as *AnnotationsService) DeleteForJob(ctx context.Context, org, pipeline, build, jobID, annotationUUID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipelines/%s/builds/%s/jobs/%s/annotations/%s", org, pipeline, build, jobID, annotationUUID)
	req, err := as.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return as.client.Do(req, nil)
}
