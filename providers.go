package buildkite

import "encoding/json"

// Provider represents a source code provider. It is read-only, but settings may be written using Pipeline.ProviderSettings.
type Provider struct {
	ID         string           `json:"id"`
	WebhookURL string           `json:"webhook_url"`
	Settings   ProviderSettings `json:"settings"`
}

// UnmarshalJSON decodes the Provider, choosing the type of the Settings from the ID.
func (p *Provider) UnmarshalJSON(data []byte) error {
	type provider Provider
	var v struct {
		provider
		Settings json.RawMessage `json:"settings"`
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*p = Provider(v.provider)

	var settings ProviderSettings
	switch v.ID {
	case "bitbucket":
		settings = &BitbucketSettings{}
	case "github":
		settings = &GitHubSettings{}
	case "github_enterprise":
		settings = &GitHubEnterpriseSettings{}
	case "gitlab", "gitlab_ee":
		settings = &GitLabSettings{}
	default:
		return nil
	}

	err = json.Unmarshal(v.Settings, settings)
	if err != nil {
		return err
	}
	p.Settings = settings

	return nil
}

// ProviderSettings represents the sum type of the settings for different source code providers.
type ProviderSettings interface {
	isProviderSettings()
}

// BitbucketSettings are settings for pipelines building from Bitbucket repositories.
type BitbucketSettings struct {
	BuildPullRequests                       bool   `json:"build_pull_requests"`
	BuildBranches                           bool   `json:"build_branches"`
	PullRequestBranchFilterEnabled          bool   `json:"pull_request_branch_filter_enabled"`
	PullRequestBranchFilterConfiguration    string `json:"pull_request_branch_filter_configuration,omitempty"`
	SkipPullRequestBuildsForExistingCommits bool   `json:"skip_pull_request_builds_for_existing_commits"`
	BuildTags                               bool   `json:"build_tags"`
	PublishCommitStatus                     bool   `json:"publish_commit_status"`
	PublishCommitStatusPerStep              bool   `json:"publish_commit_status_per_step"`

	// Read-only
	Repository string `json:"repository,omitempty"`
}

func (s *BitbucketSettings) isProviderSettings() {}

// GitHubSettings are settings for pipelines building from GitHub repositories.
type GitHubSettings struct {
	TriggerMode                             string `json:"trigger_mode,omitempty"`
	BuildPullRequests                       bool   `json:"build_pull_requests"`
	BuildBranches                           bool   `json:"build_branches"`
	PullRequestBranchFilterEnabled          bool   `json:"pull_request_branch_filter_enabled"`
	PullRequestBranchFilterConfiguration    string `json:"pull_request_branch_filter_configuration,omitempty"`
	SkipPullRequestBuildsForExistingCommits bool   `json:"skip_pull_request_builds_for_existing_commits"`
	BuildPullRequestForks                   bool   `json:"build_pull_request_forks"`
	PrefixPullRequestForkBranchNames        bool   `json:"prefix_pull_request_fork_branch_names"`
	BuildTags                               bool   `json:"build_tags"`
	PublishCommitStatus                     bool   `json:"publish_commit_status"`
	PublishCommitStatusPerStep              bool   `json:"publish_commit_status_per_step"`
	FilterEnabled                           bool   `json:"filter_enabled"`
	FilterCondition                         string `json:"filter_condition,omitempty"`
	SeparatePullRequestStatuses             bool   `json:"separate_pull_request_statuses"`
	PublishBlockedAsPending                 bool   `json:"publish_blocked_as_pending"`

	// Read-only
	Repository string `json:"repository,omitempty"`
}

func (s *GitHubSettings) isProviderSettings() {}

// GitHubEnterpriseSettings are settings for pipelines building from GitHub Enterprise repositories.
type GitHubEnterpriseSettings struct {
	TriggerMode                             string `json:"trigger_mode,omitempty"`
	BuildPullRequests                       bool   `json:"build_pull_requests"`
	BuildBranches                           bool   `json:"build_branches"`
	PullRequestBranchFilterEnabled          bool   `json:"pull_request_branch_filter_enabled"`
	PullRequestBranchFilterConfiguration    string `json:"pull_request_branch_filter_configuration,omitempty"`
	SkipPullRequestBuildsForExistingCommits bool   `json:"skip_pull_request_builds_for_existing_commits"`
	BuildTags                               bool   `json:"build_tags"`
	PublishCommitStatus                     bool   `json:"publish_commit_status"`
	PublishCommitStatusPerStep              bool   `json:"publish_commit_status_per_step"`

	// Read-only
	Repository string `json:"repository,omitempty"`
}

func (s *GitHubEnterpriseSettings) isProviderSettings() {}

// GitLabSettings are settings for pipelines building from GitLab repositories.
type GitLabSettings struct {
	// Read-only
	FilterEnabled bool   `json:"filter_enabled"`
	Repository    string `json:"repository,omitempty"`
}

func (s *GitLabSettings) isProviderSettings() {}
