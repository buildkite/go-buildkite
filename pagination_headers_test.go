package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestOrganizationsListMultipleLinkHeaders(t *testing.T) {
	for _, tc := range []struct {
		name    string
		headers []string
	}{
		{"separate", []string{`<https://api.buildkite.com/v2/organizations?page=1>; rel="first"`, `<https://api.buildkite.com/v2/organizations?page=2>; rel="prev"`, `<https://api.buildkite.com/v2/organizations?page=4>; rel="next"`, `<https://api.buildkite.com/v2/organizations?page=5>; rel="last"`}},
		{"mixed", []string{`<https://api.buildkite.com/v2/organizations?page=1>; rel="first", <https://api.buildkite.com/v2/organizations?page=2>; rel="prev"`, `<https://api.buildkite.com/v2/organizations?page=4>; rel="next", <https://api.buildkite.com/v2/organizations?page=5>; rel="last"`}},
		{"combined", []string{`<https://api.buildkite.com/v2/organizations?page=1>; rel="first", <https://api.buildkite.com/v2/organizations?page=2>; rel="prev", <https://api.buildkite.com/v2/organizations?page=4>; rel="next", <https://api.buildkite.com/v2/organizations?page=5>; rel="last"`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server, client, closeServer := newMockServerAndClient(t)
			defer closeServer()
			server.HandleFunc("/v2/organizations", func(w http.ResponseWriter, r *http.Request) {
				for _, link := range tc.headers {
					w.Header().Add("Link", link)
				}
				fmt.Fprint(w, `[]`)
			})
			_, response, err := client.Organizations.List(context.Background(), nil)
			if err != nil {
				t.Fatal(err)
			}
			if response.FirstPage != 1 || response.PrevPage != 2 || response.NextPage != 4 || response.LastPage != 5 {
				t.Errorf("pagination = first:%d prev:%d next:%d last:%d, want 1 2 4 5", response.FirstPage, response.PrevPage, response.NextPage, response.LastPage)
			}
		})
	}
}
