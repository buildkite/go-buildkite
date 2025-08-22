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

func TestClustersService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			[
				{
					"id": "528000d8-4ee1-4479-8af1-032b143185f0",
					"graphql_id": "Q2x1c3Rlci0tLTUyODAwMGQ4LTRlZTEtNDQ3OS04YWYxLTAzMmIxNDMxODVmMA==",
					"name": "Development Cluster",
					"description": "A cluster for development pipelines",
					"emoji": ":toolbox:",
					"color": "#A9CCE3",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0",
					"web_url": "https://buildkite.com/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0",
					"queues_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0/queues",
					"created_at": "2023-09-01T04:27:11.392Z",
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
					"id": "3edcecdb-5191-44f1-a5ae-370083c8f92e",
					"graphql_id": "Q2x1c3Rlci0tLTNlZGNlY2RiLTUxOTEtNDRmMS1hNWFlLTM3MDA4M2M4ZjkyZQ==",
					"name": "Production Cluster",
					"description": "A cluster for production pipelines",
					"emoji": ":toolbox:",
					"color": "#B9E3A9",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/3edcecdb-5191-44f1-a5ae-370083c8f92e",
					"web_url": "https://buildkite.com/organizations/my-great-org/clusters/3edcecdb-5191-44f1-a5ae-370083c8f92e",
					"queues_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/3edcecdb-5191-44f1-a5ae-370083c8f92e/queues",
					"created_at": "2023-09-04T04:25:55.751Z",
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

	clusters, _, err := client.Clusters.List(context.Background(), "my-great-org", nil)
	if err != nil {
		t.Errorf("TestClusters.List returned error: %v", err)
	}

	devClusterCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-09-01T04:27:11.392Z"))
	prodClusterCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-09-04T04:25:55.751Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	clusterCreator := ClusterCreator{
		ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		GraphQLID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
		Name:      "Joe Smith",
		Email:     "jsmith@example.com",
		AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := []Cluster{
		{
			ID:          "528000d8-4ee1-4479-8af1-032b143185f0",
			GraphQLID:   "Q2x1c3Rlci0tLTUyODAwMGQ4LTRlZTEtNDQ3OS04YWYxLTAzMmIxNDMxODVmMA==",
			Name:        "Development Cluster",
			Description: "A cluster for development pipelines",
			Emoji:       ":toolbox:",
			Color:       "#A9CCE3",
			URL:         "https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0",
			WebURL:      "https://buildkite.com/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0",
			QueuesURL:   "https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0/queues",
			CreatedAt:   NewTimestamp(devClusterCreatedAt),
			CreatedBy:   clusterCreator,
		},
		{
			ID:          "3edcecdb-5191-44f1-a5ae-370083c8f92e",
			GraphQLID:   "Q2x1c3Rlci0tLTNlZGNlY2RiLTUxOTEtNDRmMS1hNWFlLTM3MDA4M2M4ZjkyZQ==",
			Name:        "Production Cluster",
			Description: "A cluster for production pipelines",
			Emoji:       ":toolbox:",
			Color:       "#B9E3A9",
			URL:         "https://api.buildkite.com/v2/organizations/my-great-org/clusters/3edcecdb-5191-44f1-a5ae-370083c8f92e",
			WebURL:      "https://buildkite.com/organizations/my-great-org/clusters/3edcecdb-5191-44f1-a5ae-370083c8f92e",
			QueuesURL:   "https://api.buildkite.com/v2/organizations/my-great-org/clusters/3edcecdb-5191-44f1-a5ae-370083c8f92e/queues",
			CreatedAt:   NewTimestamp(prodClusterCreatedAt),
			CreatedBy:   clusterCreator,
		},
	}

	if diff := cmp.Diff(clusters, want); diff != "" {
		t.Errorf("TestClusters.List diff: (-got +want)\n%s", diff)
	}
}

func TestClustersService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			{
				"id": "528000d8-4ee1-4479-8af1-032b143185f0",
				"graphql_id": "Q2x1c3Rlci0tLTUyODAwMGQ4LTRlZTEtNDQ3OS04YWYxLTAzMmIxNDMxODVmMA==",
				"name": "Development Cluster",
				"description": "A cluster for development pipelines",
				"emoji": ":toolbox:",
				"color": "#A9CCE3",
				"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0",
				"web_url": "https://buildkite.com/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0",
				"queues_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0/queues",
				"created_at": "2023-09-01T04:27:11.392Z",
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

	cluster, _, err := client.Clusters.Get(context.Background(), "my-great-org", "528000d8-4ee1-4479-8af1-032b143185f0")
	if err != nil {
		t.Errorf("TestClusters.Get returned error: %v", err)
	}

	devClusterCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-09-01T04:27:11.392Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	clusterCreator := ClusterCreator{
		ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		GraphQLID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
		Name:      "Joe Smith",
		Email:     "jsmith@example.com",
		AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := Cluster{
		ID:          "528000d8-4ee1-4479-8af1-032b143185f0",
		GraphQLID:   "Q2x1c3Rlci0tLTUyODAwMGQ4LTRlZTEtNDQ3OS04YWYxLTAzMmIxNDMxODVmMA==",
		Name:        "Development Cluster",
		Description: "A cluster for development pipelines",
		Emoji:       ":toolbox:",
		Color:       "#A9CCE3",
		URL:         "https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0",
		WebURL:      "https://buildkite.com/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0",
		QueuesURL:   "https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0/queues",
		CreatedAt:   NewTimestamp(devClusterCreatedAt),
		CreatedBy:   clusterCreator,
	}

	if diff := cmp.Diff(cluster, want); diff != "" {
		t.Errorf("TestClusters.Get diff: (-got +want)\n%s", diff)
	}
}

func TestClustersService_Create(t *testing.T) {
	t.Parallel()

	type test struct {
		name        string
		input       ClusterCreate
		want        Cluster
		wantErr     bool
		wantErrText string
	}

	tests := []test{
		{
			name: "happy path",
			input: ClusterCreate{
				Name:        "Testing Cluster",
				Description: "A cluster for testing",
				Emoji:       ":construction:",
				Color:       "E5F185",
			},
			want: Cluster{
				Name:        "Testing Cluster",
				Description: "A cluster for testing",
				Emoji:       ":construction:",
				Color:       "E5F185",
			},
		},
		{
			name: "with maintainer",
			input: ClusterCreate{
				Name:        "Testing Cluster",
				Description: "A cluster for testing",
				Emoji:       ":construction:",
				Color:       "E5F185",
				Maintainers: []ClusterMaintainer{
					{
						UserID: "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
					},
				},
			},
			want: Cluster{
				Name:        "Testing Cluster",
				Description: "A cluster for testing",
				Emoji:       ":construction:",
				Color:       "E5F185",
				Maintainers: []ClusterMaintainer{
					{
						UserID: "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
					},
				},
			},
		},
		{
			name: "bad request",
			input: ClusterCreate{
				Name:        "Testing Cluster",
				Description: "A cluster for testing",
				Emoji:       ":construction:",
				Color:       "E5F185",
			},
			wantErr:     true,
			wantErrText: "bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			server.HandleFunc("/v2/organizations/my-great-org/clusters", func(w http.ResponseWriter, r *http.Request) {
				var v ClusterCreate
				err := json.NewDecoder(r.Body).Decode(&v)
				if err != nil {
					t.Fatalf("Error parsing json body: %v", err)
				}

				testMethod(t, r, "POST")

				if diff := cmp.Diff(v, tt.input); diff != "" {
					t.Errorf("Request body diff: (-got +want)\n%s", diff)
				}

				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = fmt.Fprint(w, `{"error": "bad request"}`)
					return
				}

				maintainersJSON := "null"
				if len(tt.input.Maintainers) > 0 {
					maintainersBytes, _ := json.Marshal(tt.input.Maintainers)
					maintainersJSON = string(maintainersBytes)
				}

				_, _ = fmt.Fprintf(w,
					`{
						"name": "Testing Cluster",
						"description": "A cluster for testing",
						"emoji": ":construction:",
						"color": "E5F185",
						"maintainers": %s
					}`, maintainersJSON)
			})

			cluster, _, err := client.Clusters.Create(context.Background(), "my-great-org", tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestClusters.Create returned error: %v", err)
			}

			if diff := cmp.Diff(cluster, tt.want); diff != "" {
				t.Errorf("TestClusters.Create diff: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestClustersService_Update(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	orgSlug := "my-great-org"
	clusterID := "a32cbe81-82b2-45f7-bd97-66f1ac2c0cc1"

	clustersPutEndpoint := fmt.Sprintf("/v2/organizations/%s/clusters/%s", orgSlug, clusterID)
	server.HandleFunc(clustersPutEndpoint, func(w http.ResponseWriter, r *http.Request) {
		var v ClusterUpdate
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "PATCH")

		_, _ = fmt.Fprint(w,
			`
			{
				"id": "a32cbe81-82b2-45f7-bd97-66f1ac2c0cc1",
				"name" : "Testing Cluster",
				"description": "A test cluster",
				"emoji": ":construction:",
				"color": "E5F185"
			}`)
	})

	clusterUpdate := ClusterUpdate{Description: "A test cluster"}
	cluster, _, err := client.Clusters.Update(context.Background(), orgSlug, clusterID, clusterUpdate)
	if err != nil {
		t.Errorf("TestClusters.Update returned error: %v", err)
	}

	want := Cluster{
		ID:          "a32cbe81-82b2-45f7-bd97-66f1ac2c0cc1",
		Name:        "Testing Cluster",
		Description: "A test cluster",
		Emoji:       ":construction:",
		Color:       "E5F185",
	}

	if diff := cmp.Diff(cluster, want); diff != "" {
		t.Errorf("TestClusters.Update diff: (-got +want)\n%s", diff)
	}
}

func TestClustersService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/7d2aa9b5-bf2a-4ce0-b9d7-90d3d9b8942c", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Clusters.Delete(context.Background(), "my-great-org", "7d2aa9b5-bf2a-4ce0-b9d7-90d3d9b8942c")
	if err != nil {
		t.Errorf("TestClusters.Delete returned error: %v", err)
	}
}
