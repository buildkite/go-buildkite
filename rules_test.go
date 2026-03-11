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

func TestRulesService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/rules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`[
				{
					"uuid": "a8e-4f87-a462-95b0dac9a681",
					"graphql_id": "UnVsZS0tLWE4ZS00Zjg3LWE0NjItOTViMGRhYzlhNjgx",
					"action": "trigger_build",
					"type": "pipeline.trigger_build.pipeline",
					"effect": "allow",
					"source_type": "pipeline",
					"source_uuid": "016da969-5a0c-4b22-b0c9-4e736012792a",
					"target_type": "pipeline",
					"target_uuid": "a562b76c-b990-4e0e-8046-291e8d310277",
					"organization_uuid": "a98961b7-adc1-41aa-8726-cfb2c46e42e0",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/rules/a8e-4f87-a462-95b0dac9a681",
					"created_at": "2024-01-01T10:00:00.000Z",
					"created_by": {
						"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
						"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
						"name": "Joe Smith",
						"email": "jsmith@example.com",
						"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
						"created_at": "2023-02-20T03:00:05.824Z"
					}
				},
				{
					"uuid": "b9f-5a98-b573-06c1ebd0b782",
					"graphql_id": "UnVsZS0tLWI5Zi01YTk4LWI1NzMtMDZjMWViZDBiNzgy",
					"action": "trigger_build",
					"type": "pipeline.trigger_build.pipeline",
					"effect": "allow",
					"source_type": "pipeline",
					"source_uuid": "126db970-6b1d-5c33-c1da-5f737123893b",
					"target_type": "pipeline",
					"target_uuid": "b673c87d-ca01-5f1f-9147-302f9d411388",
					"organization_uuid": "a98961b7-adc1-41aa-8726-cfb2c46e42e0",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/rules/b9f-5a98-b573-06c1ebd0b782",
					"created_at": "2024-01-02T10:00:00.000Z",
					"created_by": {
						"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
						"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
						"name": "Joe Smith",
						"email": "jsmith@example.com",
						"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
						"created_at": "2023-02-20T03:00:05.824Z"
					}
				}
			]`)
	})

	rules, _, err := client.Rules.List(context.Background(), "my-great-org", nil)
	if err != nil {
		t.Errorf("TestRules.List returned error: %v", err)
	}

	rule1CreatedAt := must(time.Parse(BuildKiteDateFormat, "2024-01-01T10:00:00.000Z"))
	rule2CreatedAt := must(time.Parse(BuildKiteDateFormat, "2024-01-02T10:00:00.000Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	creator := RuleCreator{
		ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		GraphqlID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
		Name:      "Joe Smith",
		Email:     "jsmith@example.com",
		AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := []Rule{
		{
			UUID:             "a8e-4f87-a462-95b0dac9a681",
			GraphqlID:        "UnVsZS0tLWE4ZS00Zjg3LWE0NjItOTViMGRhYzlhNjgx",
			Action:           "trigger_build",
			Type:             "pipeline.trigger_build.pipeline",
			Effect:           "allow",
			SourceType:       "pipeline",
			SourceUUID:       "016da969-5a0c-4b22-b0c9-4e736012792a",
			TargetType:       "pipeline",
			TargetUUID:       "a562b76c-b990-4e0e-8046-291e8d310277",
			OrganizationUUID: "a98961b7-adc1-41aa-8726-cfb2c46e42e0",
			URL:              "https://api.buildkite.com/v2/organizations/my-great-org/rules/a8e-4f87-a462-95b0dac9a681",
			CreatedAt:        NewTimestamp(rule1CreatedAt),
			CreatedBy:        creator,
		},
		{
			UUID:             "b9f-5a98-b573-06c1ebd0b782",
			GraphqlID:        "UnVsZS0tLWI5Zi01YTk4LWI1NzMtMDZjMWViZDBiNzgy",
			Action:           "trigger_build",
			Type:             "pipeline.trigger_build.pipeline",
			Effect:           "allow",
			SourceType:       "pipeline",
			SourceUUID:       "126db970-6b1d-5c33-c1da-5f737123893b",
			TargetType:       "pipeline",
			TargetUUID:       "b673c87d-ca01-5f1f-9147-302f9d411388",
			OrganizationUUID: "a98961b7-adc1-41aa-8726-cfb2c46e42e0",
			URL:              "https://api.buildkite.com/v2/organizations/my-great-org/rules/b9f-5a98-b573-06c1ebd0b782",
			CreatedAt:        NewTimestamp(rule2CreatedAt),
			CreatedBy:        creator,
		},
	}

	if diff := cmp.Diff(rules, want); diff != "" {
		t.Errorf("TestRules.List diff: (-got +want)\n%s", diff)
	}
}

func TestRulesService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/rules/a8e-4f87-a462-95b0dac9a681", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`{
				"uuid": "a8e-4f87-a462-95b0dac9a681",
				"graphql_id": "UnVsZS0tLWE4ZS00Zjg3LWE0NjItOTViMGRhYzlhNjgx",
				"action": "trigger_build",
				"type": "pipeline.trigger_build.pipeline",
				"effect": "allow",
				"source_type": "pipeline",
				"source_uuid": "016da969-5a0c-4b22-b0c9-4e736012792a",
				"target_type": "pipeline",
				"target_uuid": "a562b76c-b990-4e0e-8046-291e8d310277",
				"organization_uuid": "a98961b7-adc1-41aa-8726-cfb2c46e42e0",
				"url": "https://api.buildkite.com/v2/organizations/my-great-org/rules/a8e-4f87-a462-95b0dac9a681",
				"created_at": "2024-01-01T10:00:00.000Z",
				"created_by": {
					"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
					"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
					"name": "Joe Smith",
					"email": "jsmith@example.com",
					"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
					"created_at": "2023-02-20T03:00:05.824Z"
				}
			}`)
	})

	rule, _, err := client.Rules.Get(context.Background(), "my-great-org", "a8e-4f87-a462-95b0dac9a681")
	if err != nil {
		t.Errorf("TestRules.Get returned error: %v", err)
	}

	ruleCreatedAt := must(time.Parse(BuildKiteDateFormat, "2024-01-01T10:00:00.000Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	creator := RuleCreator{
		ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		GraphqlID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
		Name:      "Joe Smith",
		Email:     "jsmith@example.com",
		AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := Rule{
		UUID:             "a8e-4f87-a462-95b0dac9a681",
		GraphqlID:        "UnVsZS0tLWE4ZS00Zjg3LWE0NjItOTViMGRhYzlhNjgx",
		Action:           "trigger_build",
		Type:             "pipeline.trigger_build.pipeline",
		Effect:           "allow",
		SourceType:       "pipeline",
		SourceUUID:       "016da969-5a0c-4b22-b0c9-4e736012792a",
		TargetType:       "pipeline",
		TargetUUID:       "a562b76c-b990-4e0e-8046-291e8d310277",
		OrganizationUUID: "a98961b7-adc1-41aa-8726-cfb2c46e42e0",
		URL:              "https://api.buildkite.com/v2/organizations/my-great-org/rules/a8e-4f87-a462-95b0dac9a681",
		CreatedAt:        NewTimestamp(ruleCreatedAt),
		CreatedBy:        creator,
	}

	if diff := cmp.Diff(rule, want); diff != "" {
		t.Errorf("TestRules.Get diff: (-got +want)\n%s", diff)
	}
}

func TestRulesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	ruleCreate := RuleCreate{
		Type: "pipeline.trigger_build.pipeline",
		Value: RuleValue{
			SourcePipeline: "016da969-5a0c-4b22-b0c9-4e736012792a",
			TargetPipeline: "a562b76c-b990-4e0e-8046-291e8d310277",
		},
	}

	server.HandleFunc("/v2/organizations/my-great-org/rules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var v RuleCreate
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		if diff := cmp.Diff(v, ruleCreate); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		_, _ = fmt.Fprint(w,
			`{
				"uuid": "a8e-4f87-a462-95b0dac9a681",
				"graphql_id": "UnVsZS0tLWE4ZS00Zjg3LWE0NjItOTViMGRhYzlhNjgx",
				"action": "trigger_build",
				"type": "pipeline.trigger_build.pipeline",
				"effect": "allow",
				"source_type": "pipeline",
				"source_uuid": "016da969-5a0c-4b22-b0c9-4e736012792a",
				"target_type": "pipeline",
				"target_uuid": "a562b76c-b990-4e0e-8046-291e8d310277",
				"organization_uuid": "a98961b7-adc1-41aa-8726-cfb2c46e42e0",
				"url": "https://api.buildkite.com/v2/organizations/my-great-org/rules/a8e-4f87-a462-95b0dac9a681",
				"created_at": "2024-01-01T10:00:00.000Z",
				"created_by": {
					"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
					"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
					"name": "Joe Smith",
					"email": "jsmith@example.com",
					"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
					"created_at": "2023-02-20T03:00:05.824Z"
				}
			}`)
	})

	rule, _, err := client.Rules.Create(context.Background(), "my-great-org", ruleCreate)
	if err != nil {
		t.Errorf("TestRules.Create returned error: %v", err)
	}

	ruleCreatedAt := must(time.Parse(BuildKiteDateFormat, "2024-01-01T10:00:00.000Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	want := Rule{
		UUID:             "a8e-4f87-a462-95b0dac9a681",
		GraphqlID:        "UnVsZS0tLWE4ZS00Zjg3LWE0NjItOTViMGRhYzlhNjgx",
		Action:           "trigger_build",
		Type:             "pipeline.trigger_build.pipeline",
		Effect:           "allow",
		SourceType:       "pipeline",
		SourceUUID:       "016da969-5a0c-4b22-b0c9-4e736012792a",
		TargetType:       "pipeline",
		TargetUUID:       "a562b76c-b990-4e0e-8046-291e8d310277",
		OrganizationUUID: "a98961b7-adc1-41aa-8726-cfb2c46e42e0",
		URL:              "https://api.buildkite.com/v2/organizations/my-great-org/rules/a8e-4f87-a462-95b0dac9a681",
		CreatedAt:        NewTimestamp(ruleCreatedAt),
		CreatedBy: RuleCreator{
			ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
			GraphqlID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
			Name:      "Joe Smith",
			Email:     "jsmith@example.com",
			AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
			CreatedAt: NewTimestamp(userCreatedAt),
		},
	}

	if diff := cmp.Diff(rule, want); diff != "" {
		t.Errorf("TestRules.Create diff: (-got +want)\n%s", diff)
	}
}

func TestRulesService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/rules/a8e-4f87-a462-95b0dac9a681", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Rules.Delete(context.Background(), "my-great-org", "a8e-4f87-a462-95b0dac9a681")
	if err != nil {
		t.Errorf("TestRules.Delete returned error: %v", err)
	}
}
