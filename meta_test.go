package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMetaService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/meta", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{"webhook_ips":["192.0.2.0/24", "198.51.100.12"]}`)
	})

	meta, _, err := client.Meta.Get(context.Background())
	if err != nil {
		t.Errorf("Meta.Get returned error: %v", err)
	}

	want := Meta{WebhookIPs: []string{"192.0.2.0/24", "198.51.100.12"}}
	if diff := cmp.Diff(meta, want); diff != "" {
		t.Errorf("Meta.Get diff: (-got +want)\n%s", diff)
	}
}
