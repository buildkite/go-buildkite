package buildkite

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAPIRequestLogExportsService_Create(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name     string
		create   APIRequestLogExportCreate
		wantBody string
	}{
		{
			name: "explicit period",
			create: APIRequestLogExportCreate{
				Dataset: "rest",
				Format:  "parquet",
				From:    NewTimestamp(time.Date(2026, 11, 21, 0, 0, 0, 0, time.UTC)),
				To:      NewTimestamp(time.Date(2026, 11, 22, 0, 0, 0, 0, time.UTC)),
			},
			wantBody: `{"dataset":"rest","format":"parquet","from":"2026-11-21T00:00:00Z","to":"2026-11-22T00:00:00Z"}`,
		},
		{
			name:     "server defaults",
			create:   APIRequestLogExportCreate{Dataset: "graphql"},
			wantBody: `{"dataset":"graphql"}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			server.HandleFunc("/v2/organizations/acme/api_request_log_exports", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPost)
				assertRequestJSON(t, r, tt.wantBody)
				w.WriteHeader(http.StatusCreated)
				_, _ = fmt.Fprint(w, `{
					"uuid": "export-id",
					"url": "https://api.buildkite.com/v2/organizations/acme/api_request_log_exports/export-id",
					"dataset": "rest",
					"format": "parquet",
					"state": "pending",
					"download_url": null
				}`)
			})

			got, resp, err := client.APIRequestLogExports.Create(context.Background(), "acme", tt.create)
			if err != nil {
				t.Fatalf("Create returned error: %v", err)
			}
			if resp.StatusCode != http.StatusCreated {
				t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusCreated)
			}

			want := APIRequestLogExport{
				UUID:    "export-id",
				URL:     "https://api.buildkite.com/v2/organizations/acme/api_request_log_exports/export-id",
				Dataset: "rest",
				Format:  "parquet",
				State:   "pending",
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Create diff: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestAPIRequestLogExportsService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/acme/api_request_log_exports/export-id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"uuid": "export-id",
			"url": "https://api.buildkite.com/v2/organizations/acme/api_request_log_exports/export-id",
			"dataset": "graphql",
			"format": "ndjson",
			"state": "completed",
			"from": "2026-11-21T00:00:00.000Z",
			"to": "2026-11-22T00:00:00.000Z",
			"row_count": 0,
			"download_url": "https://api.buildkite.com/v2/organizations/acme/api_request_log_exports/export-id/download",
			"error_message": null,
			"created_by_uuid": "user-id",
			"created_at": "2026-11-22T00:20:00.000Z",
			"started_at": "2026-11-22T00:20:05.000Z",
			"finished_at": "2026-11-22T00:22:00.000Z",
			"expires_at": "2026-11-23T00:22:00.000Z"
		}`)
	})

	got, _, err := client.APIRequestLogExports.Get(context.Background(), "acme", "export-id")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	ts := func(value string) *Timestamp { return NewTimestamp(must(time.Parse(time.RFC3339, value))) }
	rowCount := 0
	downloadURL := "https://api.buildkite.com/v2/organizations/acme/api_request_log_exports/export-id/download"
	want := APIRequestLogExport{
		UUID:          "export-id",
		URL:           "https://api.buildkite.com/v2/organizations/acme/api_request_log_exports/export-id",
		Dataset:       "graphql",
		Format:        "ndjson",
		State:         "completed",
		From:          ts("2026-11-21T00:00:00Z"),
		To:            ts("2026-11-22T00:00:00Z"),
		RowCount:      &rowCount,
		DownloadURL:   &downloadURL,
		CreatedByUUID: "user-id",
		CreatedAt:     ts("2026-11-22T00:20:00Z"),
		StartedAt:     ts("2026-11-22T00:20:05Z"),
		FinishedAt:    ts("2026-11-22T00:22:00Z"),
		ExpiresAt:     ts("2026-11-23T00:22:00Z"),
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Get diff: (-want +got)\n%s", diff)
	}
}

func TestAPIRequestLogExportsService_Download(t *testing.T) {
	t.Parallel()

	// The signed link lives on another host, which must not receive the
	// Buildkite token. Go forwards Authorization to the same hostname on any
	// port, so the redirect uses localhost while the API uses 127.0.0.1.
	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("signed link request sent Authorization %q", auth)
		}
		_, _ = fmt.Fprint(w, "PAR1 export bytes")
	}))
	t.Cleanup(storage.Close)

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		if r.URL.Path != "/v2/organizations/acme/api_request_log_exports/export-id/download" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got, want := r.Header.Get("Authorization"), "Bearer token"; got != want {
			t.Errorf("Authorization = %q, want %q", got, want)
		}
		storageURL := strings.Replace(storage.URL, "127.0.0.1", "localhost", 1)
		http.Redirect(w, r, storageURL+"/export.parquet?X-Amz-Expires=600", http.StatusFound)
	}))
	t.Cleanup(api.Close)

	client, err := NewOpts(WithBaseURL(api.URL), WithTokenAuth("token"))
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if _, err := client.APIRequestLogExports.Download(context.Background(), "acme", "export-id", &buf); err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if got, want := buf.String(), "PAR1 export bytes"; got != want {
		t.Errorf("Download content = %q, want %q", got, want)
	}
}
