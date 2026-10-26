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
// Use GitHubSettingsInput for partial writes that leave unspecified settings unchanged.
type GitHubSettings struct {
	GitHubActionsMode                       *bool  `json:"github_actions_mode,omitempty"`
	GitHubWorkflowAccessTokensEnabled       *bool  `json:"github_workflow_access_tokens_enabled,omitempty"`
	IgnoreDefaultBranchPullRequests         *bool  `json:"ignore_default_branch_pull_requests,omitempty"`
	SkipBuildsForClosedPullRequests         *bool  `json:"skip_builds_for_closed_pull_requests,omitempty"`
	SkipBuildsForExistingCommits            *bool  `json:"skip_builds_for_existing_commits,omitempty"`
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

// GitHubSettingsInput configures GitHub provider settings on pipeline creation
// or update. Unset fields are omitted; Some(false) explicitly disables a setting.
// Null[string]() clears FilterCondition or PullRequestBranchFilterConfiguration.
// GitHub Actions settings require server-side availability and do not provision
// a pipeline trigger or repository webhook.
type GitHubSettingsInput struct {
	GitHubActionsMode                       Optional[bool]   `json:"github_actions_mode,omitzero"`
	GitHubWorkflowAccessTokensEnabled       Optional[bool]   `json:"github_workflow_access_tokens_enabled,omitzero"`
	IgnoreDefaultBranchPullRequests         Optional[bool]   `json:"ignore_default_branch_pull_requests,omitzero"`
	SkipBuildsForClosedPullRequests         Optional[bool]   `json:"skip_builds_for_closed_pull_requests,omitzero"`
	SkipBuildsForExistingCommits            Optional[bool]   `json:"skip_builds_for_existing_commits,omitzero"`
	TriggerMode                             Optional[string] `json:"trigger_mode,omitzero"`
	BuildPullRequests                       Optional[bool]   `json:"build_pull_requests,omitzero"`
	BuildBranches                           Optional[bool]   `json:"build_branches,omitzero"`
	PullRequestBranchFilterEnabled          Optional[bool]   `json:"pull_request_branch_filter_enabled,omitzero"`
	PullRequestBranchFilterConfiguration    Optional[string] `json:"pull_request_branch_filter_configuration,omitzero"`
	SkipPullRequestBuildsForExistingCommits Optional[bool]   `json:"skip_pull_request_builds_for_existing_commits,omitzero"`
	BuildPullRequestForks                   Optional[bool]   `json:"build_pull_request_forks,omitzero"`
	PrefixPullRequestForkBranchNames        Optional[bool]   `json:"prefix_pull_request_fork_branch_names,omitzero"`
	BuildTags                               Optional[bool]   `json:"build_tags,omitzero"`
	PublishCommitStatus                     Optional[bool]   `json:"publish_commit_status,omitzero"`
	PublishCommitStatusPerStep              Optional[bool]   `json:"publish_commit_status_per_step,omitzero"`
	FilterEnabled                           Optional[bool]   `json:"filter_enabled,omitzero"`
	FilterCondition                         Optional[string] `json:"filter_condition,omitzero"`
	SeparatePullRequestStatuses             Optional[bool]   `json:"separate_pull_request_statuses,omitzero"`
	PublishBlockedAsPending                 Optional[bool]   `json:"publish_blocked_as_pending,omitzero"`
}

func (s *GitHubSettingsInput) isProviderSettings() {}

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
	CancelDeletedBranchBuilds               bool   `json:"cancel_deleted_branch_builds"`

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
