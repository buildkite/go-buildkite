package buildkite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTestSuiteService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/analytics/organizations/my-great-org/suites", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, 
			`
			[
				{
					"id": "7c202aaa-3165-4811-9813-173c4c285463",
					"graphql_id": "N2MyMDJhYWEtMzE2NS00ODExLTk4MTMtMTczYzRjMjg1NDYz=",
					"slug": "suite-1",
					"name": "suite-1",
					"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-1",
					"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-1",
					"default_branch": "main"
				},
				{
					"id": "38ed1d73-cea9-4aba-b223-def25e66ef51",
					"graphql_id": "MzhlZDFkNzMtY2VhOS00YWJhLWIyMjMtZGVmMjVlNjZlZjUx=",
					"slug": "suite-2",
					"name": "suite-2",
					"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-2",
					"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-2",
					"default_branch": "main"
				}
			]`)
	})

	suites, _, err := client.TestSuites.List("my-great-org", nil)

	if err != nil {
		t.Errorf("TestSuites.List returned error: %v", err)
	}

	want := []TestSuite{
		{
			ID: String("7c202aaa-3165-4811-9813-173c4c285463"),
			GraphQLID: String("N2MyMDJhYWEtMzE2NS00ODExLTk4MTMtMTczYzRjMjg1NDYz="),
			Slug: String("suite-1"),
			Name: String("suite-1"),
			URL: String("https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-1"),
			WebURL: String("https://buildkite.com/organizations/my-great-org/analytics/suites/suite-1"),
			DefaultBranch: String("main"),
		}, 
		{
			ID: String("38ed1d73-cea9-4aba-b223-def25e66ef51"),
			GraphQLID: String("MzhlZDFkNzMtY2VhOS00YWJhLWIyMjMtZGVmMjVlNjZlZjUx="),
			Slug: String("suite-2"),
			Name: String("suite-2"),
			URL: String("https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-2"),
			WebURL: String("https://buildkite.com/organizations/my-great-org/analytics/suites/suite-2"),
			DefaultBranch: String("main"),
		},
	}
	if !reflect.DeepEqual(suites, want) {
		t.Errorf("TestSuites.List returned %+v, want %+v", suites, want)
	}
}

func TestTestSuiteService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, 
			`
			{
				"id": "7c202aaa-3165-4811-9813-173c4c285463",
				"graphql_id": "N2MyMDJhYWEtMzE2NS00ODExLTk4MTMtMTczYzRjMjg1NDYz=",
				"slug": "suite-1",
				"name": "suite-1",
				"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-1",
				"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suites/suite-1",
				"default_branch": "main"
			}`)
	})

	suite, _, err := client.TestSuites.Get("my-great-org", "suite-1")

	if err != nil {
		t.Errorf("TestSuites.Get returned error: %v", err)
	}

	want := &TestSuite{		
		ID: String("7c202aaa-3165-4811-9813-173c4c285463"),
		GraphQLID: String("N2MyMDJhYWEtMzE2NS00ODExLTk4MTMtMTczYzRjMjg1NDYz="),
		Slug: String("suite-1"),
		Name: String("suite-1"),
		URL: String("https://api.buildkite.com/v2/analytics/organizations/my-great-org/suites/suite-1"),
		WebURL: String("https://buildkite.com/organizations/my-great-org/analytics/suites/suite-1"),
		DefaultBranch: String("main"),
	} 
	
	if !reflect.DeepEqual(suite, want) {
		t.Errorf("TestSuites.Get returned %+v, want %+v", suite, want)
	}

}

func TestTestSuiteService_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &TestSuiteCreate{
		Name:			"Suite 3",   						
		DefaultBranch:	"main",        
		TeamUuids: 		[]string{"8369b300-fff0-4ef1-91de-010f72f4458d"},             
	}

	mux.HandleFunc("/v2/analytics/organizations/my-great-org/suites", func(w http.ResponseWriter, r *http.Request) {
		v := new(TestSuiteCreate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, 
			`
			{
				"name" : "Suite 3",
				"default_branch": "main",
				"team_ids": ["8369b300-fff0-4ef1-91de-010f72f4458d"]
			}`)
	})

	suite, _, err := client.TestSuites.Create("my-great-org", input)

	if err != nil {
		t.Errorf("TestSuites.Create returned error: %v", err)
	}

	want := &TestSuite{		
		Name: String("Suite 3"),
		DefaultBranch: String("main"),
	} 

	if !reflect.DeepEqual(suite, want) {
		t.Errorf("TestSuites.Create returned %+v, want %+v", suite, want)
	}
}

func TestTestSuiteService_Update(t *testing.T) {
	setup()
	defer teardown()

	input := &TestSuiteCreate{
		Name:			"Suite 4",   						
		DefaultBranch:	"main",        
		TeamUuids: 		[]string{"818b0849-9718-4898-8de3-42d591a7fe26"},             
	}

	mux.HandleFunc("/v2/analytics/organizations/my-great-org/suites", func(w http.ResponseWriter, r *http.Request) {
		v := new(TestSuiteCreate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, 
			`
			{
				"name" : "Suite 4",
				"default_branch": "main",
				"team_ids": ["818b0849-9718-4898-8de3-42d591a7fe26"],
				"slug": "suite-4"
			}`)
	})

	suite, _, err := client.TestSuites.Create("my-great-org", input)
	
	if err != nil {
		t.Errorf("TestSuites.Create returned error: %v", err)
	}

	// Lets update the default branch to develop
	suite.DefaultBranch = String("develop")

	mux.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-4", func(w http.ResponseWriter, r *http.Request) {
		v := new(TestSuiteCreate)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "PATCH")

		fmt.Fprint(w, 
			`
			{
				"name" : "Suite 4",
				"default_branch": "develop",
				"team_ids": ["818b0849-9718-4898-8de3-42d591a7fe26"],
				"slug": "suite-4"
			}`)
	})

	_, err = client.TestSuites.Update("my-great-org", suite)

	if err != nil {
		t.Errorf("Pipelines.Update returned error: %v", err)
	}

	want := &TestSuite{		
		Name: String("Suite 4"),
		Slug: String("suite-4"),
		DefaultBranch: String("develop"),
	}

	if !reflect.DeepEqual(suite, want) {
		t.Errorf("TestSuites.Update returned %+v, want %+v", suite, want)
	} 	
}

func TestTestSuiteService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-5", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.TestSuites.Delete("my-great-org", "suite-5")

	if err != nil {
		t.Errorf("TestSuites.Delete returned error: %v", err)
	}
}