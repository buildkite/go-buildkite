package buildkite

import (
	"context"
	"fmt"
)

// TeamPipelinesService handles communication with the team pipelines related
// methods of the buildkite API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/teams/pipelines

type TeamPipelinesService struct {
	client *Client
}

type TeamPipeline struct {
	ID          string     `json:"pipeline_id,omitempty"`
	URL         string     `json:"pipeline_url,omitempty"`
	AccessLevel string     `json:"access_level,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
}

type CreateTeamPipelines struct {
	PipelineID  string `json:"pipeline_id,omitempty"`
	AccessLevel string `json:"access_level,omitempty"`
}

type UpdateTeamPipelines struct {
	AccessLevel string `json:"access_level,omitempty"`
}

type TeamPipelinesListOptions struct {
	ListOptions
}

func (tps *TeamPipelinesService) List(ctx context.Context, org string, id string, opt *TeamPipelinesListOptions) ([]TeamPipeline, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/pipelines", org, id)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := tps.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var teamPipelines []TeamPipeline
	resp, err := tps.client.Do(req, &teamPipelines)
	if err != nil {
		return nil, resp, err
	}

	return teamPipelines, resp, err
}

func (tps *TeamPipelinesService) Get(ctx context.Context, org string, teamID string, pipelineID string) (TeamPipeline, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/pipelines/%s", org, teamID, pipelineID)

	req, err := tps.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TeamPipeline{}, nil, err
	}

	var tp TeamPipeline
	resp, err := tps.client.Do(req, &tp)
	if err != nil {
		return TeamPipeline{}, resp, err
	}

	return tp, resp, err
}

func (tps *TeamPipelinesService) Create(ctx context.Context, org string, teamID string, ctp CreateTeamPipelines) (TeamPipeline, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/pipelines", org, teamID)

	req, err := tps.client.NewRequest(ctx, "POST", u, ctp)
	if err != nil {
		return TeamPipeline{}, nil, err
	}

	var tp TeamPipeline
	resp, err := tps.client.Do(req, &tp)
	if err != nil {
		return tp, resp, err
	}

	return tp, resp, err
}

func (tps *TeamPipelinesService) Update(ctx context.Context, org string, teamID string, pipelineID string, utp UpdateTeamPipelines) (TeamPipeline, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/pipelines/%s", org, teamID, pipelineID)

	req, err := tps.client.NewRequest(ctx, "PATCH", u, utp)
	if err != nil {
		return TeamPipeline{}, nil, err
	}

	var tp TeamPipeline
	resp, err := tps.client.Do(req, &tp)
	if err != nil {
		return TeamPipeline{}, resp, err
	}

	return tp, resp, err
}

func (tps *TeamPipelinesService) Delete(ctx context.Context, org string, teamID string, pipelineID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/pipelines/%s", org, teamID, pipelineID)

	req, err := tps.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := tps.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
