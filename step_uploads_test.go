package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestStepUploadsService_ListByBuild(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		if got, want := r.URL.Query().Get("filter[source_job_id]"), "48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba"; got != want {
			t.Errorf("filter[source_job_id] query param = %q, want %q", got, want)
		}

		_, _ = fmt.Fprint(w, `{
  "items": [
    {
      "uuid": "1f1f1f1f-9c3b-4b0a-8a51-b1e59a0e5b3a",
      "graphql_id": "QnVpbGRTdGVwVXBsb2Fk",
      "state": "applied",
      "source": "job",
      "source_job_id": "48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba",
      "replace_existing_steps": false,
      "created_jobs_count": 2,
      "rejection_type": null,
      "message": null,
      "url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads/1f1f1f1f-9c3b-4b0a-8a51-b1e59a0e5b3a",
      "created_at": "2026-08-01T04:00:00Z",
      "processed_at": "2026-08-01T04:00:02Z"
    },
    {
      "uuid": "2e2e2e2e-1234-4b0a-8a51-b1e59a0e5b3a",
      "state": "rejected",
      "source": "job",
      "source_job_id": "48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba",
      "replace_existing_steps": true,
      "created_jobs_count": null,
      "rejection_type": "validation_error",
      "message": "The step upload was rejected. Check the pipeline upload command output for details.",
      "url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads/2e2e2e2e-1234-4b0a-8a51-b1e59a0e5b3a",
      "created_at": "2026-08-01T04:01:00Z",
      "processed_at": "2026-08-01T04:01:02Z"
    }
  ],
  "links": {
    "self": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads?filter[source_job_id]=48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba",
    "next": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads?after=abc123&filter[source_job_id]=48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba&per_page=2"
  }
}`)
	})

	uploads, _, err := client.StepUploads.ListByBuild(
		context.Background(),
		"my-great-org", "sup-keith", "123",
		&StepUploadsListOptions{SourceJobID: "48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba"},
	)
	if err != nil {
		t.Errorf("StepUploads.ListByBuild returned error: %v", err)
	}

	createdJobsCount := 2
	rejectionType := "validation_error"
	message := "The step upload was rejected. Check the pipeline upload command output for details."

	want := StepUploadsList{
		Items: []StepUpload{
			{
				UUID:                 "1f1f1f1f-9c3b-4b0a-8a51-b1e59a0e5b3a",
				GraphQLID:            "QnVpbGRTdGVwVXBsb2Fk",
				State:                "applied",
				Source:               "job",
				SourceJobID:          "48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba",
				ReplaceExistingSteps: false,
				CreatedJobsCount:     &createdJobsCount,
				URL:                  "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads/1f1f1f1f-9c3b-4b0a-8a51-b1e59a0e5b3a",
				CreatedAt:            NewTimestamp(time.Date(2026, 8, 1, 4, 0, 0, 0, time.UTC)),
				ProcessedAt:          NewTimestamp(time.Date(2026, 8, 1, 4, 0, 2, 0, time.UTC)),
			},
			{
				UUID:                 "2e2e2e2e-1234-4b0a-8a51-b1e59a0e5b3a",
				State:                "rejected",
				Source:               "job",
				SourceJobID:          "48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba",
				ReplaceExistingSteps: true,
				RejectionType:        &rejectionType,
				Message:              &message,
				URL:                  "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads/2e2e2e2e-1234-4b0a-8a51-b1e59a0e5b3a",
				CreatedAt:            NewTimestamp(time.Date(2026, 8, 1, 4, 1, 0, 0, time.UTC)),
				ProcessedAt:          NewTimestamp(time.Date(2026, 8, 1, 4, 1, 2, 0, time.UTC)),
			},
		},
		Links: StepUploadsListLinks{
			Self: "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads?filter[source_job_id]=48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba",
			Next: "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads?after=abc123&filter[source_job_id]=48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba&per_page=2",
		},
	}

	if diff := cmp.Diff(uploads, want); diff != "" {
		t.Errorf("StepUploads.ListByBuild diff: (-got +want)\n%s", diff)
	}
}

func TestStepUploadsService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads/1f1f1f1f-9c3b-4b0a-8a51-b1e59a0e5b3a", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{
  "uuid": "1f1f1f1f-9c3b-4b0a-8a51-b1e59a0e5b3a",
  "state": "applied",
  "source": "job",
  "source_job_id": "48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba",
  "replace_existing_steps": false,
  "created_jobs_count": 1,
  "rejection_type": null,
  "message": null,
  "url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads/1f1f1f1f-9c3b-4b0a-8a51-b1e59a0e5b3a",
  "created_at": "2026-08-01T04:00:00Z",
  "processed_at": "2026-08-01T04:00:02Z",
  "definition_bytes": 4096,
  "definition_yaml": "steps:\n- command: echo hello\n",
  "definition_yaml_omitted": false
}`)
	})

	upload, _, err := client.StepUploads.Get(
		context.Background(),
		"my-great-org", "sup-keith", "123", "1f1f1f1f-9c3b-4b0a-8a51-b1e59a0e5b3a",
	)
	if err != nil {
		t.Errorf("StepUploads.Get returned error: %v", err)
	}

	createdJobsCount := 1
	definitionBytes := 4096
	definitionYAML := "steps:\n- command: echo hello\n"
	definitionYAMLOmitted := false

	want := StepUpload{
		UUID:                  "1f1f1f1f-9c3b-4b0a-8a51-b1e59a0e5b3a",
		State:                 "applied",
		Source:                "job",
		SourceJobID:           "48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba",
		ReplaceExistingSteps:  false,
		CreatedJobsCount:      &createdJobsCount,
		URL:                   "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads/1f1f1f1f-9c3b-4b0a-8a51-b1e59a0e5b3a",
		CreatedAt:             NewTimestamp(time.Date(2026, 8, 1, 4, 0, 0, 0, time.UTC)),
		ProcessedAt:           NewTimestamp(time.Date(2026, 8, 1, 4, 0, 2, 0, time.UTC)),
		DefinitionBytes:       &definitionBytes,
		DefinitionYAML:        &definitionYAML,
		DefinitionYAMLOmitted: &definitionYAMLOmitted,
	}

	if diff := cmp.Diff(upload, want); diff != "" {
		t.Errorf("StepUploads.Get diff: (-got +want)\n%s", diff)
	}
}

func TestStepUploadsListLink_ToOptions(t *testing.T) {
	t.Parallel()

	link := StepUploadsListLink("https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads?after=abc123&filter[source_job_id]=48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba&per_page=2")

	opts, err := link.ToOptions()
	if err != nil {
		t.Fatalf("ToOptions returned error: %v", err)
	}

	want := &StepUploadsListOptions{
		SourceJobID: "48af33d8-2d4f-4f19-9a1c-e358c1e0f5ba",
		After:       "abc123",
		PerPage:     2,
	}

	if diff := cmp.Diff(opts, want); diff != "" {
		t.Errorf("ToOptions diff: (-got +want)\n%s", diff)
	}
}

func TestStepUploadsListLink_ToOptions_InvalidPerPage(t *testing.T) {
	t.Parallel()

	link := StepUploadsListLink("https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/step-uploads?per_page=abc")

	_, err := link.ToOptions()
	if err == nil {
		t.Fatal("ToOptions expected an error for per_page=abc, got nil")
	}
}
