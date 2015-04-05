// Copyright 2014 Mark Wolfe. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildkite

import "fmt"

// Emoji emoji, what else can you say?
type Emoji struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

// ListEmojis list all the emojis for a given account, including custom emojis and aliases.
//
// buildkite API docs: https://buildkite.com/docs/api/emojis
func (c *Client) ListEmojis(org string) ([]Emoji, *Response, error) {

	var u string

	u = fmt.Sprintf("v1/organizations/%s/emojis", org)

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	emoji := new([]Emoji)
	resp, err := c.Do(req, emoji)
	if err != nil {
		return nil, resp, err
	}

	return *emoji, resp, nil
}

// Token an oauth access token for the buildkite service
type Token struct {
	AccessToken *string `json:"access_token,omitempty"`
	Type        *string `json:"token_type,omitempty"`
}

// FindOrCreateToken will authenticate with the buildkite API and find or generate an access token
// with the specified client identifier.
//
// buildkite API docs: https://buildkite.com/docs/api/access_tokens
func (c *Client) FindOrCreateToken(username string, password string, clientID string) (*Token, *Response, error) {

	bt, err := NewBasicConfig(username, password)

	if err != nil {
		return nil, nil, err
	}

	c.client.Transport = bt

	params := map[string]interface{}{
		"client_id": clientID,
		"scopes":    []string{"read_user", "read_projects", "read_builds", "write_builds", "read_build_logs", "read_agents", "write_agents", "read_organizations"},
	}

	req, err := c.NewRequest("POST", "v1/access_tokens", params)

	token := new(Token)
	resp, err := c.Do(req, token)
	if err != nil {
		return nil, resp, err
	}

	tt, err := NewTokenConfig(*token.AccessToken)

	if err != nil {
		return nil, nil, err
	}

	c.client.Transport = tt

	return token, resp, err
}
