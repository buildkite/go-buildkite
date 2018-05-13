package buildkite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestPipelinesService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	pipelines, _, err := client.Pipelines.List("my-great-org", nil)
	if err != nil {
		t.Errorf("Pipelines.List returned error: %v", err)
	}

	want := []Pipeline{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(pipelines, want) {
		t.Errorf("Pipelines.List returned %+v, want %+v", pipelines, want)
	}
}

func TestPipelinesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"123"}`)
	})

	pipeline, _, err := client.Pipelines.Get("my-great-org", "my-great-pipeline")
	if err != nil {
		t.Errorf("Pipelines.Get returned error: %v", err)
	}

	want := &Pipeline{ID: String("123")}
	if !reflect.DeepEqual(pipeline, want) {
		t.Errorf("Pipelines.Get returned %+v, want %+v", pipeline, want)
	}
}

func TestPipelinesService_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &Pipeline{Name: String("My Great Pipeline")}

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		v := new(Pipeline)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":"123"}`)
	})

	pipeline, _, err := client.Pipelines.Create("my-great-org", input)
	if err != nil {
		t.Errorf("Pipelines.Create returned error: %v", err)
	}

	want := &Pipeline{ID: String("123")}
	if !reflect.DeepEqual(pipeline, want) {
		t.Errorf("Pipelines.Create returned %+v, want %+v", pipeline, want)
	}
}

func TestPipelinesService_Edit(t *testing.T) {
	setup()
	defer teardown()

	input := &Pipeline{Name: String("My Great Pipeline")}

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline", func(w http.ResponseWriter, r *http.Request) {
		v := new(Pipeline)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "PATCH")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":"123"}`)
	})

	pipeline, _, err := client.Pipelines.Edit("my-great-org", "my-great-pipeline", input)
	if err != nil {
		t.Errorf("Pipelines.Edit returned error: %v", err)
	}

	want := &Pipeline{ID: String("123")}
	if !reflect.DeepEqual(pipeline, want) {
		t.Errorf("Pipelines.Edit returned %+v, want %+v", pipeline, want)
	}
}

func TestPipelinesService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(204)
	})

	resp, err := client.Pipelines.Delete("my-great-org", "my-great-pipeline")
	if err != nil {
		t.Errorf("Pipelines.Delete returned error: %v", err)
	}

	want := 204
	if resp.StatusCode != want {
		t.Errorf("Pipelines.Delete returned HTTP %v, want HTTP %v", resp.StatusCode, want)
	}
}
