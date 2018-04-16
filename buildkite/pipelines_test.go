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

func TestPipelinesService_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &CreatePipeline{Name: *String("my-great-pipeline"),
		Repository: *String("my-great-repo"),
		Steps: []Step{Step{Type: String("script"),
			Name:    String("Build :package"),
			Command: String("script/release.sh")}},
	}

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreatePipeline)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{
						"name":"my-great-pipeline",
						"repository":"my-great-repo",
						"steps": [
							{
								"type": "script",
								"name": "Build :package:",
								"command": "script/release.sh"
							}
						]
					}`)
	})

	pipeline, _, err := client.Pipelines.Create("my-great-org", input)
	if err != nil {
		t.Errorf("Pipelines.Create returned error: %v", err)
	}

	want := &Pipeline{Name: String("my-great-pipeline"),
		Repository: String("my-great-repo"),
		Steps: []*Step{&Step{Type: String("script"),
			Name:    String("Build :package:"),
			Command: String("script/release.sh")}},
	}
	if !reflect.DeepEqual(pipeline, want) {
		t.Errorf("Pipelines.Create returned %+v, want %+v", pipeline, want)
	}

}

func TestPipelinesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"123",
						"slug":"my-great-pipeline-slug"}`)
	})

	pipeline, _, err := client.Pipelines.Get("my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.Get returned error: %v", err)
	}

	want := &Pipeline{ID: String("123"), Slug: String("my-great-pipeline-slug")}
	if !reflect.DeepEqual(pipeline, want) {
		t.Errorf("Pipelines.Get returned %+v, want %+v", pipeline, want)
	}
}

func TestPipelinesService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Pipelines.Delete("my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.Delete returned error: %v", err)
	}
}
