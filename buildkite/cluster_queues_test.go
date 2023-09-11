package buildkite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestClusterQueuesService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`
			[
				{
					"id": "da3b950b-5761-4940-900d-d9d88307c337",
					"graphql_id": "Q2x1c3RlclF1ZXVlLS0tZGEzYjk1MGItNTc2MS00OTQwLTkwMGQtZDlkODgzMDdjMzM3",
					"key": "default",
					"description": "Cluster's default queue",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/da3b950b-5761-4940-900d-d9d88307c337",
					"web_url": "https://buildkite.com/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/da3b950b-5761-4940-900d-d9d88307c337",
					"cluster_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
					"dispatch_paused": false,
					"created_at": "2023-06-06T15:02:08.951Z",
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
					"id": "46718bb6-3b2a-48da-9dcb-922c6b7ba140",
					"graphql_id": "Q2x1c3RlclF1ZXVlLS0tNDY3MThiYjYtM2IyYS00OGRhLTlkY2ItOTIyYzZiN2JhMTQw",
					"key": "development",
					"description": "Development queue",
					"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140",
					"web_url": "https://buildkite.com/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140",
					"cluster_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
					"dispatch_paused": true,
					"dispatch_paused_by": {
						"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
						"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
						"name": "Joe Smith",
						"email": "jsmith@example.com",
						"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
						"created_at": "2023-02-20T03:00:05.824Z"
					},
					"dispatch_paused_at": "2023-08-25T08:53:05.824Z",
					"dispatch_paused_note": "Weekend queue pause",
					"created_at": "2023-06-07T11:30:17.941Z",
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

	queues, _, err := client.ClusterQueues.List("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", nil)

	if err != nil {
		t.Errorf("TestClusterQueues.List returned error: %v", err)
	}

	defaultQueueCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-06-06T15:02:08.951Z")
	devQueueClusterCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-06-07T11:30:17.941Z")
	devQueuePausedAt, err := time.Parse(BuildKiteDateFormat, "2023-08-25T08:53:05.824Z")
	userCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z")

	clusterUser := &ClusterUser{
		ID:        String("7da07e25-0383-4aff-a7cf-14d1a9aa098f"),
		GraphQLID: String("VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg=="),
		Name:      String("Joe Smith"),
		Email:     String("jsmith@example.com"),
		AvatarURL: String("https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4"),
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := []ClusterQueue{
		{
			ID:             String("da3b950b-5761-4940-900d-d9d88307c337"),
			GraphQLID:      String("Q2x1c3RlclF1ZXVlLS0tZGEzYjk1MGItNTc2MS00OTQwLTkwMGQtZDlkODgzMDdjMzM3"),
			Key:            String("default"),
			Description:    String("Cluster's default queue"),
			URL:            String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/da3b950b-5761-4940-900d-d9d88307c337"),
			WebURL:         String("https://buildkite.com/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/da3b950b-5761-4940-900d-d9d88307c337"),
			ClusterURL:     String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57"),
			DispatchPaused: Bool(false),
			CreatedAt:      NewTimestamp(defaultQueueCreatedAt),
			CreatedBy:      clusterUser,
		},
		{
			ID:                 String("46718bb6-3b2a-48da-9dcb-922c6b7ba140"),
			GraphQLID:          String("Q2x1c3RlclF1ZXVlLS0tNDY3MThiYjYtM2IyYS00OGRhLTlkY2ItOTIyYzZiN2JhMTQw"),
			Key:                String("development"),
			Description:        String("Development queue"),
			URL:                String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140"),
			WebURL:             String("https://buildkite.com/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140"),
			ClusterURL:         String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57"),
			DispatchPaused:     Bool(true),
			DispatchPausedBy:   clusterUser,
			DispatchPausedAt:   NewTimestamp(devQueuePausedAt),
			DispatchPausedNote: String("Weekend queue pause"),
			CreatedAt:          NewTimestamp(devQueueClusterCreatedAt),
			CreatedBy:          clusterUser,
		},
	}

	if !reflect.DeepEqual(queues, want) {
		t.Errorf("TestClusterQueues.List returned %+v, want %+v", queues, want)
	}
}

func TestClusterQueuesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`
			{
				"id": "46718bb6-3b2a-48da-9dcb-922c6b7ba140",
				"graphql_id": "Q2x1c3RlclF1ZXVlLS0tNDY3MThiYjYtM2IyYS00OGRhLTlkY2ItOTIyYzZiN2JhMTQw",
				"key": "development",
				"description": "Development queue",
				"url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140",
				"web_url": "https://buildkite.com/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140",
				"cluster_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
				"dispatch_paused": true,
				"dispatch_paused_by": {
					"id": "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
					"graphql_id": "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
					"name": "Joe Smith",
					"email": "jsmith@example.com",
					"avatar_url": "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
					"created_at": "2023-02-20T03:00:05.824Z"
				},
				"dispatch_paused_at": "2023-08-25T08:53:05.824Z",
				"dispatch_paused_note": "Weekend queue pause",
				"created_at": "2023-06-07T11:30:17.941Z",
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

	queue, _, err := client.ClusterQueues.Get("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "46718bb6-3b2a-48da-9dcb-922c6b7ba140")

	if err != nil {
		t.Errorf("TestClusterQueues.Get returned error: %v", err)
	}

	devQueueClusterCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-06-07T11:30:17.941Z")
	devQueuePausedAt, err := time.Parse(BuildKiteDateFormat, "2023-08-25T08:53:05.824Z")
	userCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z")

	clusterUser := &ClusterUser{
		ID:        String("7da07e25-0383-4aff-a7cf-14d1a9aa098f"),
		GraphQLID: String("VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg=="),
		Name:      String("Joe Smith"),
		Email:     String("jsmith@example.com"),
		AvatarURL: String("https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4"),
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := &ClusterQueue{
		ID:                 String("46718bb6-3b2a-48da-9dcb-922c6b7ba140"),
		GraphQLID:          String("Q2x1c3RlclF1ZXVlLS0tNDY3MThiYjYtM2IyYS00OGRhLTlkY2ItOTIyYzZiN2JhMTQw"),
		Key:                String("development"),
		Description:        String("Development queue"),
		URL:                String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140"),
		WebURL:             String("https://buildkite.com/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140"),
		ClusterURL:         String("https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57"),
		DispatchPaused:     Bool(true),
		DispatchPausedBy:   clusterUser,
		DispatchPausedAt:   NewTimestamp(devQueuePausedAt),
		DispatchPausedNote: String("Weekend queue pause"),
		CreatedAt:          NewTimestamp(devQueueClusterCreatedAt),
		CreatedBy:          clusterUser,
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("TestClusterQueues.Get returned %+v, want %+v", queue, want)
	}
}

func TestClusterQueuesService_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &ClusterQueueCreate{
		Key:         String("development1"),
		Description: String("Development 1 queue"),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterQueueCreate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"key" : "development1",
				"description": "Development 1 queue"
			}`)
	})

	queue, _, err := client.ClusterQueues.Create("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", input)

	if err != nil {
		t.Errorf("TestClusterQueues.Create returned error: %v", err)
	}

	want := &ClusterQueue{
		Key:         String("development1"),
		Description: String("Development 1 queue"),
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("TestClusterQueues.Create returned %+v, want %+v", queue, want)
	}
}

func TestClusterQueuesService_Update(t *testing.T) {
	setup()
	defer teardown()

	input := &ClusterQueueCreate{
		Key:         String("development1"),
		Description: String("Development 1 queue"),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterQueueCreate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"id" : "1374ffd0-c5ed-49a5-aebe-67ce906e68ca",
				"key" : "development1",
				"description": "Development 1 queue"
			}`)
	})

	queue, _, err := client.ClusterQueues.Create("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", input)

	if err != nil {
		t.Errorf("TestClusterQueues.Update returned error: %v", err)
	}

	// Lets update the description of the cluster queue
	queue.Description = String("Development 1 Team queue")

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/1374ffd0-c5ed-49a5-aebe-67ce906e68ca", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterQueueUpdate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "PATCH")

		fmt.Fprint(w,
			`
			{
				"id" : "1374ffd0-c5ed-49a5-aebe-67ce906e68ca",
				"key" : "development1",
				"description": "Development 1 Team queue"
			}`)
	})

	queueUpdate := ClusterQueueUpdate{
		Description: String("Development 1 Team queue"),
	}

	_, err = client.ClusterQueues.Update("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "1374ffd0-c5ed-49a5-aebe-67ce906e68ca", &queueUpdate)

	if err != nil {
		t.Errorf("TestClusterQueues.Update returned error: %v", err)
	}

	want := &ClusterQueue{
		ID:          String("1374ffd0-c5ed-49a5-aebe-67ce906e68ca"),
		Key:         String("development1"),
		Description: String("Development 1 Team queue"),
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("TestClusterQueues.Update returned %+v, want %+v", queue, want)
	}
}

