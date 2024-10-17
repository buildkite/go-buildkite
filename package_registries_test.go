package buildkite

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	registry = PackageRegistry{
		Name:      "My Cool Registry",
		ID:        "0190ddd4-f249-7c47-b129-f5c4207ed0fd",
		GraphQLID: "UmVnaXN0cnktLTAxOTBkZGQ0LWYyNDktN2M0Ny1iMTI5LWFhYWFhYWFhYWFhYQ==",
		Slug:      "my-cool-registry",
		URL:       "https://api.buildkite-test.com/v2/packages/organizations/test-org/registries/my-cool-registry",
		WebURL:    "https://buildkite-test.com/organizations/test-org/packages/registries/my-cool-registry",
		Ecosystem: "ruby",
	}

	registries = []PackageRegistry{
		{
			Name:      "My Cool Registry",
			ID:        "0190ddd4-f249-7c47-b129-aaaaaaaaaaaa",
			GraphQLID: "UmVnaXN0cnktLTAxOTBkZGQ0LWYyNDktN2M0Ny1iMTI5LWFhYWFhYWFhYWFhYQ==",
			Slug:      "my-cool-registry",
			URL:       "https://api.buildkite-test.com/v2/packages/organizations/test-org/registries/my-cool-registry",
			WebURL:    "https://buildkite-test.com/organizations/test-org/packages/registries/my-cool-registry",
			Ecosystem: "ruby",
		},
		{
			Name:      "Cool Container",
			ID:        "0190718a-6081-74db-8c21-bbbbbbbbbbbb",
			GraphQLID: "UmVnaXN0cnktLTAxOTA3MThhLTYwODEtNzRkYi04YzIxLWJiYmJiYmJiYmJiYg==",
			Slug:      "cool-container",
			URL:       "https://api.buildkite-test.com/v2/packages/organizations/test-org/registries/cool-container",
			WebURL:    "https://buildkite-test.com/organizations/test-org/packages/registries/cool-container",
			Ecosystem: "container",
		},
	}
)

func TestPackageRegistryGet(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	want := registry
	server.HandleFunc("/v2/packages/organizations/test-org/registries/my-cool-registry", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		err := json.NewEncoder(w).Encode(want)
		if err != nil {
			t.Fatalf("encoding json response body: %v", err)
		}
	})

	got, _, err := client.PackageRegistriesService.Get(context.Background(), "test-org", "my-cool-registry")
	if err != nil {
		t.Fatalf("PackageRegistries.Get returned error: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("client.PackageRegistriesService.Get(context.Background(),%q, %q) diff: (-got +want)\n%s", "test-org", "my-cool-registry", diff)
	}
}

func TestPackageRegistryList(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	want := registries
	server.HandleFunc("/v2/packages/organizations/test-org/registries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		err := json.NewEncoder(w).Encode(want)
		if err != nil {
			t.Fatalf("encoding json response body: %v", err)
		}
	})

	got, _, err := client.PackageRegistriesService.List(context.Background(), "test-org")
	if err != nil {
		t.Fatalf("PackageRegistries.List returned error: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("client.PackageRegistriesService.List(context.Background(),%q) diff: (-got +want)\n%s", "test-org", diff)
	}
}

func TestPackageRegistryCreate(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	wantInput := CreatePackageRegistryInput{
		Name:        "My Cool Registry",
		Description: "A registry for all the cool things",
		Ecosystem:   "ruby",
	}

	want := registry
	server.HandleFunc("/v2/packages/organizations/test-org/registries", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		testMethod(t, r, "POST")

		var gotInput CreatePackageRegistryInput
		err := json.NewDecoder(r.Body).Decode(&gotInput)
		if err != nil {
			t.Fatalf("parsing Create Package body: %v", err)
		}

		// ensure that the input survives a roundtrip through JSON
		if diff := cmp.Diff(gotInput, wantInput); diff != "" {
			t.Fatalf("create registry input diff: (-got +want)\n%s", diff)
		}

		// send back the registry
		err = json.NewEncoder(w).Encode(want)
		if err != nil {
			t.Fatalf("encoding json response body: %v", err)
		}
	})

	got, _, err := client.PackageRegistriesService.Create(context.Background(), "test-org", wantInput)
	if err != nil {
		t.Fatalf("PackageRegistries.Create returned error: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("client.PackageRegistriesService.Create(context.Background(),%q, %#v) diff: (-got +want)\n%s", "test-org", wantInput, diff)
	}
}

func TestPackageRegistryUpdate(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	wantInput := UpdatePackageRegistryInput{
		Name:        "My Even Cooler Registry",
		Description: "A registry so cool we had to rename it",
	}

	want := registry
	want.Name = wantInput.Name
	want.Description = wantInput.Description

	server.HandleFunc("/v2/packages/organizations/test-org/registries/my-cool-registry", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		testMethod(t, r, "PATCH")

		var gotInput UpdatePackageRegistryInput
		err := json.NewDecoder(r.Body).Decode(&gotInput)
		if err != nil {
			t.Fatalf("parsing Update Package body: %v", err)
		}

		// ensure that the input survives a roundtrip through JSON
		if diff := cmp.Diff(gotInput, wantInput); diff != "" {
			t.Fatalf("update registry input diff: (-got +want)\n%s", diff)
		}

		// send back the registry
		err = json.NewEncoder(w).Encode(want)
		if err != nil {
			t.Fatalf("encoding json response body: %v", err)
		}
	})

	got, _, err := client.PackageRegistriesService.Update(context.Background(), "test-org", "my-cool-registry", wantInput)
	if err != nil {
		t.Fatalf("PackageRegistries.Update returned error: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("client.PackageRegistriesService.Update(context.Background(),%q, %q, %#v) diff: (-got +want)\n%s", "test-org", "my-cool-registry", wantInput, diff)
	}
}

func TestPackageRegistryDelete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/packages/organizations/test-org/registries/my-cool-registry", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.PackageRegistriesService.Delete(context.Background(), "test-org", "my-cool-registry")
	if err != nil {
		t.Fatalf("PackageRegistries.Delete returned error: %v", err)
	}

	if got, want := resp.StatusCode, http.StatusNoContent; got != want {
		t.Fatalf("client.PackageRegistriesService.Delete(context.Background(),%q, %q) status: %d, want %d", "test-org", "my-cool-registry", got, want)
	}
}
