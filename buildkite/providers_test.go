package buildkite

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshalBitbucketProvider(t *testing.T) {
	var provider Provider
	err := json.Unmarshal([]byte(`{"id": "bitbucket", "settings": {"repository": "my-bitbucket-repo"}}`), &provider)
	if err != nil {
		t.Errorf("Error unmarshalling Bitbucket provider: %v", err)
	}

	want := Provider{
		ID:       "bitbucket",
		Settings: &BitbucketSettings{Repository: "my-bitbucket-repo"},
	}

	if diff := cmp.Diff(provider, want); diff != "" {
		t.Errorf("Unmarshalling Bitbucket provider JSON produced unexpected output. diff: (-got +want)\n%s", diff)
	}
}

func TestUnmarshalGitHubProvider(t *testing.T) {
	var provider Provider
	err := json.Unmarshal([]byte(`{"id": "github", "settings": {"repository": "my-github-repo"}}`), &provider)
	if err != nil {
		t.Errorf("Error unmarshalling GitHub provider: %v", err)
	}

	want := Provider{
		ID:       "github",
		Settings: &GitHubSettings{Repository: "my-github-repo"},
	}

	if diff := cmp.Diff(provider, want); diff != "" {
		t.Errorf("Unmarshalling GitHub provider JSON produced unexpected output. diff: (-got +want)\n%s", diff)
	}
}

func TestUnmarshalGitHubEnterpriseProvider(t *testing.T) {
	var provider Provider
	err := json.Unmarshal([]byte(`{"id": "github_enterprise", "settings": {"repository": "my-github-enterprise-repo"}}`), &provider)
	if err != nil {
		t.Errorf("Error unmarshalling GitHub Enterprise provider: %v", err)
	}

	want := Provider{
		ID:       "github_enterprise",
		Settings: &GitHubEnterpriseSettings{Repository: "my-github-enterprise-repo"},
	}

	if diff := cmp.Diff(provider, want); diff != "" {
		t.Errorf("Unmarshalling GitHub Enterprise provider JSON produced unexpected output. diff: (-got +want)\n%s", diff)
	}
}

func TestUnmarshalGitLabProvider(t *testing.T) {
	var provider Provider
	err := json.Unmarshal([]byte(`{"id": "gitlab", "settings": {"repository": "my-gitlab-repo"}}`), &provider)
	if err != nil {
		t.Errorf("Error unmarshalling GitLab provider: %v", err)
	}

	want := Provider{
		ID:       "gitlab",
		Settings: &GitLabSettings{Repository: "my-gitlab-repo"},
	}

	if diff := cmp.Diff(provider, want); diff != "" {
		t.Errorf("Unmarshalling GitLab provider JSON produced unexpected output. diff: (-got +want)\n%s", diff)
	}
}

func TestUnmarshalUnknownProvider(t *testing.T) {
	var provider Provider
	err := json.Unmarshal([]byte(`{"id": "unknown", "settings": {"emoji": ":shrug:"}}`), &provider)
	if err != nil {
		t.Errorf("Error unmarshalling unknown provider: %v", err)
	}

	want := Provider{
		ID:       "unknown",
		Settings: nil,
	}

	if diff := cmp.Diff(provider, want); diff != "" {
		t.Errorf("Unmarshalling unknown provider JSON produced unexpected output. diff: (-got +want)\n%s", diff)
	}
}
