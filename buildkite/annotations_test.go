package buildkite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestAnnotationsService_ListByBuild(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/annotations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"id": "de0d4ab5-6360-467a-a34b-e5ef5db5320d",
			"context": "default",
			"style": "info",
			"body_html": "<h1>My Markdown Heading</h1>\n<img src=\"artifact://indy.png\" alt=\"Belongs in a museum\" height=250 />",
			"created_at": "2019-04-09T18:07:15.775Z",
			"updated_at": "2019-08-06T20:58:49.396Z"
		},
		{
			"id": "5b3ceff6-78cb-4fe9-88ae-51be5f145977",
			"context": "coverage",
			"style": "info",
			"body_html": "Read the <a href=\"artifact://coverage/index.html\">uploaded coverage report</a>",
			"created_at": "2019-04-09T18:07:16.320Z",
			"updated_at": "2019-04-09T18:07:16.320Z"
		}]`)
	})

	annotations, _, err := client.Annotations.ListByBuild("my-great-org", "sup-keith", "awesome-build", nil)
	if err != nil {
		t.Errorf("ListByBuild returned error: %v", err)
	}

	want := []Annotation{
		{
			ID:        String("de0d4ab5-6360-467a-a34b-e5ef5db5320d"),
			Context:   String("default"),
			Style:     String("info"),
			BodyHTML:  String("<h1>My Markdown Heading</h1>\n<img src=\"artifact://indy.png\" alt=\"Belongs in a museum\" height=250 />"),
			CreatedAt: NewTimestamp(time.Date(2019, 4, 9, 18, 7, 15, 775000000, time.UTC)),
			UpdatedAt: NewTimestamp(time.Date(2019, 8, 6, 20, 58, 49, 396000000, time.UTC)),
		},
		{
			ID:        String("5b3ceff6-78cb-4fe9-88ae-51be5f145977"),
			Context:   String("coverage"),
			Style:     String("info"),
			BodyHTML:  String("Read the <a href=\"artifact://coverage/index.html\">uploaded coverage report</a>"),
			CreatedAt: NewTimestamp(time.Date(2019, 4, 9, 18, 7, 16, 320000000, time.UTC)),
			UpdatedAt: NewTimestamp(time.Date(2019, 4, 9, 18, 7, 16, 320000000, time.UTC)),
		},
	}
	if !reflect.DeepEqual(annotations, want) {
		t.Errorf("ListByBuild returned %+v, want %+v", annotations, want)
	}
}

func TestAnnotationsService_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &AnnotationCreate{
		Style:       String("info"),
		Context:   	 String("default"),
		Body:        String("<h1>My Markdown Heading</h1>\n<p>An example annotation!</p>"),
		Append: 	 Bool(false),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/builds/4d8189ea-10eb-478d-8353-64d36a73f8fb/annotations", func(w http.ResponseWriter, r *http.Request) {
		v := new(AnnotationCreate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"id": "68aef727-f754-48e1-aad8-5f5da8a9960c",
				"context": "default",
				"style": "info",
				"body_html": "<h1>My Markdown Heading</h1>\n<p>An example annotation!</p>",
				"created_at": "2023-08-21T08:50:05.824Z",
				"updated_at": "2023-08-21T08:50:05.824Z"
			}`)
	})

	annotation, _, err := client.Annotations.Create("my-great-org", "4d8189ea-10eb-478d-8353-64d36a73f8fb", input)

	if err != nil {
		t.Errorf("TestAnnotations.Create returned error: %v", err)
	}

	annotationCreatedAt, err := time.Parse(BuildKiteDateFormat, "2023-08-21T08:50:05.824Z")
	annotationUpatedAt, err := time.Parse(BuildKiteDateFormat, "2023-08-21T08:50:05.824Z")

	want := &Annotation{
		ID:        String("68aef727-f754-48e1-aad8-5f5da8a9960c"),
		Context:   String("default"),
		Style:     String("info"),
		BodyHTML:  String("<h1>My Markdown Heading</h1>\n<p>An example annotation!</p>"),
		CreatedAt: NewTimestamp(annotationCreatedAt),
		UpdatedAt: NewTimestamp(annotationUpatedAt),
	}

	if !reflect.DeepEqual(annotation, want) {
		t.Errorf("TestAnnotations.Create returned %+v, want %+v", annotation, want)
	}
}



