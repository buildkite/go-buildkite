package buildkite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestClusterTokensService_List(t *testing.T) {
	setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens", func(w http.ResponseWriter, r *http.Request) {
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

	tokens, _, err := client.ClusterTokens.List("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", nil)

	if err != nil {
		t.Errorf("TestClusterTokens.List returned error: %v", err)
	}

	developmentTokenCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T08:01:02.951Z"))
	testTokenCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T08:05:00.755Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	clusterCreator := &ClusterCreator{
		ID:        String("7da07e25-0383-4aff-a7cf-14d1a9aa098f"),
		GraphQLID: String("VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg=="),
		Name:      String("Joe Smith"),
		Email:     String("jsmith@example.com"),
		AvatarURL: String("https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4"),
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := []ClusterToken{
		{
			ID:                 String("38e8fdb0-52bf-4e73-ad82-ce93cfbaa724"),
			GraphQLID:          String("Q2x1c3RlclRva2VuLS0tMzhlOGZkYjAtNTJiZi00ZTczLWFkODItY2U5M2NmYmFhNzI0"),
			Description:        String("Development cluster token"),
			URL:                String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/38e8fdb0-52bf-4e73-ad82-ce93cfbaa724"),
			ClusterURL:         String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57"),
			CreatedAt:          NewTimestamp(developmentTokenCreatedAt),
			CreatedBy:          clusterCreator,
			AllowedIPAddresses: String("99.26.83.126/24 220.189.137.145/32"),
		},
		{
			ID:          String("218baae1-3e70-44f4-9d52-2e578536693e"),
			GraphQLID:   String("Q2x1c3RlclRva2VuLS0tMjE4YmFhZTEtM2U3MC00NGY0LTlkNTItMmU1Nzg1MzY2OTNl"),
			Description: String("Test cluster token"),
			URL:         String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/218baae1-3e70-44f4-9d52-2e578536693e"),
			ClusterURL:  String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57"),
			CreatedAt:   NewTimestamp(testTokenCreatedAt),
			CreatedBy:   clusterCreator,
		},
	}

	if !reflect.DeepEqual(tokens, want) {
		t.Errorf("TestClusterTokens.List returned %+v, want %+v", tokens, want)
	}
}

func TestClusterTokensService_Get(t *testing.T) {
	setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/38e8fdb0-52bf-4e73-ad82-ce93cfbaa724", func(w http.ResponseWriter, r *http.Request) {
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

	token, _, err := client.ClusterTokens.Get("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "38e8fdb0-52bf-4e73-ad82-ce93cfbaa724")

	if err != nil {
		t.Errorf("TestClusterTokens.Get returned error: %v", err)
	}

	developmentTokenCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T08:01:02.951Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	clusterCreator := &ClusterCreator{
		ID:        String("7da07e25-0383-4aff-a7cf-14d1a9aa098f"),
		GraphQLID: String("VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg=="),
		Name:      String("Joe Smith"),
		Email:     String("jsmith@example.com"),
		AvatarURL: String("https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4"),
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := &ClusterToken{
		ID:                 String("38e8fdb0-52bf-4e73-ad82-ce93cfbaa724"),
		GraphQLID:          String("Q2x1c3RlclRva2VuLS0tMzhlOGZkYjAtNTJiZi00ZTczLWFkODItY2U5M2NmYmFhNzI0"),
		Description:        String("Development cluster token"),
		URL:                String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/38e8fdb0-52bf-4e73-ad82-ce93cfbaa724"),
		ClusterURL:         String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57"),
		CreatedAt:          NewTimestamp(developmentTokenCreatedAt),
		CreatedBy:          clusterCreator,
		AllowedIPAddresses: String("99.26.83.126/24 220.189.137.145/32"),
	}

	if !reflect.DeepEqual(token, want) {
		t.Errorf("TestClusterTokens.Get returned %+v, want %+v", token, want)
	}
}

func TestClusterTokensService_Create(t *testing.T) {
	setup(t)
	t.Cleanup(teardown)

	input := &ClusterTokenCreateUpdate{
		Description: String("Development 2 cluster token"),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterTokenCreateUpdate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"description": "Development 2 cluster token"
			}`)
	})

	token, _, err := client.ClusterTokens.Create("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", input)

	if err != nil {
		t.Errorf("TestClusterTokens.Create returned error: %v", err)
	}

	want := &ClusterToken{
		Description: String("Development 2 cluster token"),
	}

	if !reflect.DeepEqual(token, want) {
		t.Errorf("TestClusterTokens.Create returned %+v, want %+v", token, want)
	}
}

func TestClusterTokensService_Update(t *testing.T) {
	setup(t)
	t.Cleanup(teardown)

	input := &ClusterTokenCreateUpdate{
		Description: String("Development 1 Fleet Token"),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterTokenCreateUpdate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"id": "9cb33339-1c4a-4020-9aeb-3319b2e1f054",
				"description": "Development 1 agent-fleet token"
			}`)
	})

	token, _, err := client.ClusterTokens.Create("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", input)

	if err != nil {
		t.Errorf("TestClusterTokens.Update returned error: %v", err)
	}

	// Lets update the description of the cluster token
	token.Description = String("Development 1 agent token")

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/9cb33339-1c4a-4020-9aeb-3319b2e1f054", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterTokenCreateUpdate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "PATCH")

		fmt.Fprint(w,
			`
			{
				"id": "9cb33339-1c4a-4020-9aeb-3319b2e1f054",
				"description" : "Development 1 agent token"
			}`)
	})

	tokenUpdate := ClusterTokenCreateUpdate{
		Description: String("Development 1 agent token"),
	}

	_, err = client.ClusterTokens.Update("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "9cb33339-1c4a-4020-9aeb-3319b2e1f054", &tokenUpdate)

	if err != nil {
		t.Errorf("TestClusterTokens.Update returned error: %v", err)
	}

	want := &ClusterToken{
		ID:          String("9cb33339-1c4a-4020-9aeb-3319b2e1f054"),
		Description: String("Development 1 agent token"),
	}

	if !reflect.DeepEqual(token, want) {
		t.Errorf("TestClusterTokens.Update returned %+v, want %+v", token, want)
	}
}

func TestClusterTokensService_Delete(t *testing.T) {
	setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/tokens/9cb33339-1c4a-4020-9aeb-3319b2e1f054", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.ClusterTokens.Delete("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "9cb33339-1c4a-4020-9aeb-3319b2e1f054")

	if err != nil {
		t.Errorf("TestClusterTokens.Delete returned error: %v", err)
	}
}
