package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAnnotationsService_ListByBuild(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/annotations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{
			"id": "de0d4ab5-6360-467a-a34b-e5ef5db5320d",
			"context": "default",
			"style": "info",
			"scope": "build",
			"priority": 3,
			"body_html": "<h1>My Markdown Heading</h1>\n<img src=\"artifact://indy.png\" alt=\"Belongs in a museum\" height=250 />",
			"created_at": "2019-04-09T18:07:15.775Z",
			"updated_at": "2019-08-06T20:58:49.396Z"
		},
		{
			"id": "5b3ceff6-78cb-4fe9-88ae-51be5f145977",
			"context": "coverage",
			"style": "info",
			"scope": "build",
			"priority": 3,
			"body_html": "Read the <a href=\"artifact://coverage/index.html\">uploaded coverage report</a>",
			"created_at": "2019-04-09T18:07:16.320Z",
			"updated_at": "2019-04-09T18:07:16.320Z"
		}]`)
	})

	annotations, _, err := client.Annotations.ListByBuild(context.Background(), "my-great-org", "sup-keith", "awesome-build", nil)
	if err != nil {
		t.Errorf("ListByBuild returned error: %v", err)
	}

	want := []Annotation{
		{
			ID:        "de0d4ab5-6360-467a-a34b-e5ef5db5320d",
			Context:   "default",
			Style:     "info",
			Scope:     "build",
			Priority:  3,
			BodyHTML:  "<h1>My Markdown Heading</h1>\n<img src=\"artifact://indy.png\" alt=\"Belongs in a museum\" height=250 />",
			CreatedAt: NewTimestamp(time.Date(2019, 4, 9, 18, 7, 15, 775000000, time.UTC)),
			UpdatedAt: NewTimestamp(time.Date(2019, 8, 6, 20, 58, 49, 396000000, time.UTC)),
		},
		{
			ID:        "5b3ceff6-78cb-4fe9-88ae-51be5f145977",
			Context:   "coverage",
			Style:     "info",
			Scope:     "build",
			Priority:  3,
			BodyHTML:  "Read the <a href=\"artifact://coverage/index.html\">uploaded coverage report</a>",
			CreatedAt: NewTimestamp(time.Date(2019, 4, 9, 18, 7, 16, 320000000, time.UTC)),
			UpdatedAt: NewTimestamp(time.Date(2019, 4, 9, 18, 7, 16, 320000000, time.UTC)),
		},
	}
	if diff := cmp.Diff(annotations, want); diff != "" {
		t.Errorf("ListByBuild diff: (-got +want)\n%s", diff)
	}
}

func TestAnnotationsService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := AnnotationCreate{
		Style:   "info",
		Context: "default",
		Body:    "<h1>My Markdown Heading</h1>\n<p>An example annotation!</p>",
		Append:  false,
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline/builds/10/annotations", func(w http.ResponseWriter, r *http.Request) {
		var v AnnotationCreate
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		_, _ = fmt.Fprint(w,
			`
			{
				"id": "68aef727-f754-48e1-aad8-5f5da8a9960c",
				"context": "default",
				"style": "info",
				"scope": "build",
				"priority": 3,
				"body_html": "<h1>My Markdown Heading</h1>\n<p>An example annotation!</p>",
				"created_at": "2023-08-21T08:50:05.824Z",
				"updated_at": "2023-08-21T08:50:05.824Z"
			}`)
	})

	annotation, _, err := client.Annotations.Create(context.Background(), "my-great-org", "my-great-pipeline", "10", input)
	if err != nil {
		t.Errorf("TestAnnotations.Create returned error: %v", err)
	}

	annotationCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-21T08:50:05.824Z"))
	annotationUpatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-21T08:50:05.824Z"))

	want := Annotation{
		ID:        "68aef727-f754-48e1-aad8-5f5da8a9960c",
		Context:   "default",
		Style:     "info",
		Scope:     "build",
		Priority:  3,
		BodyHTML:  "<h1>My Markdown Heading</h1>\n<p>An example annotation!</p>",
		CreatedAt: NewTimestamp(annotationCreatedAt),
		UpdatedAt: NewTimestamp(annotationUpatedAt),
	}

	if diff := cmp.Diff(annotation, want); diff != "" {
		t.Errorf("TestAnnotations.Create diff: (-got +want)\n%s", diff)
	}
}

func TestAnnotationsService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline/builds/10/annotations/68aef727-f754-48e1-aad8-5f5da8a9960c", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Annotations.Delete(context.Background(), "my-great-org", "my-great-pipeline", "10", "68aef727-f754-48e1-aad8-5f5da8a9960c")
	if err != nil {
		t.Errorf("TestAnnotations.Delete returned error: %v", err)
	}
}

func TestAnnotationsService_ListByJob(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/a7c5b1d2-4f3e-4a1b-9c8d-6e2f1a3b4c5d/annotations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{
			"id": "de0d4ab5-6360-467a-a34b-e5ef5db5320d",
			"context": "test-results",
			"style": "success",
			"scope": "job",
			"priority": 3,
			"body_html": "<p>All 42 tests passed</p>\n",
			"created_at": "2024-01-15T10:30:00.000Z",
			"updated_at": "2024-01-15T10:30:00.000Z"
		}]`)
	})

	annotations, _, err := client.Annotations.ListByJob(context.Background(), "my-great-org", "sup-keith", "awesome-build", "a7c5b1d2-4f3e-4a1b-9c8d-6e2f1a3b4c5d", nil)
	if err != nil {
		t.Errorf("ListByJob returned error: %v", err)
	}

	want := []Annotation{
		{
			ID:        "de0d4ab5-6360-467a-a34b-e5ef5db5320d",
			Context:   "test-results",
			Style:     "success",
			Scope:     "job",
			Priority:  3,
			BodyHTML:  "<p>All 42 tests passed</p>\n",
			CreatedAt: NewTimestamp(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
			UpdatedAt: NewTimestamp(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
		},
	}
	if diff := cmp.Diff(annotations, want); diff != "" {
		t.Errorf("ListByJob diff: (-got +want)\n%s", diff)
	}
}

func TestAnnotationsService_CreateForJob(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := AnnotationCreate{
		Style:   "success",
		Context: "test-results",
		Body:    "Test results: 42 passed",
		Append:  false,
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline/builds/10/jobs/a7c5b1d2-4f3e-4a1b-9c8d-6e2f1a3b4c5d/annotations", func(w http.ResponseWriter, r *http.Request) {
		var v AnnotationCreate
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, `{
			"id": "a7c5b1d2-4f3e-4a1b-9c8d-6e2f1a3b4c5d",
			"context": "test-results",
			"style": "success",
			"scope": "job",
			"priority": 3,
			"body_html": "<p>Test results: 42 passed</p>\n",
			"created_at": "2024-01-15T10:30:00.000Z",
			"updated_at": "2024-01-15T10:30:00.000Z"
		}`)
	})

	annotation, _, err := client.Annotations.CreateForJob(context.Background(), "my-great-org", "my-great-pipeline", "10", "a7c5b1d2-4f3e-4a1b-9c8d-6e2f1a3b4c5d", input)
	if err != nil {
		t.Errorf("CreateForJob returned error: %v", err)
	}

	want := Annotation{
		ID:        "a7c5b1d2-4f3e-4a1b-9c8d-6e2f1a3b4c5d",
		Context:   "test-results",
		Style:     "success",
		Scope:     "job",
		Priority:  3,
		BodyHTML:  "<p>Test results: 42 passed</p>\n",
		CreatedAt: NewTimestamp(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
		UpdatedAt: NewTimestamp(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
	}

	if diff := cmp.Diff(annotation, want); diff != "" {
		t.Errorf("CreateForJob diff: (-got +want)\n%s", diff)
	}
}

func TestAnnotationsService_DeleteForJob(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline/builds/10/jobs/a7c5b1d2-4f3e-4a1b-9c8d-6e2f1a3b4c5d/annotations/68aef727-f754-48e1-aad8-5f5da8a9960c", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Annotations.DeleteForJob(context.Background(), "my-great-org", "my-great-pipeline", "10", "a7c5b1d2-4f3e-4a1b-9c8d-6e2f1a3b4c5d", "68aef727-f754-48e1-aad8-5f5da8a9960c")
	if err != nil {
		t.Errorf("DeleteForJob returned error: %v", err)
	}
}
