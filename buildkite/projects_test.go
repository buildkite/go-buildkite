package buildkite

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestProjectsServicList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/organizations/my-great-org/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	projects, _, err := client.Projects.List("my-great-org", nil)
	if err != nil {
		t.Errorf("Projects.List returned error: %v", err)
	}

	want := []Project{{ID: "123"}, {ID: "1234"}}
	if !reflect.DeepEqual(projects, want) {
		t.Errorf("Projects.List returned %+v, want %+v", projects, want)
	}
}
