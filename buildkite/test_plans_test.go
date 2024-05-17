package buildkite

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTestPlansService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-1/test_plan", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`
			{
				"identifier": "7c202aaa-3165-4811-9813-173c4c285463",
			  "parallelism": 2,
        "tasks": [
          {
            "node_number": 0,
            "tests": {
              "cases": [
                {
                  "path": "spec/models/apple_spec.rb",
                  "estimated_duration": 1
                },
                {
                  "path": "spec/models/pear_spec.rb",
                  "estimated_duration": 2
                }
              ],
              "format": "files"
            }
          },
          {
            "node_number": 1,
            "tests": {
              "cases": [
                {
                  "path": "spec/models/banana_spec.rb",
                  "estimated_duration": 3
                }
              ],
              "format": "files"
            }
          }
        ]
			}`)
	})

	suite, _, err := client.TestPlans.Get("my-great-org", "suite-1", "7c202aaa-3165-4811-9813-173c4c285463")

	if err != nil {
		t.Errorf("TestPlans.Get returned error: %v", err)
	}

	want := &TestPlan{
		Identifier:    String("7c202aaa-3165-4811-9813-173c4c285463"),
		Parallelism:   Int(2),
    Tasks:         []Task{
      {
        NodeNumber: Int(0),
        Tests: &Tests{
          Format: String("files"),
          Cases: []Case{
            {
              Path: String("spec/models/apple_spec.rb"),
              EstimatedDuration: Int(1),
            },
            {
              Path: String("spec/models/pear_spec.rb"),
              EstimatedDuration: Int(2),
            },
          },
        },
      },
      {
        NodeNumber: Int(1),
        Tests: &Tests{
          Format: String("files"),
          Cases: []Case{
            {
              Path: String("spec/models/banana_spec.rb"),
              EstimatedDuration: Int(3),
            },
          },
        },
      },
    },
	}

	if !reflect.DeepEqual(suite, want) {
		t.Errorf("TestPlans.Get returned %+v, want %+v", suite, want)
	}
}
