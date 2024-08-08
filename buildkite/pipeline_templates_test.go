package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestPipelineTemplatesService_List(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/pipeline-templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
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

	pipelineTemplateCreator := &PipelineTemplateCreator{
		ID:        String("7da07e25-0383-4aff-a7cf-14d1a9aa098f"),
		GraphQLID: String("VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg=="),
		Name:      String("Joe Smith"),
		Email:     String("jsmith@example.com"),
		AvatarURL: String("https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4"),
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := []PipelineTemplate{
		{
			UUID:          String("90333dc7-b86a-4485-98c3-9419a5dbc52e"),
			GraphQLID:     String("UGlwZWxpbmVUZW1wbG5lLS0tOTAzMzNkYzctYjg2YS00NDg1LTk4YzMtOTQxOWE1ZGJjNTJl=="),
			Name:          String("Pipeline Upload Template"),
			Description:   String("Pipeline template with basic YAML pipeline upload"),
			Configuration: String("steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload\"\n"),
			Available:     Bool(true),
			CreatedAt:     NewTimestamp(basicTemplateCreatedAt),
			CreatedBy:     pipelineTemplateCreator,
			UpdatedAt:     NewTimestamp(basicTemplateCreatedAt),
			UpdatedBy:     pipelineTemplateCreator,
			URL:           String("https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e"),
			WebURL:        String("https://buildkite.com/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e"),
		},
		{
			UUID:          String("6a25cc85-9fa2-4a00-b66c-bfe377bc5f78"),
			GraphQLID:     String("UGlwZWxpbmVUZW1wbG5lLS0tNmEyNWNjODUtOWZhMi00YTAwLWI2NmMtYmZlMzc3YmM1Zjc4=="),
			Name:          String("Pipeline-Dev Upload Template"),
			Description:   String("Pipeline template uploading buildkite-dev.yml"),
			Configuration: String("steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-dev.yml\"\n"),
			Available:     Bool(true),
			CreatedAt:     NewTimestamp(devTemplateCreatedAt),
			CreatedBy:     pipelineTemplateCreator,
			UpdatedAt:     NewTimestamp(devTemplateCreatedAt),
			UpdatedBy:     pipelineTemplateCreator,
			URL:           String("https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/6a25cc85-9fa2-4a00-b66c-bfe377bc5f78"),
			WebURL:        String("https://buildkite.com/organizations/my-great-org/pipeline-templates/6a25cc85-9fa2-4a00-b66c-bfe377bc5f78"),
		},
	}

	if !reflect.DeepEqual(pipelineTemplates, want) {
		t.Errorf("TestPipelineTemplates.List returned %+v, want %+v", pipelineTemplates, want)
	}
}

func TestPipelineTemplatesService_Get(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
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

	pipelineTemplateCreator := &PipelineTemplateCreator{
		ID:        String("7da07e25-0383-4aff-a7cf-14d1a9aa098f"),
		GraphQLID: String("VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg=="),
		Name:      String("Joe Smith"),
		Email:     String("jsmith@example.com"),
		AvatarURL: String("https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4"),
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := &PipelineTemplate{
		UUID:          String("90333dc7-b86a-4485-98c3-9419a5dbc52e"),
		GraphQLID:     String("UGlwZWxpbmVUZW1wbG5lLS0tOTAzMzNkYzctYjg2YS00NDg1LTk4YzMtOTQxOWE1ZGJjNTJl=="),
		Name:          String("Pipeline Upload Template"),
		Description:   String("Pipeline template with basic YAML pipeline upload"),
		Configuration: String("steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload\"\n"),
		Available:     Bool(true),
		CreatedAt:     NewTimestamp(basicTemplateCreatedAt),
		CreatedBy:     pipelineTemplateCreator,
		UpdatedAt:     NewTimestamp(basicTemplateCreatedAt),
		UpdatedBy:     pipelineTemplateCreator,
		URL:           String("https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e"),
		WebURL:        String("https://buildkite.com/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e"),
	}

	if !reflect.DeepEqual(pipelineTemplate, want) {
		t.Errorf("TestPipelineTemplates.Get returned %+v, want %+v", pipelineTemplate, want)
	}
}

func TestPipelineTemplatesService_Create(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := &PipelineTemplateCreateUpdate{
		Name:          String("Production Pipeline uploader"),
		Description:   String("Production pipeline upload template"),
		Configuration: String("steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n"),
		Available:     Bool(true),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/pipeline-templates", func(w http.ResponseWriter, r *http.Request) {
		v := new(PipelineTemplateCreateUpdate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
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

	want := &PipelineTemplate{
		UUID:          String("08ac0872-e3bd-41d2-b0f8-7822bb43ad41"),
		GraphQLID:     String("UGlwZWxpbmVUZW1wbG5lLS0tMDhhYzA4NzItZTNiZC00MWQyLWIwZjgtNzgyMmJiNDNhZDQxIA=="),
		Name:          String("Production Pipeline template"),
		Description:   String("Production pipeline upload template"),
		Configuration: String("steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n"),
		Available:     Bool(true),
	}

	if !reflect.DeepEqual(pipelineTemplate, want) {
		t.Errorf("TestPipelineTemplates.Create returned %+v, want %+v", pipelineTemplate, want)
	}
}

func TestPipelineTemplatesService_Update(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := &PipelineTemplateCreateUpdate{
		Name:          String("Production Pipeline uploader"),
		Description:   String("Production pipeline upload template"),
		Configuration: String("steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n"),
		Available:     Bool(true),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/pipeline-templates", func(w http.ResponseWriter, r *http.Request) {
		v := new(PipelineTemplateCreateUpdate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"uuid": "b8c2e171-1c7d-47a4-a4d1-a20d691f51d0",
				"graphql_id": "UGlwZWxpbmVUZW1wbG5lLS0tYjhjMmUxNzEtMWM3ZC00N2E0LWE0ZDEtYTIwZDY5MWY1MWQw==",
				"name" : "Production Pipeline template",
				"description": "Production pipeline upload template",
				"configuration": "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n",
				"available": true
			}`)
	})

	pipelineTemplate, _, err := client.PipelineTemplates.Create(context.Background(), "my-great-org", input)

	if err != nil {
		t.Errorf("TestPipelineTemplates.Update returned error: %v", err)
	}

	// Lets update the description of the pipeline template
	pipelineTemplate.Description = String("A pipeline template for uploading a production pipeline YAML (pipeline-production.yml)")

	mux.HandleFunc("/v2/organizations/my-great-org/pipeline-templates/b8c2e171-1c7d-47a4-a4d1-a20d691f51d0", func(w http.ResponseWriter, r *http.Request) {
		v := new(PipelineTemplateCreateUpdate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "PATCH")

		fmt.Fprint(w,
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
		Description: String("A pipeline template for uploading a production pipeline YAML (pipeline-production.yml)"),
	}

	_, err = client.PipelineTemplates.Update(context.Background(), "my-great-org", "b8c2e171-1c7d-47a4-a4d1-a20d691f51d0", &pipelineTemplateUpdate)

	if err != nil {
		t.Errorf("TestPipelineTemplates.Update returned error: %v", err)
	}

	want := &PipelineTemplate{
		UUID:          String("b8c2e171-1c7d-47a4-a4d1-a20d691f51d0"),
		GraphQLID:     String("UGlwZWxpbmVUZW1wbG5lLS0tYjhjMmUxNzEtMWM3ZC00N2E0LWE0ZDEtYTIwZDY5MWY1MWQw=="),
		Name:          String("Production Pipeline template"),
		Description:   String("A pipeline template for uploading a production pipeline YAML (pipeline-production.yml)"),
		Configuration: String("steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n"),
		Available:     Bool(true),
	}

	if !reflect.DeepEqual(pipelineTemplate, want) {
		t.Errorf("TestPipelineTemplates.Update returned %+v, want %+v", pipelineTemplate, want)
	}
}

func TestPipelineTemplatesService_Delete(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/pipeline-templates/19dbd05a-96d7-430f-bac0-14b791558562", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.PipelineTemplates.Delete(context.Background(), "my-great-org", "19dbd05a-96d7-430f-bac0-14b791558562")

	if err != nil {
		t.Errorf("TestPipelineTemplates.Delete returned error: %v", err)
	}
}
