package buildkite

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestPipelineTemplatesService_List(t *testing.T) {
	setup()
	defer teardown()

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

	pipelineTemplates, _, err := client.PipelineTemplates.List("my-great-org", nil)

	if err != nil {
		t.Errorf("TestPipelineTemplates.List returned error: %v", err)
	}

	basicTemplateCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-08-11T01:22:05.650Z")
	devTemplateCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-08-11T02:24:33.602Z")
	userCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z")

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
			UUID: 			String("90333dc7-b86a-4485-98c3-9419a5dbc52e"),
			GraphQLID:      String("UGlwZWxpbmVUZW1wbG5lLS0tOTAzMzNkYzctYjg2YS00NDg1LTk4YzMtOTQxOWE1ZGJjNTJl=="),
			Name:           String("Pipeline Upload Template"),
			Description:    String("Pipeline template with basic YAML pipeline upload"),
			Configuration:  String("steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload\"\n"),
			Available:      Bool(true),
			CreatedAt:		NewTimestamp(basicTemplateCreatedAt),
			CreatedBy:		pipelineTemplateCreator,
			UpdatedAt:      NewTimestamp(basicTemplateCreatedAt),
			UpdatedBy:      pipelineTemplateCreator,
			URL:			String("https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e"),
			WebURL:			String("https://buildkite.com/organizations/my-great-org/pipeline-templates/90333dc7-b86a-4485-98c3-9419a5dbc52e"),
		},
		{
			UUID: 			String("6a25cc85-9fa2-4a00-b66c-bfe377bc5f78"),
			GraphQLID:      String("UGlwZWxpbmVUZW1wbG5lLS0tNmEyNWNjODUtOWZhMi00YTAwLWI2NmMtYmZlMzc3YmM1Zjc4=="),
			Name:           String("Pipeline-Dev Upload Template"),
			Description:    String("Pipeline template uploading buildkite-dev.yml"),
			Configuration:  String("steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-dev.yml\"\n"),
			Available:      Bool(true),
			CreatedAt:		NewTimestamp(devTemplateCreatedAt),
			CreatedBy:		pipelineTemplateCreator,
			UpdatedAt:      NewTimestamp(devTemplateCreatedAt),
			UpdatedBy:      pipelineTemplateCreator,
			URL:			String("https://api.buildkite.com/v2/organizations/my-great-org/pipeline-templates/6a25cc85-9fa2-4a00-b66c-bfe377bc5f78"),
			WebURL:			String("https://buildkite.com/organizations/my-great-org/pipeline-templates/6a25cc85-9fa2-4a00-b66c-bfe377bc5f78"),
		},
	}

	if !reflect.DeepEqual(pipelineTemplates, want) {
		t.Errorf("TestPipelineTemplates.List returned %+v, want %+v", pipelineTemplates, want)
	}
}