package buildkite

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalBitbucketProvider(t *testing.T) {
	var provider Provider
	err := json.Unmarshal([]byte(`{"id": "bitbucket", "settings": {"repository": "my-bitbucket-repo"}}`), &provider)
	if err != nil {
		t.Errorf("Error unmarshalling Bitbucket provider: %v", err)
	}

	want := Provider{
		ID:       "bitbucket",
		Settings: &BitbucketSettings{Repository: String("my-bitbucket-repo")},
	}

	if !reflect.DeepEqual(provider, want) {
		t.Errorf("Failed to unmarshal Bitbucket provider: got %+v, want %+v", provider, want)
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
		Settings: &GitHubSettings{Repository: String("my-github-repo")},
	}

	if !reflect.DeepEqual(provider, want) {
		t.Errorf("Failed to unmarshal GitHub provider: got %+v, want %+v", provider, want)
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
		Settings: &GitHubEnterpriseSettings{Repository: String("my-github-enterprise-repo")},
	}

	if !reflect.DeepEqual(provider, want) {
		t.Errorf("Failed to unmarshal GitHub Enterprise provider: got %+v, want %+v", provider, want)
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
		Settings: &GitLabSettings{Repository: String("my-gitlab-repo")},
	}

	if !reflect.DeepEqual(provider, want) {
		t.Errorf("Failed to unmarshal GitLab provider: got %+v, want %+v", provider, want)
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

	if !reflect.DeepEqual(provider, want) {
		t.Errorf("Failed to unmarshal unknown provider: got %+v, want %+v", provider, want)
	}
}
