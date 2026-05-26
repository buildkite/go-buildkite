package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func boolPtr(b bool) *bool { return &b }

func TestPipelineSchedulesService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/pipelines/my-pipeline/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[
			{
				"id": "1bb821f4-3aaf-4d7a-bcb8-26a7fec6f7e2",
				"graphql_id": "UGlwZWxpbmVTY2hlZHVsZS0tLTFiYjgyMWY0",
				"url": "https://api.buildkite.com/v2/organizations/my-org/pipelines/my-pipeline/schedules/1bb821f4-3aaf-4d7a-bcb8-26a7fec6f7e2",
				"label": "Nightly build",
				"cronline": "@daily",
				"branch": "main",
				"commit": "HEAD",
				"message": "Nightly build",
				"env": {"FOO": "bar"},
				"enabled": true,
				"next_build_at": "2026-05-07T00:00:00Z",
				"failed_message": null,
				"failed_at": null,
				"created_at": "2026-05-01T00:00:00Z",
				"created_by": {
					"id": "user-1",
					"name": "Alice",
					"email": "alice@example.com",
					"avatar_url": "https://example.com/avatar.png",
					"created_at": "2024-01-01T00:00:00Z"
				}
			}
		]`)
	})

	got, _, err := client.PipelineSchedules.List(context.Background(), "my-org", "my-pipeline", nil)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("List returned %d schedules, want 1", len(got))
	}
	if got[0].ID != "1bb821f4-3aaf-4d7a-bcb8-26a7fec6f7e2" {
		t.Errorf("schedule id = %q, want %q", got[0].ID, "1bb821f4-3aaf-4d7a-bcb8-26a7fec6f7e2")
	}
	if got[0].Cronline != "@daily" {
		t.Errorf("schedule cronline = %q, want %q", got[0].Cronline, "@daily")
	}
	if !got[0].Enabled {
		t.Errorf("schedule enabled = false, want true")
	}
	if diff := cmp.Diff(map[string]string{"FOO": "bar"}, got[0].Env); diff != "" {
		t.Errorf("schedule env mismatch (-want +got):\n%s", diff)
	}
	if got[0].CreatedBy == nil || got[0].CreatedBy.Name != "Alice" {
		t.Errorf("schedule created_by = %+v, want name=Alice", got[0].CreatedBy)
	}
}

func TestPipelineSchedulesService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/pipelines/my-pipeline/schedules/1bb821f4-3aaf-4d7a-bcb8-26a7fec6f7e2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{
			"id": "1bb821f4-3aaf-4d7a-bcb8-26a7fec6f7e2",
			"label": "Nightly build",
			"cronline": "@daily",
			"branch": "main",
			"enabled": true
		}`)
	})

	got, _, err := client.PipelineSchedules.Get(context.Background(), "my-org", "my-pipeline", "1bb821f4-3aaf-4d7a-bcb8-26a7fec6f7e2")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	want := PipelineSchedule{
		ID:       "1bb821f4-3aaf-4d7a-bcb8-26a7fec6f7e2",
		Label:    "Nightly build",
		Cronline: "@daily",
		Branch:   "main",
		Enabled:  true,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Get mismatch (-want +got):\n%s", diff)
	}
}

func TestPipelineSchedulesService_Create(t *testing.T) {
	t.Parallel()

	in := CreatePipelineSchedule{
		Cronline: "@daily",
		Label:    "Nightly build",
		Message:  "Nightly build",
		Branch:   "main",
		Commit:   "HEAD",
		Env:      map[string]string{"FOO": "bar"},
		Enabled:  boolPtr(false),
	}

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/pipelines/my-pipeline/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var got CreatePipelineSchedule
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		if diff := cmp.Diff(in, got); diff != "" {
			t.Errorf("Create request body mismatch (-want +got):\n%s", diff)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, `{
			"id": "new-id",
			"cronline": "@daily",
			"label": "Nightly build",
			"branch": "main",
			"commit": "HEAD",
			"env": {"FOO": "bar"},
			"enabled": false
		}`)
	})

	got, _, err := client.PipelineSchedules.Create(context.Background(), "my-org", "my-pipeline", in)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if got.ID != "new-id" {
		t.Errorf("Create id = %q, want %q", got.ID, "new-id")
	}
	if got.Enabled {
		t.Errorf("Create enabled = true, want false")
	}
}

func TestPipelineSchedulesService_Update(t *testing.T) {
	t.Parallel()

	in := UpdatePipelineSchedule{
		Label:   Some("Updated label"),
		Enabled: Some(true),
	}

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/pipelines/my-pipeline/schedules/abc", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		assertRequestJSON(t, r, `{
			"label": "Updated label",
			"enabled": true
		}`)

		_, _ = fmt.Fprint(w, `{
			"id": "abc",
			"label": "Updated label",
			"enabled": true
		}`)
	})

	got, _, err := client.PipelineSchedules.Update(context.Background(), "my-org", "my-pipeline", "abc", in)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if got.Label != "Updated label" {
		t.Errorf("Update label = %q, want %q", got.Label, "Updated label")
	}
	if !got.Enabled {
		t.Errorf("Update enabled = false, want true")
	}
}

func TestPipelineSchedulesService_UpdateClearsEnv(t *testing.T) {
	t.Parallel()

	in := UpdatePipelineSchedule{
		Env: Some(map[string]string{}),
	}

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/pipelines/my-pipeline/schedules/abc", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		assertRequestJSON(t, r, `{"env":{}}`)

		_, _ = fmt.Fprint(w, `{
			"id": "abc",
			"env": {},
			"enabled": true
		}`)
	})

	got, _, err := client.PipelineSchedules.Update(context.Background(), "my-org", "my-pipeline", "abc", in)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if len(got.Env) != 0 {
		t.Errorf("Update env length = %d, want 0", len(got.Env))
	}
}

func TestPipelineSchedulesService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/pipelines/my-pipeline/schedules/abc", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.PipelineSchedules.Delete(context.Background(), "my-org", "my-pipeline", "abc")
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Delete status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}
