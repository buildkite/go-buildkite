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

func TestClusterQueuesService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues", func(w http.ResponseWriter, r *http.Request) {
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

	queues, _, err := client.ClusterQueues.List(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", nil)

	if err != nil {
		t.Errorf("TestClusterQueues.List returned error: %v", err)
	}

	defaultQueueCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-06T15:02:08.951Z"))
	devQueueClusterCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T11:30:17.941Z"))
	devQueuePausedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-25T08:53:05.824Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	clusterCreator := ClusterCreator{
		ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		GraphQLID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
		Name:      "Joe Smith",
		Email:     "jsmith@example.com",
		AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := []ClusterQueue{
		{
			ID:             "da3b950b-5761-4940-900d-d9d88307c337",
			GraphQLID:      "Q2x1c3RlclF1ZXVlLS0tZGEzYjk1MGItNTc2MS00OTQwLTkwMGQtZDlkODgzMDdjMzM3",
			Key:            "default",
			Description:    "Cluster's default queue",
			URL:            "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/da3b950b-5761-4940-900d-d9d88307c337",
			WebURL:         "https://buildkite.com/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/da3b950b-5761-4940-900d-d9d88307c337",
			ClusterURL:     "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
			DispatchPaused: false,
			CreatedAt:      NewTimestamp(defaultQueueCreatedAt),
			CreatedBy:      clusterCreator,
		},
		{
			ID:                 "46718bb6-3b2a-48da-9dcb-922c6b7ba140",
			GraphQLID:          "Q2x1c3RlclF1ZXVlLS0tNDY3MThiYjYtM2IyYS00OGRhLTlkY2ItOTIyYzZiN2JhMTQw",
			Key:                "development",
			Description:        "Development queue",
			URL:                "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140",
			WebURL:             "https://buildkite.com/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140",
			ClusterURL:         "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
			DispatchPaused:     true,
			DispatchPausedBy:   &clusterCreator,
			DispatchPausedAt:   NewTimestamp(devQueuePausedAt),
			DispatchPausedNote: "Weekend queue pause",
			CreatedAt:          NewTimestamp(devQueueClusterCreatedAt),
			CreatedBy:          clusterCreator,
		},
	}

	if diff := cmp.Diff(queues, want); diff != "" {
		t.Errorf("TestClusterQueues.List diff: (-got +want)\n%s", diff)
	}
}

func TestClusterQueuesService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140", func(w http.ResponseWriter, r *http.Request) {
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

	queue, _, err := client.ClusterQueues.Get(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "46718bb6-3b2a-48da-9dcb-922c6b7ba140")

	if err != nil {
		t.Errorf("TestClusterQueues.Get returned error: %v", err)
	}

	devQueueClusterCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-06-07T11:30:17.941Z"))
	devQueuePausedAt := must(time.Parse(BuildKiteDateFormat, "2023-08-25T08:53:05.824Z"))
	userCreatedAt := must(time.Parse(BuildKiteDateFormat, "2023-02-20T03:00:05.824Z"))

	clusterCreator := ClusterCreator{
		ID:        "7da07e25-0383-4aff-a7cf-14d1a9aa098f",
		GraphQLID: "VXNlci0tLTdkYTA3ZTI1LTAzODMtNGFmZi1hN2NmLTE0ZDFhOWFhMDk4Zg==",
		Name:      "Joe Smith",
		Email:     "jsmith@example.com",
		AvatarURL: "https://www.gravatar.com/avatar/593nf93m405mf744n3kg9456jjph9grt4",
		CreatedAt: NewTimestamp(userCreatedAt),
	}

	want := ClusterQueue{
		ID:                 "46718bb6-3b2a-48da-9dcb-922c6b7ba140",
		GraphQLID:          "Q2x1c3RlclF1ZXVlLS0tNDY3MThiYjYtM2IyYS00OGRhLTlkY2ItOTIyYzZiN2JhMTQw",
		Key:                "development",
		Description:        "Development queue",
		URL:                "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140",
		WebURL:             "https://buildkite.com/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/46718bb6-3b2a-48da-9dcb-922c6b7ba140",
		ClusterURL:         "https://api.buildkite.com/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57",
		DispatchPaused:     true,
		DispatchPausedBy:   &clusterCreator,
		DispatchPausedAt:   NewTimestamp(devQueuePausedAt),
		DispatchPausedNote: "Weekend queue pause",
		CreatedAt:          NewTimestamp(devQueueClusterCreatedAt),
		CreatedBy:          clusterCreator,
	}

	if diff := cmp.Diff(queue, want); diff != "" {
		t.Errorf("TestClusterQueues.Get diff: (-got +want)\n%s", diff)
	}
}

func TestClusterQueuesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := ClusterQueueCreate{
		Key:         "development1",
		Description: "Development 1 queue",
	}

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues", func(w http.ResponseWriter, r *http.Request) {
		var v ClusterQueueCreate
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		fmt.Fprint(w,
			`
			{
				"key" : "development1",
				"description": "Development 1 queue"
			}`)
	})

	queue, _, err := client.ClusterQueues.Create(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", input)

	if err != nil {
		t.Errorf("TestClusterQueues.Create returned error: %v", err)
	}

	want := ClusterQueue{
		Key:         "development1",
		Description: "Development 1 queue",
	}

	if diff := cmp.Diff(queue, want); diff != "" {
		t.Errorf("TestClusterQueues.Create diff: (-got +want)\n%s", diff)
	}
}

func TestClusterQueuesService_Update(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/1374ffd0-c5ed-49a5-aebe-67ce906e68ca", func(w http.ResponseWriter, r *http.Request) {
		var v ClusterQueueUpdate
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

	queueUpdate := ClusterQueueUpdate{Description: "Development 1 Team queue"}

	got, _, err := client.ClusterQueues.Update(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "1374ffd0-c5ed-49a5-aebe-67ce906e68ca", queueUpdate)

	if err != nil {
		t.Errorf("TestClusterQueues.Update returned error: %v", err)
	}

	want := ClusterQueue{
		ID:          "1374ffd0-c5ed-49a5-aebe-67ce906e68ca",
		Key:         "development1",
		Description: "Development 1 Team queue",
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestClusterQueues.Update diff: (-got +want)\n%s", diff)
	}
}

func TestClusterQueuesService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/1374ffd0-c5ed-49a5-aebe-67ce906e68ca", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.ClusterQueues.Delete(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "1374ffd0-c5ed-49a5-aebe-67ce906e68ca")

	if err != nil {
		t.Errorf("TestClusterQueues.Delete returned error: %v", err)
	}
}

func TestClusterQueuesService_Pause(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	clusterUUID := "b7c9bc4f-526f-4c18-a3be-dc854ab75d57"
	queueUUID := "b7c9bc4f-526f-4c18-a3be-dc854ab75d57"

	pauseEndpoint := fmt.Sprintf("/v2/organizations/my-great-org/clusters/%s/queues/%s/pause_dispatch", clusterUUID, queueUUID)
	server.HandleFunc(pauseEndpoint, func(w http.ResponseWriter, r *http.Request) {
		var v ClusterQueuePause
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Errorf("TestClusterQueues.Pause error decoding request: %v", err)
		}

		testMethod(t, r, "POST")

		fmt.Fprintf(w,
			`
			{
				"id" : %q,
				"key" : "development1",
				"description": "Development 1 Team queue",
				"dispatch_paused_note": "Pausing dispatch for the weekend"
			}`, queueUUID)
	})

	queuePause := ClusterQueuePause{Note: "Pausing dispatch for the weekend"}
	got, _, err := client.ClusterQueues.Pause(context.Background(), "my-great-org", clusterUUID, queueUUID, queuePause)
	if err != nil {
		t.Errorf("TestClusterQueues.Pause returned error: %v", err)
	}

	want := ClusterQueue{
		ID:                 queueUUID,
		Key:                "development1",
		Description:        "Development 1 Team queue",
		DispatchPausedNote: "Pausing dispatch for the weekend",
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestClusterQueues.Pause diff: (-got +want)\n%s", diff)
	}
}

func TestClusterQueuesService_Resume(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/queues/5cadac07-51dd-4e12-bea3-d91be4655c2f/resume_dispatch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	_, err := client.ClusterQueues.Resume(context.Background(), "my-great-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "5cadac07-51dd-4e12-bea3-d91be4655c2f")

	if err != nil {
		t.Errorf("TestClusterQueues.Resume returned error: %v", err)
	}
}
