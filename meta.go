package buildkite

import (
	"context"
	"net/http"
)

// MetaService handles communication with the meta related
// methods of the Buildkite API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/meta
type MetaService struct {
	client *Client
}

// Meta represents the meta values returned by Buildkite
type Meta struct {
	WebhookIPs []string `json:"webhook_ips,omitempty"`
}

// Get returns the meta values from the meta endpoint
func (ms *MetaService) Get(ctx context.Context) (Meta, *Response, error) {
	u := "v2/meta"
	req, err := ms.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return Meta{}, nil, err
	}

	var meta Meta
	resp, err := ms.client.Do(req, &meta)
	if err != nil {
		return Meta{}, resp, err
	}

	return meta, resp, err
}
