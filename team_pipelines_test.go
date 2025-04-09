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

func TestTeamPipelinesService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/teams/c6fa9b07-efeb-4aea-b5ad-c4aa01e91038/pipelines", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			[{
				"access_level": "manage_build_and_read",
				"created_at": "2023-08-10T05:24:08.651Z",
				"pipeline_id": "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
				"pipeline_url": "https://api.buildkite.com/v2/organizations/my-org/pipelines/pipeline-1"
			},
			{
				"access_level": "manage_build_and_read",
				"created_at": "2023-08-10T05:24:08.663Z",
				"pipeline_id": "4569ddb1-1697-4fad-a46b-372f7318432d",
				"pipeline_url": "https://api.buildkite.com/v2/organizations/my-org/pipelines/pipeline-2"
			}]
			`)
	})

	got, _, err := client.TeamPipelines.List(context.Background(), "my-org", "c6fa9b07-efeb-4aea-b5ad-c4aa01e91038", nil)
	if err != nil {
		t.Errorf("TeamPipelinesService.List returned error: %v", err)
	}

	pipeline1CreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-10T05:24:08.651Z"))
	pipeline2CreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-10T05:24:08.663Z"))

	want := []TeamPipeline{
		{
			ID:          "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
			URL:         "https://api.buildkite.com/v2/organizations/my-org/pipelines/pipeline-1",
			AccessLevel: "manage_build_and_read",
			CreatedAt:   NewTimestamp(pipeline1CreatedAt),
		},
		{
			ID:          "4569ddb1-1697-4fad-a46b-372f7318432d",
			URL:         "https://api.buildkite.com/v2/organizations/my-org/pipelines/pipeline-2",
			AccessLevel: "manage_build_and_read",
			CreatedAt:   NewTimestamp(pipeline2CreatedAt),
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TeamPipelinesService.List diff: (-got +want)\n%s", diff)
	}
}

func TestTeamPipelinesService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/teams/c6fa9b07-efeb-4aea-b5ad-c4aa01e91038/pipelines/1239d7f9-394a-4d99-badf-7c3d8577a8ff", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			{
				"access_level": "manage_build_and_read",
				"created_at": "2023-08-10T05:24:08.651Z",
				"pipeline_id": "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
				"pipeline_url": "https://api.buildkite.com/v2/organizations/my-org/pipelines/pipeline-1"
			}`)
	})

	got, _, err := client.TeamPipelines.Get(context.Background(), "my-org", "c6fa9b07-efeb-4aea-b5ad-c4aa01e91038", "1239d7f9-394a-4d99-badf-7c3d8577a8ff")
	if err != nil {
		t.Errorf("TeamPipelinesService.Get returned error: %v", err)
	}

	pipelineCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-10T05:24:08.651Z"))
	want := TeamPipeline{
		ID:          "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
		URL:         "https://api.buildkite.com/v2/organizations/my-org/pipelines/pipeline-1",
		AccessLevel: "manage_build_and_read",
		CreatedAt:   NewTimestamp(pipelineCreatedAt),
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TeamPipelinesService.Get diff: (-got +want)\n%s", diff)
	}

}

func TestTeamPipelinesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := CreateTeamPipelines{
		PipelineID:  "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
		AccessLevel: "manage_build_and_read",
	}

	server.HandleFunc("/v2/organizations/my-org/teams/c6fa9b07-efeb-4aea-b5ad-c4aa01e91038/pipelines", func(w http.ResponseWriter, r *http.Request) {
		var v CreateTeamPipelines
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Errorf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("create Team Pipelines input diff: (-got +want)\n%s", diff)
		}

		_, _ = fmt.Fprint(w,
			`
			{
				"access_level": "manage_build_and_read",
				"created_at": "2023-08-10T05:24:08.651Z",
				"pipeline_id": "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
				"pipeline_url": "https://api.buildkite.com/v2/organizations/my-org/pipelines/pipeline-1"
			}`)
	})

	got, _, err := client.TeamPipelines.Create(context.Background(), "my-org", "c6fa9b07-efeb-4aea-b5ad-c4aa01e91038", input)
	if err != nil {
		t.Errorf("TeamPipelinesService.Create returned error: %v", err)
	}

	pipelineCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-10T05:24:08.651Z"))
	want := TeamPipeline{
		ID:          "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
		URL:         "https://api.buildkite.com/v2/organizations/my-org/pipelines/pipeline-1",
		AccessLevel: "manage_build_and_read",
		CreatedAt:   NewTimestamp(pipelineCreatedAt),
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TeamPipelinesService.Create diff: (-got +want)\n%s", diff)
	}
}

func TestTeamPipelinesService_Update(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/teams/c6fa9b07-efeb-4aea-b5ad-c4aa01e91038/pipelines/1239d7f9-394a-4d99-badf-7c3d8577a8ff", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		_, _ = fmt.Fprint(w,
			`
			{
				"access_level": "build_and_read",
				"created_at": "2023-08-10T05:24:08.651Z",
				"pipeline_id": "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
				"pipeline_url": "https://api.buildkite.com/v2/organizations/my-org/pipelines/pipeline-1"
			}`)
	})

	wantUpdate := UpdateTeamPipelines{
		AccessLevel: "build_and_read",
	}

	got, _, err := client.TeamPipelines.Update(context.Background(), "my-org", "c6fa9b07-efeb-4aea-b5ad-c4aa01e91038", "1239d7f9-394a-4d99-badf-7c3d8577a8ff", wantUpdate)
	if err != nil {
		t.Errorf("TeamPipelinesService.Get returned error: %v", err)
	}

	pipelineCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-10T05:24:08.651Z"))
	want := TeamPipeline{
		ID:          "1239d7f9-394a-4d99-badf-7c3d8577a8ff",
		URL:         "https://api.buildkite.com/v2/organizations/my-org/pipelines/pipeline-1",
		AccessLevel: "build_and_read",
		CreatedAt:   NewTimestamp(pipelineCreatedAt),
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TeamPipelinesService.Get diff: (-got +want)\n%s", diff)
	}
}

func TestTeamPipelinesService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/teams/c6fa9b07-efeb-4aea-b5ad-c4aa01e91038/pipelines/1239d7f9-394a-4d99-badf-7c3d8577a8ff", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.TeamPipelines.Delete(context.Background(), "my-org", "c6fa9b07-efeb-4aea-b5ad-c4aa01e91038", "1239d7f9-394a-4d99-badf-7c3d8577a8ff")

	if err != nil {
		t.Errorf("TeamPipelinesService.Delete returned error: %v", err)
	}

}
