package buildkite

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestFlakyTestsService_List(t *testing.T) {
	setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/flaky-tests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`
			[
				{
					"id": "a915535c-a8f1-4e1a-bd6a-a5589e09f349",
					"web_url": "https://buildkite.com/organizations/my_great_org/analytics/suites/suite-example/tests/a915535c-a8f1-4e1a-bd6a-a5589e09f349",
					"scope": "User#email",
					"name": "TestExample1_Create",
					"location": "./spec/models/text_example.rb:55",
					"file_name": "./spec/models/text_example.rb",
					"instances": 1,
					"most_recent_instance_at": "2023-05-19T20:00:02.223Z"
				},
				{
					"id": "01867216-8478-7fde-a55a-0300f88bb49b",
					"web_url": "https://buildkite.com/organizations/my_great_org/analytics/suites/suite-example/tests/01867216-8478-7fde-a55a-0300f88bb49b",
					"scope": "User#email",
					"name": "TestExample1_Delete",
					"location": "./spec/models/text_example.rb:102",
					"file_name": "./spec/models/text_example.rb",
					"instances": 2,
					"most_recent_instance_at": "2023-07-10T13:14:03.214Z"
				}
			]`)
	})

	flakyTests, _, err := client.FlakyTests.List("my-great-org", "suite-example", nil)

	if err != nil {
		t.Errorf("TestSuites.List returned error: %v", err)
	}

	// Create Time instances from strings in BuildKiteDateFormat friendly format
	parsedTime1 := must(time.Parse(BuildKiteDateFormat, "2023-05-19T20:00:02.223Z"))
	parsedTime2 := must(time.Parse(BuildKiteDateFormat, "2023-07-10T13:14:03.214Z"))

	if err != nil {
		t.Errorf("TestSuites.List time.Parse error: %v", err)
	}

	want := []FlakyTest{
		{
			ID:                   String("a915535c-a8f1-4e1a-bd6a-a5589e09f349"),
			WebURL:               String("https://buildkite.com/organizations/my_great_org/analytics/suites/suite-example/tests/a915535c-a8f1-4e1a-bd6a-a5589e09f349"),
			Scope:                String("User#email"),
			Name:                 String("TestExample1_Create"),
			Location:             String("./spec/models/text_example.rb:55"),
			FileName:             String("./spec/models/text_example.rb"),
			Instances:            Int(1),
			MostRecentInstanceAt: NewTimestamp(parsedTime1),
		},
		{
			ID:                   String("01867216-8478-7fde-a55a-0300f88bb49b"),
			WebURL:               String("https://buildkite.com/organizations/my_great_org/analytics/suites/suite-example/tests/01867216-8478-7fde-a55a-0300f88bb49b"),
			Scope:                String("User#email"),
			Name:                 String("TestExample1_Delete"),
			Location:             String("./spec/models/text_example.rb:102"),
			FileName:             String("./spec/models/text_example.rb"),
			Instances:            Int(2),
			MostRecentInstanceAt: NewTimestamp(parsedTime2),
		},
	}

	if !reflect.DeepEqual(flakyTests, want) {
		t.Errorf("FlakyTests.List returned %+v, want %+v", flakyTests, want)
	}
}
