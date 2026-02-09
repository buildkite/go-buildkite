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

func TestClusterSecretsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			[
				{
					"id": "a1e2d345-6789-0abc-def1-234567890abc",
					"graphql_id": "Q2x1c3RlclNlY3JldC0tLWExZTJkMzQ1LTY3ODktMGFiYy1kZWYxLTIzNDU2Nzg5MGFiYw==",
					"key": "SSH_PRIVATE_KEY",
					"description": "SSH key for deployment",
					"policy": "any",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets/a1e2d345-6789-0abc-def1-234567890abc",
					"cluster_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
					"created_at": "2023-06-07T08:01:02.951Z",
					"created_by": {
						"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
						"name": "Joe Smith",
						"email": "jsmith@example.com"
					},
					"updated_at": "2023-06-07T08:01:02.951Z",
					"updated_by": null,
					"last_read_at": null,
					"organization": {
						"id": "0198e45b-c0d5-4a0b-8e37-e140af750d2d",
						"slug": "my-great-org",
						"url": "https://api.buildkite.com/v2/organizations/my-great-org",
						"web_url": "https://buildkite.com/my-great-org"
					}
				},
				{
					"id": "b2f3e456-7890-1bcd-ef12-345678901bcd",
					"graphql_id": "Q2x1c3RlclNlY3JldC0tLWIyZjNlNDU2LTc4OTAtMWJjZC1lZjEyLTM0NTY3ODkwMWJjZA==",
					"key": "DEPLOY_TOKEN",
					"description": "Deployment access token",
					"policy": "block",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets/b2f3e456-7890-1bcd-ef12-345678901bcd",
					"cluster_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
					"created_at": "2023-06-07T08:05:00.755Z",
					"created_by": {
						"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
						"name": "Joe Smith",
						"email": "jsmith@example.com"
					},
					"updated_at": "2023-06-07T08:05:00.755Z",
					"updated_by": null,
					"last_read_at": null,
					"organization": {
						"id": "0198e45b-c0d5-4a0b-8e37-e140af750d2d",
						"slug": "my-great-org",
						"url": "https://api.buildkite.com/v2/organizations/my-great-org",
						"web_url": "https://buildkite.com/my-great-org"
					}
				}
			]`)
	})

	secrets, _, err := client.ClusterSecrets.List(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", nil)
	if err != nil {
		t.Errorf("TestClusterSecrets.List returned error: %v", err)
	}

	sshKeyCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T08:01:02.951Z"))
	deployTokenCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T08:05:00.755Z"))

	secretCreator := SecretCreator{
		ID:    "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		Name:  "Joe Smith",
		Email: "jsmith@example.com",
	}

	org := ClusterSecretOrganization{
		ID:     "0198e45b-c0d5-4a0b-8e37-e140af750d2d",
		Slug:   "my-great-org",
		URL:    "https://api.buildkite.com/v2/organizations/my-great-org",
		WebURL: "https://buildkite.com/my-great-org",
	}

	want := []ClusterSecret{
		{
			ID:           "a1e2d345-6789-0abc-def1-234567890abc",
			GraphQLID:    "Q2x1c3RlclNlY3JldC0tLWExZTJkMzQ1LTY3ODktMGFiYy1kZWYxLTIzNDU2Nzg5MGFiYw==",
			Key:          "SSH_PRIVATE_KEY",
			Description:  "SSH key for deployment",
			Policy:       "any",
			URL:          "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets/a1e2d345-6789-0abc-def1-234567890abc",
			ClusterURL:   "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
			CreatedAt:    NewTimestamp(sshKeyCreatedAt),
			CreatedBy:    secretCreator,
			UpdatedAt:    NewTimestamp(sshKeyCreatedAt),
			Organization: org,
		},
		{
			ID:           "b2f3e456-7890-1bcd-ef12-345678901bcd",
			GraphQLID:    "Q2x1c3RlclNlY3JldC0tLWIyZjNlNDU2LTc4OTAtMWJjZC1lZjEyLTM0NTY3ODkwMWJjZA==",
			Key:          "DEPLOY_TOKEN",
			Description:  "Deployment access token",
			Policy:       "block",
			URL:          "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets/b2f3e456-7890-1bcd-ef12-345678901bcd",
			ClusterURL:   "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
			CreatedAt:    NewTimestamp(deployTokenCreatedAt),
			CreatedBy:    secretCreator,
			UpdatedAt:    NewTimestamp(deployTokenCreatedAt),
			Organization: org,
		},
	}

	if diff := cmp.Diff(secrets, want); diff != "" {
		t.Errorf("TestClusterSecrets.List diff: (-got +want)\n%s", diff)
	}
}

func TestClusterSecretsService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets/a1e2d345-6789-0abc-def1-234567890abc", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			{
				"id": "a1e2d345-6789-0abc-def1-234567890abc",
				"graphql_id": "Q2x1c3RlclNlY3JldC0tLWExZTJkMzQ1LTY3ODktMGFiYy1kZWYxLTIzNDU2Nzg5MGFiYw==",
				"key": "SSH_PRIVATE_KEY",
				"description": "SSH key for deployment",
				"policy": "any",
				"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets/a1e2d345-6789-0abc-def1-234567890abc",
				"cluster_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
				"created_at": "2023-06-07T08:01:02.951Z",
				"created_by": {
					"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
					"name": "Joe Smith",
					"email": "jsmith@example.com"
				},
				"updated_at": "2023-06-07T08:01:02.951Z",
				"updated_by": null,
				"last_read_at": null,
				"organization": {
					"id": "0198e45b-c0d5-4a0b-8e37-e140af750d2d",
					"slug": "my-great-org",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org",
					"web_url": "https://buildkite.com/my-great-org"
				}
			}`)
	})

	secret, _, err := client.ClusterSecrets.Get(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "a1e2d345-6789-0abc-def1-234567890abc")
	if err != nil {
		t.Errorf("TestClusterSecrets.Get returned error: %v", err)
	}

	secretCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T08:01:02.951Z"))

	want := ClusterSecret{
		ID:          "a1e2d345-6789-0abc-def1-234567890abc",
		GraphQLID:   "Q2x1c3RlclNlY3JldC0tLWExZTJkMzQ1LTY3ODktMGFiYy1kZWYxLTIzNDU2Nzg5MGFiYw==",
		Key:         "SSH_PRIVATE_KEY",
		Description: "SSH key for deployment",
		Policy:      "any",
		URL:         "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets/a1e2d345-6789-0abc-def1-234567890abc",
		ClusterURL:  "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
		CreatedAt:   NewTimestamp(secretCreatedAt),
		CreatedBy: SecretCreator{
			ID:    "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
			Name:  "Joe Smith",
			Email: "jsmith@example.com",
		},
		UpdatedAt: NewTimestamp(secretCreatedAt),
		Organization: ClusterSecretOrganization{
			ID:     "0198e45b-c0d5-4a0b-8e37-e140af750d2d",
			Slug:   "my-great-org",
			URL:    "https://api.buildkite.com/v2/organizations/my-great-org",
			WebURL: "https://buildkite.com/my-great-org",
		},
	}

	if diff := cmp.Diff(secret, want); diff != "" {
		t.Errorf("TestClusterSecrets.Get diff: (-got +want)\n%s", diff)
	}
}

