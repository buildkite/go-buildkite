package buildkite

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestJobsService_UnblockJob(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/unblock", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
  "id": "awesome-job-id",
  "state": "unblocked"
}`)
	})

	job, _, err := client.Jobs.UnblockJob("my-great-org", "sup-keith", "awesome-build", "awesome-job-id", nil)
	if err != nil {
		t.Errorf("UnblockJob returned error: %v", err)
	}

	want := &Job{ID: String("awesome-job-id"), State: String("unblocked")}
	if !reflect.DeepEqual(job, want) {
		t.Errorf("UnblockJob returned %+v, want %+v", job, want)
	}
}
