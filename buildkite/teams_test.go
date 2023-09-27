package buildkite

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTeamsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	teams, _, err := client.Teams.List("my-great-org", nil)
	if err != nil {
		t.Errorf("Teams.List returned error: %v", err)
	}

	want := []Team{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(teams, want) {
		t.Errorf("Teams.List returned %+v, want %+v", teams, want)
	}
}

func TestTeamsService_ListForUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"user_id": "abc",
		})
		fmt.Fprint(w, `[{"id":"123"}]`)
	})

	opt := &TeamsListOptions{UserID: "abc"}
	teams, _, err := client.Teams.List("my-great-org", opt)
	if err != nil {
		t.Errorf("Teams.List returned error: %v", err)
	}

	want := []Team{{ID: String("123")}}
	if ! reflect.DeepEqual(teams, want) {
		t.Errorf("Teams.List returned %+v, want %+v", teams, want)
	}
}