func TestClusterSecretsService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := ClusterSecretCreate{
		Key:         "SSH_PRIVATE_KEY",
		Value:       "supersecret",
		Description: "SSH key for deployment",
		Policy:      "any",
	}

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets", func(w http.ResponseWriter, r *http.Request) {
		var v ClusterSecretCreate
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
				"id": "a1e2d345-6789-0abc-def1-234567890abc",
				"key": "SSH_PRIVATE_KEY",
				"description": "SSH key for deployment",
				"policy": "any"
			}`)
	})

	secret, _, err := client.ClusterSecrets.Create(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", input)
	if err != nil {
		t.Errorf("TestClusterSecrets.Create returned error: %v", err)
	}

	want := ClusterSecret{
		ID:          "a1e2d345-6789-0abc-def1-234567890abc",
		Key:         "SSH_PRIVATE_KEY",
		Description: "SSH key for deployment",
		Policy:      "any",
	}
	if diff := cmp.Diff(secret, want); diff != "" {
		t.Errorf("TestClusterSecrets.Create diff: (-got +want)\n%s", diff)
	}
}

func TestClusterSecretsService_Update(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets/a1e2d345-6789-0abc-def1-234567890abc", func(w http.ResponseWriter, r *http.Request) {
		var v ClusterSecretUpdate
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "PUT")

		_, _ = fmt.Fprint(w,
			`
			{
				"id": "a1e2d345-6789-0abc-def1-234567890abc",
				"key": "SSH_PRIVATE_KEY",
				"description": "Updated SSH key description",
				"policy": "block"
			}`)
	})

	update := ClusterSecretUpdate{
		Description: "Updated SSH key description",
		Policy:      "block",
	}

	got, _, err := client.ClusterSecrets.Update(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "a1e2d345-6789-0abc-def1-234567890abc", update)
	if err != nil {
		t.Errorf("TestClusterSecrets.Update returned error: %v", err)
	}

	want := ClusterSecret{
		ID:          "a1e2d345-6789-0abc-def1-234567890abc",
		Key:         "SSH_PRIVATE_KEY",
		Description: "Updated SSH key description",
		Policy:      "block",
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestClusterSecrets.Update diff: (-got +want)\n%s", diff)
	}
}

func TestClusterSecretsService_UpdateValue(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets/a1e2d345-6789-0abc-def1-234567890abc/value", func(w http.ResponseWriter, r *http.Request) {
		var v ClusterSecretValueUpdate
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "PUT")

		if v.Value != "new-secret-value" {
			t.Errorf("Request body value = %q, want %q", v.Value, "new-secret-value")
		}

		w.WriteHeader(http.StatusNoContent)
	})

	input := ClusterSecretValueUpdate{Value: "new-secret-value"}

	_, err := client.ClusterSecrets.UpdateValue(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "a1e2d345-6789-0abc-def1-234567890abc", input)
	if err != nil {
		t.Errorf("TestClusterSecrets.UpdateValue returned error: %v", err)
	}
}

func TestClusterSecretsService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/secrets/a1e2d345-6789-0abc-def1-234567890abc", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.ClusterSecrets.Delete(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "a1e2d345-6789-0abc-def1-234567890abc")
	if err != nil {
		t.Errorf("TestClusterSecrets.Delete returned error: %v", err)
	}
}
