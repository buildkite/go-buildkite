package buildkite

import (
	"context"
	"fmt"
	"time"
)

// RateLimitService handles communication with the rate limit related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/organizations/rate-limits
type RateLimitService struct {
	client *Client
}

// RateLimit represents the shape of the rate limit response (a list of scopes)
type RateLimit struct {
	Scopes *RateLimitScopes `json:"scopes,omitempty"`
}

// RateLimitScopes are returned with both REST and GraphQL information
type RateLimitScopes struct {
	GraphQL *RateLimitDetails `json:"graphql,omitempty"`
	REST    *RateLimitDetails `json:"rest,omitempty"`
}

// RateLimitDetails describes the shape of a scope's response
type RateLimitDetails struct {
	Current  int64     `json:"current"`
	Enforced bool      `json:"enforced"`
	Limit    int64     `json:"limit"`
	Reset    int64     `json:"reset"`
	ResetAt  time.Time `json:"reset_at"`
}

// Get the rate limits for a given organization
//
// buildkite API docs: https://buildkite.com/docs/apis/rest-api/organizations/rate-limits#get-rate-limits
func (rls *RateLimitService) Get(ctx context.Context, org string) (RateLimit, *Response, error) {
	rl := fmt.Sprintf("v2/organizations/%s/rate_limit", org)
	req, err := rls.client.NewRequest(ctx, "GET", rl, nil)
	if err != nil {
		return RateLimit{}, nil, err
	}

	var rateLimit RateLimit
	resp, err := rls.client.Do(req, &rateLimit)
	if err != nil {
		return RateLimit{}, resp, err
	}

	return rateLimit, resp, err
}
