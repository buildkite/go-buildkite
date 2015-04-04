package buildkite

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestOrganizationsServicList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	orgs, _, err := client.Organizations.List(nil)
	if err != nil {
		t.Errorf("Organizations.List returned error: %v", err)
	}

	want := []Organization{{ID: "123"}, {ID: "1234"}}
	if !reflect.DeepEqual(orgs, want) {
		t.Errorf("Organizations.List returned %+v, want %+v", orgs, want)
	}
}
