package buildkite

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	testCacheClusterUUID  = "019a1dd3-3870-7e50-8a6d-f16efba96342"
	testCacheRegistryUUID = "019a1dd8-a49b-70c5-99b8-b8f9f0ed6824"
)

func TestCacheRegistriesServiceRegistered(t *testing.T) {
	t.Parallel()

	_, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	if client.CacheRegistries == nil {
		t.Fatal("Client.CacheRegistries is nil")
	}
}

func TestCacheRegistriesService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/acme/clusters/"+testCacheClusterUUID+"/cache-registries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{
			"after":    "ruby-gems-cursor",
			"per_page": "25",
		})
		_, _ = fmt.Fprint(w, `{
			"items": [{
				"uuid": "019a1dd8-a49b-70c5-99b8-b8f9f0ed6824",
				"slug": "ruby-gems",
				"name": "Ruby gems",
				"description": "Shared Ruby dependencies",
				"emoji": ":ruby:",
				"color": "#cc342d",
				"policy": {
					"save": {"scopes": {"branch": true}},
					"restore": {"scopes": [{}]},
					"rules": [{"effect": "allow", "action": "save"}]
				},
				"created_at": "2026-09-20T01:02:03.000Z",
				"updated_at": "2026-09-21T04:05:06.000Z",
				"url": "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342/cache-registries/019a1dd8-a49b-70c5-99b8-b8f9f0ed6824",
				"cluster_url": "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342"
			}, {
				"uuid": "019a1ddb-a210-7b03-b4a6-5fcd44828b31",
				"slug": "empty-metadata",
				"name": "Empty metadata",
				"description": null,
				"emoji": null,
				"color": null,
				"policy": null
			}],
			"links": {
				"self": "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342/cache-registries?after=ruby-gems-cursor&per_page=25",
				"first": "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342/cache-registries?per_page=25",
				"prev": "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342/cache-registries?before=ruby-gems-cursor&per_page=25",
				"next": "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342/cache-registries?after=next-cursor&per_page=25"
			}
		}`)
	})

	got, _, err := client.CacheRegistries.List(context.Background(), "acme", testCacheClusterUUID, &CacheRegistriesListOptions{
		After:   "ruby-gems-cursor",
		PerPage: 25,
	})
	if err != nil {
		t.Fatalf("CacheRegistriesService.List returned error: %v", err)
	}

	description, emoji, color := "Shared Ruby dependencies", ":ruby:", "#cc342d"
	createdAt := must(time.Parse(BuildKiteDateFormat, "2026-09-20T01:02:03.000Z"))
	updatedAt := must(time.Parse(BuildKiteDateFormat, "2026-09-21T04:05:06.000Z"))
	want := CacheRegistriesList{
		Items: []CacheRegistry{
			{
				UUID:        testCacheRegistryUUID,
				Slug:        "ruby-gems",
				Name:        "Ruby gems",
				Description: &description,
				Emoji:       &emoji,
				Color:       &color,
				Policy: CacheRegistryPolicy{
					"save":    map[string]any{"scopes": map[string]any{"branch": true}},
					"restore": map[string]any{"scopes": []any{map[string]any{}}},
					"rules":   []any{map[string]any{"effect": "allow", "action": "save"}},
				},
				CreatedAt:  NewTimestamp(createdAt),
				UpdatedAt:  NewTimestamp(updatedAt),
				URL:        "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342/cache-registries/019a1dd8-a49b-70c5-99b8-b8f9f0ed6824",
				ClusterURL: "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342",
			},
			{
				UUID: "019a1ddb-a210-7b03-b4a6-5fcd44828b31",
				Slug: "empty-metadata",
				Name: "Empty metadata",
			},
		},
		Links: CacheRegistriesListLinks{
			Self:     "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342/cache-registries?after=ruby-gems-cursor&per_page=25",
			First:    "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342/cache-registries?per_page=25",
			Previous: "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342/cache-registries?before=ruby-gems-cursor&per_page=25",
			Next:     "https://api.buildkite.com/v2/organizations/acme/clusters/019a1dd3-3870-7e50-8a6d-f16efba96342/cache-registries?after=next-cursor&per_page=25",
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("CacheRegistriesService.List diff (-want +got):\n%s", diff)
	}
}

func TestCacheRegistriesListLink_ToOptions(t *testing.T) {
	t.Parallel()

	link := CacheRegistriesListLink("https://api.buildkite.com/v2/cache-registries?before=cursor-2&per_page=50")
	got, err := link.ToOptions()
	if err != nil {
		t.Fatalf("ToOptions returned error: %v", err)
	}

	want := &CacheRegistriesListOptions{Before: "cursor-2", PerPage: 50}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ToOptions diff (-want +got):\n%s", diff)
	}
}

func TestCacheRegistriesListLink_ToOptionsRejectsInvalidPerPage(t *testing.T) {
	t.Parallel()

	_, err := CacheRegistriesListLink("https://api.buildkite.com/v2/cache-registries?per_page=many").ToOptions()
	if err == nil {
		t.Fatal("ToOptions returned nil error, want invalid per_page error")
	}
}

