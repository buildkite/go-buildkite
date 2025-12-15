package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRateLimitService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/mock-kite/rate_limit", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{
	"scopes": {
		"rest": {
			"limit": 200,
			"current": 10,
			"reset": 30,
			"reset_at": "2025-12-15T05:40:00Z",
			"enforced": true
		},
		"graphql": {
			"limit": 20000,
			"current": 10000,
			"reset": 263,
			"reset_at": "2025-12-15T05:44:33Z",
			"enforced": true
		}
	}
}
`)
	})

	limit, _, err := client.RateLimit.Get(context.Background(), "mock-kite")
	if err != nil {
		t.Errorf("RateLimit.Get returned error: %v", err)
	}

	want := RateLimit{&RateLimitScopes{
		REST: &RateLimitDetails{
			Limit:    200,
			Current:  10,
			Reset:    30,
			ResetAt:  time.Date(2025, time.December, 15, 5, 40, 0, 0, time.UTC),
			Enforced: true,
		},
		GraphQL: &RateLimitDetails{
			Limit:    20000,
			Current:  10000,
			Reset:    263,
			ResetAt:  time.Date(2025, time.December, 15, 5, 44, 33, 0, time.UTC),
			Enforced: true,
		},
	}}
	if diff := cmp.Diff(limit, want); diff != "" {
		t.Errorf("RateLImit.Get diff: (-got +want)\n%s", diff)
	}
}
