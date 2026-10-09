package buildkite

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const testNotificationServiceUUID = "9f0f9a19-1b88-4e37-98d8-c5a0cebcdb9a"

const testNotificationServiceJSON = `{
	"id": "9f0f9a19-1b88-4e37-98d8-c5a0cebcdb9a",
	"graphql_id": "Tm90aWZpY2F0aW9uU2VydmljZVdlYmhvb2stLS05ZjBmOWExOS0xYjg4LTRlMzctOThkOC1jNWEwY2ViY2RiOWE=",
	"url": "https://api.buildkite.com/v2/organizations/acme/services/9f0f9a19-1b88-4e37-98d8-c5a0cebcdb9a",
	"provider": {"id": "webhook", "name": "Webhook"},
	"description": "Deploy notifications",
	"enabled": %t,
	"scope": "some_clusters",
	"scope_uuids": ["019a1dd3-3870-7e50-8a6d-f16efba96342"],
	"branch_configuration": "main",
	"build_states": {"build_passed": false, "build_failed": true, "job_activated": false},
	"settings": {
		"url": "https://example.com/buildkite-webhook",
		"version": 3,
		"token": "xxx-yyy-zzz",
		"token_mode": "token",
		"events": ["build.finished"],
		"tls_verify": true,
		"cluster_queue_ids": []
	},
	"created_at": "2026-07-01T10:00:00.000Z",
	"created_by": {
		"id": "3d3c3bf0-7d58-4afe-8fe7-b3017d5504de",
		"graphql_id": "VXNlci0tLTNkM2MzYmYwLTdkNTgtNGFmZS04ZmU3LWIzMDE3ZDU1MDRkZQ==",
		"name": "Sam Kim",
		"email": "sam@example.com",
		"avatar_url": "https://www.gravatar.com/avatar/example",
		"created_at": "2025-01-01T00:00:00.000Z"
	}
}`

func testNotificationService(enabled bool) NotificationService {
	graphQLID := "Tm90aWZpY2F0aW9uU2VydmljZVdlYmhvb2stLS05ZjBmOWExOS0xYjg4LTRlMzctOThkOC1jNWEwY2ViY2RiOWE="
	return NotificationService{
		ID:                  testNotificationServiceUUID,
		GraphQLID:           &graphQLID,
		URL:                 "https://api.buildkite.com/v2/organizations/acme/services/9f0f9a19-1b88-4e37-98d8-c5a0cebcdb9a",
		Provider:            NotificationServiceProvider{ID: "webhook", Name: "Webhook"},
		Description:         "Deploy notifications",
		Enabled:             enabled,
		Scope:               "some_clusters",
		ScopeUUIDs:          []string{"019a1dd3-3870-7e50-8a6d-f16efba96342"},
		BranchConfiguration: "main",
		BuildStates:         map[string]bool{"build_passed": false, "build_failed": true, "job_activated": false},
		Settings: NotificationServiceSettings{
			"url":               "https://example.com/buildkite-webhook",
			"version":           float64(3),
			"token":             "xxx-yyy-zzz",
			"token_mode":        "token",
			"events":            []any{"build.finished"},
			"tls_verify":        true,
			"cluster_queue_ids": []any{},
		},
		CreatedAt: NewTimestamp(must(time.Parse(BuildKiteDateFormat, "2026-07-01T10:00:00.000Z"))),
		CreatedBy: &User{
			ID:        "3d3c3bf0-7d58-4afe-8fe7-b3017d5504de",
			Name:      "Sam Kim",
			Email:     "sam@example.com",
			CreatedAt: NewTimestamp(must(time.Parse(BuildKiteDateFormat, "2025-01-01T00:00:00.000Z"))),
		},
	}
}

