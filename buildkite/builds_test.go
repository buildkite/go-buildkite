package buildkite

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBuildsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	builds, _, err := client.Builds.List(nil)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}

func TestBuildsService_List_by_status(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"state": "running",
			"page":  "2",
		})
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{"running", "", ListOptions{Page: 2}}
	builds, _, err := client.Builds.List(opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}
