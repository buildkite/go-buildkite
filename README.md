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

## Migrating to v5 update payloads

Version 5 changes update request structs so PATCH requests can distinguish
"leave this field unchanged" from "send this field with an empty value".

Update fields that can be cleared now use `buildkite.Optional[T]`. The zero
value omits the field. Wrap values with `buildkite.Some(...)` when the field
should be sent, including empty strings, empty maps, empty slices, `0`, or
`false`.

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

Some create/update request types were split so create calls keep plain required
values while update calls use presence-aware fields:

- `ClusterTokenCreateUpdate` is now `ClusterTokenCreate` and `ClusterTokenUpdate`.
- `PipelineTemplateCreateUpdate` is now `PipelineTemplateCreate` and `PipelineTemplateUpdate`.
- `CreateTeam` is still used for team creation; `UpdateTeam` is now used for team updates.
- `TestSuitesService.Update` now takes `TestSuiteUpdate`.

See the [examples](https://github.com/buildkite/go-buildkite/tree/master/examples) directory for additional examples.

Note: not all API features are supported by `go-buildkite` just yet. If you need a feature, please make an [issue](https://github.com/buildkite/go-buildkite/issues) or submit a pull request.

# Releasing

1. Generate a changelog using [ghch](https://github.com/Songmu/ghch): `ghch --format=markdown --next-version=v<next-version-number>`, and update it in `CHANGELOG.md`
3. Commit the changelog
4. Create a release using GitHub CLI: `gh release create`, ensuring you update the release notes with the new CHANGELOG.md content

# License

This library is distributed under the BSD-style license found in the LICENSE file.