func TestClusterQueuesService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/1374ffd0-c5ed-49a5-aebe-67ce906e68ca", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.ClusterQueues.Delete("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "1374ffd0-c5ed-49a5-aebe-67ce906e68ca")

	if err != nil {
		t.Errorf("TestClusterQueues.Delete returned error: %v", err)
	}
}

func TestClusterQueuesService_Pause(t *testing.T) {
	setup()
	defer teardown()

	input := &ClusterQueueCreate{
		Key:         String("development1"),
		Description: String("Development 1 Team queue"),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterQueueCreate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"id" : "5cadac07-51dd-4e12-bea3-d91be4655c2f",
				"key" : "development1",
				"description": "Development 1 Team queue"
			}`)
	})

	queue, _, err := client.ClusterQueues.Create("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", input)

	if err != nil {
		t.Errorf("TestClusterQueues.Pause returned error: %v", err)
	}

	// Update the dispatch paused note of the queue
	queue.DispatchPausedNote = String("Pausing dispatch for the weekend")

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/5cadac07-51dd-4e12-bea3-d91be4655c2f/pause_dispatch", func(w http.ResponseWriter, r *http.Request) {
		v := new(ClusterQueueUpdate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		fmt.Fprint(w,
			`
			{
				"id" : "5cadac07-51dd-4e12-bea3-d91be4655c2f",
				"key" : "development1",
				"description": "Development 1 Team queue"
				"dispatch_paused_note": "Pausing dispatch for the weekend"",
			}`)
	})

	queuePause := ClusterQueuePause{
		Note: String("Pausing dispatch for the weekend"),
	}

	_, err = client.ClusterQueues.Pause("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "5cadac07-51dd-4e12-bea3-d91be4655c2f", &queuePause)

	if err != nil {
		t.Errorf("TestClusterQueues.Pause returned error: %v", err)
	}

	want := &ClusterQueue{
		ID:                 String("5cadac07-51dd-4e12-bea3-d91be4655c2f"),
		Key:                String("development1"),
		Description:        String("Development 1 Team queue"),
		DispatchPausedNote: String("Pausing dispatch for the weekend"),
	}

	if !reflect.DeepEqual(queue, want) {
		t.Errorf("TestClusterQueues.Pause returned %+v, want %+v", queue, want)
	}
}

func TestClusterQueuesService_Resume(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/5cadac07-51dd-4e12-bea3-d91be4655c2f/resume_dispatch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	_, err := client.ClusterQueues.Resume("my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "5cadac07-51dd-4e12-bea3-d91be4655c2f")

	if err != nil {
		t.Errorf("TestClusterQueues.Resume returned error: %v", err)
	}
}
