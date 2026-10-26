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
		ID:          "bitbucket",
		Settings:    &BitbucketSettings{Repository: "my-bitbucket-repo"},
		SettingsMap: ProviderSettingsMap{"repository": "my-bitbucket-repo"},
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
		ID:          "github",
		Settings:    &GitHubSettings{Repository: "my-github-repo"},
		SettingsMap: ProviderSettingsMap{"repository": "my-github-repo"},
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
		ID:          "github_enterprise",
		Settings:    &GitHubEnterpriseSettings{Repository: "my-github-enterprise-repo"},
		SettingsMap: ProviderSettingsMap{"repository": "my-github-enterprise-repo"},
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
		ID:          "gitlab",
		Settings:    &GitLabSettings{Repository: "my-gitlab-repo"},
		SettingsMap: ProviderSettingsMap{"repository": "my-gitlab-repo"},
	}

	if diff := cmp.Diff(provider, want); diff != "" {
		t.Errorf("Unmarshalling GitLab provider JSON produced unexpected output. diff: (-got +want)\n%s", diff)
	}
}

func TestUnmarshalGitLabEnterpriseProvider(t *testing.T) {
	var provider Provider
	err := json.Unmarshal([]byte(`{"id": "gitlab_ee", "settings": {"repository": "my-gitlab-repo"}}`), &provider)
	if err != nil {
		t.Errorf("Error unmarshalling GitLab provider: %v", err)
	}

	want := Provider{
		ID:          "gitlab_ee",
		Settings:    &GitLabSettings{Repository: "my-gitlab-repo"},
		SettingsMap: ProviderSettingsMap{"repository": "my-gitlab-repo"},
	}

	if diff := cmp.Diff(provider, want); diff != "" {
		t.Errorf("Unmarshalling GitLab Enterprise provider JSON produced unexpected output. diff: (-got +want)\n%s", diff)
	}
}

func TestUnmarshalUnknownProvider(t *testing.T) {
	var provider Provider
	err := json.Unmarshal([]byte(`{"id": "unknown", "settings": {"emoji": ":shrug:"}}`), &provider)
	if err != nil {
		t.Errorf("Error unmarshalling unknown provider: %v", err)
	}

	want := Provider{
		ID:          "unknown",
		Settings:    nil,
		SettingsMap: ProviderSettingsMap{"emoji": ":shrug:"},
	}

	if diff := cmp.Diff(provider, want); diff != "" {
		t.Errorf("Unmarshalling unknown provider JSON produced unexpected output. diff: (-got +want)\n%s", diff)
	}
}

func TestProviderSettingsDecodeErrors(t *testing.T) {
	for _, data := range []string{
		`{"id":"github","settings":{"build_branches":"invalid"}}`,
		`{"id":"github","settings":[]}`,
		`{"id":"github"}`,
		`{"id":"github","settings":{"extra":}}`,
	} {
		t.Run(data, func(t *testing.T) {
			var provider Provider
			if err := json.Unmarshal([]byte(data), &provider); err == nil {
				t.Fatal("invalid provider settings decoded without error")
			}
		})
	}
}

func TestProviderSettingsEmptyAndReused(t *testing.T) {
	var provider Provider
	for _, tt := range []struct {
		data string
		want ProviderSettingsMap
	}{
		{`{"id":"github","settings":{"build_branches":true}}`, ProviderSettingsMap{"build_branches": true}},
		{`{"id":"github","settings":{}}`, ProviderSettingsMap{}},
		{`{"id":"github","settings":null}`, nil},
		{`{"id":"unknown","settings":{"extra":false}}`, ProviderSettingsMap{"extra": false}},
		{`{"id":"unknown","settings": []}`, nil},
		{`{"id":"unknown","settings": "opaque"}`, nil},
		{`{"id":"unknown","settings":null}`, nil},
		{`{"id":"unknown"}`, nil},
	} {
		if err := json.Unmarshal([]byte(tt.data), &provider); err != nil {
			t.Fatalf("Unmarshal %s: %v", tt.data, err)
		}
		if diff := cmp.Diff(tt.want, provider.SettingsMap); diff != "" {
			t.Errorf("SettingsMap for %s (-want +got):\n%s", tt.data, diff)
		}
		if provider.ID == "unknown" && provider.Settings != nil {
			t.Errorf("unknown provider has typed settings: %T", provider.Settings)
		}
	}
}
