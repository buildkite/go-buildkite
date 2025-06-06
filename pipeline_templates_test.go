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

func TestPipelineTemplatesService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipeline-templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			[
				{
					"uuid": "90333dc7-b86a-4485-98c3-9419a5dbc52e",
					"graphql_id": "UGlwZWxpbmVUZW1wbG5lLS0tOTAzMzNkYzctYjg2YS00NDg1LTk4YzMtOTQxOWE1ZGJjNTJl==",
					"name": "Pipeline Upload Template",
					"description": "Pipeline template with basic YAML pipeline upload",
					"configuration": "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload\"\n",
					"available": true,
					"created_at": "2023-08-11T01:22:05.650Z",
					"created_by": {
						"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
						"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
						"name": "Joe Smith",
						"email": "jsmith@example.com",
						"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
						"created_at": "2023-02-20T03:00:05.824Z"
					},
					"updated_at": "2023-08-11T01:22:05.650Z",
					"updated_by": {
						"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
						"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
						"name": "Joe Smith",
						"email": "jsmith@example.com",
						"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
						"created_at": "2023-02-20T03:00:05.824Z"
					},
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e",
					"web_url": "https://buildkite.com/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e"
				},
				{
					"uuid": "6a25cc85-9fa2-4a00-b66c-bfe377bc5f78",
					"graphql_id": "UGlwZWxpbmVUZW1wbG5lLS0tNmEyNWNjODUtOWZhMi00YTAwLWI2NmMtYmZlMzc3YmM1Zjc4==",
					"name": "Pipeline-Dev Upload Template",
					"description": "Pipeline template uploading buildkite-dev.yml",
					"configuration": "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-dev.yml\"\n",
					"available": true,
					"created_at": "2023-08-11T02:24:33.602Z",
					"created_by": {
						"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
						"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
						"name": "Joe Smith",
						"email": "jsmith@example.com",
						"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
						"created_at": "2023-02-20T03:00:05.824Z"
					},
					"updated_at": "2023-08-11T02:24:33.602Z",
					"updated_by": {
						"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
						"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
						"name": "Joe Smith",
						"email": "jsmith@example.com",
						"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
						"created_at": "2023-02-20T03:00:05.824Z"
					},
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/6a25cc85-9fa2-4a00-b66c-bfe377bc5f78",
					"web_url": "https://buildkite.com/organizations/my-great-org/pipeline-templates/6a25cc85-9fa2-4a00-b66c-bfe377bc5f78"
				}
			]`)
	})

	pipelineTemplates, _, err := client.PipelineTemplates.List(context.Background(), "my-great-org", nil)

	if err != nil {
		t.Errorf("TestPipelineTemplates.List returned error: %v", err)
	}

	basicTemplateCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-11T01:22:05.650Z"))
	devTemplateCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-11T02:24:33.602Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	pipelineTemplateCreator := PipelineTemplateCreator{
		ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		GraphQLID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
		Name:      "Joe Smith",
		Email:     "jsmith@example.com",
		AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := []PipelineTemplate{
		{
			UUID:          "90333dc7-b86a-4485-98c3-9419a5dbc52e",
			GraphQLID:     "UGlwZWxpbmVUZW1wbG5lLS0tOTAzMzNkYzctYjg2YS00NDg1LTk4YzMtOTQxOWE1ZGJjNTJl==",
			Name:          "Pipeline Upload Template",
			Description:   "Pipeline template with basic YAML pipeline upload",
			Configuration: "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload\"\n",
			Available:     true,
			CreatedAt:     NewTimestamp(basicTemplateCreatedAt),
			CreatedBy:     pipelineTemplateCreator,
			UpdatedAt:     NewTimestamp(basicTemplateCreatedAt),
			UpdatedBy:     pipelineTemplateCreator,
			URL:           "https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e",
			WebURL:        "https://buildkite.com/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e",
		},
		{
			UUID:          "6a25cc85-9fa2-4a00-b66c-bfe377bc5f78",
			GraphQLID:     "UGlwZWxpbmVUZW1wbG5lLS0tNmEyNWNjODUtOWZhMi00YTAwLWI2NmMtYmZlMzc3YmM1Zjc4==",
			Name:          "Pipeline-Dev Upload Template",
			Description:   "Pipeline template uploading buildkite-dev.yml",
			Configuration: "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-dev.yml\"\n",
			Available:     true,
			CreatedAt:     NewTimestamp(devTemplateCreatedAt),
			CreatedBy:     pipelineTemplateCreator,
			UpdatedAt:     NewTimestamp(devTemplateCreatedAt),
			UpdatedBy:     pipelineTemplateCreator,
			URL:           "https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/6a25cc85-9fa2-4a00-b66c-bfe377bc5f78",
			WebURL:        "https://buildkite.com/organizations/my-great-org/pipeline-templates/6a25cc85-9fa2-4a00-b66c-bfe377bc5f78",
		},
	}

	if diff := cmp.Diff(pipelineTemplates, want); diff != "" {
		t.Errorf("TestPipelineTemplates.List diff: (-got +want)\n%s", diff)
	}
}

func TestPipelineTemplatesService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			{
				"uuid": "90333dc7-b86a-4485-98c3-9419a5dbc52e",
				"graphql_id": "UGlwZWxpbmVUZW1wbG5lLS0tOTAzMzNkYzctYjg2YS00NDg1LTk4YzMtOTQxOWE1ZGJjNTJl==",
				"name": "Pipeline Upload Template",
				"description": "Pipeline template with basic YAML pipeline upload",
				"configuration": "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload\"\n",
				"available": true,
				"created_at": "2023-08-11T01:22:05.650Z",
				"created_by": {
					"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
					"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
					"name": "Joe Smith",
					"email": "jsmith@example.com",
					"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
					"created_at": "2023-02-20T03:00:05.824Z"
				},
				"updated_at": "2023-08-11T01:22:05.650Z",
				"updated_by": {
					"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
					"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
					"name": "Joe Smith",
					"email": "jsmith@example.com",
					"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
					"created_at": "2023-02-20T03:00:05.824Z"
				},
				"url": "https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e",
				"web_url": "https://buildkite.com/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e"
			}`)
	})

	pipelineTemplate, _, err := client.PipelineTemplates.Get(context.Background(), "my-great-org", "90333dc7-b86a-4485-98c3-9419a5dbc52e")

	if err != nil {
		t.Errorf("TestPipelineTemplates.Get returned error: %v", err)
	}

	basicTemplateCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-11T01:22:05.650Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	pipelineTemplateCreator := PipelineTemplateCreator{
		ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		GraphQLID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
		Name:      "Joe Smith",
		Email:     "jsmith@example.com",
		AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := PipelineTemplate{
		UUID:          "90333dc7-b86a-4485-98c3-9419a5dbc52e",
		GraphQLID:     "UGlwZWxpbmVUZW1wbG5lLS0tOTAzMzNkYzctYjg2YS00NDg1LTk4YzMtOTQxOWE1ZGJjNTJl==",
		Name:          "Pipeline Upload Template",
		Description:   "Pipeline template with basic YAML pipeline upload",
		Configuration: "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload\"\n",
		Available:     true,
		CreatedAt:     NewTimestamp(basicTemplateCreatedAt),
		CreatedBy:     pipelineTemplateCreator,
		UpdatedAt:     NewTimestamp(basicTemplateCreatedAt),
		UpdatedBy:     pipelineTemplateCreator,
		URL:           "https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e",
		WebURL:        "https://buildkite.com/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e",
	}

	if diff := cmp.Diff(pipelineTemplate, want); diff != "" {
		t.Errorf("TestPipelineTemplates.Get diff: (-got +want)\n%s", diff)
	}
}

func TestPipelineTemplatesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := PipelineTemplateCreateUpdate{
		Name:          "Production Pipeline uploader",
		Description:   "Production pipeline upload template",
		Configuration: "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n",
		Available:     true,
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipeline-templates", func(w http.ResponseWriter, r *http.Request) {
		var v PipelineTemplateCreateUpdate
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
				"uuid": "08ac0872-e3bd-41d2-b0f8-7822bb43ad41",
				"graphql_id": "UGlwZWxpbmVUZW1wbG5lLS0tMDhhYzA4NzItZTNiZC00MWQyLWIwZjgtNzgyMmJiNDNhZDQxIA==",
				"name" : "Production Pipeline template",
				"description": "Production pipeline upload template",
				"configuration": "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n",
				"available": true
			}`)
	})

	pipelineTemplate, _, err := client.PipelineTemplates.Create(context.Background(), "my-great-org", input)

	if err != nil {
		t.Errorf("TestPipelineTemplates.Create returned error: %v", err)
	}

	want := PipelineTemplate{
		UUID:          "08ac0872-e3bd-41d2-b0f8-7822bb43ad41",
		GraphQLID:     "UGlwZWxpbmVUZW1wbG5lLS0tMDhhYzA4NzItZTNiZC00MWQyLWIwZjgtNzgyMmJiNDNhZDQxIA==",
		Name:          "Production Pipeline template",
		Description:   "Production pipeline upload template",
		Configuration: "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n",
		Available:     true,
	}

	if diff := cmp.Diff(pipelineTemplate, want); diff != "" {
		t.Errorf("TestPipelineTemplates.Create diff: (-got +want)\n%s", diff)
	}
}

func TestPipelineTemplatesService_Update(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipeline-templates/b8c2e171-1c7d-47a4-a4d1-a20d691f51d0", func(w http.ResponseWriter, r *http.Request) {
		var v PipelineTemplateCreateUpdate
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "PATCH")

		_, _ = fmt.Fprint(w,
			`
			{
				"uuid": "b8c2e171-1c7d-47a4-a4d1-a20d691f51d0",
				"graphql_id": "UGlwZWxpbmVUZW1wbG5lLS0tYjhjMmUxNzEtMWM3ZC00N2E0LWE0ZDEtYTIwZDY5MWY1MWQw==",
				"name" : "Production Pipeline template",
				"description": "A pipeline template for uploading a production pipeline YAML (pipeline-production.yml)",
				"configuration": "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n",
				"available": true
			}`)
	})

	pipelineTemplateUpdate := PipelineTemplateCreateUpdate{
		Description: "A pipeline template for uploading a production pipeline YAML (pipeline-production.yml",
	}

	got, _, err := client.PipelineTemplates.Update(context.Background(), "my-great-org", "b8c2e171-1c7d-47a4-a4d1-a20d691f51d0", pipelineTemplateUpdate)
	if err != nil {
		t.Errorf("TestPipelineTemplates.Update returned error: %v", err)
	}

	want := PipelineTemplate{
		UUID:          "b8c2e171-1c7d-47a4-a4d1-a20d691f51d0",
		GraphQLID:     "UGlwZWxpbmVUZW1wbG5lLS0tYjhjMmUxNzEtMWM3ZC00N2E0LWE0ZDEtYTIwZDY5MWY1MWQw==",
		Name:          "Production Pipeline template",
		Description:   "A pipeline template for uploading a production pipeline YAML (pipeline-production.yml)",
		Configuration: "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n",
		Available:     true,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestPipelineTemplates.Update diff: (-got +want)\n%s", diff)
	}
}

func TestPipelineTemplatesService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipeline-templates/19dbd05a-96d7-430f-bac0-14b791558562", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.PipelineTemplates.Delete(context.Background(), "my-great-org", "19dbd05a-96d7-430f-bac0-14b791558562")

	if err != nil {
		t.Errorf("TestPipelineTemplates.Delete returned error: %v", err)
	}
}
