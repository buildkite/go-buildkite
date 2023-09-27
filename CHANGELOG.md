## Unreleased
* Add support for `user_id` query string parameter when loading teams

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
