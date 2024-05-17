package buildkite

import "fmt"

type TestPlansService struct {
	client *Client
}

type Case struct {
	Path              *string `json:"path,omitempty" yaml:"path,omitempty"`
	EstimatedDuration *int `json:"estimated_duration,omitempty" yaml:"estimated_duration,omitempty"`
}

type Tests struct {
	Format *string `json:"format,omitempty" yaml:"format,omitempty"`
	Cases  []Case  `json:"cases,omitempty" yaml:"cases,omitempty"`
}

type Task struct {
	NodeNumber *int   `json:"node_number,omitempty" yaml:"node_number,omitempty"`
	Tests      *Tests `json:"tests,omitempty" yaml:"tests,omitempty"`
}

type TestPlan struct {
	Identifier  *string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	Parallelism *int    `json:"parallelism,omitempty" yaml:"paralellism,omitempty"`
	Tasks       []Task  `json:"tasks,omitempty" yaml:"tasks,omitempty"`
}

func (tps *TestPlansService) Get(org, slug string, identifier string) (*TestPlan, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/test_plan?identifier=%s", org, slug, identifier)

	req, err := tps.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, err
	}

	testPlan := new(TestPlan)

	resp, err := tps.client.Do(req, testPlan)

	if err != nil {
		return nil, resp, err
	}

	return testPlan, resp, err
}
