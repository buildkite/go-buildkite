package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

var pkg = Package{
	ID:       uuid.NewString(),
	URL:      "https://example.com/my-package",
	WebURL:   "https://buildkite.com/my-org/my-registry/my-package",
	Registry: registry,
	Organization: Organization{
		ID:   "my-org",
		Slug: "my-org",
		Name: "My Org",
	},
}

func TestGetPackage(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	want := pkg
	endpoint := fmt.Sprintf("/v2/packages/organizations/my-org/registries/my-registry/packages/%s", pkg.ID)
	server.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		err := json.NewEncoder(w).Encode(pkg)
		if err != nil {
			t.Fatalf("marshalling package to json: %v", err)
		}
	})

	p, _, err := client.PackagesService.Get(context.Background(), "my-org", "my-registry", pkg.ID)
	if err != nil {
		t.Fatalf("Packages.Get returned error: %v", err)
	}

	if diff := cmp.Diff(p, want); diff != "" {
		t.Fatalf("client.PackagesService.Get(context.Background(),%q, %q, %q) diff: (-got +want)\n%s", "test-org", "my-cool-registry", pkg.ID, diff)
	}
}

func TestCopyPackage(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	want := pkg
	endpoint := fmt.Sprintf("/v2/packages/organizations/my-org/registries/source-registry/packages/%s/copy", pkg.ID)
	server.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"to": "dest-registry"})

		err := json.NewEncoder(w).Encode(pkg)
		if err != nil {
			t.Fatalf("marshalling package to json: %v", err)
		}
	})

	p, _, err := client.PackagesService.Copy(context.Background(), "my-org", "source-registry", pkg.ID, "dest-registry")
	if err != nil {
		t.Fatalf("Packages.Copy returned error: %v", err)
	}

	if diff := cmp.Diff(p, want); diff != "" {
		t.Fatalf("client.PackagesService.Copy(context.Background(), %q, %q, %q, %q) diff: (-got +want)\n%s", "my-org", "source-registry", pkg.ID, "dest-registry", diff)
	}
}
