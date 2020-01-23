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

func TestJobsService_GetJobLog(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/log", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
  "url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sub-keith/builds/awesome-build/jobs/awesome-job-id/log",
  "content": "This is the job's log output",
	"size": 28,
	"header_times": [1563337899810051000,1563337899811015000,1563337905336878000,1563337906589603000,156333791038291900]
}`)
	})

	job, _, err := client.Jobs.GetJobLog("my-great-org", "sup-keith", "awesome-build", "awesome-job-id")
	if err != nil {
		t.Errorf("GetJobLog returned error: %v", err)
	}

	want := &JobLog{
		URL:         String("https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sub-keith/builds/awesome-build/jobs/awesome-job-id/log"),
		Content:     String("This is the job's log output"),
		Size:        Int(28),
		HeaderTimes: []int64{1563337899810051000, 1563337899811015000, 1563337905336878000, 1563337906589603000, 156333791038291900},
	}
	if !reflect.DeepEqual(job, want) {
		t.Errorf("GetJobLog returned %+v, want %+v", job, want)
	}
}
