package buildkite

import (
	"context"
	"fmt"
)

// EmojisService handles communication with the emoji related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api/emojis
type EmojisService struct {
	client *Client
}

type Emoji struct {
	Name    string   `json:"name,omitempty"`
	URL     string   `json:"url,omitempty"`
	Aliases []string `json:"aliases,omitempty"`
}

// List all the emojis for a given account, including custom emojis and aliases.
//
// buildkite API docs: https://buildkite.com/docs/api/emojis
func (es *EmojisService) List(ctx context.Context, org string) ([]Emoji, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/emojis", org)
	req, err := es.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var emojis []Emoji
	resp, err := es.client.Do(req, &emojis)
	if err != nil {
		return nil, resp, err
	}

	return emojis, resp, nil
}
