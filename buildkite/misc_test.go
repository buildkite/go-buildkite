package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestListEmojis(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/emojis", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"name":"rocket","url":"https://a.buildboxassets.com/assets/emoji2/unicode/1f680.png?v2"}]`)
	})

	emoji, _, err := client.ListEmojis(context.Background(), "my-great-org")
	if err != nil {
		t.Errorf("ListEmojis returned error: %v", err)
	}

	want := []Emoji{{Name: String("rocket"), URL: String("https://a.buildboxassets.com/assets/emoji2/unicode/1f680.png?v2")}}
	if diff := cmp.Diff(want, emoji); diff != "" {
		t.Errorf("ListEmojis diff: (-got +want)\n%s", diff)
	}
}
