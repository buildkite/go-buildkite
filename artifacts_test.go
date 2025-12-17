package buildkite

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestArtifactsService_ListByBuild(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{"id":"art-1","job_id":"job-1","path":"artifact.txt"},{"id":"art-2","job_id":"job-2","path":"artifact2.txt"}]`)
	})

	artifacts, _, err := client.Artifacts.ListByBuild(context.Background(), "my-great-org", "sup-keith", "123", nil)
	if err != nil {
		t.Errorf("ListByBuild returned error: %v", err)
	}

	want := []Artifact{
		{ID: "art-1", JobID: "job-1", Path: "artifact.txt"},
		{ID: "art-2", JobID: "job-2", Path: "artifact2.txt"},
	}
	if diff := cmp.Diff(artifacts, want); diff != "" {
		t.Errorf("ListByBuild diff: (-got +want)\n%s", diff)
	}
}

func TestArtifactsService_ListByBuild_WithPagination(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "2",
			"per_page": "10",
		})
		_, _ = fmt.Fprint(w, `[{"id":"art-1"}]`)
	})

	opt := &ArtifactListOptions{
		ListOptions: ListOptions{Page: 2, PerPage: 10},
	}
	artifacts, _, err := client.Artifacts.ListByBuild(context.Background(), "my-great-org", "sup-keith", "123", opt)
	if err != nil {
		t.Errorf("ListByBuild returned error: %v", err)
	}

	want := []Artifact{{ID: "art-1"}}
	if diff := cmp.Diff(artifacts, want); diff != "" {
		t.Errorf("ListByBuild diff: (-got +want)\n%s", diff)
	}
}

func TestArtifactsService_ListByJob(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-456/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{"id":"art-1","job_id":"job-456","path":"output.log"}]`)
	})

	artifacts, _, err := client.Artifacts.ListByJob(context.Background(), "my-great-org", "sup-keith", "123", "job-456", nil)
	if err != nil {
		t.Errorf("ListByJob returned error: %v", err)
	}

	want := []Artifact{{ID: "art-1", JobID: "job-456", Path: "output.log"}}
	if diff := cmp.Diff(artifacts, want); diff != "" {
		t.Errorf("ListByJob diff: (-got +want)\n%s", diff)
	}
}

func TestArtifactsService_ListByJob_WithPagination(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-456/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "3",
		})
		_, _ = fmt.Fprint(w, `[{"id":"art-1"}]`)
	})

	opt := &ArtifactListOptions{
		ListOptions: ListOptions{Page: 3},
	}
	artifacts, _, err := client.Artifacts.ListByJob(context.Background(), "my-great-org", "sup-keith", "123", "job-456", opt)
	if err != nil {
		t.Errorf("ListByJob returned error: %v", err)
	}

	want := []Artifact{{ID: "art-1"}}
	if diff := cmp.Diff(artifacts, want); diff != "" {
		t.Errorf("ListByJob diff: (-got +want)\n%s", diff)
	}
}

func TestArtifactsService_GetArtifact(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-456/artifacts/art-789", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{
			"id": "art-789",
			"job_id": "job-456",
			"url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-456/artifacts/art-789",
			"download_url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-456/artifacts/art-789/download",
			"state": "finished",
			"path": "logs/output.log",
			"dirname": "logs",
			"filename": "output.log",
			"mime_type": "text/plain",
			"file_size": 1024,
			"sha1sum": "abc123def456"
		}`)
	})

	artifact, _, err := client.Artifacts.GetArtifact(context.Background(), "my-great-org", "sup-keith", "123", "job-456", "art-789")
	if err != nil {
		t.Errorf("GetArtifact returned error: %v", err)
	}

	want := Artifact{
		ID:          "art-789",
		JobID:       "job-456",
		URL:         "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-456/artifacts/art-789",
		DownloadURL: "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-456/artifacts/art-789/download",
		State:       "finished",
		Path:        "logs/output.log",
		Dirname:     "logs",
		Filename:    "output.log",
		MimeType:    "text/plain",
		FileSize:    1024,
		SHA1:        "abc123def456",
	}
	if diff := cmp.Diff(artifact, want); diff != "" {
		t.Errorf("GetArtifact diff: (-got +want)\n%s", diff)
	}
}

func TestArtifactsService_DownloadArtifactByURL(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	expectedContent := "This is the artifact content"
	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-456/artifacts/art-789/download", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprint(w, expectedContent)
	})

	var buf bytes.Buffer
	_, err := client.Artifacts.DownloadArtifactByURL(context.Background(), "v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-456/artifacts/art-789/download", &buf)
	if err != nil {
		t.Errorf("DownloadArtifactByURL returned error: %v", err)
	}

	if got := buf.String(); got != expectedContent {
		t.Errorf("DownloadArtifactByURL content = %q, want %q", got, expectedContent)
	}
}