func TestNotificationServicesService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/acme/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{
			"after":    "cursor-1",
			"per_page": "30",
		})
		_, _ = fmt.Fprintf(w, `{
			"items": [`+testNotificationServiceJSON+`, {
				"id": "019a1ddb-a210-7b03-b4a6-5fcd44828b31",
				"graphql_id": null,
				"provider": {"id": "datadog_pipeline_visibility", "name": "Datadog Pipeline Visibility"},
				"enabled": true,
				"scope": "all",
				"scope_uuids": [],
				"branch_configuration": "",
				"settings": {"api_key": "XXXXXXXXXXXXcdef", "datadog_site": "datadoghq.com", "datadog_tags": null},
				"created_by": null
			}],
			"links": {
				"self": "https://api.buildkite.com/v2/organizations/acme/services?after=cursor-1&per_page=30",
				"next": "https://api.buildkite.com/v2/organizations/acme/services?after=cursor-2&per_page=30"
			}
		}`, true)
	})

	got, _, err := client.NotificationServices.List(context.Background(), "acme", &NotificationServicesListOptions{
		After:   "cursor-1",
		PerPage: 30,
	})
	if err != nil {
		t.Fatalf("NotificationServicesService.List returned error: %v", err)
	}

	want := NotificationServicesList{
		Items: []NotificationService{
			testNotificationService(true),
			{
				ID:         "019a1ddb-a210-7b03-b4a6-5fcd44828b31",
				Provider:   NotificationServiceProvider{ID: "datadog_pipeline_visibility", Name: "Datadog Pipeline Visibility"},
				Enabled:    true,
				Scope:      "all",
				ScopeUUIDs: []string{},
				Settings: NotificationServiceSettings{
					"api_key":      "XXXXXXXXXXXXcdef",
					"datadog_site": "datadoghq.com",
					"datadog_tags": nil,
				},
			},
		},
		Links: NotificationServicesListLinks{
			Self: "https://api.buildkite.com/v2/organizations/acme/services?after=cursor-1&per_page=30",
			Next: "https://api.buildkite.com/v2/organizations/acme/services?after=cursor-2&per_page=30",
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NotificationServicesService.List diff (-want +got):\n%s", diff)
	}

	next, err := got.Links.Next.ToOptions()
	if err != nil {
		t.Fatalf("ToOptions returned error: %v", err)
	}
	if diff := cmp.Diff(&NotificationServicesListOptions{After: "cursor-2", PerPage: 30}, next); diff != "" {
		t.Errorf("ToOptions diff (-want +got):\n%s", diff)
	}
}

func TestNotificationServicesService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/acme/services/"+testNotificationServiceUUID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, testNotificationServiceJSON, true)
	})

	got, _, err := client.NotificationServices.Get(context.Background(), "acme", testNotificationServiceUUID)
	if err != nil {
		t.Fatalf("NotificationServicesService.Get returned error: %v", err)
	}
	if diff := cmp.Diff(testNotificationService(true), got); diff != "" {
		t.Errorf("NotificationServicesService.Get diff (-want +got):\n%s", diff)
	}
}

func TestNotificationServicesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/acme/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{
			"provider": "webhook",
			"description": "Deploy notifications",
			"branch_configuration": "main",
			"scope": "some_clusters",
			"scope_uuids": ["019a1dd3-3870-7e50-8a6d-f16efba96342"],
			"build_states": {"build_failed": true},
			"settings": {
				"url": "https://example.com/buildkite-webhook",
				"token": "xxx-yyy-zzz",
				"token_mode": "signature",
				"events": ["build.finished"]
			}
		}`)
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, testNotificationServiceJSON, true)
	})

	got, resp, err := client.NotificationServices.Create(context.Background(), "acme", NotificationServiceCreate{
		Provider:            "webhook",
		Description:         "Deploy notifications",
		BranchConfiguration: "main",
		Scope:               "some_clusters",
		ScopeUUIDs:          []string{"019a1dd3-3870-7e50-8a6d-f16efba96342"},
		BuildStates:         map[string]bool{"build_failed": true},
		Settings: NotificationServiceSettings{
			"url":        "https://example.com/buildkite-webhook",
			"token":      "xxx-yyy-zzz",
			"token_mode": "signature",
			"events":     []string{"build.finished"},
		},
	})
	if err != nil {
		t.Fatalf("NotificationServicesService.Create returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Create response status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if diff := cmp.Diff(testNotificationService(true), got); diff != "" {
		t.Errorf("NotificationServicesService.Create diff (-want +got):\n%s", diff)
	}
}

func TestNotificationServicesService_UpdatePatchPresence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input NotificationServiceUpdate
		want  string
	}{
		{
			name:  "omits unset fields",
			input: NotificationServiceUpdate{},
			want:  `{}`,
		},
		{
			name: "sends values",
			input: NotificationServiceUpdate{
				Provider:            Some("webhook"),
				Description:         Some("Production deploy notifications"),
				BranchConfiguration: Some("main release/*"),
				Scope:               Some("some_projects"),
				ScopeUUIDs:          Some([]string{"16f3b56f-4934-4546-923c-287859851332"}),
				BuildStates:         Some(map[string]bool{"build_failed": true, "build_fixed": true}),
				Settings:            Some(NotificationServiceSettings{"events": []string{"build.finished", "job.finished"}}),
			},
			want: `{
				"provider": "webhook",
				"description": "Production deploy notifications",
				"branch_configuration": "main release/*",
				"scope": "some_projects",
				"scope_uuids": ["16f3b56f-4934-4546-923c-287859851332"],
				"build_states": {"build_failed": true, "build_fixed": true},
				"settings": {"events": ["build.finished", "job.finished"]}
			}`,
		},
		{
			name: "sends explicit empty values",
			input: NotificationServiceUpdate{
				Description:         Some(""),
				BranchConfiguration: Some(""),
				Scope:               Some("all"),
				ScopeUUIDs:          Some([]string{}),
			},
			want: `{"description":"","branch_configuration":"","scope":"all","scope_uuids":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			server.HandleFunc("/v2/organizations/acme/services/"+testNotificationServiceUUID, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPatch)
				assertRequestJSON(t, r, tt.want)
				_, _ = fmt.Fprintf(w, testNotificationServiceJSON, true)
			})

			got, _, err := client.NotificationServices.Update(context.Background(), "acme", testNotificationServiceUUID, tt.input)
			if err != nil {
				t.Fatalf("NotificationServicesService.Update returned error: %v", err)
			}
			if got.ID != testNotificationServiceUUID {
				t.Errorf("Update service ID = %q, want %q", got.ID, testNotificationServiceUUID)
			}
		})
	}
}

func TestNotificationServicesService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/acme/services/"+testNotificationServiceUUID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.NotificationServices.Delete(context.Background(), "acme", testNotificationServiceUUID)
	if err != nil {
		t.Fatalf("NotificationServicesService.Delete returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Delete response status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestNotificationServicesService_EnableDisable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		action  string
		enabled bool
		call    func(*Client) (NotificationService, *Response, error)
	}{
		{
			action:  "enable",
			enabled: true,
			call: func(c *Client) (NotificationService, *Response, error) {
				return c.NotificationServices.Enable(context.Background(), "acme", testNotificationServiceUUID)
			},
		},
		{
			action:  "disable",
			enabled: false,
			call: func(c *Client) (NotificationService, *Response, error) {
				return c.NotificationServices.Disable(context.Background(), "acme", testNotificationServiceUUID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			t.Parallel()

			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			server.HandleFunc("/v2/organizations/acme/services/"+testNotificationServiceUUID+"/"+tt.action, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPut)
				if body := must(io.ReadAll(r.Body)); len(body) != 0 {
					t.Errorf("request body = %q, want empty", body)
				}
				_, _ = fmt.Fprintf(w, testNotificationServiceJSON, tt.enabled)
			})

			got, _, err := tt.call(client)
			if err != nil {
				t.Fatalf("NotificationServicesService %s returned error: %v", tt.action, err)
			}
			if diff := cmp.Diff(testNotificationService(tt.enabled), got); diff != "" {
				t.Errorf("NotificationServicesService %s diff (-want +got):\n%s", tt.action, diff)
			}
		})
	}
}
