package buildkite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestClustersService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/clusters", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
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

	clusters, _, err := client.Clusters.List("my-great-org", nil)

	if err != nil {
		t.Errorf("TestClusters.List returned error: %v", err)
	}

	devClusterCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-09-01T04:27:11.392Z")
	prodClusterCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-09-04T04:25:55.751Z")
	userCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z")

	want := []Cluster{
		{
			ID:          String("528000d8-4ee1-4479-8af1-032b143185f0"),
			GraphQLID:   String("Q2x1c3Rlci0tLTUyODAwMGQ4LTRlZTEtNDQ3OS04YWYxLTAzMmIxNDMxODVmMA=="),
			Name:        String("Development Cluster"),
			Description: String("A cluster for development pipelines"),
			Emoji:       String(":toolbox:"),
			Color:       String("#A9CCE3"),
			URL:         String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0"),
			WebURL:      String("https://buildkite.com/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0"),
			QueuesURL:   String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0/queues"),
			CreatedAt:   NewTimestamp(devClusterCreatedAt),
			CreatedBy: &ClusterCreator{
				ID:        String("7da07e25-0383-4aff-a7cf-14d1a9aa098f"),
				GraphQLID: String("VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg=="),
				Name:      String("Joe Smith"),
				Email:     String("jsmith@example.com"),
				AvatarURL: String("https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4"),
				CreatedAt: NewTimestamp(userCreatedAt),
			},
		},
		{
			ID:          String("3edcecdb-5191-44f1-a5ae-370083c8f92e"),
			GraphQLID:   String("Q2x1c3Rlci0tLTNlZGNlY2RiLTUxOTEtNDRmMS1hNWFlLTM3MDA4M2M4ZjkyZQ=="),
			Name:        String("Production Cluster"),
			Description: String("A cluster for production pipelines"),
			Emoji:       String(":toolbox:"),
			Color:       String("#B9E3A9"),
			URL:         String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/3edcecdb-5191-44f1-a5ae-370083c8f92e"),
			WebURL:      String("https://buildkite.com/organizations/my-great-org/clusters/3edcecdb-5191-44f1-a5ae-370083c8f92e"),
			QueuesURL:   String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/3edcecdb-5191-44f1-a5ae-370083c8f92e/queues"),
			CreatedAt:   NewTimestamp(prodClusterCreatedAt),
			CreatedBy: &ClusterCreator{
				ID:        String("7da07e25-0383-4aff-a7cf-14d1a9aa098f"),
				GraphQLID: String("VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg=="),
				Name:      String("Joe Smith"),
				Email:     String("jsmith@example.com"),
				AvatarURL: String("https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4"),
				CreatedAt: NewTimestamp(userCreatedAt),
			},
		},
	}

	if !reflect.DeepEqual(clusters, want) {
		t.Errorf("TestClusters.List returned %+v, want %+v", clusters, want)
	}
}

func TestClustersService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
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

	cluster, _, err := client.Clusters.Get("my-great-org", "528000d8-4ee1-4479-8af1-032b143185f0")

	if err != nil {
		t.Errorf("TestClusters.Get returned error: %v", err)
	}

	devClusterCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-09-01T04:27:11.392Z")
	userCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z")

	want := &Cluster{
		ID:          String("528000d8-4ee1-4479-8af1-032b143185f0"),
		GraphQLID:   String("Q2x1c3Rlci0tLTUyODAwMGQ4LTRlZTEtNDQ3OS04YWYxLTAzMmIxNDMxODVmMA=="),
		Name:        String("Development Cluster"),
		Description: String("A cluster for development pipelines"),
		Emoji:       String(":toolbox:"),
		Color:       String("#A9CCE3"),
		URL:         String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0"),
		WebURL:      String("https://buildkite.com/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0"),
		QueuesURL:   String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/528000d8-4ee1-4479-8af1-032b143185f0/queues"),
		CreatedAt:   NewTimestamp(devClusterCreatedAt),
		CreatedBy: &ClusterCreator{
			ID:        String("7da07e25-0383-4aff-a7cf-14d1a9aa098f"),
			GraphQLID: String("VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg=="),
			Name:      String("Joe Smith"),
			Email:     String("jsmith@example.com"),
			AvatarURL: String("https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4"),
			CreatedAt: NewTimestamp(userCreatedAt),
		},
	}

	if !reflect.DeepEqual(cluster, want) {
		t.Errorf("TestClusters.Get returned %+v, want %+v", cluster, want)
	}
}

func TestClustersService_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &ClusterCreate{
		Name:        "Testing Cluster",
		Description: String("A cluster for testing"),
		Emoji:       String(":construction:"),
		Color:       String("E5F185"),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/clusters", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterCreate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"name" : "Testing Cluster",
				"description": "A cluster for testing",
				"emoji": ":construction:",
				"color": "E5F185"
			}`)
	})

	cluster, _, err := client.Clusters.Create("my-great-org", input)

	if err != nil {
		t.Errorf("TestClusters.Create returned error: %v", err)
	}

	want := &Cluster{
		Name:        String("Testing Cluster"),
		Description: String("A cluster for testing"),
		Emoji:       String(":construction:"),
		Color:       String("E5F185"),
	}

	if !reflect.DeepEqual(cluster, want) {
		t.Errorf("TestClusters.Create returned %+v, want %+v", cluster, want)
	}
}

func TestClustersService_Update(t *testing.T) {
	setup()
	defer teardown()

	input := &ClusterCreate{
		Name:          "Testing Cluster",
		Description:   String("A cluster for testing"),
		Emoji: 	       String(":construction:"),
		Color:         String("E5F185"),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/clusters", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterCreate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"id": "a32cbe81-82b2-45f7-bd97-66f1ac2c0cc1",
				"name" : "Testing Cluster",
				"description": "A cluster for testing",
				"emoji": ":construction:",
				"color": "E5F185"
			}`)
	})

	cluster, _, err := client.Clusters.Create("my-great-org", input)

	if err != nil {
		t.Errorf("TestClusters.Create returned error: %v", err)
	}

	// Lets update the description of the cluster
	cluster.Description = String("A test cluster")

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/a32cbe81-82b2-45f7-bd97-66f1ac2c0cc1", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterUpdate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "PATCH")

		fmt.Fprint(w,
			`
			{
				"id": "a32cbe81-82b2-45f7-bd97-66f1ac2c0cc1",
				"name" : "Testing Cluster",
				"description": "A test cluster",
				"emoji": ":construction:",
				"color": "E5F185"
			}`)
	})

	clusterUpdate := ClusterUpdate{
		Description: String("A test cluster"),
	}

	_, err = client.Clusters.Update("my-great-org", *cluster.ID, &clusterUpdate)

	if err != nil {
		t.Errorf("TestClusters.Update returned error: %v", err)
	}
	
	want := &Cluster{
		ID:			   String("a32cbe81-82b2-45f7-bd97-66f1ac2c0cc1"),
		Name:          String("Testing Cluster"),
		Description:   String("A test cluster"),
		Emoji: 	       String(":construction:"),
		Color:         String("E5F185"),
	}
	
	if !reflect.DeepEqual(cluster, want) {
		t.Errorf("TestClusters.Update returned %+v, want %+v", cluster, want)
	}
	
	
}

func TestClustersService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/7d2aa9b5-bf2a-4ce0-b9d7-90d3d9b8942c", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Clusters.Delete("my-great-org", "7d2aa9b5-bf2a-4ce0-b9d7-90d3d9b8942c")

	if err != nil {
		t.Errorf("TestClusters.Delete returned error: %v", err)
	}
}
