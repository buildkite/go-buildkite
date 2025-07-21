package buildkite

import (
	"context"
)

type AccessTokensService struct {
	client *Client
}

type AccessToken struct {
	UUID        string     `json:"uuid"`        // The token's UUID
	Scopes      []string   `json:"scopes"`      // The scopes the token has access to
	Description string     `json:"description"` // A description for the access token
	CreatedAt   *Timestamp `json:"created_at"`  // The date and time the access token was created
	User        struct {
		Name  string `json:"name"`  // The name of the user who the access token belongs to
		Email string `json:"email"` // The email of the user who the access token belongs to
	} `json:"user"` // The user who created the access token
}

// Get gets the current token which was used to authenticate the request
//
// buildkite API docs: https://buildkite.com/docs/rest-api/access-token
func (ats *AccessTokensService) Get(ctx context.Context) (AccessToken, *Response, error) {
	u := "v2/access-token"
	req, err := ats.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return AccessToken{}, nil, err
	}

	var accessToken AccessToken
	resp, err := ats.client.Do(req, &accessToken)
	if err != nil {
		return AccessToken{}, resp, err
	}

	return accessToken, resp, nil
}

// Revokes the current token which was used to authenticate the request
//
// buildkite API docs: https://buildkite.com/docs/rest-api/access-token
func (ats *AccessTokensService) Revoke(ctx context.Context) (*Response, error) {
	u := "v2/access-token"
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
