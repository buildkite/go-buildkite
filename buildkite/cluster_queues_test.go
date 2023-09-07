package buildkite

import (
	//	"encoding/json"
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
