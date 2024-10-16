package buildkite

import (
	"context"
	"fmt"
)

type AccessTokensService struct {
	client *Client
}

type AccessToken struct {
	UUID   *string   `json:"uuid,omitempty"`
	Scopes *[]string `json:"scopes,omitempty"`
}

// Get gets the current token which was used to authenticate the request
//
// buildkite API docs: https://buildkite.com/docs/rest-api/access-token
func (ats *AccessTokensService) Get(ctx context.Context) (*AccessToken, *Response, error) {

	var u string

	u = fmt.Sprintf("v2/access-token")

	req, err := ats.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	accessToken := new(AccessToken)

	resp, err := ats.client.Do(req, accessToken)
	if err != nil {
		return nil, resp, err
	}

	return accessToken, resp, nil
}

// Revokes the current token which was used to authenticate the request
//
// buildkite API docs: https://buildkite.com/docs/rest-api/access-token
func (ats *AccessTokensService) Revoke(ctx context.Context) (*Response, error) {

	var u string

	u = fmt.Sprintf("v2/access-token")

	req, err := ats.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := ats.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
