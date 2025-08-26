package buildkite

import (
	"context"
	"fmt"
)

// PipelineTemplatesService handles communication with pipeline template related
// methods of the Buildkite API.
//
// Buildkite API docs: <to-fill>
type PipelineTemplatesService struct {
	client *Client
}

type PipelineTemplate struct {
	UUID          string `json:"uuid,omitempty"`
	GraphQLID     string `json:"graphql_id,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Configuration string `json:"configuration,omitempty"`
	Available     bool   `json:"available"`
	URL           string `json:"url,omitempty"`
	WebURL        string `json:"web_url,omitempty"`

	CreatedAt *Timestamp              `json:"created_at,omitempty"`
	CreatedBy PipelineTemplateCreator `json:"created_by,omitempty"`

	UpdatedAt *Timestamp              `json:"updated_at,omitempty"`
	UpdatedBy PipelineTemplateCreator `json:"updated_by,omitempty"`
}

type PipelineTemplateCreateUpdate struct {
	Name          string `json:"name,omitempty"`
	Configuration string `json:"configuration,omitempty"`
	Description   string `json:"description,omitempty"`
	Available     bool   `json:"available"`
}

type PipelineTemplateCreator struct {
	ID        string     `json:"id,omitempty"`
	GraphQLID string     `json:"graphql_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

type PipelineTemplateListOptions struct {
	ListOptions
}

func (pts *PipelineTemplatesService) List(ctx context.Context, org string, opt *PipelineTemplateListOptions) ([]PipelineTemplate, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipeline-templates", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := pts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var templates []PipelineTemplate
	resp, err := pts.client.Do(req, &templates)
	if err != nil {
		return nil, resp, err
	}

	return templates, resp, err
}

func (pts *PipelineTemplatesService) Get(ctx context.Context, org, templateUUID string) (PipelineTemplate, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipeline-templates/%s", org, templateUUID)
	req, err := pts.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return PipelineTemplate{}, nil, err
	}

	var template PipelineTemplate
	resp, err := pts.client.Do(req, &template)
	if err != nil {
		return PipelineTemplate{}, resp, err
	}

	return template, resp, err
}

func (pts *PipelineTemplatesService) Create(ctx context.Context, org string, ptc PipelineTemplateCreateUpdate) (PipelineTemplate, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipeline-templates", org)
	req, err := pts.client.NewRequest(ctx, "POST", u, ptc)
	if err != nil {
		return PipelineTemplate{}, nil, err
	}

	var template PipelineTemplate
	resp, err := pts.client.Do(req, &template)
	if err != nil {
		return PipelineTemplate{}, resp, err
	}

	return template, resp, err
}

func (pts *PipelineTemplatesService) Update(ctx context.Context, org, templateUUID string, ptu PipelineTemplateCreateUpdate) (PipelineTemplate, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipeline-templates/%s", org, templateUUID)
	req, err := pts.client.NewRequest(ctx, "PATCH", u, ptu)
	if err != nil {
		return PipelineTemplate{}, nil, err
	}

	var template PipelineTemplate
	resp, err := pts.client.Do(req, &template)
	if err != nil {
		return PipelineTemplate{}, resp, err
	}

	return template, resp, err
}

func (pts *PipelineTemplatesService) Delete(ctx context.Context, org, templateUUID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/pipeline-templates/%s", org, templateUUID)
	req, err := pts.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return pts.client.Do(req, nil)
}
