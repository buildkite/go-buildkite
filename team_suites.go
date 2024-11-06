package buildkite

import (
	"context"
	"fmt"
)

// TeamSuitesService handles communication with the team pipelines related
// methods of the buildkite API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/teams/suites

type TeamSuitesService struct {
	client *Client
}

type TeamSuites struct {
	ID          string     `json:"suite_id,omitempty"`
	URL         string     `json:"suite_url,omitempty"`
	AccessLevel []string   `json:"access_level,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
}

type CreateTeamSuites struct {
	SuiteID     string   `json:"suite_id,omitempty"`
	AccessLevel []string `json:"access_level,omitempty"`
}

type UpdateTeamSuites struct {
	AccessLevel []string `json:"access_level,omitempty"`
}

type TeamSuitesListOptions struct {
	ListOptions
}

func (tss *TeamSuitesService) List(ctx context.Context, org string, id string, opt *TeamSuitesListOptions) ([]TeamSuites, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/suites", org, id)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := tss.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var TeamSuites []TeamSuites
	resp, err := tss.client.Do(req, &TeamSuites)
	if err != nil {
		return nil, resp, err
	}

	return TeamSuites, resp, err
}

func (tss *TeamSuitesService) Get(ctx context.Context, org string, teamID string, suiteID string) (TeamSuites, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/suites/%s", org, teamID, suiteID)

	req, err := tss.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TeamSuites{}, nil, err
	}

	var ts TeamSuites
	resp, err := tss.client.Do(req, &ts)
	if err != nil {
		return TeamSuites{}, resp, err
	}

	return ts, resp, err
}

func (tss *TeamSuitesService) Create(ctx context.Context, org string, teamID string, cts CreateTeamSuites) (TeamSuites, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/suite", org, teamID)

	req, err := tss.client.NewRequest(ctx, "POST", u, cts)
	if err != nil {
		return TeamSuites{}, nil, err
	}

	var ts TeamSuites
	resp, err := tss.client.Do(req, &ts)
	if err != nil {
		return TeamSuites{}, resp, err
	}

	return ts, resp, err
}

func (tss *TeamSuitesService) Update(ctx context.Context, org string, teamID string, pipelineID string, utp UpdateTeamSuites) (TeamSuites, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/suites/%s", org, teamID, pipelineID)

	req, err := tss.client.NewRequest(ctx, "PATCH", u, utp)
	if err != nil {
		return TeamSuites{}, nil, err
	}

	var ts TeamSuites
	resp, err := tss.client.Do(req, &ts)
	if err != nil {
		return TeamSuites{}, resp, err
	}

	return ts, resp, err
}

func (tss *TeamSuitesService) Delete(ctx context.Context, org string, teamID string, suiteID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/teams/%s/suites/%s", org, teamID, suiteID)

	req, err := tss.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := tss.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
