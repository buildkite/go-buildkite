package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEmojisService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/emojis", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[
			{"name":"rocket","url":"https://a.buildboxassets.com/assets/emoji2/unicode/1f680.png?v2"},
			{"name":"shipit","url":"https://a.buildboxassets.com/assets/emoji2/shipit.png?v2","aliases":["squirrel"]}
		]`)
	})

	emoji, _, err := client.Emojis.List(context.Background(), "my-great-org")
	if err != nil {
		t.Errorf("Emojis.List returned error: %v", err)
	}

	want := []Emoji{
		{Name: "rocket", URL: "https://a.buildboxassets.com/assets/emoji2/unicode/1f680.png?v2"},
		{Name: "shipit", URL: "https://a.buildboxassets.com/assets/emoji2/shipit.png?v2", Aliases: []string{"squirrel"}},
	}
	if diff := cmp.Diff(emoji, want); diff != "" {
		t.Errorf("Emojis.List diff: (-got +want): \n%s", diff)
	}
}

// TestClient_ListEmojis covers the deprecated Client.ListEmojis shim, which
// should delegate to Emojis.List.
func TestClient_ListEmojis(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/emojis", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{"name":"rocket","url":"https://a.buildboxassets.com/assets/emoji2/unicode/1f680.png?v2"}]`)
	})

	emoji, _, err := client.ListEmojis(context.Background(), "my-great-org")
	if err != nil {
		t.Errorf("ListEmojis returned error: %v", err)
	}

	want := []Emoji{{Name: "rocket", URL: "https://a.buildboxassets.com/assets/emoji2/unicode/1f680.png?v2"}}
	if diff := cmp.Diff(emoji, want); diff != "" {
		t.Errorf("ListEmojis diff: (-got +want): \n%s", diff)
	}
}
