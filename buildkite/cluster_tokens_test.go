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

func TestClusterTokensService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`
			[
				{
					"id": "38e8fdb0-52bf-4e73-ad82-ce93cfbaa724",
					"graphql_id": "Q2x1c3RlclRva2VuLS0tMzhlOGZkYjAtNTJiZi00ZTczLWFkODItY2U5M2NmYmFhNzI0",
					"description": "Development cluster token",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/38e8fdb0-52bf-4e73-ad82-ce93cfbaa724",
					"cluster_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
					"created_at": "2023-06-07T08:01:02.951Z",
					"created_by": {
						"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
						"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
						"name": "Joe Smith",
						"email": "jsmith@example.com",
						"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
						"created_at": "2023-02-20T03:00:05.824Z"
					},
					"allowed_ip_addresses": "99.26.83.126/24 220.189.137.145/32"
				},
				{
					"id": "218baae1-3e70-44f4-9d52-2e578536693e",
					"graphql_id": "Q2x1c3RlclRva2VuLS0tMjE4YmFhZTEtM2U3MC00NGY0LTlkNTItMmU1Nzg1MzY2OTNl",
					"description": "Test cluster token",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/218baae1-3e70-44f4-9d52-2e578536693e",
					"cluster_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
					"created_at": "2023-06-07T08:05:00.755Z",
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

	tokens, _, err := client.ClusterTokens.List(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", nil)

	if err != nil {
		t.Errorf("TestClusterTokens.List returned error: %v", err)
	}

	developmentTokenCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T08:01:02.951Z"))
	testTokenCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T08:05:00.755Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	clusterCreator := ClusterCreator{
		ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		GraphQLID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
		Name:      "Joe Smith",
		Email:     "jsmith@example.com",
		AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := []ClusterToken{
		{
			ID:                 "38e8fdb0-52bf-4e73-ad82-ce93cfbaa724",
			GraphQLID:          "Q2x1c3RlclRva2VuLS0tMzhlOGZkYjAtNTJiZi00ZTczLWFkODItY2U5M2NmYmFhNzI0",
			Description:        "Development cluster token",
			URL:                "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/38e8fdb0-52bf-4e73-ad82-ce93cfbaa724",
			ClusterURL:         "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
			CreatedAt:          NewTimestamp(developmentTokenCreatedAt),
			CreatedBy:          clusterCreator,
			AllowedIPAddresses: "99.26.83.126/24 220.189.137.145/32",
		},
		{
			ID:          "218baae1-3e70-44f4-9d52-2e578536693e",
			GraphQLID:   "Q2x1c3RlclRva2VuLS0tMjE4YmFhZTEtM2U3MC00NGY0LTlkNTItMmU1Nzg1MzY2OTNl",
			Description: "Test cluster token",
			URL:         "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/218baae1-3e70-44f4-9d52-2e578536693e",
			ClusterURL:  "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
			CreatedAt:   NewTimestamp(testTokenCreatedAt),
			CreatedBy:   clusterCreator,
		},
	}

	if diff := cmp.Diff(tokens, want); diff != "" {
		t.Errorf("TestClusterTokens.List diff: (-got +want)\n%s", diff)
	}
}

func TestClusterTokensService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/38e8fdb0-52bf-4e73-ad82-ce93cfbaa724", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`
			{
				"id": "38e8fdb0-52bf-4e73-ad82-ce93cfbaa724",
				"graphql_id": "Q2x1c3RlclRva2VuLS0tMzhlOGZkYjAtNTJiZi00ZTczLWFkODItY2U5M2NmYmFhNzI0",
				"description": "Development cluster token",
				"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/38e8fdb0-52bf-4e73-ad82-ce93cfbaa724",
				"cluster_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
				"created_at": "2023-06-07T08:01:02.951Z",
				"created_by": {
					"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
					"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
					"name": "Joe Smith",
					"email": "jsmith@example.com",
					"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
					"created_at": "2023-02-20T03:00:05.824Z"
				},
				"allowed_ip_addresses": "99.26.83.126/24 220.189.137.145/32"
			}`)
	})

	token, _, err := client.ClusterTokens.Get(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "38e8fdb0-52bf-4e73-ad82-ce93cfbaa724")
	if err != nil {
		t.Errorf("TestClusterTokens.Get returned error: %v", err)
	}

	developmentTokenCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T08:01:02.951Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	clusterCreator := ClusterCreator{
		ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		GraphQLID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
		Name:      "Joe Smith",
		Email:     "jsmith@example.com",
		AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := ClusterToken{
		ID:                 "38e8fdb0-52bf-4e73-ad82-ce93cfbaa724",
		GraphQLID:          "Q2x1c3RlclRva2VuLS0tMzhlOGZkYjAtNTJiZi00ZTczLWFkODItY2U5M2NmYmFhNzI0",
		Description:        "Development cluster token",
		URL:                "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/38e8fdb0-52bf-4e73-ad82-ce93cfbaa724",
		ClusterURL:         "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
		CreatedAt:          NewTimestamp(developmentTokenCreatedAt),
		CreatedBy:          clusterCreator,
		AllowedIPAddresses: "99.26.83.126/24 220.189.137.145/32",
	}

	if diff := cmp.Diff(token, want); diff != "" {
		t.Errorf("TestClusterTokens.Get diff: (-got +want)\n%s", diff)
	}
}

func TestClusterTokensService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := ClusterTokenCreateUpdate{Description: "Development 2 cluster token"}

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens", func(w http.ResponseWriter, r *http.Request) {
		var v ClusterTokenCreateUpdate
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		fmt.Fprint(w,
			`
			{
				"description": "Development 2 cluster token"
			}`)
	})

	token, _, err := client.ClusterTokens.Create(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", input)

	if err != nil {
		t.Errorf("TestClusterTokens.Create returned error: %v", err)
	}

	want := ClusterToken{Description: "Development 2 cluster token"}
	if diff := cmp.Diff(token, want); diff != "" {
		t.Errorf("TestClusterTokens.Create diff: (-got +want)\n%s", diff)
	}
}

func TestClusterTokensService_Update(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/9cb33339-1c4a-4020-9aeb-3319b2e1f054", func(w http.ResponseWriter, r *http.Request) {
		var v ClusterTokenCreateUpdate
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "PATCH")

		fmt.Fprint(w,
			`
			{
				"id": "9cb33339-1c4a-4020-9aeb-3319b2e1f054",
				"description" : "Development 1 agent token"
			}`)
	})

	tokenUpdate := ClusterTokenCreateUpdate{Description: "Development 1 agent token"}

	got, _, err := client.ClusterTokens.Update(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "9cb33339-1c4a-4020-9aeb-3319b2e1f054", tokenUpdate)
	if err != nil {
		t.Errorf("TestClusterTokens.Update returned error: %v", err)
	}

	want := ClusterToken{
		ID:          "9cb33339-1c4a-4020-9aeb-3319b2e1f054",
		Description: "Development 1 agent token",
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestClusterTokens.Update diff: (-got +want)\n%s", diff)
	}
}

func TestClusterTokensService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/9cb33339-1c4a-4020-9aeb-3319b2e1f054", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.ClusterTokens.Delete(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "9cb33339-1c4a-4020-9aeb-3319b2e1f054")
	if err != nil {
		t.Errorf("TestClusterTokens.Delete returned error: %v", err)
	}
}
