## [v4.9.1](https://github.com/buildkite/go-buildkite/compare/v4.9.0...v4.9.1) (2025-10-20)

* fix: update cluster maintainers to match REST api for create and list [#251](https://github.com/buildkite/go-buildkite/pull/251) ([wolfeidau](https://github.com/wolfeidau))

## [v4.9.0](https://github.com/buildkite/go-buildkite/compare/v4.8.0...v4.9.0) (2025-10-14)

* Fix URL resolution to preserve base URL path prefixes [#250](https://github.com/buildkite/go-buildkite/pull/250) ([bearcage-dayjob](https://github.com/bearcage-dayjob))

## [v4.8.0](https://github.com/buildkite/go-buildkite/compare/v4.7.1...v4.8.0) (2025-09-22)
[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v4.7.1...v4.8.0)

* feat: unmarshall step data when present in api responses [#245](https://github.com/buildkite/go-buildkite/pull/245) ([mtibben](https://github.com/mtibben))

## [v4.7.1](https://github.com/buildkite/go-buildkite/compare/v4.7.0...v4.7.1) (2025-09-19)

* Setting min go version to 1.23 with toolchain version to 1.24 [#244](https://github.com/buildkite/go-buildkite/pull/244) ([PriyaSudip](https://github.com/PriyaSudip))

## [v4.7.0](https://github.com/buildkite/go-buildkite/compare/v4.6.0...v4.7.0) (2025-08-27)

* chore: code formatting via gofumpt [#242](https://github.com/buildkite/go-buildkite/pull/242) ([mcncl](https://github.com/mcncl))
* Add agent pause and resume functionality [#243](https://github.com/buildkite/go-buildkite/pull/243) ([JoeColeman95](https://github.com/JoeColeman95))

## [v4.6.0](https://github.com/buildkite/go-buildkite/compare/v4.5.1...v4.6.0) (2025-08-22)

* Support Cluster maintainers on create [#240](https://github.com/buildkite/go-buildkite/pull/240) ([mcncl](https://github.com/mcncl))
* support X-Buildkite-Token in headers [#239](https://github.com/buildkite/go-buildkite/pull/239) ([mcncl](https://github.com/mcncl))


## [v4.5.1](https://github.com/buildkite/go-buildkite/tree/v4.5.1) (2025-07-24)

* feat: we are adding filters to the list pipelines REST call [#237](https://github.com/buildkite/go-buildkite/pull/237) (@wolfeidau)

## [v4.5.0](https://github.com/buildkite/go-buildkite/compare/v4.4.0...v4.5.0) (2025-07-22)

* Update access-tokens response type with new data [#235](https://github.com/buildkite/go-buildkite/pull/235) ([moskyb](https://github.com/moskyb))

## [v4.4.0](https://github.com/buildkite/go-buildkite/compare/v4.3.0...v4.4.0) (2025-06-05)

* Fix tests being backwards [#233](https://github.com/buildkite/go-buildkite/pull/233) ([blaknite](https://github.com/blaknite))
* feat: add run result to test runs [#228](https://github.com/buildkite/go-buildkite/pull/228) ([mcncl](https://github.com/mcncl))
* Allow fetching optional failure_expanded from failed_executions API [#232](https://github.com/buildkite/go-buildkite/pull/232) ([blaknite](https://github.com/blaknite))
* Account for fetching test engine data with an optional param [#231](https://github.com/buildkite/go-buildkite/pull/231) ([blaknite](https://github.com/blaknite))
* Add support for test engine run ids [#229](https://github.com/buildkite/go-buildkite/pull/229) ([blaknite](https://github.com/blaknite))
* Add support for getting failed executions [#230](https://github.com/buildkite/go-buildkite/pull/230) ([blaknite](https://github.com/blaknite))

## [v4.3.0](https://github.com/buildkite/go-buildkite/compare/v4.2.0...v4.3.0) (2025-06-04)

* Add support for excluding jobs and pipeline data from List Builds [#225](https://github.com/buildkite/go-buildkite/pull/225) ([blaknite](https://github.com/blaknite))

## [v4.2.0](https://github.com/buildkite/go-buildkite/compare/v4.1.1...v4.2.0) (2025-06-03)

* Export ClientOpt type for programmatic option composition [#217](https://github.com/buildkite/go-buildkite/pull/217) ([prateek](https://github.com/prateek))
* chore(release): CHANGELOG entries for new version [#220](https://github.com/buildkite/go-buildkite/pull/220) ([mcncl](https://github.com/mcncl))

## [v4.1.1](https://github.com/buildkite/go-buildkite/compare/v4.1.0...v4.1.1) (2025-05-21)

* Fix for build author either being a string or a structure [#218](https://github.com/buildkite/go-buildkite/pull/218) ([wolfeidau](https://github.com/wolfeidau))
* List registry packages + delete packages [#216](https://github.com/buildkite/go-buildkite/pull/216) ([moskyb](https://github.com/moskyb))
* Add OIDC policy scopes [#214](https://github.com/buildkite/go-buildkite/pull/214) ([moskyb](https://github.com/moskyb))
* Upgrade Go to 1.24 [#215](https://github.com/buildkite/go-buildkite/pull/215) ([moskyb](https://github.com/moskyb))
* fix: fixes gitlab enterprise which appears in webhooks [#213](https://github.com/buildkite/go-buildkite/pull/213) ([mtibben](https://github.com/mtibben))
* Update golang.org/x/net [#210](https://github.com/buildkite/go-buildkite/pull/210) ([yob](https://github.com/yob))
* Don't use omitempty for bool fields [#209](https://github.com/buildkite/go-buildkite/pull/209) ([moskyb](https://github.com/moskyb))

## [v4.1.0](https://github.com/buildkite/go-buildkite/compare/v4.0.1...v4.1.0) (2024-11-06)

* Support Team Suites and Team Pipelines API [#205](https://github.com/buildkite/go-buildkite/pull/205) ([lizrabuya](https://github.com/lizrabuya))
* Expand the Teams service to support additional endpoints  [#204](https://github.com/buildkite/go-buildkite/pull/204) ([mstifflin](https://github.com/mstifflin))

## [v4.0.0](https://github.com/buildkite/go-buildkite/tree/v4.0.0) (2024-10-21)

Our first major release in a little while! This release is mostly a cleanup release, which __does not contain any major feature releases__.

### Notable Breaking Changes

#### Import Path
The module import path has changed from `github.com/buildkite/go-buildkite/v3/buildkite` to `github.com/buildkite/go-buildkite/v4`. Note that the version number has changed, and that the `buildkite` package is now at the root of the module.

#### Contexts
All methods that interact with Buildkite API now require a context parameter. This is to allow for better control over timeouts and cancellations. For example, to spend a maximum of five seconds trying to get a list of builds for a pipeline, you could run:
```go
ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
builds, _, err := client.Builds.ListByPipeline(ctx, "my-org", "my-pipeline", nil)
```

#### Removal of pointers
Previously, this library had a lot of pointers in its method signatures and struct values. Where practical, we've removed these pointers to make the API clearer and more idiomatic. This means that you'll need to update your code to pass values directly instead of pointers. For example, previously you might have done:
```go
pipeline := buildkite.Pipeline{
  Name: buildkite.String("My Pipeline"),
}
_, _, err := client.Pipelines.Create(ctx, "my-org", &pipeline)
```

Now, you should do:
```go
pipeline := buildkite.Pipeline{
  Name: "My Pipeline",
}
_, _, err := client.Pipelines.Create(ctx, "my-org", pipeline)
```

Along with this change, we've removed the `buildkite.String`, `buildkite.Bool`, and `buildkite.Int` helper functions. You should now pass values directly to the struct fields.

One notable difference in API after this change is that many (but not all!) `Update` methods for various API resources previously unmarshalled their response into a parameter that was passed into them. This is no longer the case, and the response is now returned directly from the method. For example, previously you might have done:
```go
updatePipeline := buildkite.UpdatePipeline{
  Name: buildkite.String("My Pipeline"),
}

_, err := client.Pipelines.Update("my-org", "my-pipeline", &updatePipeline)
// Result of update is unmarshalled into updatePipeline, with surprising results in some cases
```
now, you should do:
```go
updatePipeline := buildkite.UpdatePipeline{
  Name: "My Pipeline",
}
updated, _, err := client.Pipelines.Update(ctx, "my-org", "my-pipeline", updatePipeline)
// updated is the result of the update
// updatePipeline is unchanged
```

#### New Client Creation
We've changed the `buildkite.NewClient` method so that it includes functional args. Previously, you would have done something like this to create a client:
```go
config, err := buildkite.NewTokenConfig(apiToken, debug)
if err != nil {
  return fmt.Errorf("client config failed: %w", err)
}

client := buildkite.NewClient(config.Client())
```
Config creation has been moved inside the `NewClient` method and its functional arguments, so the new equivalent is:
```go
client, err := buildkite.NewClient(
  buildkite.WithTokenAuth(apitoken),
  buildkite.WithHTTPDebug(debug),
)
```

For a full list of functional arguments to `NewClient`, see the [godoc](https://pkg.go.dev/github.com/buildkite/go-buildkite/v4@v4.0.0).

The `NewOpts` method, which was introduced in v3.12.0, remains in place as an alias for `NewClient`.

#### Removal of YAML Bindings
We've removed the YAML bindings for exported data types. This means that marshalling the data types to and from YAML will no longer work. Code like the following will no longer work, and will produce undefined behaviour:
```go
pipeline := buildkite.Pipeline{
  Name: "My Pipeline",
  Repository: "https://github.com/buildkite/go-buildkite",
}

yamlBytes, err := yaml.Marshal(pipeline)
```

These YAML bindings weren't used within the library itself, and so aren't necessary for the operation of this library. If you were marshalling this library's data types to YAML, you should be able to use the `encoding/json` package instead, as all JSON is a subset of YAML.

### Other Changes

#### Changed
- Create packages using presigned uploads, rather than transiting them through the buildkite backend [#194](https://github.com/buildkite/go-buildkite/pull/194) [#195](https://github.com/buildkite/go-buildkite/pull/195) (@moskyb)

#### Internal
- Use golang-ci lint in CI [#199](https://github.com/buildkite/go-buildkite/pull/199) (@moskyb)
- Update README with NewOpts client/kingpin example [#192](https://github.com/buildkite/go-buildkite/pull/192) (@james2791)
- Update tests to use cmp.Diff instead of reflect.DeepEqual [#198](https://github.com/buildkite/go-buildkite/pull/198) (@moskyb)


## [v3.13.0](https://github.com/buildkite/go-buildkite/compare/v3.12.0...v3.13.0) (2024-08-27)
* Add `Name` field to `buildkite.Package` struct [#190](https://github.com/buildkite/go-buildkite/pull/190) ([moskyb](ttps://github.com/moskyb))


## [v3.12.0](https://github.com/buildkite/go-buildkite/compare/v3.11.0...v3.12.0) (2024-08-19)
* Deprecate [`buildkite.NewClient`](https://pkg.go.dev/github.com/buildkite/go-buildkite/v3@v3.11.0/buildkite#NewClient) and its associated authenticating roundtrippers, and replace them with a new function `buildkite.NewOpts` [#185](https://github.com/buildkite/go-buildkite/pull/185) ([moskyb](https://github.com/moskyb))
* Add bindings for Buildkite Packages APIs [#187](https://github.com/buildkite/go-buildkite/pull/186) ([moskyb](https://github.com/moskyb))


## [v3.11.0](https://github.com/buildkite/go-buildkite/compare/v3.10.0...v3.11.0) (2024-02-08)
* Expose retry_source and retry_type on jobs [#171](https://github.com/buildkite/go-buildkite/pull/171) ([drcapulet](https://github.com/drcapulet))
* Expose additional webhook fields on builds and jobs [#173](https://github.com/buildkite/go-buildkite/pull/173) ([mstifflin](https://github.com/mstifflin))
* SUP-1697 Add username to Author struct [#174](https://github.com/buildkite/go-buildkite/pull/174) ([lizrabuya](https://github.com/lizrabuya))
* SUP-1681: Webhook event type definition for struct literal creation [#175](https://github.com/buildkite/go-buildkite/pull/175) ([james2791](https://github.com/james2791))
* Adding job priority as part of job struct [#176](https://github.com/buildkite/go-buildkite/pull/176) ([pankti11](https://github.com/pankti11))

## [v3.10.0](https://github.com/buildkite/go-buildkite/compare/v3.9.0...v3.10.0) (2023-11-15)
* cluster_tokens: expose token string [#168](https://github.com/buildkite/go-buildkite/pull/168) ([gmichelo](https://github.com/gmichelo))

## [v3.9.0](https://github.com/buildkite/go-buildkite/compare/v3.8.0...v3.9.0) (2023-11-14)
* SUP-1448: Pipeline Templates integration [#165](https://github.com/buildkite/go-buildkite/pull/165) ([james2791](https://github.com/james2791))

## [v3.8.0](https://github.com/buildkite/go-buildkite/compare/v3.7.0...v3.8.0) (2023-11-02)
* SUP-1537: Annotations create endpoint interaction [#163](https://github.com/buildkite/go-buildkite/pull/163) ([james2791](https://github.com/james2791))

## [v3.7.0](https://github.com/buildkite/go-buildkite/compare/v3.6.0...v3.7.0) (2023-10-31)
* Add support for `user_id` query string parameter when loading teams [#159](https://github.com/buildkite/go-buildkite/pull/159) ([mwgamble](https://github.com/mwgamble))
* jobs: add cluster_id and cluster_queue_id to job [#160](https://github.com/buildkite/go-buildkite/pull/160) ([gmichelo](https://github.com/gmichelo))

## [v3.6.0](https://github.com/buildkite/go-buildkite/compare/v3.5.0...v3.6.0) (2023-09-18)
* SUP-1353: Clusters integration [#154](https://github.com/buildkite/go-buildkite/pull/154) ([james2791](https://github.com/james2791))
* SUP-1354: Cluster Queues integration [#156](https://github.com/buildkite/go-buildkite/pull/156) ([james2791](https://github.com/james2791))
* SUP-1355: Cluster Tokens integration [#157](https://github.com/buildkite/go-buildkite/pull/157) ([james2791](https://github.com/james2791))
* Stricter type for build meta_data [#155](https://github.com/buildkite/go-buildkite/pull/155) ([erichiggins0](https://github.com/erichiggins0))

## [v3.5.0](https://github.com/buildkite/go-buildkite/compare/v3.4.0...v3.5.0) (2023-09-01)
* Addition of issue/new feature/release templates, Codeowners [#148](https://github.com/buildkite/go-buildkite/pull/148) ([james2791](https://github.com/james2791))
* Add tag fields to pipeline object [#147](https://github.com/buildkite/go-buildkite/pull/147) ([DeanBruntThirdfort](https://github.com/DeanBruntThirdfort))
* SUP-1053: Ability to set `trigger_mode` for Pipelines with GitHub Enterprise ProviderSettings [#149](https://github.com/buildkite/go-buildkite/pull/149) ([james2791](https://github.com/james2791))
* SUP-1410: Addition of unblocked_at to Jobs struct [#152](https://github.com/buildkite/go-buildkite/pull/152) ([james2791](https://github.com/james2791))

## [v3.4.0](https://github.com/buildkite/go-buildkite/compare/v3.3.1...v3.4.0) (2023-08-10)
* Support build.failing events [#141](https://github.com/buildkite/go-buildkite/pull/141) ([mcncl](https://github.com/mcncl))
* SUP-1314: Test Analytics Integration [#142](https://github.com/buildkite/go-buildkite/pull/142) ([james2791](https://github.com/james2791))
* SUP-1321: Pipeline updates utilising dedicated PipelineUpdate struct [#145](https://github.com/buildkite/go-buildkite/pull/145) ([james2791](https://github.com/james2791))

### Notice
As part of this release, the properties that can be set when updating a pipeline using the [PipelinesService](https://github.com/buildkite/go-buildkite/blob/main/buildkite/pipelines.go) have changed to reflect the permitted request body properties in the pipeline update [REST API endpoint](https://buildkite.com/docs/apis/rest-api/pipelines#update-a-pipeline).

## [v3.3.1](https://github.com/buildkite/go-buildkite/compare/v3.3.0...v3.3.1) (2023-06-08)
* Resolved issue on 500 error when request body is null [#137](https://github.com/buildkite/go-buildkite/pull/137) ([lizrabuya](ttps://github.com/lizrabuya))
* Bump to v3.3.1 [#138](https://github.com/buildkite/go-buildkite/pull/138) ([lizrabuya](ttps://github.com/lizrabuya))

## [v3.3.0](https://github.com/buildkite/go-buildkite/compare/v3.2.1...v3.3.0) (2023-06-08)

* update teams docs [#102](https://github.com/buildkite/go-buildkite/pull/102) ([sfunkhouser](https://github.com/sfunkhouser))
* Add clean_checkout to CreateBuild struct [#95](https://github.com/buildkite/go-buildkite/pull/95) ([justingallardo-okta](https://github.com/justingallardo-okta))
* Support filtering list builds call by multiple branches [#96](https://github.com/buildkite/go-buildkite/pull/96) ([gsavit](https://github.com/gsavit))
* Improve README grammar [#110](https://github.com/buildkite/go-buildkite/pull/110) ([mdb](https://github.com/mdb))
* Upgrade to latest go (merge to main) [#125](https://github.com/buildkite/go-buildkite/pull/125) ([lizrabuya](https://github.com/lizrabuya))
* Bump golang.org/x/net from 0.0.0-20190404232315-eb5bcb51f2a3 to 0.7.0 [#124](https://github.com/buildkite/go-buildkite/pull/124) ([dependabot](https://github.com/dependabot))
* Add `rebuilt_from` field to Build struct [#129](https://github.com/buildkite/go-buildkite/pull/129) ([matthiasr](https://github.com/matthiasr))
* add group_key to job struct if available [#132](https://github.com/buildkite/go-buildkite/pull/132) ([alexnguyennn](ttps://github.com/alexnguyennn))
* add visibility field to pipeline create [#133](https://github.com/buildkite/go-buildkite/pull/133) ([sfunkhouser](https://github.com/sfunkhouser))
* Allow filtering builds by meta_data [#127](https://github.com/buildkite/go-buildkite/pull/127) ([andrewhamon](https://github.com/andrewhamon))
* Bump to 3.3.0 (release) [#134](https://github.com/buildkite/go-buildkite/pull/134) ([james2791](https://github.com/james2791))

## [v3.2.1](https://github.com/buildkite/go-buildkite/compare/v3.2.0...v3.2.1) (2023-03-16)

* Change parallel indices type in job struct to int (bug fix) [#121](https://github.com/buildkite/go-buildkite/pull/121) ([yqiu24](https://github.com/yqiu24))

## [v3.2.0](https://github.com/buildkite/go-buildkite/compare/v3.1.0...v3.2.0) (2023-03-14)

* Add GraphQLID to more api objects [#116](https://github.com/buildkite/go-buildkite/pull/116) ([benmoss](https://github.com/benmoss))
* fix "more than one error-wrapping directive %w" [#112](https://github.com/buildkite/go-buildkite/pull/112) ([mdb](https://github.com/mdb))
* Add Label to steps [#117](https://github.com/buildkite/go-buildkite/pull/117) ([benmoss](https://github.com/benmoss))
* Add missing parallel job fields to struct [#119](https://github.com/buildkite/go-buildkite/pull/119) ([yqiu24](https://github.com/yqiu24))
* Update docs to reference pkg.go.dev [#92](https://github.com/buildkite/go-buildkite/pull/92) ([y-yagi](https://github.com/y-yagi))

## [v3.1.0](https://github.com/buildkite/go-buildkite/compare/v3.0.1...v3.1.0) (2023-01-11)

- Add support for plugins as an array [#106](https://github.com/buildkite/go-buildkite/pull/106) ([benmoss](https://github.com/benmoss))
- Add job get environment variables api [#107](https://github.com/buildkite/go-buildkite/pull/107) ([gu-kevin](https://github.com/gu-kevin))

## [v3.0.1](https://github.com/buildkite/go-buildkite/compare/v3.0.0...v3.0.1) (2021-10-26)

- add build_branches param to all provider settings properties [#94](https://github.com/buildkite/go-buildkite/pull/94) ([carol-he](https://github.com/carol-he))

## [v3.0.0](https://github.com/buildkite/go-buildkite/compare/v2.9.0...v3.0.0) (2021-08-20)

- Add block-job label [#90](https://github.com/buildkite/go-buildkite/pull/90) ([BEvgeniyS](https://github.com/BEvgeniyS))
- Fix pagination for teams [#89](https://github.com/buildkite/go-buildkite/pull/89) ([atishpatel](https://github.com/atishpatel))
- Add support for plugins to Pipeline steps [#82](https://github.com/buildkite/go-buildkite/pull/82) ([luma](https://github.com/luma))

## [v2.9.0](https://github.com/buildkite/go-buildkite/compare/v2.8.1...v2.9.0) (2021-08-03)

- Implement Access Tokens API [#87](https://github.com/buildkite/go-buildkite/pull/87) ([alam0rt](https://github.com/alam0rt))
- implement visibility pipeline field [#88](https://github.com/buildkite/go-buildkite/pull/88) ([alam0rt](https://github.com/alam0rt))

## [v2.8.1](https://github.com/buildkite/go-buildkite/compare/v2.8.0...v2.8.1) (2021-07-14)

- add missing organization fields as per api docs [#86](https://github.com/buildkite/go-buildkite/pull/86) ([alam0rt](https://github.com/alam0rt))

## [v2.8.0](https://github.com/buildkite/go-buildkite/compare/v2.7.1...v2.8.0) (2021-06-16)

- Support for Arching/Unarchiving pipelines [#85](https://github.com/buildkite/go-buildkite/pull/85) ([ksepehr](https://github.com/ksepehr))

## [v2.7.1](https://github.com/buildkite/go-buildkite/compare/v2.7.0...v2.7.1) (2021-06-09)

- Adding cluster_id fields to Pipeline and CreatePipeline structs [#84](https://github.com/buildkite/go-buildkite/pull/84) ([ksepehr](https://github.com/ksepehr))

## [v2.7.0](https://github.com/buildkite/go-buildkite/compare/v2.6.0...v2.7.0) (2021-06-08)

- Updated PipelinesService to include AddWebhook [#83](https://github.com/buildkite/go-buildkite/pull/83) ([ksepehr](https://github.com/ksepehr))

## [v2.6.1](https://github.com/buildkite/go-buildkite/compare/v2.6.0...v2.6.1) (2021-03-25)

- Add Author filed to Builds [#80](https://github.com/buildkite/go-buildkite/pull/80) ([dan-embark](https://github.com/dan-embark))

## [v2.6.0](https://github.com/buildkite/go-buildkite/compare/v2.5.1...v2.6.0) (2021-01-05)

- Add soft_failed attribute for job [#77](https://github.com/buildkite/go-buildkite/pull/77) ([qinjin](https://github.com/qinjin))
- add Job unblocker/unblocking fields [#76](https://github.com/buildkite/go-buildkite/pull/76) ([jsleeio](https://github.com/jsleeio))

## [v2.5.1](https://github.com/buildkite/go-buildkite/compare/v2.5.0...v2.5.1) (2020-10-28)

### Changed

- Add missing fields to CreatePipeline, Pipeline, and GitHubSettings structs [#74](https://github.com/buildkite/go-buildkite/pull/74) ([kushmansingh](https://github.com/kushmansingh))

## [v2.5.0](https://github.com/buildkite/go-buildkite/compare/v2.4.0...v2.5.0) (2020-08-20)

### Changed

- Support GitHub Enterprise pipeline settings [#71](https://github.com/buildkite/go-buildkite/pull/71) ([niceking](https://github.com/niceking))
- Bump Go to version 1.15 [#72](https://github.com/buildkite/go-buildkite/pull/72) ([niceking](https://github.com/niceking))

## [v2.4.0](https://github.com/buildkite/go-buildkite/tree/v2.4.0) (2019-02-15)

[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v2.1.0...v2.4.0)

### Changed

- Added support for rebuilding builds [#27](https://github.com/buildkite/go-buildkite/pull/27) (@kevsmith)
- Add typed provider settings [#23](https://github.com/buildkite/go-buildkite/pull/23) (@haines)
- Add blocked to build struct [#30](https://github.com/buildkite/go-buildkite/pull/30) (@jsm)
- Add JobsService with unblockJob method [#31](https://github.com/buildkite/go-buildkite/pull/31) (@dschofie)
- Support Job.runnable_at [#26](https://github.com/buildkite/go-buildkite/pull/26) (@lox)
- Parse AgentQueryRules and TimeoutInMinutes for Steps [#25](https://github.com/buildkite/go-buildkite/pull/25) (@lox)
- Parse the timestamp format webhooks use [#19](https://github.com/buildkite/go-buildkite/pull/19) (@euank)

## [v2.3.1](https://github.com/buildkite/go-buildkite/tree/v2.3.1) (2020-02-19)

[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v2.3.0...v2.3.1)

### Changed

- Go Modules: /v2 suffix on module path [#63](https://github.com/buildkite/go-buildkite/pull/63) (@pda)

## [v2.3.0](https://github.com/buildkite/go-buildkite/tree/v2.3.0) (2020-02-18)

[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v2.2.0...v2.3.0)

### Changed

- Add Annotations API support [#62](https://github.com/buildkite/go-buildkite/pull/62) (@toolmantim)
- Add the `step_key` field to Jobs [#61](https://github.com/buildkite/go-buildkite/pull/61) (@toolmantim)
- Add header times to job log [#59](https://github.com/buildkite/go-buildkite/pull/59) (@kushmansingh)
- run tests on go 1.13 [#58](https://github.com/buildkite/go-buildkite/pull/58) (@yob)
- Fix tests: compile error, segfault. [#57](https://github.com/buildkite/go-buildkite/pull/57) (@pda)
- Go modules; delete vendor/ [#56](https://github.com/buildkite/go-buildkite/pull/56) (@pda)
- License/Copyright tidy [#55](https://github.com/buildkite/go-buildkite/pull/55) (@pda)
- Add functionality for teams endpoint [#54](https://github.com/buildkite/go-buildkite/pull/54) (@pyasi)
- Add client func for the Artifacts ListByJob endpoint [#53](https://github.com/buildkite/go-buildkite/pull/53) (@aashishkapur)
- Add support for the Job ArtifactsURL field [#52](https://github.com/buildkite/go-buildkite/pull/52) (@aashishkapur)
- Added include_retried_jobs for BuildsListOptions [#51](https://github.com/buildkite/go-buildkite/pull/51) (@NorseGaud)
- Add support for getting a job's log output [#48](https://github.com/buildkite/go-buildkite/pull/48) (@y-yagi)
- Fixed spelling errors in comment [#46](https://github.com/buildkite/go-buildkite/pull/46) (@eightseventhreethree)
- Add YAML support to types [#44](https://github.com/buildkite/go-buildkite/pull/44) (@gdhagger)
- Include `source` field for the builds API response [#49](https://github.com/buildkite/go-buildkite/pull/49) (@angulito)
- Add missing AgentListOptions. [#50](https://github.com/buildkite/go-buildkite/pull/50) (@philwo)
- Add filtering option for listing builds by commit [#47](https://github.com/buildkite/go-buildkite/pull/47) (@srmocher)
- Add cancel build function [#45](https://github.com/buildkite/go-buildkite/pull/45) (@dschofie)
- Add pull request build information to CreateBuild [#43](https://github.com/buildkite/go-buildkite/pull/43) (@jradtilbrook)
- Allow updating pipeline provider settings [#41](https://github.com/buildkite/go-buildkite/pull/41) (@gdhagger)
- Explicitly set JSON content-type on API requests [#40](https://github.com/buildkite/go-buildkite/pull/40) (@matthewd)

## [v2.2.0](https://github.com/buildkite/go-buildkite/tree/v2.2.0) (2019-02-20)

[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v2.1.0...v2.2.0)

### Changed

- Update comments to have func name, run gofmt and golint [#36](https://github.com/buildkite/go-buildkite/pull/36) (@charlottestjohn)
- Added support for rebuilding builds [#27](https://github.com/buildkite/go-buildkite/pull/27) (@kevsmith)
- Add typed provider settings [#23](https://github.com/buildkite/go-buildkite/pull/23) (@haines)
- Add blocked to build struct [#30](https://github.com/buildkite/go-buildkite/pull/30) (@jsm)
- Add JobsService with unblockJob method [#31](https://github.com/buildkite/go-buildkite/pull/31) (@dschofie)
- Support Job.runnable_at [#26](https://github.com/buildkite/go-buildkite/pull/26) (@lox)
- Parse AgentQueryRules and TimeoutInMinutes for Steps [#25](https://github.com/buildkite/go-buildkite/pull/25) (@lox)
- Parse the timestamp format webhooks use [#19](https://github.com/buildkite/go-buildkite/pull/19) (@euank)
- Pipeline CRUD [#22](https://github.com/buildkite/go-buildkite/pull/22) (@mubeta06)
- Actually bump version to 2.1.1 [#17](https://github.com/buildkite/go-buildkite/pull/17) (@lox)

## [v2.1.0](https://github.com/buildkite/go-buildkite/tree/v2.1.0) (2017-11-23)

[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v2.0.0...v2.1.0)

### Changed

- Add retry with exp backoff for GET/429 responses [#16](https://github.com/buildkite/go-buildkite/pull/16) (@lox)
- Add Create() to BuildsService [#11](https://github.com/buildkite/go-buildkite/pull/11) (@bshi)
- Update README.md [#15](https://github.com/buildkite/go-buildkite/pull/15) (@lox)
- Fix the build and bump to 1.7.3 [#13](https://github.com/buildkite/go-buildkite/pull/13) (@toolmantim)
- Add Artifacts Service [#10](https://github.com/buildkite/go-buildkite/pull/10) (@pquerna)
- Removed the deprecated/dead featured_build pipeline property [#7](https://github.com/buildkite/go-buildkite/pull/7) (@tobyjoe)
- Added BadgeURL property to Pipeline [#6](https://github.com/buildkite/go-buildkite/pull/6) (@tobyjoe)
- Changed the wolfeidau refs to buildkite [#9](https://github.com/buildkite/go-buildkite/pull/9) (@tobyjoe)

## [v2.0.0](https://github.com/buildkite/go-buildkite/tree/v2.0.0) (2016-03-31)

[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v1.0.0...v2.0.0)

### Changed

- Changes for v2 [#5](https://github.com/buildkite/go-buildkite/pull/5) (@lox)

## [v1.0.0](https://github.com/buildkite/go-buildkite/tree/v1.0.0) (2016-01-10)

[Full Changelog](https://github.com/buildkite/go-buildkite/compare/44957ec...v1.0.0)

- Initial implementation
