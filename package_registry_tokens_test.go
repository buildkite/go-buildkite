package buildkite

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var (
	registryToken = PackageRegistryToken{
		ID:          "0191b6a2-aa51-70d0-8a5f-aabce115b0fd",
		GraphQLID:   "UmVnaXN0cnlUb2tlbi0tMDE5MWI2YTItYWE1MS03MGQwLThhNWYtYWFiY2UxMTViMGZk",
		Description: "CI deploy token",
		URL:         "https://api.buildkite-test.com/v2/packages/organizations/test-org/registries/my-cool-registry/tokens/0191b6a2-aa51-70d0-8a5f-aabce115b0fd",
		CreatedAt:   NewTimestamp(time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)),
		CreatedBy: PackageRegistryTokenCreator{
			ID:        "0191b6a2-aa51-70d0-8a5f-cccccccccccc",
			GraphQLID: "VXNlci0tMDE5MWI2YTItYWE1MS03MGQwLThhNWYtY2NjY2NjY2NjY2Nj",
			Name:      "Tom Watt",
			Email:     "tom@buildkite.com",
			AvatarURL: "https://www.gravatar.com/avatar/abc123",
			CreatedAt: NewTimestamp(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		Organization: PackageRegistryTokenOrg{
			ID:     "0191b6a2-aa51-70d0-8a5f-dddddddddddd",
			Slug:   "test-org",
			URL:    "https://api.buildkite-test.com/v2/organizations/test-org",
			WebURL: "https://buildkite-test.com/test-org",
		},
		Registry: PackageRegistryTokenRegistry{
			ID:        "0190ddd4-f249-7c47-b129-f5c4207ed0fd",
			GraphQLID: "UmVnaXN0cnktLTAxOTBkZGQ0LWYyNDktN2M0Ny1iMTI5LWFhYWFhYWFhYWFhYQ==",
			Slug:      "my-cool-registry",
			URL:       "https://api.buildkite-test.com/v2/packages/organizations/test-org/registries/my-cool-registry",
			WebURL:    "https://buildkite-test.com/organizations/test-org/packages/registries/my-cool-registry",
		},
	}

	registryTokens = []PackageRegistryToken{
		registryToken,
		{
			ID:           "0191b6a2-aa51-70d0-8a5f-eeeeeeeeeeee",
			GraphQLID:    "UmVnaXN0cnlUb2tlbi0tMDE5MWI2YTItYWE1MS03MGQwLThhNWYtZWVlZWVlZWVlZWVl",
			Description:  "Publish token",
			URL:          "https://api.buildkite-test.com/v2/packages/organizations/test-org/registries/my-cool-registry/tokens/0191b6a2-aa51-70d0-8a5f-eeeeeeeeeeee",
			CreatedAt:    NewTimestamp(time.Date(2026, 7, 2, 9, 30, 0, 0, time.UTC)),
			CreatedBy:    registryToken.CreatedBy,
			Organization: registryToken.Organization,
			Registry:     registryToken.Registry,
		},
	}
)

func TestPackageRegistryTokenGet(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	want := registryToken
	server.HandleFunc("/v2/packages/organizations/test-org/registries/my-cool-registry/tokens/0191b6a2-aa51-70d0-8a5f-aabce115b0fd", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		err := json.NewEncoder(w).Encode(want)
		if err != nil {
			t.Fatalf("encoding json response body: %v", err)
		}
	})

	got, _, err := client.PackageRegistryTokensService.Get(context.Background(), "test-org", "my-cool-registry", "0191b6a2-aa51-70d0-8a5f-aabce115b0fd")
	if err != nil {
		t.Fatalf("PackageRegistryTokens.Get returned error: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("client.PackageRegistryTokensService.Get(...) diff: (-got +want)\n%s", diff)
	}
}

func TestPackageRegistryTokenList(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	want := registryTokens
	server.HandleFunc("/v2/packages/organizations/test-org/registries/my-cool-registry/tokens", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		err := json.NewEncoder(w).Encode(want)
		if err != nil {
			t.Fatalf("encoding json response body: %v", err)
		}
	})

	got, _, err := client.PackageRegistryTokensService.List(context.Background(), "test-org", "my-cool-registry")
	if err != nil {
		t.Fatalf("PackageRegistryTokens.List returned error: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("client.PackageRegistryTokensService.List(...) diff: (-got +want)\n%s", diff)
	}
}