func TestCacheRegistriesService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/acme/clusters/"+testCacheClusterUUID+"/cache-registries/"+testCacheRegistryUUID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{"uuid":%q,"slug":"ruby-gems","name":"Ruby gems","policy":{},"default":true}`, testCacheRegistryUUID)
	})

	got, _, err := client.CacheRegistries.Get(context.Background(), "acme", testCacheClusterUUID, testCacheRegistryUUID)
	if err != nil {
		t.Fatalf("CacheRegistriesService.Get returned error: %v", err)
	}
	if diff := cmp.Diff(CacheRegistry{
		UUID:    testCacheRegistryUUID,
		Slug:    "ruby-gems",
		Name:    "Ruby gems",
		Policy:  CacheRegistryPolicy{},
		Default: true,
	}, got); diff != "" {
		t.Errorf("CacheRegistriesService.Get diff (-want +got):\n%s", diff)
	}
}

func TestCacheRegistriesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	description, emoji, color := "Shared Ruby dependencies", ":ruby:", "#cc342d"
	input := CacheRegistryCreate{
		Name:        "Ruby gems",
		Description: Some(&description),
		Emoji:       Some(&emoji),
		Color:       Some(&color),
		Policy: Some(CacheRegistryPolicy{
			"save": map[string]any{"scopes": map[string]any{"branch": true}},
		}),
	}

	server.HandleFunc("/v2/organizations/acme/clusters/"+testCacheClusterUUID+"/cache-registries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{
			"name": "Ruby gems",
			"description": "Shared Ruby dependencies",
			"emoji": ":ruby:",
			"color": "#cc342d",
			"policy": {"save": {"scopes": {"branch": true}}}
		}`)
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"uuid":%q,"slug":"ruby-gems","name":"Ruby gems"}`, testCacheRegistryUUID)
	})

	got, resp, err := client.CacheRegistries.Create(context.Background(), "acme", testCacheClusterUUID, input)
	if err != nil {
		t.Fatalf("CacheRegistriesService.Create returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Create response status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got.UUID != testCacheRegistryUUID {
		t.Errorf("Create registry UUID = %q, want %q", got.UUID, testCacheRegistryUUID)
	}
}

func TestCacheRegistriesService_UpdatePatchPresence(t *testing.T) {
	t.Parallel()

	description, emoji, color := "Updated dependencies", ":package:", "#123456"
	tests := []struct {
		name  string
		input CacheRegistryUpdate
		want  string
	}{
		{
			name:  "omits unset fields",
			input: CacheRegistryUpdate{},
			want:  `{}`,
		},
		{
			name: "sends values",
			input: CacheRegistryUpdate{
				Name:        Some("Updated gems"),
				Description: Some(&description),
				Emoji:       Some(&emoji),
				Color:       Some(&color),
				Policy:      Some(CacheRegistryPolicy{"restore": map[string]any{"scopes": []any{map[string]any{}}}}),
			},
			want: `{
				"name": "Updated gems",
				"description": "Updated dependencies",
				"emoji": ":package:",
				"color": "#123456",
				"policy": {"restore": {"scopes": [{}]}}
			}`,
		},
		{
			name: "sends explicit null to clear nullable fields",
			input: CacheRegistryUpdate{
				Description: Some((*string)(nil)),
				Emoji:       Some((*string)(nil)),
				Color:       Some((*string)(nil)),
				Policy:      Some(CacheRegistryPolicy(nil)),
			},
			want: `{"description":null,"emoji":null,"color":null,"policy":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			server.HandleFunc("/v2/organizations/acme/clusters/"+testCacheClusterUUID+"/cache-registries/"+testCacheRegistryUUID, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPatch)
				assertRequestJSON(t, r, tt.want)
				_, _ = fmt.Fprintf(w, `{"uuid":%q,"slug":"updated-gems","name":"Updated gems"}`, testCacheRegistryUUID)
			})

			got, _, err := client.CacheRegistries.Update(context.Background(), "acme", testCacheClusterUUID, testCacheRegistryUUID, tt.input)
			if err != nil {
				t.Fatalf("CacheRegistriesService.Update returned error: %v", err)
			}
			if got.UUID != testCacheRegistryUUID {
				t.Errorf("Update registry UUID = %q, want %q", got.UUID, testCacheRegistryUUID)
			}
		})
	}
}

func TestCacheRegistriesService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/acme/clusters/"+testCacheClusterUUID+"/cache-registries/"+testCacheRegistryUUID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.CacheRegistries.Delete(context.Background(), "acme", testCacheClusterUUID, testCacheRegistryUUID)
	if err != nil {
		t.Fatalf("CacheRegistriesService.Delete returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Delete response status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestCacheRegistriesService_APIErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{name: "permission denied", statusCode: http.StatusForbidden, message: "Forbidden"},
		{name: "default registry cannot be deleted", statusCode: http.StatusUnprocessableEntity, message: "Cannot destroy default cache registry"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/clusters/"+testCacheClusterUUID+"/cache-registries/"+testCacheRegistryUUID, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = fmt.Fprintf(w, `{"message":%q}`, tt.message)
			})

			resp, err := client.CacheRegistries.Delete(context.Background(), "acme", testCacheClusterUUID, testCacheRegistryUUID)
			if err == nil {
				t.Fatal("CacheRegistriesService.Delete returned nil error, want API error")
			}
			if resp == nil || resp.StatusCode != tt.statusCode {
				t.Fatalf("Delete response = %#v, want status %d", resp, tt.statusCode)
			}

			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) {
				t.Fatalf("Delete error type = %T, want *ErrorResponse", err)
			}
			if apiErr.Message != tt.message {
				t.Errorf("API error message = %q, want %q", apiErr.Message, tt.message)
			}
		})
	}
}
