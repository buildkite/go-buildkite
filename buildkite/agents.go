// Copyright 2014 Mark Wolfe. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildkite

import "fmt"

// AgentsService handles communication with the agent related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api/agents
type AgentsService struct {
	client *Client
}

// Agent represents a buildkite build agent.
type Agent struct {
	ID             string     `json:"id"`
	URL            string     `json:"url"`
	Name           string     `json:"name"`
	ConnectedState string     `json:"connection_state"`
	AgentToken     string     `json:"access_token"`
	Hostname       string     `json:"hostname"`
	IPAddress      string     `json:"ip_address"`
	UserAgent      string     `json:"user_agent"`
	CreatedAt      *Timestamp `json:"created_at,omitempty"`

	// the user which created the agent
	Creator *User `json:"creator"`
}

// AgentListOptions specifies the optional parameters to the
// AgentService.List method.
type AgentListOptions struct {
	ListOptions
}

// List the agents for a given orginisation.
//
// buildkite API docs: https://buildkite.com/docs/api/agents#list-agents
func (as *AgentsService) List(org string, opt *AgentListOptions) ([]Agent, *Response, error) {
	var u string

	u = fmt.Sprintf("v1/organizations/%s/agents", org)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := as.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	agents := new([]Agent)
	resp, err := as.client.Do(req, agents)
	if err != nil {
		return nil, resp, err
	}

	return *agents, resp, err
}

// Get fetches an agent.
//
// buildkite API docs: https://buildkite.com/docs/api/agents#get-an-agent
func (as *AgentsService) Get(org string, id string) (*Agent, *Response, error) {

	u := fmt.Sprintf("v1/organizations/%s/agents/%s", org, id)

	req, err := as.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	agent := new(Agent)
	resp, err := as.client.Do(req, agent)
	if err != nil {
		return nil, resp, err
	}

	return agent, resp, err
}

// Create a new buildkite agent.
//
// buildkite API docs: https://buildkite.com/docs/api/agents#create-an-agent
func (as *AgentsService) Create(org string, name string) (*Agent, *Response, error) {

	var u string

	u = fmt.Sprintf("v1/organizations/%s/agents", org)

	params := map[string]string{
		"name": name,
	}

	req, err := as.client.NewRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	agent := new(Agent)
	resp, err := as.client.Do(req, agent)
	if err != nil {
		return nil, resp, err
	}

	return agent, resp, err
}
