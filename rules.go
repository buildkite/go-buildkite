package buildkite

import (
	"context"
	"fmt"
	"net/http"
)

type RulesService struct {
	client *Client
}

type Rule struct {
	Action           string      `json:"action,omitempty"`
	CreatedAt        *Timestamp  `json:"created_at,omitempty"`
	CreatedBy        RuleCreator `json:"created_by"`
	Effect           string      `json:"effect,omitempty"`
	GraphqlID        string      `json:"graphql_id,omitempty"`
	OrganizationUUID string      `json:"organization_uuid,omitempty"`
	SourceType       string      `json:"source_type,omitempty"`
	SourceUUID       string      `json:"source_uuid,omitempty"`
	TargetType       string      `json:"target_type,omitempty"`
	TargetUUID       string      `json:"target_uuid,omitempty"`
	Type             string      `json:"type,omitempty"`
	URL              string      `json:"url,omitempty"`
	UUID             string      `json:"uuid,omitempty"`
}

type RuleCreator struct {
	AvatarURL string     `json:"avatar_url,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	Email     string     `json:"email,omitempty"`
	GraphqlID string     `json:"graphql_id,omitempty"`
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
}

type RuleCreate struct {
	Type  string    `json:"type,omitempty"`
	Value RuleValue `json:"value"`
}

type RuleValue struct {
	SourcePipeline string `json:"source_pipeline,omitempty"`
	TargetPipeline string `json:"target_pipeline,omitempty"`
}

type RulesListOptions struct {
	ListOptions
}

func (rs *RulesService) List(ctx context.Context, org string, opt *RulesListOptions) ([]Rule, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/rules", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := rs.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var rules []Rule
	resp, err := rs.client.Do(req, &rules)
	if err != nil {
		return nil, resp, err
	}

	return rules, resp, err
}

func (rs *RulesService) Get(ctx context.Context, org, ruleUUID string) (Rule, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/rules/%s", org, ruleUUID)
	req, err := rs.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return Rule{}, nil, err
	}

	var rule Rule
	resp, err := rs.client.Do(req, &rule)
	if err != nil {
		return Rule{}, resp, err
	}

	return rule, resp, err
}

func (rs *RulesService) Create(ctx context.Context, org string, rc RuleCreate) (Rule, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/rules", org)
	req, err := rs.client.NewRequest(ctx, http.MethodPost, u, rc)
	if err != nil {
		return Rule{}, nil, err
	}

	var rule Rule
	resp, err := rs.client.Do(req, &rule)
	if err != nil {
		return Rule{}, resp, err
	}

	return rule, resp, err
}

func (rs *RulesService) Delete(ctx context.Context, org, ruleUUID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/rules/%s", org, ruleUUID)
	req, err := rs.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return rs.client.Do(req, nil)
}
