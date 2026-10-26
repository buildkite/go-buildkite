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

## Provider settings maps

Use `ProviderSettingsMap` for settings not represented by a typed input. It
implements `ProviderSettings`, so existing request fields accept either form:

```go
settings := buildkite.ProviderSettingsMap{
    "build_branches":   false,
    "trigger_mode":     "",
    "filter_condition": nil,
}

created, _, err := client.Pipelines.Create(ctx, org, buildkite.CreatePipeline{
    Name:             "Example",
    Repository:       "https://github.com/example/repo.git",
    ClusterID:        clusterID,
    ProviderSettings: settings,
})

updated, _, err := client.Pipelines.Update(ctx, org, pipelineSlug, buildkite.UpdatePipeline{
    ProviderSettings: buildkite.Some[buildkite.ProviderSettings](settings),
})
```

Only included keys are sent. Omission, `false`, `""`, and `nil` (JSON null)
remain distinct. Values can include nested maps, slices, `json.Number`, and
valid `json.RawMessage` for JSON passthrough. Use `json.Number` or integer types
instead of `float64` when writing numbers that need exact precision. Unsupported
Go values and invalid raw JSON return an encoding error before sending a request.

Leave the request field unset to omit `provider_settings`. An initialized empty
map sends `{}`; a nil map supplied to either request sends `null`. These values
do not imply that the API clears all settings. The server determines which
properties, values, and providers it accepts; the map does not bypass validation
or make a server feature available.

Pipeline responses expose complete settings objects in `Provider.SettingsMap`,
including unknown properties and unknown providers. All map numbers, including
numbers inside nested objects and arrays, decode as `json.Number`. Missing keys
remain absent, JSON null becomes `nil`, and empty objects remain empty maps.
For unknown providers, `Provider.Settings` remains nil. Known providers retain
their existing typed settings and decoding errors. Unknown providers with
non-object settings remain accepted, but their `SettingsMap` is nil.

`Settings` and `SettingsMap` are independent snapshots. Editing one does not
update the other. `SettingsMap` is excluded from JSON serialization; marshaling
a provider still uses only the legacy typed `Settings` view and does not retain
unknown properties. Marshal `SettingsMap` itself if you need the complete object.
The client does not merge these views or automatically resend response settings.
Build a write map with only intended writable keys, not the full response map,
which can contain read-only properties such as `repository`.

See the [examples](https://github.com/buildkite/go-buildkite/tree/master/examples) directory for additional examples.

Note: not all API features are supported by `go-buildkite` just yet. If you need a feature, please make an [issue](https://github.com/buildkite/go-buildkite/issues) or submit a pull request.

# Releasing

1. Generate a changelog using [ghch](https://github.com/Songmu/ghch): `ghch --format=markdown --next-version=v<next-version-number>`, and update it in `CHANGELOG.md`
3. Commit the changelog
4. Create a release using GitHub CLI: `gh release create`, ensuring you update the release notes with the new CHANGELOG.md content

# License

This library is distributed under the BSD-style license found in the LICENSE file.
