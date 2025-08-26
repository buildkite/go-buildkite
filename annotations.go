package buildkite

import (
	"context"
	"fmt"
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
	BodyHTML  string     `json:"body_html,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	UpdatedAt *Timestamp `json:"updated_at,omitempty"`
}

type AnnotationCreate struct {
	Body    string `json:"body,omitempty"`
	Context string `json:"context,omitempty"`
	Style   string `json:"style,omitempty"`
	Append  bool   `json:"append"`
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
