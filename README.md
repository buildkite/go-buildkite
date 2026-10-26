# buildkite-go [![Go Reference](https://pkg.go.dev/badge/github.com/buildkite/go-buildkite.svg)](https://pkg.go.dev/github.com/buildkite/go-buildkite/v5) [![Build status](https://badge.buildkite.com/b16a0d730b8732a1cfba06068f8450aa7cc4b2cf40eb6e6717.svg?branch=master)](https://buildkite.com/buildkite/go-buildkite)

A [Go](http://golang.org) library and client for the [Buildkite API](https://buildkite.com/docs/api). This project draws a lot of its structure and testing methods from [go-github](https://github.com/google/go-github).

# Usage

To get the package, execute:

```
go get github.com/buildkite/go-buildkite/v5
```

Simple shortened example for listing all pipelines:

```go
import (
    "github.com/buildkite/go-buildkite/v5"
    "github.com/alecthomas/kingpin/v2"
)

var (
    apiToken = kingpin.Flag("token", "API token").Required().String()
    org = kingpin.Flag("org", "Organization slug").Required().String()
)

client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))

if err != nil {
	log.Fatalf("client config failed: %s", err)
}

pipelines, _, err := client.Pipelines.List(*org, nil)
```

## Test Engine workflows

Use `client.TestWorkflows` to list and inspect workflows within a suite.
Suite paths accept a slug or UUID. Requests require `read_suites`.

Response types include monitor conditions, filters, alarm/recovery action
configuration, identity, and links. Pointer fields distinguish absent fields
from zero thresholds or null labels. Webhook endpoints retain the API's redaction.
Unknown response types remain readable, and unknown response fields are ignored.

```go
workflow, _, err := client.TestWorkflows.Get(ctx, org, suite, workflowID)
```

See the [workflow example](test_workflows_example_test.go) for listing
and inspection. Unavailable or unauthorized requests return errors rather than
an empty workflow list. Workflow writes and event history are not included.

## Migrating to v5 update payloads

Version 5 changes update request structs so PATCH requests can distinguish
"leave this field unchanged" from "send this field with an empty value".

Update fields that can be cleared now use `buildkite.Optional[T]`. The zero
value omits the field. Wrap values with `buildkite.Some(...)` when the field
should be sent, including empty strings, empty maps, empty slices, `0`, or
`false`.

Use initialized empty slices and maps when the API expects `[]` or `{}`. Nil
slices and maps wrapped with `buildkite.Some(...)` are encoded as `null`.

```go
_, _, err := client.Pipelines.Update(ctx, org, pipelineSlug, buildkite.UpdatePipeline{
    Description:              buildkite.Some(""),
    Tags:                     buildkite.Some([]string{}),
    SkipQueuedBranchBuilds:   buildkite.Some(false),
})
```

```go
_, _, err := client.PipelineSchedules.Update(ctx, org, pipelineSlug, scheduleID, buildkite.UpdatePipelineSchedule{
    Env:     buildkite.Some(map[string]string{}),
    Enabled: buildkite.Some(false),
})
```

Leaving an `Optional[T]` unset omits that field from the PATCH body:

```go
_, _, err := client.Clusters.Update(ctx, org, clusterID, buildkite.ClusterUpdate{
    Name: buildkite.Some("macOS builders"),
    // Description is unchanged.
})
```

Use `buildkite.Null[T]()` to send explicit JSON `null` for a nullable field.
For example, `BranchConfiguration: buildkite.Null[string]()` clears the stored
pipeline branch filter. `buildkite.Some("")` still sends an empty string, not null.
Only use `Null` where the API accepts null.

Some create/update request types were split so create calls keep plain required
values while update calls use presence-aware fields:

- `ClusterTokenCreateUpdate` is now `ClusterTokenCreate` and `ClusterTokenUpdate`.
- `PipelineTemplateCreateUpdate` is now `PipelineTemplateCreate` and `PipelineTemplateUpdate`.
- `CreateTeam` is still used for team creation; `UpdateTeam` is now used for team updates.
- `TestSuitesService.Update` now takes `TestSuiteUpdate`.

## GitHub pipeline settings

Use `GitHubSettingsInput` in `CreatePipeline.ProviderSettings` or
`UpdatePipeline.ProviderSettings` to send only the settings you specify.
Unset fields are omitted, `Some(false)` sends false, and `Null[string]()` clears
`FilterCondition` or `PullRequestBranchFilterConfiguration`.

```go
_, _, err := client.Pipelines.Update(ctx, org, pipelineSlug, buildkite.UpdatePipeline{
    BranchConfiguration:       buildkite.Null[string](),
    SkipQueuedBranchBuilds:    buildkite.Some(false),
    CancelRunningBranchBuilds: buildkite.Some(false),
    ProviderSettings: buildkite.Some[buildkite.ProviderSettings](&buildkite.GitHubSettingsInput{
        GitHubActionsMode:                 buildkite.Some(true),
        GitHubWorkflowAccessTokensEnabled: buildkite.Some(true),
        PublishCommitStatus:              buildkite.Some(false),
        TriggerMode:                      buildkite.Some("none"),
    }),
})
```

Pipeline responses still decode GitHub settings as `*GitHubSettings`. The new
GHAC and baseline boolean fields use pointers so absent fields stay absent and
explicit false values survive JSON readback. Existing fields and legacy
`GitHubSettings` writes keep their previous behavior, including sending every
existing boolean. Use `GitHubSettingsInput` for partial writes instead of
resending response settings. Existing response string fields still decode null
as an empty string; use the request type to choose omission, null, or a string.

On create, an omitted `CreatePipeline.BranchConfiguration` uses the API's null
default. On update, use `Null[string]()` to clear an existing branch filter.

These fields represent configuration, not server feature availability. They do
not create a GitHub Actions pipeline trigger or repository webhook. Follow the
[native GHAC setup](https://buildkite.com/docs/pipelines/migration/run-github-actions-workflows)
for the complete settings and separate trigger request. GitHub Enterprise
Server is not supported by native GHAC setup or workflow access tokens.

## Pipeline template assignment

Set `CreatePipeline.PipelineTemplateUUID` to create a pipeline backed by an
existing template. Leave `Configuration` and `Steps` unset; no placeholder YAML
is required. Without a template, supply YAML configuration or visual steps as
appropriate for the pipeline. The server validates template availability and
configuration requirements.

```go
pipeline, _, err := client.Pipelines.Create(ctx, org, buildkite.CreatePipeline{
    Name:                 "Example",
    Repository:           "https://github.com/example/repo.git",
    ClusterID:            clusterID,
    PipelineTemplateUUID: templateUUID,
})
```

On update, use `PipelineTemplateUUID: buildkite.Some(templateUUID)` to assign a
template, or `PipelineTemplateUUID: buildkite.Null[string]()` to remove the
assignment. Leaving the field unset preserves the existing assignment. If
needed, supply replacement `Configuration` in the same update. Pipeline
responses expose `PipelineTemplateUUID` as `*string`; null or absent values
decode as nil.

## Traditional GitHub webhook processing

Use `client.Pipelines.GetGitHubWebhookProcessing(ctx, org, pipelineSlug)` to
read processing state. `EnableGitHubWebhookProcessing` and
`DisableGitHubWebhookProcessing` enable and disable it. All three return
`(GitHubWebhookProcessing, *Response, error)`. The state includes `URL`,
`Enabled`, nullable `DisabledAt` (`*Timestamp`), and nullable `DisabledBy`
(`*string`, the user's name). Repeated enable or disable requests return the
current state.

These methods use GET, PUT, and DELETE on the pipeline's `github-webhooks`
endpoint, with no request body. They require Full Access to the pipeline and
`read_pipelines` for get or `write_pipelines` for enable/disable. Availability
depends on organization enrollment and a GitHub or GitHub Enterprise provider.
Permission/scope failures and availability/provider 404 responses remain API
errors; the client does not treat them as disabled state.

Disabling processing does not delete repository webhooks or disable dedicated
GitHub Actions pipeline triggers. `AddWebhook` remains the separate POST
`webhook` operation that creates a repository webhook.

See the [examples](https://github.com/buildkite/go-buildkite/tree/master/examples) directory for additional examples.

Note: not all API features are supported by `go-buildkite` just yet. If you need a feature, please make an [issue](https://github.com/buildkite/go-buildkite/issues) or submit a pull request.

# Releasing

1. Generate a changelog using [ghch](https://github.com/Songmu/ghch): `ghch --format=markdown --next-version=v<next-version-number>`, and update it in `CHANGELOG.md`
3. Commit the changelog
4. Create a release using GitHub CLI: `gh release create`, ensuring you update the release notes with the new CHANGELOG.md content

# License

This library is distributed under the BSD-style license found in the LICENSE file.
